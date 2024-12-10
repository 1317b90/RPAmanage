package main

import (
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// 测试环境

// var dsn = "olive:wangrong@2024@tcp(127.0.0.1:3306)/rpa?charset=utf8mb4&parseTime=True&loc=Local"

// 正式环境
var dsn = "olive:wangrong@2024@tcp(172.18.59.37:3306)/rpa?charset=utf8mb4&parseTime=True&loc=Local"

var db, _ = gorm.Open(mysql.Open(dsn), &gorm.Config{})

// 定义路由
var R = gin.Default()

func main() {
	// 全局设置时区
	os.Setenv("TZ", "Asia/Shanghai")

	// 迁移数据表
	_ = db.AutoMigrate(&RPA{})
	_ = db.AutoMigrate(&Task{})
	_ = db.AutoMigrate(&Var{})
	_ = db.AutoMigrate(&RPAGroup{})
	_ = db.AutoMigrate(&Log{})

	S.Start()

	// 允许所有跨域请求
	R.Use(cors.Default())

	// 打个招呼
	R.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "hello",
		})
	})
	// --------- 任务类 --------- 任务类 --------- 任务类 --------- 任务类 --------- 任务类 --------- 任务类 --------- 任务类
	// 添加任务
	R.POST("/task", create_task)

	// 查询任务 从内存中
	R.GET("/task/memory", get_task_memory)

	// 查询任务 从数据库中
	R.GET("/task/db", get_task_db)

	// 更新任务状态为ing
	R.GET("/task/ing", ing_task)

	// 任务完成后，修改数据库中任务状态
	R.PUT("/task/done", done_task)

	// 删除任务，从数据库
	R.DELETE("/task", del_task)

	// 删除任务，从内存
	R.DELETE("/task/memory", del_task_memory)

	// 统计任务数据
	R.GET("/task/count", count_task)

	// --------- RPA组 --------- RPA组 --------- RPA组 --------- RPA组 --------- RPA组 --------- RPA组 --------- RPA组

	// 新增和编辑RPA组
	R.PUT("/rpa/group", put_rpa_group)

	// 获取所有的RPA组
	R.GET("/rpa/group", get_rpa_group)

	// 删除RPA组
	R.DELETE("/rpa/group", del_rpa_group)

	// 获取RPA组名称对应表
	R.GET("/rpa/group/dict", get_rpa_group_dict)

	// --------- RPA表 --------- RPA表 --------- RPA表 --------- RPA表 --------- RPA表 --------- RPA表 --------- RPA表

	// 新增和编辑RPA
	R.PUT("/rpa", put_rpa)

	// 获取所有的RPA
	R.GET("/rpa", get_rpa)

	// 删除RPA
	R.DELETE("/rpa", del_rpa)

	// 获取RPA名称对应表
	R.GET("/rpa/dict", get_rpa_dict)

	// --------- 变量类 --------- 变量类 --------- 变量类 --------- 变量类 --------- 变量类 --------- 变量类 --------- 变量类
	// 查询变量
	R.GET("/var", get_var)

	// 修改或新增某个变量的值
	R.PUT("/var", put_var)

	// 根据id删除变量
	R.DELETE("/var", del_var)

	// 根据输入的多个RPAName获取变量
	R.GET("/var/rpa", get_var_by_rpa_name_s)

	// --------- 文件 --------- 文件 --------- 文件 --------- 文件 --------- 文件 --------- 文件 --------- 文件 --------- 文件

	// 设置file文件夹为静态文件目录
	R.Static("/file", "./file")

	// 上传通用的文件
	R.POST("/upfile/common", up_file_common)

	// 上传文件
	R.POST("/upfile/xlsx", up_file_xlsx)

	// 上传xlsx并返回二维表格
	R.POST("/upfile/batch", up_file_batch)

	// 上传企业微信数据
	R.POST("/upfile/wecom", up_file_wecom)

	// 下载模板文件
	R.POST("/downfile/template", down_file_template)

	// --------- 微搭 --------- 微搭 --------- 微搭 --------- 微搭 --------- 微搭 --------- 微搭 --------- 微搭 --------- 微搭

	// 获取所有微搭数据
	R.GET("/weda", get_weda_api)

	// 新增微搭数据
	R.POST("/weda", add_weda_api)

	// 修改微搭数据
	R.PUT("/weda", set_weda_api)

	// 删除微搭数据
	R.DELETE("/weda", del_weda_api)

	// --------- 定时任务 --------- 定时任务 --------- 定时任务 --------- 定时任务 --------- 定时任务 --------- 定时任务 --------- 定时任务 --------- 定时任务

	// 新增定时任务
	R.POST("/cron", create_cron)

	// 查询定时任务
	R.GET("/cron", get_cron)

	// 删除定时任务
	R.DELETE("/cron", del_cron)

	// ---------- 企业微信 ---------- 企业微信 ---------- 企业微信 ---------- 企业微信 ---------- 企业微信 ---------- 企业微信 ---------- 企业微信 ---------- 企业微信
	// 新增发送消息的队列
	R.POST("/wecom", add_message)

	// 读取消息的列
	R.GET("/wecom", get_message)

	// 获取消息列表
	R.GET("/wecom/list", get_message_list)

	// 删除消息
	R.DELETE("/wecom/:id", del_message)

	// ---------- 日志 ---------- 日志 ---------- 日志 ---------- 日志 ---------- 日志 ---------- 日志 ---------- 日志 ---------- 日志
	// 查询所有日志
	R.GET("/log", get_log)

	// 新增日志
	R.POST("/log", add_log)

	// 删除日志
	R.DELETE("/log", del_log)

	_ = R.Run()
}
