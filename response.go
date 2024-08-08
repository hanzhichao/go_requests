package requests

import (
	"encoding/json"
	"fmt"
)

// 请求配置

// 响应结构体
type Response struct {
	StatusCode int               `json:"status_code"` // 状态码
	Reason     string            `json:"reason"`      // 状态码说明
	Elapsed    float64           `json:"elapsed"`     // 请求耗时(秒)
	Content    []byte            `json:"content"`     // 响应二进制内容
	Text       string            `json:"text"`        // 响应文本
	Headers    map[string]string `json:"headers"`     // 响应头
	Cookies    map[string]string `json:"cookies"`     // 响应Cookies
	Request    *Request          `json:"request"`     // 原始请求
}

func (res *Response) Json() map[string]interface{} {
	data := make(map[string]interface{})
	err := json.Unmarshal(res.Content, &data)
	if err != nil {
		fmt.Println("响应文本转map失败")
	}
	return data
}

func (res *Response) JsonArray() []map[string]interface{} {
	var data []map[string]interface{}
	err := json.Unmarshal(res.Content, &data)
	if err != nil {
		fmt.Println("响应文本转map失败")
	}
	return data
}
