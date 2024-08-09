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
		updateMap(req.Cookies, s.cookies)
	}
	resp := req.Send()
	if resp.Cookies != nil {
		updateMap(s.cookies, resp.Cookies)
		//s.cookies = resp.Cookies // TODO Merge cookies
	}
	return resp
}

func (s *Session) Get(url string) *Response {
	req := NewRequestWithConfig(s.Config).
		SetMethod("GET").
		SetUrl(url)
	return s.SendRequest(req)
}

func (s *Session) GetWithParams(url string, params map[string]string) *Response {
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

func (s *Session) PostAsMultipartForm(url string, data map[string]string, files map[string]string) *Response {
	req := NewRequestWithConfig(s.Config).
		SetMethod("POST").
		SetUrl(url).
		SetFormData(data).
		SetUploadFiles(files)
	return s.SendRequest(req)
}

func (s *Session) PostAsJson(url string, json string) *Response {
	req := NewRequestWithConfig(s.Config).
		SetMethod("POST").
		SetUrl(url).
		SetJsonData(json)
	return req.Send()
}

func (s *Session) PostAsRaw(url, raw, contentType string) *Response {
	req := NewRequestWithConfig(s.Config).
		SetMethod("POST").
		SetUrl(url).
		SetRawData(raw).
		SetContentType(contentType)
	return req.Send()
}
