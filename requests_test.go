package requests

import (
	"fmt"
	"testing"
)

// 发送GET 请求
func TestGet(t *testing.T) {
	resp := Get("https://httpbin.org/get?name=张三&age=12")
	fmt.Printf("状态码: %d\n", resp.StatusCode)
	fmt.Printf("原因: %s\n", resp.Reason)
	fmt.Printf("响应时间: %f秒\n", resp.Elapsed)
	fmt.Printf("响应文本: %s\n", resp.Text)
}

// 发送GET 带单独Query参数请求
func TestGetWithParams(t *testing.T) {
	resp := GetWithParams("https://httpbin.org/get", map[string]string{"name": "张三", "age": "12"})
	fmt.Printf("响应文本: %s\n", resp.Text)
}

// 发送POST 表单请求 带自定义Headers
func TestPostForm(t *testing.T) {
	resp := PostAsForm("https://httpbin.org/post", map[string]string{"name": "张三", "age": "12"})
	fmt.Printf("响应文本: %s\n", resp.Text)
}

// 发送POST JSON请求
func TestPostJson(t *testing.T) {
	resp := PostAsJson("https://httpbin.org/post", `{"name": "张三", "age": "12"}`)
	fmt.Printf("姓名: %s\n", resp.Get("json.name"))
	fmt.Printf("年龄: %s\n", resp.Get("json.age"))
}

// 发送POST XML请求
func TestPostXML(t *testing.T) {
	resp := PostAsRaw("https://httpbin.org/post", `<xml>hello</xml>`, "application/xml")
	fmt.Printf("状态码: %d\n", resp.StatusCode)
	fmt.Printf("原因: %s\n", resp.Reason)
	fmt.Printf("响应时间: %f秒\n", resp.Elapsed)
	fmt.Printf("响应文本: %s\n", resp.Text)
}

func TestPostMultipartFormData(t *testing.T) {
	resp := PostAsMultipartForm("https://httpbin.org/post", map[string]string{"name": "张三", "age": "12"}, map[string]string{"pic": "./testdata/logo.png"})
	fmt.Printf("响应文本: %s\n", resp.Text)
}
