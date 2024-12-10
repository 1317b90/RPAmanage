package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Log struct {
	gorm.Model
	Type    string `gorm:"size:30"`
	Title   string `gorm:"size:255"`
	State   bool   `gorm:"default:true"`
	Content string `gorm:"type:text;"`
}

// 查询所有记录
func get_log(c *gin.Context) {
	var data []Log
	query := db.Model(&Log{})

	// 根据输入的参数进行动态查询
	if Type := c.Query("Type"); Type != "" {
		query = query.Where("type = ?", Type)
	}
	if Title := c.Query("Title"); Title != "" {
		query = query.Where("title LIKE ?", "%"+Title+"%")
	}
	if State := c.Query("State"); State != "" {
		query = query.Where("state = ?", State)
	}
	if Content := c.Query("Content"); Content != "" {
		query = query.Where("content LIKE ?", "%"+Content+"%")
	}

	// 执行查询
	result := query.Find(&data)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "查询日志失败",
			"error":   result.Error.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": data,
	})
}

// 新增记录，本地使用
func add_log_func(Type, Title, Content string, State bool) error {
	result := db.Create(&Log{Type: Type, Title: Title, Content: Content, State: State})
	if result.Error != nil {
		return fmt.Errorf("创建日志记录失败: %v", result.Error)
	}
	return nil
}

// 新增记录
func add_log(c *gin.Context) {
	type inputType struct {
		Type    string `form:"Type" json:"Type" binding:"required"`
		Title   string `form:"Title" json:"Title" binding:"required"`
		Content string `form:"Content" json:"Content"`
		State   bool   `form:"State" json:"State"`
	}
	var inputData inputType
	if err := c.ShouldBind(&inputData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	err := add_log_func(inputData.Type, inputData.Title, inputData.Content, inputData.State)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}

// 删除记录
func del_log(c *gin.Context) {
	id := c.DefaultQuery("id", "")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "id不能为空"})
		return
	}
	err := db.Delete(&Log{}, id).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}
