package requests

import (
	"fmt"
	"testing"
)

// 发送GET 请求
func TestGet(t *testing.T) {
	r := Request{
		Method: "get",
		Url:    "https://httpbin.org/get?name=张三&age=12"}
	resp := r.Send()
	fmt.Printf("状态码: %d\n", resp.StatusCode)
	fmt.Printf("原因: %s\n", resp.Reason)
	fmt.Printf("响应时间: %f秒\n", resp.Elapsed)
	fmt.Printf("响应文本: %s\n", resp.Text)
}

// 发送GET 带单独Query参数请求
func TestGetWithParams(t *testing.T) {
	r := Request{
		Method:  "GET",
		Url:     "https://httpbin.org/get",
		Params:  map[string]string{"name": "张三", "age": "12"},
		Headers: map[string]string{"Cookie": "abc", "Token": "123"}}
	resp := r.Send()
	fmt.Printf("状态码: %d\n", resp.StatusCode)
	fmt.Printf("原因: %s\n", resp.Reason)
	fmt.Printf("响应时间: %f秒\n", resp.Elapsed)
	fmt.Printf("响应文本: %s\n", resp.Text)
}

// 发送POST 表单请求 带自定义Headers
func TestPostForm(t *testing.T) {
	r := Request{
		Method:  "POST",
		Url:     "https://httpbin.org/post",
		Data:    map[string]string{"name": "张三", "age": "12"},
		Headers: map[string]string{"Cookie": "abc", "Token": "123"}}
	resp := r.Send()
	fmt.Printf("状态码: %d\n", resp.StatusCode)
	fmt.Printf("原因: %s\n", resp.Reason)
	fmt.Printf("响应时间: %f秒\n", resp.Elapsed)
	fmt.Printf("响应文本: %s\n", resp.Text)
}

// 发送POST JSON请求
func TestPostJson(t *testing.T) {
	r := Request{
		Method: "POST",
		Url:    "https://httpbin.org/post",
		Json:   `{"name": "张三", "age": "12"}`}
	resp := r.Send()
	fmt.Printf("状态码: %d\n", resp.StatusCode)
	fmt.Printf("原因: %s\n", resp.Reason)
	fmt.Printf("响应时间: %f秒\n", resp.Elapsed)
	fmt.Printf("响应文本: %s\n", resp.Text)
}

// 发送POST XML请求
func TestPostXML(t *testing.T) {
	r := Request{
		Method:  "POST",
		Url:     "https://httpbin.org/post",
		Raw:     `<xml>hello</xml>`,
		Headers: map[string]string{"Content-Type": "application/xml"}}
	resp := r.Send()
	fmt.Printf("状态码: %d\n", resp.StatusCode)
	fmt.Printf("原因: %s\n", resp.Reason)
	fmt.Printf("响应时间: %f秒\n", resp.Elapsed)
	fmt.Printf("响应文本: %s\n", resp.Text)
}

// 测试MultipartFormData
func TestPostMultipartFormData(t *testing.T) {
	r := Request{
		Method: "POST",
		Url:    "https://httpbin.org/post",
		Data:   map[string]string{"name": "张三", "age": "12"},
		Files:  map[string]string{"pic": "./testdata/logo.png"},
	}
	resp := r.Send()
	fmt.Printf("状态码: %d\n", resp.StatusCode)
	fmt.Printf("原因: %s\n", resp.Reason)
	fmt.Printf("响应时间: %f秒\n", resp.Elapsed)
	fmt.Printf("响应文本: %s\n", resp.Text)
}

// 测试BasicAuth
func TestRequestWithAuth(t *testing.T) {
	r := Request{
		Method: "GET",
		Url:    "https://httpbin.org/get",
		Auth:   []string{"Kevin", "123456"},
	}
	resp := r.Send()
	fmt.Printf("响应文本：%s\n", resp.Text)
}

// 测试关闭重定向及响应Cookies
func TestNotAllowRedirests(t *testing.T) {
	url := "https://newecshop.longtest.cn/admin/privilege.php"
	r := Request{
		Method: "POST",
		Url:    url,
		Data: map[string]string{
			"username": "***", // todo 修改为正确的用户名
			"password": "***",  // todo 修改为正确的密码
			"act":      "signin"},
		NoRedirects: true,
	}
	resp := r.Send()
	fmt.Printf("状态码: %d\n", resp.StatusCode)
	for key, value := range resp.Cookies {
		fmt.Printf("响应Cookies项：%s=%s\n", key, value)
	}
}

// 测试请求Timeout
func TestRequestTimeout(t *testing.T) {
	r := Request{
		Method:  "get",
		Url:     "https://httpbin.org/get",
		Timeout: 1, // 毫秒
	}
	r.Send()
}

// 测试请求默认配置
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

// 读取JSON文件发送请求
func TestRequestFromJsonFile(t *testing.T) {
	r := GetRequestFromJsonFile("./testdata/data.json")
	resp := r.Send()
	fmt.Printf("状态码: %d\n", resp.StatusCode)
	fmt.Printf("原因: %s\n", resp.Reason)
	fmt.Printf("响应时间: %f秒\n", resp.Elapsed)
	fmt.Printf("响应文本: %s\n", resp.Text)
}

// 默认请求方法
func TestRequestWithDefaultMethod(t *testing.T) {
	r := Request{Url: "https://httpbin.org/get"}
	resp := r.Send()
	fmt.Printf("状态码: %d\n", resp.StatusCode)

	r = Request{Url: "https://httpbin.org/post", Data: map[string]string{"name": "Kevin"}}
	resp = r.Send()
	fmt.Printf("状态码: %d\n", resp.StatusCode)
}

// 使用HTTP2及关闭TLS验证
func TestRequestWithHttp2(t *testing.T) { // todo 换其他方式验证
	r := Request{Url: "https://stackoverflow.com", HTTP2: true, NoVerify: true}
	resp := r.Send()
	fmt.Printf("响应头: %v\n", resp.Headers)
}

// 测试使用HTTP代理
func TestRequestWithProxy(t *testing.T) { // todo 换其他方式验证
	r := Request{Url: "https://httpbin.org/get", Proxy: "http://localhost:8888", NoVerify: true}
	resp := r.Send()
	fmt.Printf("状态码: %d\n", resp.StatusCode)
}

// 异步发送请求
func TestAsyncSendRequest(t *testing.T) {
	r := Request{Url: "https://www.baidu.com"}
	for i := 0; i < 10; i++ {
		r.AsyncSend()
		resp := <- Ch
		fmt.Println(resp.StatusCode)
	}
}
