package telegram_api

import (
	"context"
	"fmt"
	"time"

	"github.com/bytedance/sonic"
	"github.com/go-resty/resty/v2"
)

type contextKey string

const (
	requestStartTimeKey contextKey = "request_start_time"
)

func setRequestStartTime(ctx context.Context, t time.Time) context.Context {
	return context.WithValue(ctx, requestStartTimeKey, t)
}

func getRequestStartTime(ctx context.Context) (time.Time, bool) {
	val := ctx.Value(requestStartTimeKey)
	if val == nil {
		return time.Time{}, false
	}

	t, ok := val.(time.Time)
	if !ok {
		return time.Time{}, false
	}

	return t, true
}

func setupRestyClient(cfg *Config) *resty.Client {
	client := resty.New()

	// 基础配置
	client.SetTimeout(cfg.GetTimeout())
	client.SetHeader("User-Agent", cfg.GetUserAgent())

	// 重试配置
	client.SetRetryCount(3)
	client.SetRetryWaitTime(1 * time.Second)
	client.AddRetryCondition(func(resp *resty.Response, err error) bool {
		if err != nil {
			return true
		}
		return resp.StatusCode() >= 500
	})

	// OnBeforeRequest 钩子 - 请求发送前处理
	client.OnBeforeRequest(func(c *resty.Client, req *resty.Request) error {
		// 记录请求开始时间
		req.SetContext(setRequestStartTime(req.Context(), time.Now()))

		// 如果有调试模式，记录请求信息
		if cfg.Debug {
			logger.Debug(fmt.Sprintf("Request: %s %s", req.Method, req.URL))
			if req.Body != nil {
				logger.Debug(fmt.Sprintf("Request Body:\n%s", formatRequestBody(req.Body)))
			}
		}

		return nil
	})

	// OnAfterResponse 钩子 - 请求响应后处理
	client.OnAfterResponse(func(c *resty.Client, resp *resty.Response) error {
		// 计算请求耗时
		if startTime, ok := getRequestStartTime(resp.Request.Context()); ok {
			elapsed := time.Since(startTime)

			// 记录请求信息
			if cfg.Debug {
				logger.Debug(fmt.Sprintf("Response: %d %s (took %v)", resp.StatusCode(), resp.Request.URL, elapsed))
				if resp.Body() != nil && len(resp.Body()) > 0 {
					logger.Debug(fmt.Sprintf("Response Body: %s", formatRequestBody(resp.Body())))
				}
			}

			// 如果响应状态码 >= 400，记录错误日志
			if resp.StatusCode() >= 400 {
				logger.Errorf("HTTP Error: %d, URL: %s, Duration: %v", resp.StatusCode(), resp.Request.URL, elapsed)
			}
		}

		return nil
	})

	// OnError 钩子 - 请求错误处理
	client.OnError(func(req *resty.Request, err error) {
		if cfg.Debug {
			logger.Errorf("Request Error: %v, URL: %s", err, req.URL)
		}
	})

	// 设置日志级别
	if cfg.Debug {
		client.SetLogger(logger)
	}

	return client
}

func formatRequestBody(body interface{}) string {
	switch v := body.(type) {
	case string:
		// 尝试解析并格式化 JSON 字符串
		var raw interface{}
		if err := sonic.Unmarshal([]byte(v), &raw); err == nil {
			if formatted, err := sonic.MarshalIndent(raw, "", "  "); err == nil {
				return string(formatted)
			}
		}
		return v

	case []byte:
		var raw interface{}
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
