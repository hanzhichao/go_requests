package requests

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestBuildTestSuite(t *testing.T) {
	name := "测试套件1"
	config := Config{BaseUrl: "https://httpbin.org", Headers: map[string]string{"token": "abc123"}}
	variables := map[string]string{"a": "100", "b": "200"} // todo 支持字符串以外其他类型
	testcase1_step1 := Step{
		Name:     "步骤1-GET请求",
		Request:  &Request{Method: "GET", Url: "https://httpbin.org/get"},
		Register: map[string]string{"u": "content.url"},
		Verify:   []map[string][]string{{"eq": []string{"status_code", "200"}}},
	}

	testcase1_step2 := Step{
		Name: "步骤1-GET请求",
		Request: &Request{Method: "POST", Url: "https://httpbin.org/post",
			Data: map[string]string{"a": "$a", "b": "2", "c": "$c"}},
	}

	testcase1 := TestCase{
		Name:     "测试用例1",
		Priority: "P0",
		Tags:     []string{"demo"},
		Steps:    []Step{testcase1_step1, testcase1_step2},
	}

	testsuite := TestSuite{
		Name:      name,
		Config:    &config,
		Variables: variables,
		TestCases: []TestCase{testcase1},
	}
	data, _ := json.MarshalIndent(testsuite, "", "  ")
	fmt.Println(string(data))
}

