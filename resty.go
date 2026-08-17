package telegram_api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/bytedance/sonic"
	"github.com/go-resty/resty/v2"
)

func newRestyClient(cfg *Config) *resty.Client {
	debug := cfg.Debug

	client := resty.New()

	// 基础配置
	client.SetTimeout(cfg.GetTimeout())
	client.SetHeader("User-Agent", cfg.GetUserAgent())
	client.SetHeader("Content-Type", "application/json")

	// 重试配置
	client.SetRetryCount(3)
	client.SetRetryWaitTime(1 * time.Second)
	client.AddRetryCondition(func(resp *resty.Response, err error) bool {
		if err != nil {
			return true
		}
		code := resp.StatusCode()
		return code >= 500 || code == http.StatusTooManyRequests
	})

	// OnBeforeRequest 钩子 - 请求发送前处理
	client.OnBeforeRequest(func(c *resty.Client, req *resty.Request) error {
		// 如果有调试模式，记录请求信息
		if debug {
			logger.Debug(fmt.Sprintf("Request: %s %s", req.Method, req.URL))
			if req.Body != nil {
				logger.Debug(fmt.Sprintf("Request Body:\n%s", formatRequestBody(req.Body)))
			}
		}

		return nil
	})

	// OnAfterResponse 钩子 - 请求响应后处理
	client.OnAfterResponse(func(c *resty.Client, resp *resty.Response) error {
		elapsed := time.Since(resp.Request.Time)

		// 记录请求信息
		if debug {
			logger.Debug(fmt.Sprintf("Response: %d %s (took %v)", resp.StatusCode(), resp.Request.URL, elapsed))
			if resp.Body() != nil && len(resp.Body()) > 0 {
				logger.Debug(fmt.Sprintf("Response Body: %s", formatRequestBody(resp.Body())))
			}
		}

		// 如果响应状态码 >= 500，记录错误日志 400 错误多半是业务逻辑
		if resp.StatusCode() >= 500 {
			logger.Errorf("HTTP Error: %d, URL: %s, Duration: %v", resp.StatusCode(), resp.Request.URL, elapsed)
		}

		return nil
	})

	// OnError 钩子 - 请求错误处理
	client.OnError(func(req *resty.Request, err error) {
		if debug {
			logger.Errorf("Request Error: %v, URL: %s", err, req.URL)
		}
	})

	// 设置日志级别
	if debug {
		client.SetLogger(logger)
	}

	return client
}

func formatRequestBody(body any) string {
	switch v := body.(type) {
	case string:
		// 尝试解析并格式化 JSON 字符串
		var raw any
		if err := sonic.Unmarshal([]byte(v), &raw); err == nil {
			if formatted, err := sonic.MarshalIndent(raw, "", "  "); err == nil {
				return string(formatted)
			}
		}
		return v

	case []byte:
		var raw any
		if err := sonic.Unmarshal(v, &raw); err == nil {
			if formatted, err := sonic.MarshalIndent(raw, "", "  "); err == nil {
				return string(formatted)
			}
		}
		return string(v)

	default:
		// 使用 sonic 序列化
		if prettyJSON, err := sonic.MarshalIndent(v, "", "  "); err == nil {
			return string(prettyJSON)
		}
		return fmt.Sprintf("%+v", v)
	}
}
