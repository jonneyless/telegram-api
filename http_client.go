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
)

// httpClient 封装 HTTP 客户端
type httpClient struct {
	client      *http.Client
	baseURL     string
	baseURLTemp string
	headers     map[string]string
	timeout     time.Duration
	retryConfig *retryConfig // 重试配置
	backoffFunc backoffFunc  // 退避算法函数
}

// retryConfig 重试配置
type retryConfig struct {
	maxRetries          int           // 最大重试次数
	retryInterval       time.Duration // 基础重试间隔
	maxRetryInterval    time.Duration // 最大重试间隔
	retryOnStatus       []int         // 针对哪些状态码重试（默认为5xx）
	retryOnTimeout      bool          // 是否在超时时重试
	retryOnNetworkError bool          // 是否在网络错误时重试
}

// backoffFunc 退避算法函数类型
type backoffFunc func(retryCount int, retryConfig *retryConfig) time.Duration

// httpError HTTP 错误
type httpError struct {
	statusCode int
	status     string
	body       string
	retries    int // 重试次数
	baseURL    string
	path       string
}

func (e *httpError) Error() string {
	return fmt.Sprintf("Path: %s%s, HTTP %d: %s, Body: %s (Retries: %d)", e.baseURL, e.path, e.statusCode, e.status, e.body, e.retries)
}

// newHTTPClient 创建新的 HTTP 客户端
func newHTTPClient(baseURL string, options ...option) *httpClient {
	client := &httpClient{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
		baseURL: baseURL,
		headers: make(map[string]string),
		retryConfig: &retryConfig{
			maxRetries:          3,                         // 默认重试3次
			retryInterval:       1 * time.Second,           // 默认1秒间隔
			maxRetryInterval:    30 * time.Second,          // 最大30秒间隔
			retryOnStatus:       []int{500, 502, 503, 504}, // 默认对5xx重试
			retryOnTimeout:      true,                      // 超时重试
			retryOnNetworkError: true,                      // 网络错误重试
		},
		backoffFunc: exponentialBackoff, // 默认使用指数退避
	}

	// 应用选项
	for _, option := range options {
		option(client)
	}

	return client
}

// ========================== 重试相关方法==========================

// setRetryConfig 设置重试配置
func (hc *httpClient) setRetryConfig(config *retryConfig) {
	if config != nil {
		hc.retryConfig = config
	}
}

// setBackoffFunc 设置退避算法函数
func (hc *httpClient) setBackoffFunc(fn backoffFunc) {
	if fn != nil {
		hc.backoffFunc = fn
	}
}

// defaultRetryConfig 获取默认重试配置
func defaultRetryConfig() *retryConfig {
	return &retryConfig{
		maxRetries:          3,
		retryInterval:       1 * time.Second,
		maxRetryInterval:    30 * time.Second,
		retryOnStatus:       []int{500, 502, 503, 504},
		retryOnTimeout:      true,
		retryOnNetworkError: true,
	}
}

// noRetry 禁用重试
func noRetry() *retryConfig {
	return &retryConfig{
		maxRetries: 0,
	}
}

// exponentialBackoff 指数退避算法
func exponentialBackoff(retryCount int, config *retryConfig) time.Duration {
	if retryCount <= 0 {
		return 0
	}

	// 计算退避时间：base * 2^(retry-1)
	backoff := float64(config.retryInterval) * math.Pow(2, float64(retryCount-1))

	// 添加抖动（10%的随机波动）
	backoff = backoff * (0.9 + 0.2*(float64(time.Now().UnixNano()%100)/100))

	// 限制最大退避时间
	if backoff > float64(config.maxRetryInterval) {
		backoff = float64(config.maxRetryInterval)
	}

	return time.Duration(backoff)
}

// linearBackoff 线性退避算法
func linearBackoff(retryCount int, config *retryConfig) time.Duration {
	if retryCount <= 0 {
		return 0
	}

	backoff := float64(config.retryInterval) * float64(retryCount)

	// 添加抖动
	backoff = backoff * (0.9 + 0.2*(float64(time.Now().UnixNano()%100)/100))

	if backoff > float64(config.maxRetryInterval) {
		backoff = float64(config.maxRetryInterval)
	}

	return time.Duration(backoff)
}

// fixedBackoff 固定间隔退避算法
func fixedBackoff(retryCount int, config *retryConfig) time.Duration {
	if retryCount <= 0 {
		return 0
	}

	// 添加抖动
	backoff := float64(config.retryInterval) * (0.9 + 0.2*(float64(time.Now().UnixNano()%100)/100))

	return time.Duration(backoff)
}

// shouldRetry 判断是否应该重试
func (hc *httpClient) shouldRetry(err error, statusCode int) bool {
	if hc.retryConfig == nil || hc.retryConfig.maxRetries <= 0 {
		return false
	}

	// 检查网络错误
	if err != nil && hc.retryConfig.retryOnNetworkError {
		// 判断是否为网络超时错误
		if hc.retryConfig.retryOnTimeout {
			if strings.Contains(err.Error(), "timeout") ||
				strings.Contains(err.Error(), "deadline") ||
				strings.Contains(err.Error(), "context deadline") {
				return true
			}
		}
		if statusCode < 500 {
			return false
		}
		// 其他网络错误
		return true
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

// ========================== 带重试的请求方法==========================

// doRequestWithRetry 执行带重试的 HTTP 请求
func (hc *httpClient) doRequestWithRetry(method, path string, queryParams map[string]string, data interface{}, result interface{}) error {
	return hc.doRequestWithHeadersAndRetry(method, path, queryParams, data, nil, result, hc.retryConfig)
}

// doRequestWithHeadersAndRetry 执行带自定义头部和重试的 HTTP 请求
func (hc *httpClient) doRequestWithHeadersAndRetry(method, path string, queryParams map[string]string, data interface{}, customHeaders map[string]string, result interface{}, retryConfig *retryConfig) error {
	// 使用提供的重试配置或默认配置
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

	for retryCount = 0; retryCount <= config.maxRetries; retryCount++ {
		// 如果不是第一次重试，等待退避时间
		if retryCount > 0 {
			waitTime := hc.backoffFunc(retryCount, config)
			if waitTime > 0 {
				time.Sleep(waitTime)
			}
		}

		// 执行请求
		err := hc.doRequestSingle(method, path, queryParams, data, customHeaders, &lastResp)

		// 获取状态码
		if lastResp != nil {
			statusCode = lastResp.StatusCode
		}

		// 如果没有错误，处理响应
		if err == nil && statusCode >= 200 && statusCode < 500 {
			return hc.handleResponse(lastResp, result)
		}

		// 保存最后一次错误
		if err != nil {
			lastErr = err
		} else if lastResp != nil {
			bodyBytes, _ := io.ReadAll(lastResp.Body)
			lastErr = &httpError{
				statusCode: statusCode,
				status:     lastResp.Status,
				body:       string(bodyBytes),
				retries:    retryCount,
				baseURL:    hc.baseURL,
				path:       path,
			}
		}

		// 检查是否应该重试
		if retryCount < config.maxRetries && hc.shouldRetry(lastErr, statusCode) {
			// 准备下一次重试
			if lastResp != nil {
				err := lastResp.Body.Close()
				if err != nil {
					return err
				}
				lastResp = nil
			}
			continue
		}

		// 不需要或不能重试，返回错误
		break
	}

	// 清理资源
	if lastResp != nil {
		err := lastResp.Body.Close()
		if err != nil {
			return err
		}
	}

	if retryCount < 1 {
		return fmt.Errorf("请求失败: %w", lastErr)
	}
	return fmt.Errorf("请求失败 (重试 %d 次): %w", config.maxRetries, lastErr)
}

// doRequestSingle 执行单次请求（不包含重试逻辑）
func (hc *httpClient) doRequestSingle(method, path string, queryParams map[string]string, data interface{}, customHeaders map[string]string, resp **http.Response) error {
	// 构建 URL
	requestURL := hc.buildURL(path, queryParams)

	// 准备请求体
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
			// 默认为 JSON
			jsonData, err := json.Marshal(v)
			if err != nil {
				return fmt.Errorf("序列化 JSON 失败: %w", err)
			}
			body = bytes.NewBuffer(jsonData)
			contentType = "application/json"
		}
	}

	// 创建请求
	req, err := http.NewRequest(method, requestURL, body)
	if err != nil {
		return fmt.Errorf("创建请求失败: %w", err)
	}

	// 设置头部
	hc.setHeaders(req, contentType, customHeaders)

	// 发送请求
	*resp, err = hc.client.Do(req)
	if err != nil {
		return fmt.Errorf("发送请求失败: %w", err)
	}

	return nil
}

// ========================== 原有方法的包装==========================

// get 发送 GET 请求（带重试）
func (hc *httpClient) get(path string, queryParams map[string]string, result interface{}) error {
	return hc.doRequestWithRetry("GET", path, queryParams, nil, result)
}

// getWithRetry 发送带自定义重试配置的 GET 请求
func (hc *httpClient) getWithRetry(path string, queryParams map[string]string, result interface{}, retryConfig *retryConfig) error {
	return hc.doRequestWithHeadersAndRetry("GET", path, queryParams, nil, nil, result, retryConfig)
}

// getWithHeaders 发送带自定义头部的 GET 请求（带重试）
func (hc *httpClient) getWithHeaders(path string, queryParams map[string]string, headers map[string]string, result interface{}) error {
	return hc.doRequestWithHeadersAndRetry("GET", path, queryParams, nil, headers, result, hc.retryConfig)
}

// post 发送 POST 请求（带重试）
func (hc *httpClient) post(path string, data interface{}, result interface{}) error {
	return hc.doRequestWithRetry("POST", path, nil, data, result)
}

// postWithRetry 发送带自定义重试配置的 POST 请求
func (hc *httpClient) postWithRetry(path string, data interface{}, result interface{}, retryConfig *retryConfig) error {
	return hc.doRequestWithHeadersAndRetry("POST", path, nil, data, nil, result, retryConfig)
}

// postWithHeaders 发送带自定义头部的 POST 请求（带重试）
func (hc *httpClient) postWithHeaders(path string, data interface{}, headers map[string]string, result interface{}) error {
	return hc.doRequestWithHeadersAndRetry("POST", path, nil, data, headers, result, hc.retryConfig)
}

// postForm 发送表单 POST 请求（带重试）
func (hc *httpClient) postForm(path string, formData map[string]string, result interface{}) error {
	values := url.Values{}
	for key, value := range formData {
		values.Add(key, value)
	}

	return hc.doRequestWithRetry("POST", path, nil, values.Encode(), result)
}

// ========================== 保持原有方法不变==========================

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

// setHeader 设置默认头部
func (hc *httpClient) setHeader(key, value string) {
	hc.headers[key] = value
}

// removeHeader 移除头部
func (hc *httpClient) removeHeader(key string) {
	delete(hc.headers, key)
}

// clearHeaders 清空所有头部
func (hc *httpClient) clearHeaders() {
	hc.headers = make(map[string]string)
}

// option 客户端配置选项
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
		hc.retryConfig = config
	}
}

// withBackoffFunc 设置退避算法函数
func withBackoffFunc(fn backoffFunc) option {
	return func(hc *httpClient) {
		hc.backoffFunc = fn
	}
}

// withHeaders 设置默认头部
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

// getRaw 发送 GET 请求并返回原始响应（带重试）
func (hc *httpClient) getRaw(path string, queryParams map[string]string) ([]byte, error) {
	var result []byte
	err := hc.doRequestWithRetry("GET", path, queryParams, nil, &result)
	return result, err
}

// postRaw 发送 POST 请求并返回原始响应（带重试）
func (hc *httpClient) postRaw(path string, data interface{}) ([]byte, error) {
	var result []byte
	err := hc.doRequestWithRetry("POST", path, nil, data, &result)
	return result, err
}

// doRequest 执行 HTTP 请求（兼容原有方法，带重试）
func (hc *httpClient) doRequest(method, path string, queryParams map[string]string, data interface{}, result interface{}) error {
	return hc.doRequestWithRetry(method, path, queryParams, data, result)
}

// doRequestWithHeaders 执行带自定义头部的 HTTP 请求（兼容原有方法，带重试）
func (hc *httpClient) doRequestWithHeaders(method, path string, queryParams map[string]string, data interface{}, customHeaders map[string]string, result interface{}) error {
	return hc.doRequestWithHeadersAndRetry(method, path, queryParams, data, customHeaders, result, hc.retryConfig)
}

// buildURL 构建完整的 URL
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

// setHeaders 设置请求头部
func (hc *httpClient) setHeaders(req *http.Request, contentType string, customHeaders map[string]string) {
	// 设置默认头部
	for key, value := range hc.headers {
		req.Header.Set(key, value)
	}

	// 设置内容类型
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}

	// 设置自定义头部
	for key, value := range customHeaders {
		req.Header.Set(key, value)
	}

	// 如果没有设置 User-Agent，设置一个默认的
	if req.Header.Get("User-Agent") == "" {
		req.Header.Set("User-Agent", "Go-HTTP-Client/1.0")
	}
}

// handleResponse 处理响应
func (hc *httpClient) handleResponse(resp *http.Response, result interface{}) error {
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("读取响应失败: %w", err)
	}

	// 如果不需要解析结果，直接返回
	if result == nil {
		return nil
	}

	// 如果结果是字节切片，直接返回
	if byteResult, ok := result.(*[]byte); ok {
		*byteResult = bodyBytes
		return nil
	}

	// 尝试解析为 JSON
	if len(bodyBytes) > 0 {
		if err := json.Unmarshal(bodyBytes, result); err != nil {
			return fmt.Errorf("解析 JSON 响应失败: %w, 响应内容: %s", err, string(bodyBytes))
		}
	}

	return nil
}

// uploadFile 文件上传（带重试）
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

	for retryCount := 0; retryCount <= config.maxRetries; retryCount++ {
		// 如果不是第一次重试，等待退避时间
		if retryCount > 0 {
			waitTime := hc.backoffFunc(retryCount, config)
			if waitTime > 0 {
				time.Sleep(waitTime)
			}
		}

		// 执行单次上传
		err := hc.uploadFileSingle(path, fileField, filename, fileContent, formData, &lastResp)

		if err == nil {
			return hc.handleResponse(lastResp, result)
		}

		lastErr = err

		// 检查是否应该重试
		if retryCount < config.maxRetries && hc.shouldRetry(lastErr, 0) {
			if lastResp != nil {
				err := lastResp.Body.Close()
				if err != nil {
					return err
				}
				lastResp = nil
			}
			continue
		}

		break
	}

	if lastResp != nil {
		err := lastResp.Body.Close()
		if err != nil {
			return err
		}
	}

	return fmt.Errorf("文件上传失败 (重试 %d 次): %w", config.maxRetries, lastErr)
}

// uploadFileSingle 单次文件上传
func (hc *httpClient) uploadFileSingle(path string, fileField string, filename string, fileContent []byte, formData map[string]string, resp **http.Response) error {
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	// 添加文件
	part, err := writer.CreateFormFile(fileField, filename)
	if err != nil {
		return fmt.Errorf("创建表单文件失败: %w", err)
	}
	if _, err := part.Write(fileContent); err != nil {
		return fmt.Errorf("写入文件内容失败: %w", err)
	}

	// 添加其他表单字段
	for key, value := range formData {
		if err := writer.WriteField(key, value); err != nil {
			return fmt.Errorf("写入表单字段失败: %w", err)
		}
	}

	if err := writer.Close(); err != nil {
		return fmt.Errorf("关闭 multipart writer 失败: %w", err)
	}

	// 发送请求
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

// downloadFile 文件下载（带重试）
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

	for retryCount := 0; retryCount <= config.maxRetries; retryCount++ {
		if retryCount > 0 {
			waitTime := hc.backoffFunc(retryCount, config)
			if waitTime > 0 {
				time.Sleep(waitTime)
			}
		}

		err := hc.downloadFileSingle(path, queryParams, &lastResp)

		if err == nil && lastResp != nil && lastResp.StatusCode == http.StatusOK {
			content, err := io.ReadAll(lastResp.Body)
			if err != nil {
				err := lastResp.Body.Close()
				if err != nil {
					return nil, "", err
				}
				return nil, "", fmt.Errorf("读取响应失败: %w", err)
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

			err = lastResp.Body.Close()
			if err != nil {
				return nil, "", err
			}
			return content, filename, nil
		}

		lastErr = err

		if retryCount < config.maxRetries && hc.shouldRetry(lastErr, 0) {
			if lastResp != nil {
				err := lastResp.Body.Close()
				if err != nil {
					return nil, "", err
				}
				lastResp = nil
			}
			continue
		}

		break
	}

	if lastResp != nil {
		err := lastResp.Body.Close()
		if err != nil {
			return nil, "", err
		}
	}

	return nil, "", fmt.Errorf("文件下载失败 (重试 %d 次): %w", config.maxRetries, lastErr)
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
