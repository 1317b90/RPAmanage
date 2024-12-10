package main

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
)

func up_file_xlsx(c *gin.Context) {
	// 接收上传的文件
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(400, gin.H{"message": "文件上传失败"})
		return
	}

	// 检查文件扩展名是否为.xlsx
	if filepath.Ext(file.Filename) != ".xlsx" {
		c.JSON(400, gin.H{"message": "只允许上传.xlsx文件"})
		return
	}

	// 确保目标文件夹存在
	uploadDir := "file/xlsx"

	// 生成目标文件路径
	dst := filepath.Join(uploadDir, file.Filename)

	// 保存文件
	if err := c.SaveUploadedFile(file, dst); err != nil {
		c.JSON(500, gin.H{"message": "保存文件失败"})
		return
	}

	// 确保在成功时返回JSON响应
	c.JSON(200, gin.H{"message": "文件上传成功", "filename": file.Filename})
}

// 上传xlsx并返回二维表格
func up_file_batch(c *gin.Context) {
	// 接收上传的文件
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(400, gin.H{"message": "文件上传失败"})
		return
	}

	// 检查文件扩展名是否为.xlsx
	if filepath.Ext(file.Filename) != ".xlsx" {
		c.JSON(400, gin.H{"message": "只允许上传.xlsx文件"})
		return
	}

	// 确保目标文件夹存在
	uploadDir := "file/xlsx"

	// 生成目标文件路径
	dst := filepath.Join(uploadDir, file.Filename)

	// 保存文件
	if err := c.SaveUploadedFile(file, dst); err != nil {
		c.JSON(500, gin.H{"message": "保存文件失败"})
		return
	}

	// 使用excelize库读取文件
	f, err := excelize.OpenFile(dst)
	if err != nil {
		c.JSON(500, gin.H{"message": "无法打开Excel文件"})
		return
	}
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println("关闭Excel文件时出错:", err)
		}
	}()

	// 获取第一个工作表
	sheetName := f.GetSheetName(0)
	rows, err := f.GetRows(sheetName)
	if err != nil {
		c.JSON(500, gin.H{"message": "无法读取Excel数据"})
		return
	}

	// 获取表头
	if len(rows) < 4 {
		c.JSON(400, gin.H{"message": "Excel文件数据不完整，请检查数据是否正确"})
		return
	}
	headers := rows[0]

	// 获取符合条件的数据
	var data []map[string]string
	for _, row := range rows[3:] { // 跳过表头
		rowData := make(map[string]string)

		for i, cell := range row {
			if i < len(headers) {
				rowData[headers[i]] = strings.TrimSpace(cell)
			}
		}
		data = append(data, rowData)
	}

	// 返回结果
	c.JSON(200, gin.H{
		"message": "数据获取成功",
		"data":    data,
	})

}

// 上传文件 统一
func up_file_common(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(500, gin.H{"message": "文件上传失败"})
		return
	}

	// 保存文件到 "file/wecom" 路径下
	dst := filepath.Join("file", "common", file.Filename)
	if err := c.SaveUploadedFile(file, dst); err != nil {
		c.JSON(500, gin.H{"message": "保存文件失败"})
		return
	}

	// 将Windows路径分隔符转换为URL路径分隔符
	urlPath := strings.ReplaceAll(dst, "\\", "/")
	c.JSON(200, gin.H{"message": "ok", "url": fmt.Sprintf("https://test.g4b.cn/rpa/%s", urlPath)})
}

// 上传企业微信数据
func up_file_wecom(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(500, gin.H{"message": "文件上传失败"})
		return
	}

	// 保存文件到 "file/wecom" 路径下
	dst := filepath.Join("file", "wecom", file.Filename)
	if err := c.SaveUploadedFile(file, dst); err != nil {
		c.JSON(500, gin.H{"message": "保存文件失败"})
		return
	}

	// 将Windows路径分隔符转换为URL路径分隔符
	urlPath := strings.ReplaceAll(dst, "\\", "/")
	c.JSON(200, gin.H{"message": "ok", "url": fmt.Sprintf("https://test.g4b.cn/rpa/%s", urlPath)})

}

// 下载模板文件
func down_file_template(c *gin.Context) {
	var headData []VarType

	if err := c.ShouldBind(&headData); err != nil {
		c.JSON(400, gin.H{"message": "无效的表头数据", "error": err.Error()})
		return
	}

	// 创建一个新的Excel文件
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	// 创建一个新的工作表
	sheetName := "Sheet1"
	index, err := f.NewSheet(sheetName)
	if err != nil {
		c.JSON(500, gin.H{"message": "创建工作表失败"})
		return
	}

	// 设置活动工作表
	f.SetActiveSheet(index)

	// 写入表头数据
	for i, data := range headData {
		// 第一行：VarName
		cell1 := fmt.Sprintf("%c1", 'A'+i)
		f.SetCellValue(sheetName, cell1, data.VarName)

		// 第二行：VarRemark
		cell2 := fmt.Sprintf("%c2", 'A'+i)
		f.SetCellValue(sheetName, cell2, data.VarRemark)

		// 第三行：必填或空字符串
		cell3 := fmt.Sprintf("%c3", 'A'+i)
		if data.Required {
			f.SetCellValue(sheetName, cell3, "必填")
		} else {
			f.SetCellValue(sheetName, cell3, "")
		}
	}

	// 保存文件
	filename := "template.xlsx"
	filePath := filepath.Join("file", "mod", filename)
	if err := f.SaveAs(filePath); err != nil {
		c.JSON(500, gin.H{"message": "保存模板文件失败"})
		return
	}

	c.JSON(200, gin.H{"message": "ok", "url": "https://test.g4b.cn/rpa/" + filePath})
	// c.JSON(200, gin.H{"message": "ok", "url": "http://127.0.0.1:8080/" + filePath})
}
