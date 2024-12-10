package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// 变量表
type Var struct {
	gorm.Model
	// 对应的RPA组
	RPAGroup string `gorm:"size:255"`
	// 对应的任务类型
	RPAName string `gorm:"size:255"`
	// 变量名称
	VarName string `gorm:"size:255"`
	// 变量的备注
	VarRemark string `gorm:"size:255"`
	// 变量在别的地方的名称 例如微搭
	AsName string `gorm:"size:255"`
	// 变量的数据类型
	VarType string `gorm:"size:255" default:"string"`
	// 变量的验证类型
	VerifyType string `gorm:"size:255"`
	// 变量的默认值
	Default string `gorm:"size:255"`
	// 变量是否必填
	Required bool
}

// 获取变量
func get_var(c *gin.Context) {
	RPAName := c.DefaultQuery("RPAName", "")
	RPAGroupName := c.DefaultQuery("RPAGroupName", "")

	var data []Var
	query := db

	if RPAName != "" {
		query = query.Where("rpa_name = ?", RPAName)
	}

	if RPAGroupName != "" {
		query = query.Where("rpa_group = ?", RPAGroupName)
	}

	if err := query.Find(&data).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "请求变量失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": data,
	})
}

// 根据输入的多个RPAName获取变量
func get_var_by_rpa_name_s(c *gin.Context) {
	RPANameList := c.DefaultQuery("RPANameList", "")
	if RPANameList == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "至少输入一个RPAName"})
		return
	} else {
		// 将RPANameList按逗号分割成切片
		rpaNames := strings.Split(RPANameList, ",")

		// 查询符合条件的数据
		var data []Var
		if err := db.Where("rpa_name IN ?", rpaNames).Find(&data).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "查询变量失败"})
			return
		}

		// 返回查询结果
		c.JSON(http.StatusOK, gin.H{
			"data": data,
		})
	}

}

// 根据id查询变量
func get_var_by_id(ID int) (*Var, error) {
	var varr Var
	result := db.First(&varr, ID)
	if result.Error != nil {
		return nil, result.Error
	}
	return &varr, nil
}

type VarType struct {
	ID         int    `form:"ID" json:"ID" binding:"required"`
	RPAGroup   string `form:"RPAGroup" json:"RPAGroup" binding:"required"`
	RPAName    string `form:"RPAName" json:"RPAName" binding:"required"`
	VarName    string `form:"VarName" json:"VarName" binding:"required"`
	VarRemark  string `form:"VarRemark" json:"VarRemark" binding:"required"`
	AsName     string `form:"AsName" json:"AsName" `
	VarType    string `form:"VarType" json:"VarType" `
	VerifyType string `form:"VerifyType" json:"VerifyType" `
	Default    string `form:"Default" json:"Default" `
	Required   bool   `form:"Required" json:"Required" `
}

// 修改或新增变量
func put_var(c *gin.Context) {

	var inputData VarType
	if err := c.ShouldBind(&inputData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	// 删除inputType中string字段的空格
	inputData.RPAName = strings.TrimSpace(inputData.RPAName)
	inputData.VarName = strings.TrimSpace(inputData.VarName)
	inputData.VarRemark = strings.TrimSpace(inputData.VarRemark)
	inputData.VarType = strings.TrimSpace(inputData.VarType)
	inputData.VerifyType = strings.TrimSpace(inputData.VerifyType)
	inputData.AsName = strings.TrimSpace(inputData.AsName)
	inputData.Default = strings.TrimSpace(inputData.Default)

	// 如果是新增变量
	if inputData.ID == -1 {
		data := &Var{
			RPAGroup:   inputData.RPAGroup,
			RPAName:    inputData.RPAName,
			VarName:    inputData.VarName,
			VarRemark:  inputData.VarRemark,
			AsName:     inputData.AsName,
			VarType:    inputData.VarType,
			VerifyType: inputData.VerifyType,
			Default:    inputData.Default,
			Required:   inputData.Required,
		}

		if result := db.Create(data); result.Error != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": result.Error,
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"message": "添加成功",
			})
		}

		// 如果是修改变量
	} else {
		data, err := get_var_by_id(inputData.ID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "该记录不存在",
			})
		} else {
			data.RPAGroup = inputData.RPAGroup
			data.RPAName = inputData.RPAName
			data.VarName = inputData.VarName
			data.VarRemark = inputData.VarRemark
			data.AsName = inputData.AsName
			data.VarType = inputData.VarType
			data.VerifyType = inputData.VerifyType
			data.Default = inputData.Default
			data.Required = inputData.Required

			result := db.Save(&data)
			if result.Error != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"message": result.Error,
				})
			} else {
				c.JSON(http.StatusOK, gin.H{
					"message": "修改成功",
				})
			}
		}
	}
}

// 删除变量
func del_var(c *gin.Context) {
	id := c.Query("id")
	intid, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "id是数字",
		})
	} else {
		varr, errr := get_var_by_id(intid)
		if errr != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "该变量不存在",
			})
		} else {
			db.Delete(varr, intid)
			c.JSON(http.StatusOK, gin.H{
				"message": "删除成功",
			})
		}
	}
}

// -------- 内部使用 ---------

// 根据rpa名称获取变量
func get_var_by_rpa_name(rpa_name string) ([]Var, error) {
	var data []Var
	query := db.Where("rpa_name = ?", rpa_name)
	if err := query.Find(&data).Error; err != nil {
		return nil, err
	}
	return data, nil
}

// 修改变量的group
func set_var_group(oldGroup string, newGroup string) error {
	result := db.Model(&Var{}).Where("`group` = ?", oldGroup).Update("group", newGroup)
	if result.Error != nil {
		return fmt.Errorf("修改变量的group时发生错误: %v", result.Error)
	}
	return nil
}

// 修改变量的rpa名称
func set_var_rpa_name(oldRPAName string, newRPAName string) error {
	result := db.Model(&Var{}).Where("rpa_name = ?", oldRPAName).Update("rpa_name", newRPAName)
	if result.Error != nil {
		return fmt.Errorf("修改变量的rpa名称时发生错误: %v", result.Error)
	}
	return nil
}
