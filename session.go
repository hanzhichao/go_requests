package requests

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

func (s *Session) Get(url string, headers map[string]string) *Response {
	req := NewRequestWithConfig(s.Config, "GET", url).
		SetHeaders(headers)
	return s.SendRequest(req)
}

func (s *Session) Post(url, data string, headers map[string]string) *Response {
	req := NewRequestWithConfig(s.Config, "POST", url).
		SetRawData(data).
		SetHeaders(headers)
	return req.Send()
}

func (s *Session) Put(url, data string, headers map[string]string) *Response {
	req := NewRequestWithConfig(s.Config, "Put", url).
		SetRawData(data).
		SetHeaders(headers)
	return req.Send()
}

func (s *Session) Delete(url string, headers map[string]string) *Response {
	req := NewRequestWithConfig(s.Config, "DELETE", url).
		SetHeaders(headers)
	return s.SendRequest(req)
}

func (s *Session) Head(url string, headers map[string]string) *Response {
	req := NewRequestWithConfig(s.Config, "HEAD", url).
		SetHeaders(headers)
	return s.SendRequest(req)
}

func (s *Session) Options(url string, headers map[string]string) *Response {
	req := NewRequestWithConfig(s.Config, "OPTIONS", url).
		SetHeaders(headers)
	return s.SendRequest(req)
}

func Get(url string, headers map[string]string) *Response {
	return NewSession(nil).Get(url, headers)
}

func Post(url, data string, headers map[string]string) *Response {
	return NewSession(nil).Post(url, data, headers)
}

func Put(url, data string, headers map[string]string) *Response {
	return NewSession(nil).Post(url, data, headers)
}

func Delete(url string, headers map[string]string) *Response {
	return NewSession(nil).Get(url, headers)
}

func Head(url string, headers map[string]string) *Response {
	return NewSession(nil).Head(url, headers)
}

func Options(url string, headers map[string]string) *Response {
	return NewSession(nil).Options(url, headers)
}
