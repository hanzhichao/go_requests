package requests

func Get(url string) *Response {
	return NewSession(nil).Get(url)
}

func GetWithParams(url string, params map[string]string) *Response {
	return NewSession(nil).GetWithParams(url, params)
}

func PostAsForm(url string, data map[string]string) *Response {
	return NewSession(nil).PostAsForm(url, data)
}

func PostAsJson(url string, json string) *Response {
	return NewSession(nil).PostAsJson(url, json)
}

func PostAsRaw(url, raw, contentType string) *Response {
	return NewSession(nil).PostAsRaw(url, raw, contentType)
}

func PostAsMultipartForm(url string, data map[string]string, files map[string]string) *Response {
	return NewSession(nil).PostAsMultipartForm(url, data, files)
}
