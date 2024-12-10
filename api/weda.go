package main

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
)

const (
	envID          = os.Getenv("WEDA_ENV_ID")          // 环境ID
	envType        = os.Getenv("WEDA_ENV_TYPE")        // 环境类型
	datasourceName = os.Getenv("WEDA_DATASOURCE_NAME") // 数据表名
	basicKey       = os.Getenv("WEDA_BASIC_KEY")       // 密钥
)

var (
	host    = fmt.Sprintf("https://%s.ap-shanghai.tcb-api.tencentcloudapi.com", envID)
	wedaUrl = fmt.Sprintf("https://%s.ap-shanghai.tcb-api.tencentcloudapi.com/weda/odata/v1/%s/%s", envID, envType, datasourceName)
	headers = map[string]string{
		"Authorization": "",
		"Content-Type":  "application/json",
	}
)

// 更新微搭Token
func getToken() error {
	url := fmt.Sprintf("%s/auth/v1/token/clientCredential", host)
	authorizationKey := "Basic " + base64.StdEncoding.EncodeToString([]byte(basicKey))

	client := resty.New()

	resp, err := client.R().
		SetHeader("Authorization", authorizationKey).
		SetHeader("Content-Type", "application/json").
		SetBody(map[string]string{
			"grant_type": "client_credentials",
		}).
		Post(url)

	if err != nil {
		return err
	}

	var responseData map[string]interface{}
	err = json.Unmarshal(resp.Body(), &responseData)
	if err != nil {
		return err
	}

	// 更改全局的token
	headers["Authorization"] = "Bearer " + responseData["access_token"].(string)
	return nil
}

// 根据id获取单条数据
func get_weda_by_id(id string) (map[string]interface{}, error) {
	client := resty.New()
	url := fmt.Sprintf("%s('%s')", wedaUrl, id)
	resp, err := client.R().
		SetHeaders(headers).
		Get(url)

	if err != nil {
		return nil, fmt.Errorf("发送请求失败: %v", err)
	}

	if resp.IsError() {
		return nil, fmt.Errorf("请求失败: %s", string(resp.Body()))
	}

	var result map[string]interface{}
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		return nil, fmt.Errorf("JSON解析失败: %v, 响应体: %s", err, string(resp.Body()))
	}

	if result["code"] == "PERMISSION_DENIED" || result["code"] == "INVALID_ACCESS_TOKEN" {
		// 刷新token
		err := getToken()
		if err != nil {
			return nil, fmt.Errorf("更新token失败: %v", err)
		}

		// 重新执行请求
		return get_weda_by_id(id)
	}

	if result["_id"] == nil {
		return nil, fmt.Errorf("该用户不存在")
	}

	return result, nil
}

// 获取所有数据
func get_weda_all(filter, selects, phone string) ([]interface{}, error) {
	client := resty.New()
	url := wedaUrl

	// 设置查询参数
	queryParams := map[string]string{}

	// 默认获取所有数据
	queryParams["$top"] = "9999"

	// 筛选
	if phone != "" {
		queryParams["$filter"] = fmt.Sprintf("contains(sjhm,'%s')", phone)

	} else if filter != "" {
		queryParams["$filter"] = filter

		// 指定获取某些字段
	} else if selects != "" {
		queryParams["$select"] = selects
	}

	// 发送请求
	resp, err := client.R().
		SetQueryParams(queryParams).
		SetHeaders(headers).
		Get(url)

	if err != nil {
		return nil, fmt.Errorf("发送请求失败: %v", err)
	}

	if resp.IsError() {
		return nil, fmt.Errorf("请求失败: %s", string(resp.Body()))
	}

	if len(resp.Body()) == 0 {
		return nil, fmt.Errorf("响应体为空")
	}

	var result map[string]interface{}
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		return nil, fmt.Errorf("JSON解析失败: %v, 响应体: %s", err, string(resp.Body()))
	}

	// token失效
	if result["code"] == "PERMISSION_DENIED" || result["code"] == "INVALID_ACCESS_TOKEN" {
		// 刷新token
		err := getToken()
		if err != nil {
			return nil, fmt.Errorf("更新token失败: %v", err)
		}

		// 重新执行请求
		return get_weda_all(filter, selects, phone)
	}

	dataList, ok := result["value"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("数据获取失败")
	}

	return dataList, nil
}

// 新增数据
func add_weda(data map[string]interface{}) (map[string]interface{}, error) {
	fmt.Println(wedaUrl)

	// 创建resty客户端
	client := resty.New()

	// 发送POST请求
	resp, err := client.R().
		SetHeaders(headers).
		SetBody(data).
		Post(wedaUrl)

	fmt.Println(wedaUrl)
	fmt.Println(data)

	if err != nil {
		return nil, fmt.Errorf("发送请求失败: %v", err)
	}

	if resp.IsError() {
		return nil, fmt.Errorf("请求失败: %s", string(resp.Body()))
	}

	// 处理响应
	var result map[string]interface{}
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		return nil, fmt.Errorf("JSON解析失败: %v, 响应体: %s", err, string(resp.Body()))
	}

	// 检查是否需要刷新token
	if result["code"] == "PERMISSION_DENIED" || result["code"] == "INVALID_ACCESS_TOKEN" {
		err := getToken()
		if err != nil {
			return nil, fmt.Errorf("更新token失败: %v", err)
		}
		// 重新执行请求
		return add_weda(data)
	}

	return result, nil
}

// 修改数据
func set_weda(data map[string]interface{}) error {
	url := fmt.Sprintf("%s('%s')", wedaUrl, data["_id"])

	_, err := get_weda_by_id(data["_id"].(string))
	if err != nil {
		return fmt.Errorf("该用户不存在")
	}

	// 创建resty客户端
	client := resty.New()

	// 发送POST请求
	resp, err := client.R().
		SetHeaders(headers).
		SetBody(data).
		Patch(url)

	if err != nil {
		return fmt.Errorf("发送请求失败: %v", err)
	}

	// 检查响应状态码是否为204
	if resp.StatusCode() == 200 {
		var result map[string]interface{}
		_ = json.Unmarshal(resp.Body(), &result)

		// 检查是否需要刷新token
		if result["code"] == "PERMISSION_DENIED" || result["code"] == "INVALID_ACCESS_TOKEN" {
			err := getToken()
			if err != nil {
				return fmt.Errorf("更新token失败: %v", err)
			}
			// 重新执行请求
			return set_weda(data)
		}
		// 请求成功
	} else if resp.StatusCode() == 204 {
		return nil
	}
	return errors.New(string(resp.Body()))
}

// 删除数据
func del_weda(id string) error {
	url := fmt.Sprintf("%s('%s')", wedaUrl, id)
	// 创建resty客户端
	client := resty.New()

	// 发送POST请求
	resp, err := client.R().
		SetHeaders(headers).
		Delete(url)

	if err != nil {
		return fmt.Errorf("发送请求失败: %v", err)
	}

	// 检查响应状态码是否为200
	if resp.StatusCode() == 200 {
		var result map[string]interface{}
		_ = json.Unmarshal(resp.Body(), &result)

		// 检查是否需要刷新token
		if result["code"] == "PERMISSION_DENIED" || result["code"] == "INVALID_ACCESS_TOKEN" {
			err := getToken()
			if err != nil {
				return fmt.Errorf("更新token失败: %v", err)
			}
			// 重新执行请求
			return del_weda(id)
		}
		// 请求成功
	} else if resp.StatusCode() == 204 {
		return nil
	}
	return errors.New(string(resp.Body()))
}

// --------- 接口 --------- 接口 --------- 接口 --------- 接口 --------- 接口 --------- 接口 --------- 接口
// 获取数据接口
func get_weda_api(c *gin.Context) {
	// 指定id
	id := c.DefaultQuery("id", "")

	if id != "" {
		response, err := get_weda_by_id(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"data": response})
	} else {
		// 指定手机号码
		phone := c.DefaultQuery("phone", "")
		// 筛选条件
		filter := c.DefaultQuery("filter", "")
		// 选择字段
		selects := c.DefaultQuery("selects", "")

		response, err := get_weda_all(filter, selects, phone)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"data": response})
	}
}

// 新增数据
func add_weda_api(c *gin.Context) {
	type inputS struct {
		Data json.RawMessage `form:"data" json:"data" binding:"required"`
	}

	var inputData inputS

	if err := c.ShouldBind(&inputData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	var data map[string]interface{}
	if err := json.Unmarshal(inputData.Data, &data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "解析输入数据失败"})
		return
	}

	// 在手机号码不为空的情况下，检测是否有重复的
	if data["sjhm"] != nil {
		response, err := get_weda_all("", "", data["sjhm"].(string))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		if len(response) > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"message": "手机号码已存在"})
			return
		}
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"message": "手机号码不能为空"})
		return
	}

	response, err := add_weda(data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": response})
}

// 修改数据接口
func set_weda_api(c *gin.Context) {
	type inputS struct {
		Data json.RawMessage `form:"data" json:"data" binding:"required"`
	}

	var inputData inputS

	if err := c.ShouldBind(&inputData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	var data map[string]interface{}
	if err := json.Unmarshal(inputData.Data, &data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "解析输入数据失败"})
		return
	}

	err := set_weda(data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}

// 删除数据接口
func del_weda_api(c *gin.Context) {
	id := c.DefaultQuery("id", "")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"messsage": "id不能为空"})
	} else {
		response, _ := get_weda_by_id(id)
		if response == nil {
			c.JSON(http.StatusBadRequest, gin.H{"messsage": "该用户不存在"})
			return
		}
		err := del_weda(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"messsage": err.Error()})
		} else {
			c.JSON(http.StatusOK, gin.H{"message": "ok"})
		}
	}

}
