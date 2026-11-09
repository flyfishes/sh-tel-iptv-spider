package http_client

import (
	"fmt"
	"iptv-spider-sh/global"
	"net"
	"net/http"
	"net/url"
	"strings"

	"github.com/go-resty/resty/v2"
)

type HttpClient struct {
	client    *resty.Client
	userAgent string
	resp      *resty.Response
}

// buildCurl 构造可完整复现本次请求的 curl 命令
func (c *HttpClient) buildCurl(uri, method string, form map[string]string) string {
	values := url.Values{}
	for k, v := range form {
		values.Set(k, v)
	}

	var sb strings.Builder
	sb.WriteString("curl")

	// 方法
	if method == "GET" {
		// GET 时参数拼到 URL query，-G 让 --data-urlencode 走 query
		sb.WriteString(" -G")
	} else {
		sb.WriteString(" -X " + method)
	}

	// URL（单引号包裹，避免 shell 解析特殊字符）
	sb.WriteString(" '" + uri + "'")

	// User-Agent
	sb.WriteString(" -H 'User-Agent: " + c.userAgent + "'")

	// Cookie
	if len(c.client.Cookies) > 0 {
		var cookieParts []string
		for _, ck := range c.client.Cookies {
			cookieParts = append(cookieParts, ck.Name+"="+ck.Value)
		}
		sb.WriteString(" -H 'Cookie: " + strings.Join(cookieParts, "; ") + "'")
	}

	// 参数
	encoded := values.Encode()
	if encoded != "" {
		if method == "GET" {
			// GET：--data-urlencode 配合 -G 自动拼到 query
			for k, vs := range values {
				for _, v := range vs {
					sb.WriteString(fmt.Sprintf(" --data-urlencode '%s=%s'", k, v))
				}
			}
		} else {
			// POST：form body
			sb.WriteString(" -H 'Content-Type: application/x-www-form-urlencoded'")
			sb.WriteString(" --data-raw '" + encoded + "'")
		}
	}

	return strings.ReplaceAll(sb.String(), `'`, `"`)
}

func (c *HttpClient) Request(uri, method string, form map[string]string) *HttpClient {
	r := c.client.R()
	method = strings.ToUpper(method)
	switch method {
	case "GET":
		r.SetQueryParams(form)
	case "POST":
		r.SetFormData(form)
	}
	//global.LOG.Debug(fmt.Sprintf("%s %s Data: %v", method, uri, form))
	global.LOG.Debug(c.buildCurl(uri, method, form))
	var err error
	c.resp, err = r.Execute(method, uri)
	if err != nil {
		global.LOG.Error(fmt.Sprintf("%s", err.Error()))
	}
	//global.LOG.Debug(fmt.Sprintf("Resp Body: %s", string(c.resp.Body())))
	return c
}

func (c *HttpClient) GetResp() *resty.Response {
	return c.resp
}

func (c *HttpClient) GetRespBytes() []byte {
	return c.resp.Body()
}

func NewHttpClient(opts ...HttpClientOption) *HttpClient {
	c := &HttpClient{
		userAgent: "webkit;Resolution(PAL,720P,1080P,2106P,4K)",
	}
	for _, opt := range opts {
		opt(c)
	}
	if c.client == nil {
		c.client = resty.New()
	}
	c.afterAction()
	return c
}

func (c *HttpClient) afterAction() {
	// setUserAgent
	c.client.SetHeader("User-Agent", c.userAgent)
}

func (c *HttpClient) SetCookies(cookies ...*http.Cookie) {
	if len(cookies) <= 0 {
		return
	}
	c.client.SetCookies(cookies)
}

func (c *HttpClient) Cookies() []*http.Cookie {
	return c.client.Cookies
}

type HttpClientOption func(client *HttpClient)

func WithUserAgent(ua string) HttpClientOption {
	return func(c *HttpClient) {
		c.userAgent = ua
	}
}

func WithLocalAddr(addr string) HttpClientOption {
	return func(c *HttpClient) {
		tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
		if err != nil {
			global.LOG.Warn("ResolveTCPAddr failed: " + addr)
		} else {
			global.LOG.Warn("ResolveTCPAddr success: " + addr)
		}
		c.client = resty.NewWithLocalAddr(tcpAddr)
	}
}
