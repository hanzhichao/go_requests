package requests

import (
	"fmt"
	"testing"
)

func TestGetWithConfig(t *testing.T) {
	config := Config{
		BaseUrl: "https://httpbin.org",
		Headers: map[string]string{"token": "token123"},
	}
	r := Request{
		Config: &config,
		Method: "get",
		Url:    "/get?name=张三&age=12"}
	resp := r.Send()
	fmt.Printf("状态码: %d\n", resp.StatusCode)
	fmt.Printf("原因: %s\n", resp.Reason)
	fmt.Printf("响应时间: %f秒\n", resp.Elapsed)
	fmt.Printf("响应文本: %s\n", resp.Text)
}

func TestRequestConfig(t *testing.T) {
	GlobalConfig.BaseUrl = "https://httpbin.org"
	GlobalConfig.Params = map[string]string{"Token": "abc"}
	GlobalConfig.Headers = map[string]string{"Test": "123"}
	GlobalConfig.Cookies = map[string]string{"sid": "hhh"}
	GlobalConfig.Timeout = 10000 // 10秒

	r := Request{
		Method: "get",
		Url:    "/get",
	}
	resp := r.Send()
	fmt.Printf("响应文本：%s\n", resp.Text)
}
