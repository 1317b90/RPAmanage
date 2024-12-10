package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type RPA struct {
	gorm.Model
	Name   string `gorm:"size:30;unique"`
	Remark string `gorm:"size:255;unique"`
	Group  string `gorm:"size:30"`
	Now    bool
	Spont  bool
}

// 查询所有
func get_rpa(c *gin.Context) {
	group := c.DefaultQuery("group", "")
	var data []RPA
	query := db.Model(&RPA{})

	if group != "" {
		query = query.Where("`group` = ?", group)
	}
	query.Find(&data)

	c.JSON(http.StatusOK, gin.H{
		"data": data,
	})
}

// 修改或新增
func put_rpa(c *gin.Context) {
	type inputType struct {
		ID     int    `form:"ID" json:"ID" binding:"required"`
		Name   string `form:"Name" json:"Name"`
		Remark string `form:"Remark" json:"Remark"`
		Group  string `form:"Group" json:"Group"`
		Now    bool   `form:"Now" json:"Now" default:"false"`
		Spont  bool   `form:"Spont" json:"Spont" default:"false"`
	}

	// Spont（自发启动）和Now（立刻执行）的区别
	// Spont，直接启动，不需要影刀或其他第三方监听
	// Now，执行后持续监听影子的执行情况，然后立即返回

	var inputData inputType
	if err := c.ShouldBind(&inputData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	// 删除inputType中string字段的空格
	inputData.Name = strings.TrimSpace(inputData.Name)
	inputData.Remark = strings.TrimSpace(inputData.Remark)

	// 如果是新增变量
	if inputData.ID == -1 {
		data := &RPA{
			Name:   inputData.Name,
			Remark: inputData.Remark,
			Group:  inputData.Group,
			Now:    inputData.Now,
			Spont:  inputData.Spont,
		}

		if result := db.Create(data); result.Error != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": result.Error.Error(),
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"message": "添加成功",
			})
		}

		// 如果是修改变量
	} else {
		var data RPA
		result := db.First(&data, inputData.ID)
		if result.Error != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "该记录不存在",
			})
		} else {
			data.Remark = inputData.Remark
			data.Group = inputData.Group
			data.Now = inputData.Now
			data.Spont = inputData.Spont

			// 如果RPA名称发生变化，对应的var表中的rpa名称也要修改
			if data.Name != inputData.Name {
				set_var_rpa_name(data.Name, inputData.Name)
				data.Name = inputData.Name
			}

			// 保存更新
			result := db.Save(&data)
			if result.Error != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"message": result.Error.Error(),
				})
			} else {
				c.JSON(http.StatusOK, gin.H{
					"message": "修改成功",
				})
			}
		}
	}
}

// 删除
func del_rpa(c *gin.Context) {
	// 获取要删除的RPA名称
	id := c.Query("id")

	// 永久删除
	result := db.Unscoped().Delete(&RPA{}, id)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "删除失败",
			"error":   result.Error.Error(),
		})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "该记录不存在",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "删除成功",
	})

}

// 获取Name与Remark字典
func get_rpa_dict(c *gin.Context) {
	name := c.DefaultQuery("name", "")
	group := c.DefaultQuery("group", "")
	var data []RPA

	query := db.Model(&RPA{})

	if name != "" {
		query = query.Where("name = ?", name)
	}

	if group != "" {
		query = query.Where("`group` = ?", group)
	}

	query.Find(&data)

	dict := make(map[string]string)
	for _, item := range data {
		dict[item.Name] = item.Remark
	}
	c.JSON(http.StatusOK, gin.H{
		"data": dict,
	})
}

// -------- 内部使用 ---------

// 根据rpa名称获取rpa
func get_rpa_by_name(name string) (RPA, error) {
	var data RPA
	result := db.Where("name = ?", name).First(&data)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return RPA{}, fmt.Errorf("未找到名为 %s 的RPA", name)
		}
		return RPA{}, fmt.Errorf("查询RPA时发生错误: %v", result.Error)
	}
	return data, nil
}

// 修改rpa的group
func set_rpa_group(oldGroup string, newGroup string) error {
	result := db.Model(&RPA{}).Where("`group` = ?", oldGroup).Update("group", newGroup)
	if result.Error != nil {
		return fmt.Errorf("修改RPA的group时发生错误: %v", result.Error)
	}
	return nil
}
