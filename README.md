# requests
Golang版人性化HTTP请求库，灵感来自Python版requests库

## 特性
- 支持GET、POST等各种请求方法，支持默认请求方法
- 支持独立Query参数、自定义Headers及自定义Cookies
- 支持JSON、表单、`mutipart/form-data`及Raw格式数据
- 支持HTTP2.0及跳过TLS服务端证书验证
- 支持HTTP请求代理
- 支持请求Timeout
- 支持NoRedirects禁止重定向
- 支持BasicAuth基础授权
- 支持GlobalConfig配置默认BasicUrl、默认Params、Headers、Cookies及Auth
- 响应支持状态码、原因、响应二进制内容、响应文本、响应头、Cookies、请求耗时及JSON转map
- 支持从JSON及JSON文件中读取请求配置并发送
- 支持异步请求及并发


## 安装方法
```shell
$ go get -u "github.com/hanzhichao/requests"
```

## 结构模型
### 请求
```go
// 请求结构体
type Request struct {
    Method      string            `json:"method"`          // 请求方法
    Url         string            `json:"url"`             // 请求url
    Params      map[string]string `json:"params"`          // Query参数
    Headers     map[string]string `json:"headers"`         // 请求头
    Cookies     map[string]string `json:"cookies"`         // Cookies
    Data        map[string]string `json:"data"`            // 表单格式请求数据
    Json        string            `json:"json"`            // JSON格式请求数据
    Files       map[string]string `json:"files"`           // mutipart/form-data需要上传的文体 Files
    Raw         string            `json:"raw"`             // 原始请求数据
    Auth        []string          `json:"auth"`            // BaseAuth授权用户名及密码
    Proxy       string            `json:"proxy"`           // 代理地址 例如  "http://127.0.0.1:8888"
    Timeout     int               `json:"timeout"`         // 超时时间，单位 毫秒
    NoRedirects bool              `json:"allow_redirects"` // 关闭重定向, 默认开启
    NoVerify    bool              `json:"no_verify"`       // 跳过TLS证书验证，默认不跳过
    HTTP2       bool              `json:"http_2"`          // 是否启用HTTP2，默认不启用，受GlobalConfig影响
}
```

### 响应
```go
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
```


### 全局请求配置
```go
// 请求配置
type Config struct {
    BaseUrl string            `json:"base_url"` // 基础url
    Params  map[string]string `json:"params"`   // 默认Query参数
    Headers map[string]string `json:"headers"`  // 默认请求头
    Cookies map[string]string `json:"cookies"`  // 默认请求头
    Auth    []string          `json:"auth"`     // 默认BasicAuth授权用户名及密码
    Timeout int               `json:"timeout"`  // 默认超时时间，单位 毫秒
    HTTP2   bool              `json:"http_2"`   // 是否默认启用HTTP2，默认不启用
    Proxy   string            `json:"proxy"`    // 默认代理地址 例如  "http://127.0.0.1:8888"
    // todo 暴露跟多 http.Transport 所需配置
}
```


## 使用示例
> 需要`import "github.com/hanzhichao/requests"`

### 默认请求方法
> Method可以省略，有数据时默认请求方法为POST，否则默认请方法为GET

```go
// 默认请求方法
func TestRequestWithDefaultMethod(t *testing.T) {
    // 发送GET请求
    r := requests.Request{Url: "https://httpbin.org/get"}
    resp := r.Send()
    fmt.Printf("状态码: %d\n", resp.StatusCode)

    // 发送POST请求
    r = requests.Request{Url: "https://httpbin.org/post", Data: map[string]string{"name": "Kevin"}}
    resp = r.Send()
    fmt.Printf("状态码: %d\n", resp.StatusCode)
}
```

### 发送GET请求
```go
func TestGet(t *testing.T) {
    r := requests.Request{
        Method:  "get", 
        Url:     "https://httpbin.org/get?name=张三&age=12"}
    
    resp := r.Send()
    fmt.Printf("状态码: %d\n", resp.StatusCode)
    fmt.Printf("原因: %s\n", resp.Reason)
    fmt.Printf("响应时间: %f秒\n", resp.Elapsed)
    fmt.Printf("响应文本: %s\n", resp.Text)
}
```
### 发送GET 带单独Query参数请求
```go
func TestGetWithParams(t *testing.T) {
    r := requests.Request{
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
```


### 发送POST 表单请求 并携带自定义Headers
```go
func TestPostForm(t *testing.T) {
    r := requests.Request{
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
```


### 发送POST JSON请求
```go
func TestPostJson(t *testing.T) {
    r := Request{
        Method:  "POST",
        Url:     "https://httpbin.org/post",
        Json:    `{"name": "张三", "age": "12"}`}
    
    resp := r.Send()
    fmt.Printf("状态码: %d\n", resp.StatusCode)
    fmt.Printf("原因: %s\n", resp.Reason)
    fmt.Printf("响应时间: %f秒\n", resp.Elapsed)
    fmt.Printf("响应文本: %s\n", resp.Text)
}
```

### 发送POST XML请求
```go
func TestPostXML(t *testing.T) {
    r := requests.Request{
        Method:  "POST",
        Url:     "https://httpbin.org/post",
        Raw:    `<xml>hello</xml>`,
        Headers: map[string]string{"Content-Type": "application/xml"}}
    
    resp := r.Send()
    fmt.Printf("状态码: %d\n", resp.StatusCode)
    fmt.Printf("原因: %s\n", resp.Reason)
    fmt.Printf("响应时间: %f秒\n", resp.Elapsed)
    fmt.Printf("响应文本: %s\n", resp.Text)
}
```


### 发送multipart/form-data请求
```go
func TestPostMultipartFormData(t *testing.T) {
    r := requets.Request{
        Method:  "POST",
        Url:     "https://httpbin.org/post",
        Data:    map[string]string{"name": "张三", "age": "12"},
        Files:   map[string]string{"pic": "./testdata/logo.png"},
    }
    
    resp := r.Send()
    fmt.Printf("状态码: %d\n", resp.StatusCode)
    fmt.Printf("原因: %s\n", resp.Reason)
    fmt.Printf("响应时间: %f秒\n", resp.Elapsed)
    fmt.Printf("响应文本: %s\n", resp.Text)
}

```


## 发送BasicAuth请求
```go
func TestRequestWithAuth(t *testing.T){
    r := requests.Request{
        Method: "GET",
        Url: "https://httpbin.org/get",
        Auth: []string{"Kevin", "123456"},
    }
    
    resp := r.Send()
    fmt.Printf("响应文本：%s\n", resp.Text)
}

```


## 发送禁止重定向请求并获取响应Cookies
```go
func TestNotAllowRedirests(t *testing.T) {
    r := requests.Request{
        Method: "POST",
        Url: "https://newecshop.longtest.cn/admin/privilege.php",
        Data: map[string]string{
                "username": "***",
                "password": "***",
                "act":      "signin"},
        NoRedirects: true,
    }
    
    resp := r.Send()
    fmt.Printf("状态码: %d\n", resp.StatusCode)
    for key, value := range(resp.Cookies){
    fmt.Printf("响应Cookies项：%s=%s\n", key, value)
    }
}
```

### 请求超时时间设置
```go
func TestRequestTimeout(t *testing.T){
    r := requests.Request{
        Method: "get",
        Url: "https://httpbin.org/get",
        Timeout: 1, // 毫秒
    }
    
    r.Send()
}

```

### 使用请求默认配置
```go
func TestRequestConfig(t *testing.T){
    requests.GlobalConfig.BaseUrl = "https://httpbin.org"
    requests.GlobalConfig.Params = map[string]string{"Token": "abc"}
    requests.GlobalConfig.Headers = map[string]string{"Test": "123"}
    requests.GlobalConfig.Cookies = map[string]string{"sid": "hhh"}
    requests.GlobalConfig.Timeout = 10000 // 10秒

    r := requests.Request{
        Method: "get",
        Url: "/get",
    }
    resp := r.Send()
    fmt.Printf("响应文本：%s\n", resp.Text)
}

```

### 读取JSON文件发送请求
**testdata/data.json**内容
```json
{
  "method": "post",
  "url": "https://httpbin.org/post",
  "data": {"name": "Kevin", "age":  "12"}
}
```
> 注意：data中的value值如`age`必须是string类型，不然会反序列化失败


```go
func TestRequestFromJsonFile(t *testing.T) {
    r := requets.GetRequestFromJsonFile("./testdata/data.json")
    resp := r.Send()
    fmt.Printf("状态码: %d\n", resp.StatusCode)
    fmt.Printf("原因: %s\n", resp.Reason)
    fmt.Printf("响应时间: %f秒\n", resp.Elapsed)
    fmt.Printf("响应文本: %s\n", resp.Text)
}
```

### 发送HTTP2请求
```go
// 使用HTTP2及关闭TLS验证
func TestRequestWithHttp2(t *testing.T){  // todo 换其他方式验证
    r := requests.Request{Url: "https://stackoverflow.com", HTTP2: true, NoVerify: true}
    resp := r.Send()
    fmt.Printf("响应头: %v\n", resp.Headers)
}
```
> 可以通过调试在原始http.Resposne对象res.Proto属性中查看到请求协议为HTTP/2.0

### 使用HTTP代理
```go
func TestRequestWithProxy(t *testing.T){  // todo 换其他方式验证
    r := requests.Request{Url: "https://httpbin.org/get", Proxy: "http://localhost:8888", NoVerify: true}
    resp := r.Send()
    fmt.Printf("状态码: %d\n", resp.StatusCode)
}
```
### 使用异步请求
```go
// 异步发送请求
func TestAsyncSendRequest(t *testing.T) {
	r := requests.Request{Url: "https://www.baidu.com"}
	for i := 0; i < 10; i++ {
		r.AsyncSend()
		resp := <- requests.Ch
		fmt.Println(resp.StatusCode)
	}
}
```

## ToDo
- [ ] 异常处理
- [x] 异步请求
- [x] 并发请求
- [ ] SSL验证/关闭验证
- [ ] Runner实现
- [ ] 性能测试及指标计算
- [ ] HTTP3
- [ ] 支持WebSocket
- [ ] 异步请求并发配置

## 已知问题
- [ ] 不支持流式发送单文件binary数据
- [ ] 无法自定义Transport配置，无法添加个人TLS证书及密钥
- [ ] 无法获取响应HTTP版本
- [ ] 不支持国密TLS

## 参考
- <https://pkg.go.dev/golang.org/x/net/http2#ConfigureTransport>
- <https://httpwg.org/specs/rfc7540.html>
