package requests

import "fmt"

type Step struct {
	Name string `json:"name"` // 步骤名称
	Request *Request `json:"request"`   // 请求数据
	Register map[string]string `json:"register"`  // 注册提取的变量
	Verify  []map[string][]string `json:"verify"`  // 断言表达式列表
}

type TestCase struct {
	Name string `json:"name"` // 用例名称
	Priority string 	`json:"priority"` // 用例优先级
	Tags []string `json:"tags"` // 用例标签
	Steps []Step `json:"steps"` // 测试步骤
	Setups []Step `json:"setups"`  // 准备铺在
	TearDowns []Step `json:"tear_downs"` // 清理步骤
}

type TestSuite struct {
	Name string `json:"name"` // 测试套件名称
	Config *Config `json:"config"` // 全局请求配置
	Variables map[string]string  `json:"variables"`  // 全局变量
	TestCases []TestCase `json:"test_cases"`  // 用例列表
}


// 运行步骤
func (s *Step) RunStep(){
	fmt.Printf("  执行步骤：%s\n", s.Name)
	s.Request.Send()
}


// 运行测试用例
func (tc *TestCase) RunTestCase(){
	fmt.Printf(" 执行用例：%s\n", tc.Name)
	for _, step := range(tc.Steps){
		step.RunStep()
	}
}

// 运行测试套件
func (ts *TestSuite) RunTestSuite(){
	fmt.Printf("执行测试套件：%s\n", ts.Name)
	for _, testcase := range(ts.TestCases){
		testcase.RunTestCase()
	}

}