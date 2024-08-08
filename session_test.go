package requests

import (
	"fmt"
	"testing"
)

func TestKeepCookies(t *testing.T) {
	s := NewSession(nil)
	resp := s.PostAsForm("http://127.0.0.1:5000/api/user/login/", map[string]string{"name": "张三", "password": "123456"})

	fmt.Printf("状态码: %d\n", resp.StatusCode)
	fmt.Printf("原因: %s\n", resp.Reason)
	fmt.Printf("响应时间: %f秒\n", resp.Elapsed)
	fmt.Printf("响应文本: %s\n", resp.Text)
	fmt.Printf("响应Cookies: %s\n", resp.Cookies)

	resp2 := s.Get("http://127.0.0.1:5000/api/user/logout/", map[string]string{"name": "张三"})

	fmt.Printf("状态码: %d\n", resp2.StatusCode)
	fmt.Printf("原因: %s\n", resp2.Reason)
	fmt.Printf("响应时间: %f秒\n", resp2.Elapsed)
	fmt.Printf("响应文本: %s\n", resp2.Text)

}
