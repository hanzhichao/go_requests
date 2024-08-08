package requests

// Config 请求配置
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

func NewConfig() *Config {
	return &Config{}
}

func (conf *Config) SetBaseUrl(baseUrl string) *Config {
	conf.BaseUrl = baseUrl
	return conf
}

func (conf *Config) SetParams(params map[string]string) *Config {
	conf.Params = params
	return conf
}

func (conf *Config) SetHeaders(headers map[string]string) *Config {
	conf.Headers = headers
	return conf
}

func (conf *Config) SetCookies(cookies map[string]string) *Config {
	conf.Cookies = cookies
	return conf
}

func (conf *Config) SetBasicAuth(auth []string) *Config {
	conf.Auth = auth
	return conf
}

func (conf *Config) SetTimeout(timeout int) *Config {
	conf.Timeout = timeout
	return conf
}

func (conf *Config) SetProxy(proxy string) *Config {
	conf.Proxy = proxy
	return conf
}

func (conf *Config) EnableHTTP2(enable bool) *Config {
	conf.HTTP2 = enable
	return conf
}
