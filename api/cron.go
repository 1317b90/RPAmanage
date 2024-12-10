package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"

	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-co-op/gocron/v2"
	"github.com/google/uuid"
)

var S, _ = gocron.NewScheduler()

// 具体执行定时的操作
// 逻辑：
// 1. 如果RPAName为clean_task，则执行clean_task_func()
// 2. 如果Input为weda，则从微搭获取数据

func doFunc(RPAName string, Input string) {
	var inputMap map[string]interface{}

	// 从微搭获取数据
	if Input == "weda" {
		req, _ := http.NewRequest("GET", "/weda/dsj", nil)
		result := httptest.NewRecorder()

		// 使用Gin实例处理请求
		R.ServeHTTP(result, req)

		// 解析响应体
		var response struct {
			Data []map[string]interface{} `json:"data"`
		}
		if err := json.Unmarshal(result.Body.Bytes(), &response); err != nil {
			fmt.Println("解析响应失败:", err)
			return
		}
	} else if Input == "" {
		// 不做任何操作
	} else {
		// 解析json
		err := json.Unmarshal([]byte(Input), &inputMap)
		if err != nil {
			add_log_func("定时任务", RPAName, "输入JSON数据有误:"+err.Error(), false)
			return
		}
	}

	// 使用Gin给自己发送请求
	// 准备请求数据
	requestData := map[string]interface{}{
		"RPAName": RPAName,
		"Input":   inputMap,
	}
	jsonData, err := json.Marshal(requestData)
	if err != nil {
		add_log_func("定时任务", RPAName, "输入JSON数据有误:"+err.Error(), false)
		return
	}

	// 创建任务
	req, err := http.NewRequest("POST", "/task", bytes.NewBuffer(jsonData))
	if err != nil {
		add_log_func("定时任务", RPAName, "创建请求失败:"+err.Error(), false)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	result := httptest.NewRecorder()
	R.ServeHTTP(result, req)

	// 解析响应
	var response map[string]interface{}
	json.Unmarshal(result.Body.Bytes(), &response)

	// 获取message和状态
	message := response["message"].(string)
	state := result.Code == http.StatusOK

	// 记录日志
	add_log_func("定时任务", RPAName, message, state)

}

// 创建定时任务
func create_cron(c *gin.Context) {
	type inputS struct {
		Name     string `form:"Name" json:"Name" binding:"required"`
		RPAName  string `form:"RPAName" json:"RPAName" binding:"required"`
		Type     string `form:"Type" json:"Type" default:"cron"`
		CronData string `form:"CronData" json:"CronData" default:"1 * * * *"` // 指定cron：通配符 每个月执行：几号执行 只执行一次：“”
		Input    string `form:"Input" json:"Input"`
	}

	var inputData inputS

	if err := c.ShouldBind(&inputData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	var cronType gocron.JobDefinition

	// 指定cron
	if inputData.Type == "cron" {
		cronType = gocron.CronJob(
			inputData.CronData,
			false,
		)
		// 每个月某天执行
	} else if inputData.Type == "month" {
		day, err := strconv.Atoi(inputData.CronData)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "无效的日期"})
			return
		}

		cronType = gocron.MonthlyJob(
			1,
			gocron.NewDaysOfTheMonth(day),
			gocron.NewAtTimes(
				gocron.NewAtTime(0, 0, 0),
			),
		)
		// 只执行一次
	} else {
		cronType = gocron.OneTimeJob(
			gocron.OneTimeJobStartImmediately(),
		)
	}

	j, err := S.NewJob(
		cronType,
		gocron.NewTask(
			func() { doFunc(inputData.RPAName, inputData.Input) },
		),
		gocron.WithName(inputData.Name+"@"+inputData.RPAName+"@"+inputData.Type+"@"+inputData.CronData),
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	} else {
		nextRun, _ := j.NextRun()

		c.JSON(http.StatusOK, gin.H{
			"message": "ok",
			"ID":      j.ID(),
			"name":    inputData.Name,
			"nextRun": nextRun,
		})
		return
	}
}

// 查询定时任务
func get_cron(c *gin.Context) {
	idStr := c.Query("id")
	var jobs []map[string]interface{}

	for _, job := range S.Jobs() {
		nextRun, _ := job.NextRun()
		remark := strings.Split(job.Name(), "@")
		jobName := remark[0]
		RPAName := remark[1]
		Type := remark[2]
		CronData := remark[3]
		jobInfo := map[string]interface{}{
			"ID":       job.ID(),
			"Name":     jobName,
			"RPAName":  RPAName,
			"Type":     Type,
			"CronData": CronData,
			"NextRun":  nextRun.Format("2006-01-02 15:04:05"),
		}

		if idStr == "" || job.ID() == uuid.MustParse(idStr) {
			jobs = append(jobs, jobInfo)
		}
	}

	if len(jobs) == 0 && idStr != "" {
		c.JSON(http.StatusOK, gin.H{
			"code":    404,
			"message": "未找到定时任务",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "ok",
		"data":    jobs,
	})
}

// 删除定时任务
func del_cron(c *gin.Context) {
	idStr := c.Query("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的 UUID"})
		return
	}

	S.RemoveJob(id)
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "ok",
	})
}
