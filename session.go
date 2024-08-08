package requests

const (
	ContentTypeFormUrlEncoded = iota
	ContentTypeApplicationJson
)

type Session struct {
	Config  *Config `json:"config"` // 请求配置
	cookies map[string]string
}

func NewSession(config *Config) *Session {
	return &Session{Config: config}
}

func (s *Session) SendRequest(req *Request) *Response {
	if s.cookies != nil {
		req.Cookies = s.cookies // TODO Merge cookies
	}
	resp := req.Send()
	if resp.Cookies != nil {
		s.cookies = resp.Cookies // TODO Merge cookies
	}
	return resp
}

func (s *Session) Get(url string, params map[string]string) *Response {
	req := NewRequestWithConfig(s.Config).
		SetMethod("GET").
		SetUrl(url).
		SetParams(params)
	return s.SendRequest(req)
}

func (s *Session) PostAsForm(url string, data map[string]string) *Response {
	req := NewRequestWithConfig(s.Config).
		SetMethod("POST").
		SetUrl(url).
		SetFormData(data)
	return s.SendRequest(req)
}

func (s *Session) PostAsJson(method, url string, data map[string]interface{}) *Response {
	req := NewRequestWithConfig(s.Config).
		SetMethod(method).
		SetUrl(url).
		SetJsonData(data)
	return req.Send()
}
