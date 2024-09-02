package go_requests

import (
	"fmt"
	"testing"
)

// 发送GET请求, 带Query参数
func TestGet(t *testing.T) {
	resp := Get("https://httpbin.org/get?name=张三&age=12", nil)
	fmt.Printf("状态码: %d\n", resp.StatusCode)
	fmt.Printf("原因: %s\n", resp.Reason)
	fmt.Printf("响应时间: %f秒\n", resp.Elapsed)
	fmt.Printf("响应文本: %s\n", resp.Text)
}

// 发送POST 表单请求 带自定义Headers
func TestPostForm(t *testing.T) {
	data := "name=张三&age=12"
	headers := map[string]string{"Content-Type": "application/x-www-form-urlencoded"}
	resp := Post("https://httpbin.org/post", data, headers)
	fmt.Printf("响应文本: %s\n", resp.Text)
}

// 发送POST JSON请求
func TestPostJson(t *testing.T) {
	data := `{"name": "张三", "age": "12"}`
	headers := map[string]string{"Content-Type": "application/json"}
	resp := Post("https://httpbin.org/post", data, headers)
	fmt.Printf("姓名: %s\n", resp.Get("json.name"))
	fmt.Printf("年龄: %s\n", resp.Get("json.age"))
}

// 发送POST XML请求
func TestPostXML(t *testing.T) {
	data := `<xml>hello</xml>`
	headers := map[string]string{"Content-Type": "application/xml"}
	resp := Post("https://httpbin.org/post", data, headers)
	fmt.Printf("状态码: %d\n", resp.StatusCode)
	fmt.Printf("原因: %s\n", resp.Reason)
	fmt.Printf("响应时间: %f秒\n", resp.Elapsed)
	fmt.Printf("响应文本: %s\n", resp.Text)
}

func TestKeepCookies(t *testing.T) {
	s := NewSession(nil)
	data := "name=张三&password=123456"
	resp := s.Post("http://127.0.0.1:5000/api/user/login/", data, nil)

	fmt.Printf("状态码: %d\n", resp.StatusCode)
	fmt.Printf("原因: %s\n", resp.Reason)
	fmt.Printf("响应时间: %f秒\n", resp.Elapsed)
	fmt.Printf("响应文本: %s\n", resp.Text)
	fmt.Printf("响应Cookies: %s\n", resp.Cookies)

	resp2 := s.Get("http://127.0.0.1:5000/api/user/logout/?name=张三", nil)

	fmt.Printf("状态码: %d\n", resp2.StatusCode)
	fmt.Printf("原因: %s\n", resp2.Reason)
	fmt.Printf("响应时间: %f秒\n", resp2.Elapsed)
	fmt.Printf("响应文本: %s\n", resp2.Text)

}
