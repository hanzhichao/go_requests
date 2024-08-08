package requests

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestParseJsonResponse(t *testing.T) {
	r := Request{
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

	r := Request{
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
