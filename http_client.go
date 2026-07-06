package telegram_api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/bytedance/sonic"
)

// backoffFunc 退避算法函数类型，用于计算重试等待时间
type backoffFunc func(retryCount int, retryConfig *retryConfig, retryAfter time.Duration) time.Duration

// TelegramAPIError Telegram API 错误响应结构体
type TelegramAPIError struct {
	OK          bool   `json:"ok"`          // 请求是否成功
	ErrorCode   int    `json:"error_code"`  // 错误码
	Description string `json:"description"` // 错误描述
	Parameters  *struct {
		RetryAfter      int   `json:"retry_after"`        // 需要等待的秒数
		MigrateToChatID int64 `json:"migrate_to_chat_id"` // 迁移到的群组ID
	} `json:"parameters"` // 额外参数
}

// Error 实现 error 接口
func (e *TelegramAPIError) Error() string {
	if e.Parameters != nil && e.Parameters.RetryAfter > 0 {
		return fmt.Sprintf("Telegram API error: %s (retry_after: %d seconds)", e.Description, e.Parameters.RetryAfter)
	}
	if e.Parameters != nil && e.Parameters.MigrateToChatID > 0 {
		return fmt.Sprintf("Telegram API error: %s (migrate_to_chat_id: %d)", e.Description, e.Parameters.MigrateToChatID)
	}
	return fmt.Sprintf("Telegram API error: %s", e.Description)
}

// NeedRetry 判断是否需要重试
func (e *TelegramAPIError) NeedRetry() bool {
	if e.Parameters == nil {
		return false
	}
	return e.Parameters.RetryAfter > 0 || e.Parameters.MigrateToChatID > 0
}

// GetRetryAfter 获取重试等待时间
func (e *TelegramAPIError) GetRetryAfter() time.Duration {
	if e.Parameters != nil && e.Parameters.RetryAfter > 0 {
		return time.Duration(e.Parameters.RetryAfter) * time.Second
	}
	return 0
}

// IsMigrateError 判断是否是群组迁移错误
func (e *TelegramAPIError) IsMigrateError() bool {
	return e.Parameters != nil && e.Parameters.MigrateToChatID > 0
}

// GetMigrateToChatID 获取迁移目标群组ID
func (e *TelegramAPIError) GetMigrateToChatID() int64 {
	if e.Parameters != nil {
		return e.Parameters.MigrateToChatID
	}
	return 0
}

// httpError HTTP 请求错误结构体
type httpError struct {
	statusCode  int           // HTTP状态码
	status      string        // HTTP状态描述
	body        string        // 响应体内容
	retries     int           // 已重试次数
	baseURL     string        // 基础URL
	path        string        // 请求路径
	retryAfter  time.Duration // 需要等待的时间
	migrateToID int64         // 迁移目标群组ID
}

// Error 实现 error 接口
func (e *httpError) Error() string {
	var parts []string
	parts = append(parts, fmt.Sprintf("Path: %s%s", e.baseURL, e.path))
	if e.statusCode > 0 {
		parts = append(parts, fmt.Sprintf("HTTP %d: %s", e.statusCode, e.status))
	}
	if e.retryAfter > 0 {
		parts = append(parts, fmt.Sprintf("RetryAfter: %v", e.retryAfter))
	}
	if e.migrateToID > 0 {
		parts = append(parts, fmt.Sprintf("MigrateToID: %d", e.migrateToID))
	}
	if e.body != "" {
		parts = append(parts, fmt.Sprintf("Body: %s", e.body))
	}
	parts = append(parts, fmt.Sprintf("Retries: %d", e.retries))
	return strings.Join(parts, ", ")
}

// retryConfig 重试配置结构体
type retryConfig struct {
	maxRetries          int           // 最大重试次数
	retryInterval       time.Duration // 基础重试间隔
	maxRetryInterval    time.Duration // 最大重试间隔
	retryOnStatus       []int         // 需要重试的HTTP状态码列表
	retryOnTimeout      bool          // 是否在超时时重试
	retryOnNetworkError bool          // 是否在网络错误时重试
	respectRetryAfter   bool          // 是否遵守Telegram的retry_after
}

// httpClient HTTP客户端结构体
type httpClient struct {
	client      *http.Client      // HTTP客户端
	baseURL     string            // 基础URL
	baseURLTemp string            // 临时基础URL
	headers     map[string]string // 默认请求头
	timeout     time.Duration     // 超时时间
	retryConfig *retryConfig      // 重试配置
	backoffFunc backoffFunc       // 退避算法函数
}

// exponentialBackoff 指数退避算法
func exponentialBackoff(retryCount int, config *retryConfig, retryAfter time.Duration) time.Duration {
	if retryCount <= 0 {
		return 0
	}
	// 如果有 retry_after，优先使用
	if retryAfter > 0 {
		extraJitter := time.Duration(100+time.Now().UnixNano()%400) * time.Millisecond
		return retryAfter + extraJitter
	}
	// 计算退避时间：base * 2^(retry-1)
	backoff := float64(config.retryInterval) * math.Pow(2, float64(retryCount-1))
	// 添加抖动（10%的随机波动）
	jitter := 0.9 + 0.2*(float64(time.Now().UnixNano()%100)/100)
	backoff = backoff * jitter
	// 限制最大退避时间
	if backoff > float64(config.maxRetryInterval) {
		backoff = float64(config.maxRetryInterval)
	}
	return time.Duration(backoff)
}

// linearBackoff 线性退避算法
func linearBackoff(retryCount int, config *retryConfig, retryAfter time.Duration) time.Duration {
	if retryCount <= 0 {
		return 0
	}
	if retryAfter > 0 {
		extraJitter := time.Duration(100+time.Now().UnixNano()%400) * time.Millisecond
		return retryAfter + extraJitter
	}
	backoff := float64(config.retryInterval) * float64(retryCount)
	jitter := 0.9 + 0.2*(float64(time.Now().UnixNano()%100)/100)
	backoff = backoff * jitter
	if backoff > float64(config.maxRetryInterval) {
		backoff = float64(config.maxRetryInterval)
	}
	return time.Duration(backoff)
}

// fixedBackoff 固定间隔退避算法
func fixedBackoff(retryCount int, config *retryConfig, retryAfter time.Duration) time.Duration {
	if retryCount <= 0 {
		return 0
	}
	if retryAfter > 0 {
		extraJitter := time.Duration(100+time.Now().UnixNano()%400) * time.Millisecond
		return retryAfter + extraJitter
	}
	jitter := 0.9 + 0.2*(float64(time.Now().UnixNano()%100)/100)
	backoff := float64(config.retryInterval) * jitter
	return time.Duration(backoff)
}

// newHTTPClient 创建新的HTTP客户端
func newHTTPClient(baseURL string, options ...option) *httpClient {
	client := &httpClient{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
		baseURL: baseURL,
		headers: make(map[string]string),
		retryConfig: &retryConfig{
			maxRetries:          3,
			retryInterval:       1 * time.Second,
			maxRetryInterval:    60 * time.Second,
			retryOnStatus:       []int{429, 500, 502, 503, 504},
			retryOnTimeout:      true,
			retryOnNetworkError: true,
			respectRetryAfter:   true,
		},
		backoffFunc: exponentialBackoff,
	}
	for _, option := range options {
		option(client)
	}
	return client
}

// defaultRetryConfig 获取默认重试配置
func defaultRetryConfig() *retryConfig {
	return &retryConfig{
		maxRetries:          3,
		retryInterval:       1 * time.Second,
		maxRetryInterval:    60 * time.Second,
		retryOnStatus:       []int{429, 500, 502, 503, 504},
		retryOnTimeout:      true,
		retryOnNetworkError: true,
		respectRetryAfter:   true,
	}
}

// noRetry 返回禁用重试的配置
func noRetry() *retryConfig {
	return &retryConfig{
		maxRetries: 0,
	}
}

// option 客户端配置选项函数类型
type option func(*httpClient)

// withTimeout 设置超时时间
func withTimeout(timeout time.Duration) option {
	return func(hc *httpClient) {
		hc.client.Timeout = timeout
		hc.timeout = timeout
	}
}

// withRetryConfig 设置重试配置
func withRetryConfig(config *retryConfig) option {
	return func(hc *httpClient) {
		if config != nil {
			hc.retryConfig = config
		}
	}
}

// withBackoffFunc 设置退避算法函数
func withBackoffFunc(fn backoffFunc) option {
	return func(hc *httpClient) {
		if fn != nil {
			hc.backoffFunc = fn
		}
	}
}

// withHeaders 设置默认请求头
func withHeaders(headers map[string]string) option {
	return func(hc *httpClient) {
		for k, v := range headers {
			hc.headers[k] = v
		}
	}
}

// withTransport 设置自定义传输层
func withTransport(transport http.RoundTripper) option {
	return func(hc *httpClient) {
		hc.client.Transport = transport
	}
}

// shouldRetry 判断是否应该重试
func (hc *httpClient) shouldRetry(err error, statusCode int, retryAfter time.Duration) bool {
	if hc.retryConfig == nil || hc.retryConfig.maxRetries <= 0 {
		return false
	}
	// 如果有 retry_after，应该重试
	if retryAfter > 0 {
		return true
	}
	// 检查错误类型
	if err != nil {
		errMsg := err.Error()
		if hc.retryConfig.retryOnTimeout {
			if strings.Contains(errMsg, "timeout") ||
				strings.Contains(errMsg, "deadline") ||
				strings.Contains(errMsg, "context deadline") ||
				strings.Contains(errMsg, "connection refused") ||
				strings.Contains(errMsg, "connection reset") {
				return true
			}
		}
		if hc.retryConfig.retryOnNetworkError {
			if strings.Contains(errMsg, "network") ||
				strings.Contains(errMsg, "tls") ||
				strings.Contains(errMsg, "handshake") {
				return true
			}
		}
	}
	// 检查状态码
	if statusCode > 0 {
		for _, retryCode := range hc.retryConfig.retryOnStatus {
			if statusCode == retryCode {
				return true
			}
		}
	}
	return false
}

// doRequestWithHeadersAndRetry 执行带重试的HTTP请求
func (hc *httpClient) doRequestWithHeadersAndRetry(method, path string, queryParams map[string]string, data interface{}, customHeaders map[string]string, result interface{}, retryConfig *retryConfig) error {
	config := retryConfig
	if config == nil {
		config = hc.retryConfig
	}
	if config == nil {
		config = defaultRetryConfig()
	}

	var lastErr error
	var lastResp *http.Response
	var statusCode int
	var retryCount int
	var retryAfter time.Duration
	var migrateToID int64

	for retryCount = 0; retryCount <= config.maxRetries; retryCount++ {
		// 如果不是第一次重试，等待退避时间
		if retryCount > 0 {
			waitTime := hc.backoffFunc(retryCount, config, retryAfter)
			if waitTime > 0 {
				time.Sleep(waitTime)
			}
			retryAfter = 0
		}

		// 执行请求
		err := hc.doRequestSingle(method, path, queryParams, data, customHeaders, &lastResp)

		// 获取状态码
		if lastResp != nil {
			statusCode = lastResp.StatusCode
		}

		// 如果请求成功，处理响应
		if err == nil && statusCode >= 200 && statusCode < 400 {
			return hc.handleResponse(lastResp, result)
		}

		// 解析 Telegram API 错误，获取 retry_after 和 migrate_to_chat_id
		if lastResp != nil && (statusCode == 429 || statusCode == 400) {
			bodyBytes, readErr := io.ReadAll(lastResp.Body)
			if readErr == nil {
				// 恢复 Body 以便后续处理
				lastResp.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

				var telegramErr TelegramAPIError
				if jsonErr := json.Unmarshal(bodyBytes, &telegramErr); jsonErr == nil {
					if telegramErr.Parameters != nil {
						if telegramErr.Parameters.RetryAfter > 0 {
							retryAfter = time.Duration(telegramErr.Parameters.RetryAfter) * time.Second
						}
						if telegramErr.Parameters.MigrateToChatID > 0 {
							migrateToID = telegramErr.Parameters.MigrateToChatID
						}
						lastErr = &telegramErr
					}
				}
			}
		}

		// 保存最后一次错误
		if err != nil {
			lastErr = err
		} else if lastResp != nil {
			bodyBytes, readErr := io.ReadAll(lastResp.Body)
			if readErr != nil {
				bodyBytes = []byte{}
			}
			lastErr = &httpError{
				statusCode:  statusCode,
				status:      lastResp.Status,
				body:        string(bodyBytes),
				retries:     retryCount,
				baseURL:     hc.baseURL,
				path:        path,
				retryAfter:  retryAfter,
				migrateToID: migrateToID,
			}
		}

		// 检查是否应该重试
		if retryCount < config.maxRetries && hc.shouldRetry(lastErr, statusCode, retryAfter) {
			if lastResp != nil {
				_ = lastResp.Body.Close()
				lastResp = nil
			}
			continue
		}

		break
	}

	// 清理资源
	if lastResp != nil {
		_ = lastResp.Body.Close()
	}

	if lastErr != nil {
		if retryCount > 0 {
			return fmt.Errorf("请求失败 (重试 %d 次): %w", retryCount, lastErr)
		}
		return fmt.Errorf("请求失败: %w", lastErr)
	}

	return nil
}

// doRequestSingle 执行单次HTTP请求
func (hc *httpClient) doRequestSingle(method, path string, queryParams map[string]string, data interface{}, customHeaders map[string]string, resp **http.Response) error {
	requestURL := hc.buildURL(path, queryParams)

	var body io.Reader
	var contentType string

	if data != nil {
		switch v := data.(type) {
		case string:
			body = bytes.NewBufferString(v)
		case []byte:
			body = bytes.NewBuffer(v)
		case url.Values:
			body = bytes.NewBufferString(v.Encode())
			contentType = "application/x-www-form-urlencoded"
		default:
			jsonData, err := sonic.Marshal(v)
			if err != nil {
				return fmt.Errorf("序列化 JSON 失败: %w", err)
			}
			body = bytes.NewBuffer(jsonData)
			contentType = "application/json"
		}
	}

	req, err := http.NewRequest(method, requestURL, body)
	if err != nil {
		return fmt.Errorf("创建请求失败: %w", err)
	}

	hc.setHeaders(req, contentType, customHeaders)

	*resp, err = hc.client.Do(req)
	if err != nil {
		return fmt.Errorf("发送请求失败: %w", err)
	}

	return nil
}

// get 发送 GET 请求
func (hc *httpClient) get(path string, queryParams map[string]string, result interface{}) error {
	return hc.doRequestWithHeadersAndRetry("GET", path, queryParams, nil, nil, result, hc.retryConfig)
}

// post 发送 POST 请求
func (hc *httpClient) post(path string, data interface{}, result interface{}) error {
	return hc.doRequestWithHeadersAndRetry("POST", path, nil, data, nil, result, hc.retryConfig)
}

// postWithHeaders 发送带自定义头部的 POST 请求
func (hc *httpClient) postWithHeaders(path string, data interface{}, headers map[string]string, result interface{}) error {
	return hc.doRequestWithHeadersAndRetry("POST", path, nil, data, headers, result, hc.retryConfig)
}

// postForm 发送表单 POST 请求
func (hc *httpClient) postForm(path string, formData map[string]string, result interface{}) error {
	values := url.Values{}
	for key, value := range formData {
		values.Add(key, value)
	}
	return hc.doRequestWithHeadersAndRetry("POST", path, nil, values.Encode(), nil, result, hc.retryConfig)
}

// getRaw 发送 GET 请求并返回原始字节
func (hc *httpClient) getRaw(path string, queryParams map[string]string) ([]byte, error) {
	var result []byte
	err := hc.doRequestWithHeadersAndRetry("GET", path, queryParams, nil, nil, &result, hc.retryConfig)
	return result, err
}

// postRaw 发送 POST 请求并返回原始字节
func (hc *httpClient) postRaw(path string, data interface{}) ([]byte, error) {
	var result []byte
	err := hc.doRequestWithHeadersAndRetry("POST", path, nil, data, nil, &result, hc.retryConfig)
	return result, err
}

// doRequest 执行 HTTP 请求（兼容方法）
func (hc *httpClient) doRequest(method, path string, queryParams map[string]string, data interface{}, result interface{}) error {
	return hc.doRequestWithHeadersAndRetry(method, path, queryParams, data, nil, result, hc.retryConfig)
}

// doRequestWithHeaders 执行带自定义头部的 HTTP 请求（兼容方法）
func (hc *httpClient) doRequestWithHeaders(method, path string, queryParams map[string]string, data interface{}, customHeaders map[string]string, result interface{}) error {
	return hc.doRequestWithHeadersAndRetry(method, path, queryParams, data, customHeaders, result, hc.retryConfig)
}

// getBaseUrl 获取基础URL
func (hc *httpClient) getBaseUrl() string {
	return hc.baseURL
}

// setBaseUrl 设置基础URL
func (hc *httpClient) setBaseUrl(url string) {
	hc.baseURL = url
}

// setBaseUrlTemp 设置临时基础URL
func (hc *httpClient) setBaseUrlTemp(url string) {
	hc.baseURLTemp = url
}

// setHeader 设置默认请求头
func (hc *httpClient) setHeader(key, value string) {
	hc.headers[key] = value
}

// removeHeader 移除请求头
func (hc *httpClient) removeHeader(key string) {
	delete(hc.headers, key)
}

// clearHeaders 清空所有请求头
func (hc *httpClient) clearHeaders() {
	hc.headers = make(map[string]string)
}

// buildURL 构建完整的URL
func (hc *httpClient) buildURL(path string, queryParams map[string]string) string {
	fullURL := hc.baseURL + path

	if hc.baseURLTemp != "" {
		fullURL = hc.baseURLTemp + path
		hc.baseURLTemp = ""
	}

	if len(queryParams) > 0 {
		params := url.Values{}
		for key, value := range queryParams {
			params.Add(key, value)
		}
		fullURL += "?" + params.Encode()
	}

	return fullURL
}

// setHeaders 设置请求头
func (hc *httpClient) setHeaders(req *http.Request, contentType string, customHeaders map[string]string) {
	for key, value := range hc.headers {
		req.Header.Set(key, value)
	}

	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}

	for key, value := range customHeaders {
		req.Header.Set(key, value)
	}

	if req.Header.Get("User-Agent") == "" {
		req.Header.Set("User-Agent", "Go-HTTP-Client/1.0")
	}
}

// handleResponse 处理HTTP响应
func (hc *httpClient) handleResponse(resp *http.Response, result interface{}) error {
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("读取响应失败: %w", err)
	}

	if result == nil {
		return nil
	}

	if byteResult, ok := result.(*[]byte); ok {
		*byteResult = bodyBytes
		return nil
	}

	if len(bodyBytes) > 0 {
		if err := sonic.Unmarshal(bodyBytes, result); err != nil {
			return fmt.Errorf("解析 JSON 响应失败: %w, 响应内容: %s", err, string(bodyBytes))
		}
	}

	return nil
}

// uploadFile 上传文件
func (hc *httpClient) uploadFile(path string, fileField string, filename string, fileContent []byte, formData map[string]string, result interface{}) error {
	return hc.uploadFileWithRetry(path, fileField, filename, fileContent, formData, result, hc.retryConfig)
}

// uploadFileWithRetry 带重试的文件上传
func (hc *httpClient) uploadFileWithRetry(path string, fileField string, filename string, fileContent []byte, formData map[string]string, result interface{}, retryConfig *retryConfig) error {
	config := retryConfig
	if config == nil {
		config = hc.retryConfig
	}

	var lastErr error
	var lastResp *http.Response
	var retryAfter time.Duration

	for retryCount := 0; retryCount <= config.maxRetries; retryCount++ {
		if retryCount > 0 {
			waitTime := hc.backoffFunc(retryCount, config, retryAfter)
			if waitTime > 0 {
				time.Sleep(waitTime)
			}
			retryAfter = 0
		}

		err := hc.uploadFileSingle(path, fileField, filename, fileContent, formData, &lastResp)

		if err == nil && lastResp != nil && lastResp.StatusCode >= 200 && lastResp.StatusCode < 400 {
			return hc.handleResponse(lastResp, result)
		}

		if lastResp != nil && lastResp.StatusCode == 429 {
			bodyBytes, readErr := io.ReadAll(lastResp.Body)
			if readErr == nil {
				var telegramErr TelegramAPIError
				if jsonErr := json.Unmarshal(bodyBytes, &telegramErr); jsonErr == nil {
					if telegramErr.Parameters != nil && telegramErr.Parameters.RetryAfter > 0 {
						retryAfter = time.Duration(telegramErr.Parameters.RetryAfter) * time.Second
						lastErr = &telegramErr
					}
				}
				lastResp.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
			}
		}

		if err != nil {
			lastErr = err
		} else if lastResp != nil {
			bodyBytes, readErr := io.ReadAll(lastResp.Body)
			if readErr != nil {
				bodyBytes = []byte{}
			}
			lastErr = &httpError{
				statusCode: lastResp.StatusCode,
				status:     lastResp.Status,
				body:       string(bodyBytes),
				retries:    retryCount,
				baseURL:    hc.baseURL,
				path:       path,
				retryAfter: retryAfter,
			}
		}

		if retryCount < config.maxRetries && hc.shouldRetry(lastErr, 0, retryAfter) {
			if lastResp != nil {
				_ = lastResp.Body.Close()
				lastResp = nil
			}
			continue
		}

		break
	}

	if lastResp != nil {
		_ = lastResp.Body.Close()
	}

	if lastErr != nil {
		return fmt.Errorf("文件上传失败 (重试 %d 次): %w", config.maxRetries, lastErr)
	}

	return nil
}

// uploadFileSingle 单次文件上传
func (hc *httpClient) uploadFileSingle(path string, fileField string, filename string, fileContent []byte, formData map[string]string, resp **http.Response) error {
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	part, err := writer.CreateFormFile(fileField, filename)
	if err != nil {
		return fmt.Errorf("创建表单文件失败: %w", err)
	}

	if _, err := part.Write(fileContent); err != nil {
		return fmt.Errorf("写入文件内容失败: %w", err)
	}

	for key, value := range formData {
		if err := writer.WriteField(key, value); err != nil {
			return fmt.Errorf("写入表单字段失败: %w", err)
		}
	}

	if err := writer.Close(); err != nil {
		return fmt.Errorf("关闭 multipart writer 失败: %w", err)
	}

	requestURL := hc.baseURL + path
	req, err := http.NewRequest("POST", requestURL, &buf)
	if err != nil {
		return fmt.Errorf("创建请求失败: %w", err)
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	hc.setHeaders(req, "", nil)

	*resp, err = hc.client.Do(req)
	if err != nil {
		return fmt.Errorf("发送请求失败: %w", err)
	}

	return nil
}

// downloadFile 下载文件
func (hc *httpClient) downloadFile(path string, queryParams map[string]string) ([]byte, string, error) {
	return hc.downloadFileWithRetry(path, queryParams, hc.retryConfig)
}

// downloadFileWithRetry 带重试的文件下载
func (hc *httpClient) downloadFileWithRetry(path string, queryParams map[string]string, retryConfig *retryConfig) ([]byte, string, error) {
	config := retryConfig
	if config == nil {
		config = hc.retryConfig
	}

	var lastErr error
	var lastResp *http.Response
	var retryAfter time.Duration

	for retryCount := 0; retryCount <= config.maxRetries; retryCount++ {
		if retryCount > 0 {
			waitTime := hc.backoffFunc(retryCount, config, retryAfter)
			if waitTime > 0 {
				time.Sleep(waitTime)
			}
			retryAfter = 0
		}

		err := hc.downloadFileSingle(path, queryParams, &lastResp)

		if err == nil && lastResp != nil && lastResp.StatusCode == http.StatusOK {
			content, readErr := io.ReadAll(lastResp.Body)
			if readErr != nil {
				_ = lastResp.Body.Close()
				return nil, "", fmt.Errorf("读取响应失败: %w", readErr)
			}

			filename := ""
			contentDisposition := lastResp.Header.Get("Content-Disposition")
			if contentDisposition != "" {
				if strings.Contains(contentDisposition, "filename=") {
					parts := strings.Split(contentDisposition, "filename=")
					if len(parts) > 1 {
						filename = strings.Trim(parts[1], `"`)
					}
				}
			}

			_ = lastResp.Body.Close()
			return content, filename, nil
		}

		if lastResp != nil && lastResp.StatusCode == 429 {
			bodyBytes, readErr := io.ReadAll(lastResp.Body)
			if readErr == nil {
				var telegramErr TelegramAPIError
				if jsonErr := json.Unmarshal(bodyBytes, &telegramErr); jsonErr == nil {
					if telegramErr.Parameters != nil && telegramErr.Parameters.RetryAfter > 0 {
						retryAfter = time.Duration(telegramErr.Parameters.RetryAfter) * time.Second
						lastErr = &telegramErr
					}
				}
				lastResp.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
			}
		}

		if err != nil {
			lastErr = err
		} else if lastResp != nil {
			bodyBytes, readErr := io.ReadAll(lastResp.Body)
			if readErr != nil {
				bodyBytes = []byte{}
			}
			lastErr = &httpError{
				statusCode: lastResp.StatusCode,
				status:     lastResp.Status,
				body:       string(bodyBytes),
				retries:    retryCount,
				baseURL:    hc.baseURL,
				path:       path,
				retryAfter: retryAfter,
			}
		}

		if retryCount < config.maxRetries && hc.shouldRetry(lastErr, 0, retryAfter) {
			if lastResp != nil {
				_ = lastResp.Body.Close()
				lastResp = nil
			}
			continue
		}

		break
	}

	if lastResp != nil {
		_ = lastResp.Body.Close()
	}

	if lastErr != nil {
		return nil, "", fmt.Errorf("文件下载失败 (重试 %d 次): %w", config.maxRetries, lastErr)
	}

	return nil, "", nil
}

// downloadFileSingle 单次文件下载
func (hc *httpClient) downloadFileSingle(path string, queryParams map[string]string, resp **http.Response) error {
	requestURL := hc.buildURL(path, queryParams)

	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		return fmt.Errorf("创建请求失败: %w", err)
	}

	hc.setHeaders(req, "", nil)

	*resp, err = hc.client.Do(req)
	if err != nil {
		return fmt.Errorf("发送请求失败: %w", err)
	}

	return nil
}
