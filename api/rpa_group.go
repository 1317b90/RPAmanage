package main

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// RPA组表
type RPAGroup struct {
	gorm.Model
	Name   string `gorm:"size:30;unique"`
	Remark string `gorm:"size:255;unique"`
	IP     string `gorm:"size:20;default:'127.0.0.1'"`
}

// 查询所有rpa组
func get_rpa_group(c *gin.Context) {
	var data []RPAGroup
	db.Find(&data)

	c.JSON(http.StatusOK, gin.H{
		"data": data,
	})
}

// 修改或新增rpa组
func put_rpa_group(c *gin.Context) {
	type inputType struct {
		ID     int    `form:"ID" json:"ID" binding:"required"`
		Name   string `form:"Name" json:"Name"`
		Remark string `form:"Remark" json:"Remark"`
		IP     string `form:"IP" json:"IP"`
	}

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
		newData := &RPAGroup{
			Name:   inputData.Name,
			Remark: inputData.Remark,
			IP:     inputData.IP,
		}

		if result := db.Create(newData); result.Error != nil {
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
		var oldData RPAGroup
		result := db.First(&oldData, inputData.ID)
		if result.Error != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "该记录不存在",
			})
		} else {
			oldData.Remark = inputData.Remark
			oldData.IP = inputData.IP

			// 如果RPA组名发生变化，对应的rpa表和var表中的group也要修改
			if oldData.Name != inputData.Name {
				set_rpa_group(oldData.Name, inputData.Name)
				set_var_group(oldData.Name, inputData.Name)
				oldData.Name = inputData.Name
			}
			// 保存更新
			result := db.Save(&oldData)
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

// 删除rpa
func del_rpa_group(c *gin.Context) {
	// 获取要删除的RPA名称
	id := c.Query("id")

	// 永久删除
	result := db.Unscoped().Delete(&RPAGroup{}, id)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "删除失败",
			"error":   result.Error.Error(),
		})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "未找到该记录",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "删除成功",
	})

}

// 获取Name与Remark字典
func get_rpa_group_dict(c *gin.Context) {
	name := c.DefaultQuery("name", "")
	var data []RPAGroup
	if name == "" {
		db.Find(&data)
	} else {
		db.Where("name = ?", name).Find(&data)
	}

	dict := make(map[string]string)
	for _, item := range data {
		dict[item.Name] = item.Remark
	}
	c.JSON(http.StatusOK, gin.H{
		"data": dict,
	})
}
