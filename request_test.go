package requests

import (
	"fmt"
	"testing"
)

// 发送GET 请求
func TestRequestGet(t *testing.T) {
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
func TestRequestGetWithParams(t *testing.T) {
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
func TestRequestPostForm(t *testing.T) {
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
func TestRequestPostJson(t *testing.T) {
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
func TestRequestPostXML(t *testing.T) {
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
func TestRequestPostMultipartFormData(t *testing.T) {
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
			"password": "***", // todo 修改为正确的密码
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
		resp := <-Ch
		fmt.Println(resp.StatusCode)
	}
}

func TestBuildRequest(t *testing.T) {
	r := NewRequest().
		SetMethod("get").
		SetUrl("https://httpbin.org/get").
		SetParams(map[string]string{"name": "张三", "age": "12"}).
		SetTimeout(3000)
	resp := r.Send()
	fmt.Printf("状态码: %d\n", resp.StatusCode)
	fmt.Printf("原因: %s\n", resp.Reason)
	fmt.Printf("响应时间: %f秒\n", resp.Elapsed)
	fmt.Printf("响应文本: %s\n", resp.Text)
}

func TestRequestWithPrevResponseCookies(t *testing.T) {
	r := Request{
		Method: "POST",
		Url:    "http://127.0.0.1:5000/api/user/login/",
		Data:   map[string]string{"name": "张三", "password": "123456"}}
	resp := r.Send()
	fmt.Printf("状态码: %d\n", resp.StatusCode)
	fmt.Printf("原因: %s\n", resp.Reason)
	fmt.Printf("响应时间: %f秒\n", resp.Elapsed)
	fmt.Printf("响应文本: %s\n", resp.Text)
	fmt.Printf("响应Cookies: %s\n", resp.Cookies)

	r2 := Request{
		Method:  "GET",
		Url:     "http://127.0.0.1:5000/api/user/logout/",
		Params:  map[string]string{"name": "张三"},
		Cookies: resp.Cookies,
	}
	resp2 := r2.Send()

	fmt.Printf("状态码: %d\n", resp2.StatusCode)
	fmt.Printf("原因: %s\n", resp2.Reason)
	fmt.Printf("响应时间: %f秒\n", resp2.Elapsed)
	fmt.Printf("响应文本: %s\n", resp2.Text)
}
