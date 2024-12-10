package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Message struct {
	ID      string `json:"ID" form:"ID"`
	UserID  string `json:"userID" form:"userID" binding:"required"`
	TaskID  string `json:"taskID" form:"taskID" binding:"required"`
	Group   string `json:"Group" form:"Group" binding:"required"`
	Content string `json:"Content" form:"Content" default:""`
	Type    string `json:"Type" form:"Type" default:"text"`
	Data    string `json:"Data" form:"Data" default:""`
}

var messages = make([]Message, 0)

// 新增发送消息的队列
func add_message(c *gin.Context) {
	var message Message
	if err := c.ShouldBind(&message); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	message.ID = uuid.New().String()
	if message.Type == "" {
		message.Type = "text"
	}
	messages = append(messages, message)

	c.JSON(http.StatusOK, gin.H{
		"message": "success",
		"id":      message.ID,
	})
}

// 获取最新需要发送的消息
func get_message(c *gin.Context) {
	group := c.Query("group")
	if len(messages) > 0 {
		// 单独取出第一条信息并返回
		found := false
		for i, message := range messages {
			if message.Group == group {
				messages = append(messages[:i], messages[i+1:]...)
				found = true
				c.JSON(http.StatusOK, message)
				break
			}
		}
		if !found {
			c.JSON(http.StatusNotFound, nil)
		}
	} else {
		c.JSON(http.StatusNotFound, nil)
	}
}

// 获取消息列表
func get_message_list(c *gin.Context) {
	c.JSON(http.StatusOK, messages)
}

// 删除队列中的指定消息
func del_message(c *gin.Context) {
	id := c.Param("id")

	found := false
	for i, message := range messages {
		if message.ID == id {
			messages = append(messages[:i], messages[i+1:]...)
			found = true
			break
		}
	}
	if !found {
		c.JSON(http.StatusBadRequest, gin.H{"message": "消息ID不存在"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}
