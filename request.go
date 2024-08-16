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

var Wait sync.WaitGroup
var Ch chan Response = make(chan Response, 10)

// Request 请求结构体
type Request struct {
	Config *Config `json:"config"` // 请求配置

	Method      string            `json:"method"`          // 请求方法
	Url         string            `json:"url"`             // 请求url
	Params      map[string]string `json:"params"`          // Query参数
	Headers     map[string]string `json:"headers"`         // 请求头
	Cookies     map[string]string `json:"cookies"`         // Cookies
	Data        map[string]string `json:"data"`            // 表单格式请求数据
	Json        string            `json:"json"`            // JSON格式请求数据
	Files       map[string]string `json:"files"`           // mutipart/form-data需要上传的文体 Files
	Raw         string            `json:"raw"`             // 原始请求数据
	Auth        []string          `json:"auth"`            // BasicAuth授权用户名及密码
	Proxy       string            `json:"proxy"`           // 代理地址 例如  "http://127.0.0.1:8888"
	Timeout     int               `json:"timeout"`         // 超时时间，单位 毫秒
	NoRedirects bool              `json:"allow_redirects"` // 关闭重定向, 默认开启
	NoVerify    bool              `json:"no_verify"`       // 跳过TLS证书验证，默认不跳过
	HTTP2       bool              `json:"http_2"`          // 是否启用HTTP2，默认不启用，受Config影响
}

func NewRequest(method, url string) *Request {
	return &Request{Method: method, Url: url}
}

func NewRequestWithConfig(config *Config, method, url string) *Request {
	return &Request{Config: config, Method: method, Url: url}
}

func (req *Request) SetParams(params map[string]string) *Request {
	req.Params = params
	return req
}

func (req *Request) SetHeaders(headers map[string]string) *Request {
	if req.Headers == nil {
		req.Headers = headers
	}
	updateMap(req.Headers, headers)
	return req
}

func (req *Request) SetCookies(cookies map[string]string) *Request {
	req.Cookies = cookies
	return req
}

func (req *Request) SetFormData(data map[string]string) *Request {
	req.Data = data
	return req
}

func (req *Request) SetJsonData(json string) *Request {
	req.Json = json
	return req
}

func (req *Request) SetRawData(raw string) *Request {
	req.Raw = raw
	return req
}

func (req *Request) SetUploadFiles(files map[string]string) *Request {
	req.Files = files
	return req
}

func (req *Request) SetBasicAuth(auth []string) *Request {
	req.Auth = auth
	return req
}

func (req *Request) SetTimeout(timeout int) *Request {
	req.Timeout = timeout
	return req
}

func (req *Request) SetProxy(proxy string) *Request {
	req.Proxy = proxy
	return req
}

func (req *Request) SetNoRedirects(enable bool) *Request {
	req.NoRedirects = enable
	return req
}

func (req *Request) SetNoVerify(enable bool) *Request {
	req.NoVerify = enable
	return req
}

func (req *Request) EnableHttp2(enable bool) *Request {
	req.HTTP2 = enable
	return req
}

func (req *Request) SetContentType(contentType string) *Request {
	req.SetHeaders(map[string]string{"Content-Type": contentType})
	return req
}

func (req *Request) SetBearerToken(token string) *Request {
	req.SetHeaders(map[string]string{"Authorization": fmt.Sprintf("Bearer %s", token)})
	return req
}

// 处理Config配置
func (req *Request) handleConfig() {
	// 处理BaseUrl
	config := req.Config
	if config == nil {
		return
	}
	if config.BaseUrl != "" && !strings.HasPrefix(req.Url, "http") {
		req.Url = fmt.Sprintf("%s%s", config.BaseUrl, req.Url)
	}
	// 处理默认Params
	if config.Params != nil {
		if req.Params != nil {
			for key, value := range config.Params {
				if _, ok := req.Params[key]; !ok {
					req.Params[key] = value
				}
			}
		}
		req.Params = config.Params
	}
	// 处理默认请求头
	if config.Headers != nil {
		if req.Headers != nil {
			for key, value := range config.Headers {
				if _, ok := req.Headers[key]; !ok {
					req.Headers[key] = value
				}
			}
		}
		req.Headers = config.Headers
	}
	// 处理默认Cookies
	if config.Cookies != nil {
		if req.Cookies != nil {
			for key, value := range config.Cookies {
				if _, ok := req.Cookies[key]; !ok {
					req.Cookies[key] = value
				}
			}
		}
		req.Cookies = config.Cookies
	}
	// 处理默认Timeout
	if config.Timeout > 0 && req.Timeout == 0 {
		req.Timeout = config.Timeout
	}
	// 处理默认BasicAuth
	if config.Auth != nil && len(config.Auth) == 2 {
		req.Auth = config.Auth
	}
	// 处理默认是否开启HTTP2
	if config.HTTP2 == true {
		req.HTTP2 = true
	}
	// 处理默认Proxy配置
	if config.Proxy != "" {
		req.Proxy = config.Proxy
	}
}

// 处理请求方法
func (req *Request) getMethod() string {
	if req.Method == "" {
		if req.Raw == "" && req.Json == "" && (req.Data == nil || len(req.Data) == 0) && (req.Files == nil || len(req.Files) == 0) {
			req.Method = "GET" // 无任何数据是默认请求方法GET
		} else {
			req.Method = "POST" // 有数据是默认请求方法是POST
		}
	}
	return strings.ToUpper(req.Method) // 必须转为全部大写
}

// 组装URL
func (req *Request) getUrl() string {
	if req.Params != nil && len(req.Params) > 0 {
		urlValues := url.Values{}
		Url, err := url.Parse(req.Url) // todo 处理err
		if err != nil {
			fmt.Printf("解析URL\"%s\"失败: %s\n", req.Url, err)
		}
		for key, value := range req.Params {
			urlValues.Set(key, value)
		}
		Url.RawQuery = urlValues.Encode()
		return Url.String()
	}
	return req.Url
}

// 组装请求数据
func (req *Request) getData() io.Reader {
	var reqBody string
	if req.Headers == nil {
		req.Headers = map[string]string{}
	}
	// 处理Raw格式请求，需要自行添加请求头Content-Type
	if req.Raw != "" {
		reqBody = req.Raw
		return strings.NewReader(reqBody)
	}
	// 处理multipart/formdata
	if req.Files != nil && len(req.Files) > 0 {
		// 实例化multipart
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)

		// 处理r.Data中的字符串参数
		if req.Data != nil && len(req.Data) > 0 {
			for key, value := range req.Data {
				_ = writer.WriteField(key, value)
			}
		}
		// 处理r.Files中的文件路径
		for key, filePath := range req.Files {
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
		req.Headers["Content-Type"] = writer.FormDataContentType()
		return body
	}

	// 处理application/x-www-form-urlencoded
	if req.Data != nil && len(req.Data) > 0 {
		urlValues := url.Values{}
		for key, value := range req.Data {
			urlValues.Add(key, value)
		}
		reqBody = urlValues.Encode()
		req.Headers["Content-Type"] = "application/x-www-form-urlencoded"
		return strings.NewReader(reqBody)
	}

	// 处理application/json
	if req.Json != "" {
		//bytesData, _ := json.Marshal(req.Json)
		reqBody = req.Json
		req.Headers["Content-Type"] = "application/json"
		return strings.NewReader(reqBody)
	}

	return strings.NewReader(reqBody)
}

// 添加请求头-需要在getData后使用
func (req *Request) addHeaders(r *http.Request) {
	if req.Headers != nil {
		for key, value := range req.Headers {
			r.Header.Add(key, value)
		}
	}
}

// 添加Cookies
func (req *Request) addCookies(r *http.Request) {
	if req.Cookies != nil {
		for key, value := range req.Cookies {
			r.AddCookie(&http.Cookie{Name: key, Value: value})
		}
	}
}

// 处理Auth
func (req *Request) setAuth(r *http.Request) {
	if req.Auth != nil && len(req.Auth) == 2 {
		username := req.Auth[0]
		password := req.Auth[1]
		r.SetBasicAuth(username, password)
	}
}

// 准备请求
func (req *Request) prepare() *http.Request {
	req.handleConfig()
	Method := req.getMethod()
	Url := req.getUrl()
	Data := req.getData()
	r, err := http.NewRequest(Method, Url, Data)
	if err != nil {
		fmt.Printf("构造请求失败: %s\n", err)
	}
	req.addHeaders(r)
	req.addCookies(r)
	req.setAuth(r)
	return r
}

// 组装响应对象
func (req *Request) buildResponse(res *http.Response, elapsed float64) Response {
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

func (req *Request) getClient() *http.Client {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: req.NoVerify}, // 是否跳过验服务端证书
	}
	// 处理Proxy
	if req.Proxy != "" {
		proxy, err := url.Parse(req.Proxy)
		if err != nil {
			fmt.Printf("解析代理地址 \"%s\" 出错: %s\n", req.Proxy, err)
		}
		transport.Proxy = http.ProxyURL(proxy)
	}

	// 处理是否HTTP2
	if req.HTTP2 == true {
		err := http2.ConfigureTransport(transport)
		if err != nil {
			fmt.Printf("HTTP2传输配置出错: %s\n", err)
		}
		//transport.AllowHTTP = true
	}

	client := &http.Client{Transport: transport}
	if req.Timeout > 0 {
		client.Timeout = time.Duration(req.Timeout) * time.Millisecond
	}
	if req.NoRedirects == true {
		client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}
	}
	return client
}

// Send 发送请求
func (req *Request) Send() *Response {
	r := req.prepare()
	client := req.getClient()
	start := time.Now()
	res, err := client.Do(r)
	if err != nil {
		fmt.Printf("发送请求失败: %s\n", err)
	}
	defer res.Body.Close()
	elapsed := time.Since(start).Seconds()
	Wait.Add(1)
	resp := req.buildResponse(res, elapsed)
	Ch <- resp
	Wait.Done()

	//// 处理Set-Cookies 修改 Request
	//if resp.Cookies != nil {
	//	req.SetCookies(resp.Cookies)
	//}

	return &resp
}

// AsyncSend 发送异步请求
func (req *Request) AsyncSend() {
	go req.Send()
	Wait.Wait()
}

// GetRequestFromJson 从JSON字符串得到Request结构体
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

// GetRequestFromJsonFile 从JSON文件得到Request结构体
func GetRequestFromJsonFile(jsonFilePath string) Request {
	jsonData, err := ioutil.ReadFile(jsonFilePath)
	if err != nil {
		fmt.Printf("读取JSON文件出错: %s\n", err)
	}
	return GetRequestFromJson(jsonData)
}
