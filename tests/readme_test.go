package tests

import (
	"encoding/json"
	"fmt"
	"github.com/hanzhichao/go_requests"
	"testing"
)

func TestGet(t *testing.T) {
	resp := go_requests.Get("https://httpbin.org/get?name=张三&age=12", nil)
	fmt.Printf("状态码: %d\n", resp.StatusCode)
	fmt.Printf("原因: %s\n", resp.Reason)
	fmt.Printf("响应时间: %f秒\n", resp.Elapsed)
	fmt.Printf("响应文本: %s\n", resp.Text)
}

func TestPostForm(t *testing.T) {
	resp := go_requests.Post("https://httpbin.org/post", "name=张三&age=12",
		map[string]string{"Content-Type": "application/x-www-form-urlencoded"})
	fmt.Printf("响应文本: %s\n", resp.Text)
}

func TestPostJson(t *testing.T) {
	resp := go_requests.Post("https://httpbin.org/post", `{"name": "张三", "age": "12"}`,
		map[string]string{"Content-Type": "application/json"})
	// JSON响应解析
	fmt.Printf("姓名: %s\n", resp.Get("json.name"))
	fmt.Printf("年龄: %s\n", resp.Get("json.age"))
}

func TestPostXML(t *testing.T) {
	resp := go_requests.Post("https://httpbin.org/post", `<xml>hello</xml>`,
		map[string]string{"Content-Type": "application/xml"})
	fmt.Printf("响应文本: %s\n", resp.Text)
}

func TestPostMultipartFormData(t *testing.T) {
	r := go_requests.Request{
		Method: "POST",
		Url:    "https://httpbin.org/get",
		Data:   map[string]string{"name": "张三", "age": "12"},
		Files:  map[string]string{"pic": "./testdata/logo.png"}}
	resp := r.Send()
	fmt.Printf("响应文本: %s\n", resp.Text)
}

func TestRequestWithAuth(t *testing.T) {
	r := go_requests.Request{
		Method: "GET",
		Url:    "https://httpbin.org/get",
		Auth:   []string{"Kevin", "123456"},
	}

	resp := r.Send()
	fmt.Printf("响应文本：%s\n", resp.Text)
}

func TestRequestTimeout(t *testing.T) {
	r := go_requests.Request{
		Method:  "get",
		Url:     "https://httpbin.org/get",
		Timeout: 1, // 毫秒
	}

	r.Send()
}

func TestRequestConfig(t *testing.T) {
	go_requests.GlobalConfig.BaseUrl = "https://httpbin.org"
	go_requests.GlobalConfig.Params = map[string]string{"Token": "abc"}
	go_requests.GlobalConfig.Headers = map[string]string{"Test": "123"}
	go_requests.GlobalConfig.Cookies = map[string]string{"sid": "hhh"}
	go_requests.GlobalConfig.Timeout = 10000 // 10秒

	r := go_requests.Request{
		Method: "get",
		Url:    "/get",
	}
	resp := r.Send()
	fmt.Printf("响应文本：%s\n", resp.Text)
}

func TestRequestFromJsonFile(t *testing.T) {
	r := go_requests.GetRequestFromJsonFile("./testdata/data.json")
	resp := r.Send()
	fmt.Printf("状态码: %d\n", resp.StatusCode)
	fmt.Printf("原因: %s\n", resp.Reason)
	fmt.Printf("响应时间: %f秒\n", resp.Elapsed)
	fmt.Printf("响应文本: %s\n", resp.Text)
}

func TestRequestWithHttp2(t *testing.T) { // todo 换其他方式验证
	r := go_requests.Request{Url: "https://stackoverflow.com", HTTP2: true, NoVerify: true}
	resp := r.Send()
	fmt.Printf("响应头: %v\n", resp.Headers)
}

func TestRequestWithProxy(t *testing.T) { // todo 换其他方式验证
	r := go_requests.Request{Url: "https://httpbin.org/get", Proxy: "http://localhost:8888", NoVerify: true}
	resp := r.Send()
	fmt.Printf("状态码: %d\n", resp.StatusCode)
}

func TestAsyncSendRequest(t *testing.T) {
	r := go_requests.Request{Url: "https://www.baidu.com"}
	for i := 0; i < 10; i++ {
		r.AsyncSend()
		resp := <-go_requests.Ch
		fmt.Println(resp.StatusCode)
	}
}

func TestParseJsonResponse(t *testing.T) {
	r := go_requests.Request{
		Method: "get",
		Url:    "https://httpbin.org/get?name=张三&age=12"}
	resp := r.Send()
	respJson := resp.Json()
	//fmt.Println(resp.Text)
	args := respJson["args"].(map[string]interface{})
	name := args["name"].(string)
	age := args["age"].(string)
	fmt.Println(name, age)
}

func TestParseJsonResponseToStruct(t *testing.T) {
	type Args struct {
		Name string `json:"name"`
		Age  string `json:"age"`
	}

	type Headers struct {
		Accept_Encoding string `json:"Accept-Encoding"`
		Host            string `json:"Host"`
		User_Agent      string `json:"User-Agent"`
		X_Amzn_Trace_Id string `json:"X-Amzn-Trace-Id"`
	}

	type MyResponse struct {
		Args    Args    `json:"args"`
		Headers Headers `json:"headers"`
		Origin  string  `json:"origin"`
		Url     string  `json:"url"`
	}

	r := go_requests.Request{
		Method: "get",
		Url:    "https://httpbin.org/get?name=张三&age=12"}
	resp := r.Send()
	fmt.Println(resp.Text)
	//
	var respObj MyResponse
	err := json.Unmarshal(resp.Content, &respObj)
	if err != nil {
		fmt.Println("JSON反序列化失败")
	}
	name := respObj.Args.Name
	age := respObj.Args.Age
	fmt.Println(name, age)
}
