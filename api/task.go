// 任务接口

package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// 任务表
type Task struct {
	gorm.Model
	// 等同于RPA type
	RPAName string `gorm:"size:255"`
	Input   json.RawMessage
	Output  json.RawMessage
	// 可以是waiting（等待执行） running（执行中）success（执行成功） error（执行失败）
	State string `gorm:"type:varchar(10);default:'waiting'"`
}

// 内存任务
var tasks = make([]Task, 0)

// 创建任务
func create_task(c *gin.Context) {
	type inputS struct {
		RPAName string          `form:"RPAName" json:"RPAName" binding:"required"`
		Input   json.RawMessage `form:"Input" json:"Input" binding:"required"`
	}

	var inputData inputS

	if err := c.ShouldBind(&inputData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	// 将input转换为Map
	var inputMap map[string]interface{}
	if err := json.Unmarshal(inputData.Input, &inputMap); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "解析输入数据失败"})
		return
	}

	// 1. 获取该RPA的相关数据
	rpaData, err := get_rpa_by_name(inputData.RPAName)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "创建任务失败", "error": err.Error()})
		return
	}

	// 2. 检查变量
	varData, err := get_var_by_rpa_name(inputData.RPAName)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "创建任务失败，未能正常读取对应变量"})
		return
	} else {
		if varData == nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "创建任务失败，该任务类型不存在"})
			return
		} else {

			lackVars := ""
			// 需要的所有变量
			for _, item := range varData {
				// 如果该变量非必填，则跳过
				if !item.Required {
					continue
				}
				found := false
				// 输入的所有变量
				for key, value := range inputMap {
					// 找到对应的字段
					if key == item.VarName && value != nil && value != "" {
						found = true
						break
					}
				}

				// 如果用户输入的该字段为空
				if !found {
					// 如果该变量有默认值，则使用默认值
					if item.Default != "" {
						inputMap[item.VarName] = item.Default
					} else {
						// 如果该变量没有默认值，则报错
						lackVars += item.VarRemark + ","
					}
				}
			}
			if lackVars != "" {
				c.JSON(http.StatusBadRequest, gin.H{"message": lackVars + "不可为空"})
				return
			}
		}
	}

	// 3. 如果任务自发启动，不需要监听
	if rpaData.Spont {
		do_spont_task(inputData.RPAName, inputMap)
		return
	}

	// 4. 创建任务到数据库
	task := &Task{
		RPAName: inputData.RPAName,
		Input:   inputData.Input,
	}

	result := db.Create(task)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "创建任务失败"})
		return
	}

	// 5. 将任务添加到内存
	tasks = append(tasks, *task)

	// 6. 将更新后的数据重新编码为JSON
	updatedInput, err := json.Marshal(inputMap)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "更新输入数据失败"})
		return
	}

	inputData.Input = updatedInput

	// 7. 判断任务是否即刻执行
	// 如果任务需要即刻执行
	if rpaData.Now {
		var watchTask Task
		forTime := time.Now()
		// 循环检查任务状态
		for {
			result := db.First(&watchTask, task.ID)
			if result.Error == nil {
				// 直到任务完成
				if watchTask.State == "success" {
					c.JSON(200, watchTask.Output)
					return
				} else if watchTask.State == "error" {
					c.JSON(400, watchTask.Output)
					return
				}
			}
			// 假如任务超时 100秒
			if time.Since(forTime).Seconds() > 100 {
				c.JSON(http.StatusRequestTimeout, gin.H{
					"message": "任务执行超时",
				})
				return
			}
			time.Sleep(time.Second * 1)
		}

		// 如果任务排队执行
	} else {
		c.JSON(http.StatusOK, gin.H{
			"message": "任务创建完毕",
			"taskID":  task.ID,
		})
	}
}

// 查询任务，从内存
func get_task_memory(c *gin.Context) {
	taskID := c.Query("id")
	// 如果id为空，从内存，返回全部任务
	if taskID == "" {
		c.JSON(http.StatusOK, gin.H{
			"message": "ok",
			"data":    tasks,
		})

	} else {
		// 将taskID转换为uint
		id, err := strconv.ParseUint(taskID, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "任务ID必须是有效的数字",
			})
			return
		}

		// 在tasks切片中查找匹配的任务
		var foundTask *Task
		for _, task := range tasks {
			if task.ID == uint(id) {
				foundTask = &task
				break
			}
		}

		if foundTask != nil {
			c.JSON(http.StatusOK, gin.H{
				"message": "ok",
				"data":    foundTask,
			})
		} else {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "未找到指定ID的任务",
			})
		}
	}

}

// 查询任务，从数据库
func get_task_db(c *gin.Context) {
	taskID := c.Query("id")
	// 如果id为空，从内存，返回全部任务
	if taskID == "" {
		state := c.Query("state")
		var tasks []Task
		var result *gorm.DB
		if state == "" {
			result = db.Find(&tasks)
		} else {
			result = db.Where("state = ?", state).Find(&tasks)
		}
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "获取任务列表失败",
				"error":   result.Error,
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message": "ok",
			"data":    tasks,
		})
	} else {
		var task Task

		result := db.First(&task, taskID)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "任务不存在",
				"error":   result.Error,
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message": "ok",
			"data":    task,
		})
	}
}

// 任务开始执行，修改数据状态
func ing_task(c *gin.Context) {
	taskID := c.Query("id")
	// 将taskID转换为uint
	id, err := strconv.ParseUint(taskID, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "无效的任务ID",
			"error":   err.Error(),
		})
		return
	}

	// 从内存中删除任务
	for i, task := range tasks {
		if task.ID == uint(id) {
			tasks = append(tasks[:i], tasks[i+1:]...)
			break
		}
	}

	// 在数据库中更新任务状态
	result := db.Model(&Task{}).Where("id = ?", id).Update("state", "running")
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "更新任务状态失败",
			"error":   result.Error.Error(),
		})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "未找到指定ID的任务",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "任务状态已更新为进行中",
	})
}

type OutputData struct {
	Code int
	Msg  string
	Data string
}

// 完成任务时，修改数据库状态
func done_task(c *gin.Context) {
	type inputType struct {
		TaskID uint64 `form:"id" json:"id" binding:"required"`
		Code   int    `form:"code" json:"code" default:"500"`
		Msg    string `form:"msg" json:"msg" `
		Data   string `form:"data" json:"data" `
	}
	var inputData inputType

	if err := c.ShouldBind(&inputData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	var output = &OutputData{Code: inputData.Code, Msg: inputData.Msg, Data: inputData.Data}

	// 更新任务的Output字段和状态
	state := "error"
	if inputData.Code == 200 {
		state = "success"
	}

	outputJSON, err := json.Marshal(output)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "序列化输出失败", "error": err.Error()})
		return
	}

	result := db.Model(&Task{}).Where("id = ?", inputData.TaskID).Updates(map[string]interface{}{
		"output": json.RawMessage(outputJSON),
		"state":  state,
	})

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "更新数据库状态失败",
			"error":   result.Error.Error(),
		})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "未找到指定ID的任务",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "任务Output已成功更新",
	})
}

// 删除任务，从数据库
func del_task(c *gin.Context) {
	taskID := c.Query("id")
	// 将taskID转换为uint
	id, err := strconv.ParseUint(taskID, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "无效的任务ID",
			"error":   err.Error(),
		})
		return
	}

	// 从数据库中删除任务
	result := db.Unscoped().Delete(&Task{}, id)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "删除失败，" + result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

// 删除任务，从内存
func del_task_memory(c *gin.Context) {
	taskID := c.DefaultQuery("id", "")
	taskIDUint, err := strconv.ParseUint(taskID, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "请输入正确的任务ID"})
		return
	}

	// 从内存中删除任务
	for i, task := range tasks {
		if task.ID == uint(taskIDUint) {
			tasks = append(tasks[:i], tasks[i+1:]...)
			break
		}
	}
	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

// 统计任务数据
func count_task(c *gin.Context) {
	var totalCount int64
	var waitingCount int64
	var runningCount int64
	var successCount int64
	var errorCount int64

	// 统计总任务数
	db.Model(&Task{}).Count(&totalCount)

	// 统计不同状态的任务数
	db.Model(&Task{}).Where("state = ?", "waiting").Count(&waitingCount)
	db.Model(&Task{}).Where("state = ?", "running").Count(&runningCount)
	db.Model(&Task{}).Where("state = ?", "success").Count(&successCount)
	db.Model(&Task{}).Where("state = ?", "error").Count(&errorCount)

	c.JSON(http.StatusOK, gin.H{
		"message": "任务统计数据",
		"data": gin.H{
			"all":     totalCount,
			"waiting": waitingCount,
			"running": runningCount,
			"success": successCount,
			"error":   errorCount,
		},
	})
}

// 定时清理任务记录 本地执行
// 为什么不做成接口：做成接口谁来读取？
func clean_task_func() (string, error) {
	// 删除超过24小时且状态为waiting的任务记录
	deleteResult := db.Where("state = ? AND created_at < ?", "waiting", time.Now().Add(-24*time.Hour)).Delete(&Task{})
	deletedCount := deleteResult.RowsAffected

	// 将超过24小时且状态为running的任务记录更新为error状态
	updateResult := db.Model(&Task{}).
		Where("state = ? AND created_at < ?", "running", time.Now().Add(-24*time.Hour)).
		Updates(map[string]interface{}{
			"state":  "error",
			"output": json.RawMessage(`{"Code":500,"Msg":"任务执行超时","Data":null}`),
		})
	updatedCount := updateResult.RowsAffected

	// 记录日志
	logMessage := fmt.Sprintf("清理任务记录：删除了 %d 条等待状态的过期任务，更新了 %d 条执行中的超时任务", deletedCount, updatedCount)
	return logMessage, nil
}

// 自发执行任务
func do_spont_task(rpaName string, inputMap map[string]interface{}) {
	var output OutputData
	var state string
	if rpaName == "clean_task" {
		message, err := clean_task_func()
		if err != nil {
			output.Code = 500
			output.Msg = err.Error()
			state = "error"
		}
		output.Code = 200
		output.Msg = message
		state = "success"
	}

	// 将结果保存到数据库
	inputJSON, _ := json.Marshal(inputMap)
	outputJSON, _ := json.Marshal(output)

	task := &Task{
		RPAName: rpaName,
		Input:   json.RawMessage(inputJSON),
		Output:  json.RawMessage(outputJSON),
		State:   state,
	}

	db.Create(task)
}
