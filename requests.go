package requests

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"golang.org/x/net/http2"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

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
	// todo 暴露跟多 http.Transport 所需配置/**/
}

var GlobalConfig = &Config{} // 全局配置
var Wait sync.WaitGroup
var Ch chan Response = make(chan Response, 10)

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

func (rs *Response) Json() map[string]interface{} {
	data := make(map[string]interface{})
	err := json.Unmarshal(rs.Content, &data)
	if err != nil {
		fmt.Println("响应文本转map失败")
	}
	return data
}

func (rs *Response) JsonArray() []map[string]interface{} {
	var data []map[string]interface{}
	err := json.Unmarshal(rs.Content, &data)
	if err != nil {
		fmt.Println("响应文本转map失败")
	}
	return data
}

// 处理Config配置
func (r *Request) handleConfig() {
	// 处理BaseUrl
	if GlobalConfig.BaseUrl != "" && !strings.HasPrefix(r.Url, "http") {
		r.Url = fmt.Sprintf("%s%s", GlobalConfig.BaseUrl, r.Url)
	}
	// 处理默认Params
	if GlobalConfig.Params != nil {
		if r.Params != nil {
			for key, value := range GlobalConfig.Params {
				if _, ok := r.Params[key]; !ok {
					r.Params[key] = value
				}
			}
		}
		r.Params = GlobalConfig.Params
	}
	// 处理默认请求头
	if GlobalConfig.Headers != nil {
		if r.Headers != nil {
			for key, value := range GlobalConfig.Headers {
				if _, ok := r.Headers[key]; !ok {
					r.Headers[key] = value
				}
			}
		}
		r.Headers = GlobalConfig.Headers
	}
	// 处理默认Cookies
	if GlobalConfig.Cookies != nil {
		if r.Cookies != nil {
			for key, value := range GlobalConfig.Cookies {
				if _, ok := r.Cookies[key]; !ok {
					r.Cookies[key] = value
				}
			}
		}
		r.Cookies = GlobalConfig.Cookies
	}
	// 处理默认Timeout
	if GlobalConfig.Timeout > 0 && r.Timeout == 0 {
		r.Timeout = GlobalConfig.Timeout
	}
	// 处理默认BasicAuth
	if GlobalConfig.Auth != nil && len(GlobalConfig.Auth) == 2 {
		r.Auth = GlobalConfig.Auth
	}
	// 处理默认是否开启HTTP2
	if GlobalConfig.HTTP2 == true {
		r.HTTP2 = true
	}
	// 处理默认Proxy配置
	if GlobalConfig.Proxy != "" {
		r.Proxy = GlobalConfig.Proxy
	}
}

// 处理请求方法
func (r *Request) getMethod() string {
	if r.Method == "" {
		if r.Raw == "" && r.Json == "" && (r.Data == nil || len(r.Data) == 0) && (r.Files == nil || len(r.Files) == 0) {
			r.Method = "GET" // 无任何数据是默认请求方法GET
		} else {
			r.Method = "POST" // 有数据是默认请求方法是POST
		}
	}
	return strings.ToUpper(r.Method) // 必须转为全部大写
}

// 组装URL
func (r *Request) getUrl() string {
	if r.Params != nil && len(r.Params) > 0 {
		urlValues := url.Values{}
		Url, err := url.Parse(r.Url) // todo 处理err
		if err != nil {
			fmt.Printf("解析URL\"%s\"失败: %s\n", r.Url, err)
		}
		for key, value := range r.Params {
			urlValues.Set(key, value)
		}
		Url.RawQuery = urlValues.Encode()
		return Url.String()
	}
	return r.Url
}

// 组装请求数据
func (r *Request) getData() io.Reader {
	var reqBody string
	if r.Headers == nil {
		r.Headers = map[string]string{}
	}
	// 处理Raw格式请求，需要自行添加请求头Content-Type
	if r.Raw != "" {
		reqBody = r.Raw
		return strings.NewReader(reqBody)
	}
	// 处理multipart/formdata
	if r.Files != nil && len(r.Files) > 0 {
		// 实例化multipart
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)

		// 处理r.Data中的字符串参数
		if r.Data != nil && len(r.Data) > 0 {
			for key, value := range r.Data {
				_ = writer.WriteField(key, value)
			}
		}
		// 处理r.Files中的文件路径
		for key, filePath := range r.Files {
			item, _ := writer.CreateFormFile(key, filepath.Base(filePath))
			file, _ := os.Open(filePath)
			defer file.Close()
			_, err := io.Copy(item, file)
			if err != nil {
				fmt.Println("复制文件内容失败")
			}
		}
		// 关闭writer
		err := writer.Close()
		if err != nil {
			fmt.Println("关闭writer出错")
		}
		r.Headers["Content-Type"] = writer.FormDataContentType()
		return body
	}

	// 处理application/x-www-form-urlencoded
	if r.Data != nil && len(r.Data) > 0 {
		urlValues := url.Values{}
		for key, value := range r.Data {
			urlValues.Add(key, value)
		}
		reqBody = urlValues.Encode()
		r.Headers["Content-Type"] = "application/x-www-form-urlencoded"
		return strings.NewReader(reqBody)
	}

	// 处理application/json
	if r.Json != "" {
		//bytesData, _ := json.Marshal(r.Json)
		reqBody = r.Json
		r.Headers["Content-Type"] = "application/json"
		return strings.NewReader(reqBody)
	}

	return strings.NewReader(reqBody)
}

// 添加请求头-需要在getData后使用
func (r *Request) addHeaders(req *http.Request) {
	if r.Headers != nil {
		for key, value := range r.Headers {
			req.Header.Add(key, value)
		}
	}
}

// 添加Cookies
func (r *Request) addCookies(req *http.Request) {
	if r.Cookies != nil {
		for key, value := range r.Cookies {
			req.AddCookie(&http.Cookie{Name: key, Value: value})
		}
	}
}

// 处理Auth
func (r *Request) setAuth(req *http.Request) {
	if r.Auth != nil && len(r.Auth) == 2 {
		username := r.Auth[0]
		password := r.Auth[1]
		req.SetBasicAuth(username, password)
	}
}

// 准备请求
func (r *Request) prepare() *http.Request {
	r.handleConfig()
	Method := r.getMethod()
	Url := r.getUrl()
	Data := r.getData()
	req, err := http.NewRequest(Method, Url, Data)
	if err != nil {
		fmt.Printf("构造请求失败: %s\n", err)
	}
	r.addHeaders(req)
	r.addCookies(req)
	r.setAuth(req)
	return req
}

// 组装响应对象
func (r *Request) buildResponse(res *http.Response, elapsed float64) Response {
	var resp Response
	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("读取响应数据失败: %s\n", err)
	}
	resp.Content = resBody
	resp.Text = string(resBody)
	resp.StatusCode = res.StatusCode
	resp.Reason = strings.Split(res.Status, " ")[1]
	resp.Elapsed = elapsed
	resp.Headers = map[string]string{}
	for key, value := range res.Header {
		resp.Headers[key] = strings.Join(value, ";")
	}
	resp.Cookies = map[string]string{}
	for _, item := range res.Cookies() {
		resp.Cookies[item.Name] = item.Value
	}
	return resp
}

func (r *Request) getClient() *http.Client {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: r.NoVerify}, // 是否跳过验服务端证书
	}
	// 处理Proxy
	if r.Proxy != "" {
		proxy, err := url.Parse(r.Proxy)
		if err != nil {
			fmt.Printf("解析代理地址 \"%s\" 出错: %s\n", r.Proxy, err)
		}
		transport.Proxy = http.ProxyURL(proxy)
	}

	// 处理是否HTTP2
	if r.HTTP2 == true {
		err := http2.ConfigureTransport(transport)
		if err != nil {
			fmt.Printf("HTTP2传输配置出错: %s\n", err)
		}
		//transport.AllowHTTP = true
	}

	client := &http.Client{Transport: transport}
	if r.Timeout > 0 {
		client.Timeout = time.Duration(r.Timeout) * time.Millisecond
	}
	if r.NoRedirects == true {
		client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}
	}
	return client
}

// 发送请求
func (r *Request) Send() Response {
	req := r.prepare()
	client := r.getClient()
	start := time.Now()
	res, err := client.Do(req)
	if err != nil {
		fmt.Printf("发送请求失败: %s\n", err)
	}
	defer res.Body.Close()
	elapsed := time.Since(start).Seconds()
	Wait.Add(1)
	resp := r.buildResponse(res, elapsed)
	Ch <- resp
	Wait.Done()
	return resp
}

// 发送异步请求
func (r *Request) AsyncSend() {
	go r.Send()
	Wait.Wait()
}

// 从JSON字符串得到Request结构体
func GetRequestFromJson(jsonData []byte) Request {
	r := Request{
		Params:  map[string]string{},
		Headers: map[string]string{},
		Data:    map[string]string{},
		Files:   map[string]string{}}
	err := json.Unmarshal(jsonData, &r)
	if err != nil {
		fmt.Printf("反序列化出错: %s\n", err)
	}
	return r
}

// 从JSON文件得到Request结构体
func GetRequestFromJsonFile(jsonFilePath string) Request {
	jsonData, err := ioutil.ReadFile(jsonFilePath)
	if err != nil {
		fmt.Printf("读取JSON文件出错: %s\n", err)
	}
	return GetRequestFromJson(jsonData)
}
