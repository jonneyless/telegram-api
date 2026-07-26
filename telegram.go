package telegram_api

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/jonneyless/telegram-api/models"
	"github.com/jonneyless/telegram-api/requests"
)

// 多实例模式存储
var (
	instances = make(map[int64]*Telegram)
	instance  *Telegram
	mu        sync.RWMutex
	once      sync.Once
)

// Telegram Telegram API 客户端
type Telegram struct {
	client *resty.Client // HTTP客户端
	debug  bool          // 是否开启调试模式
	botId  int64         // Bot ID
	token  string        // Bot Token
}

func (t *Telegram) GetBotId() int64 {
	return t.botId
}

func (t *Telegram) GetToken() string {
	return t.token
}

// MessageDeleteOptions 消息删除选项
type MessageDeleteOptions struct {
	Delay       time.Duration // 延迟删除时间
	DeleteReply bool          // 是否删除回复消息
	MessageIds  []int64       // 要删除的消息ID列表
}

func (t *Telegram) post(path string, params any, result any) error {
	if t.token == "" {
		return fmt.Errorf("bot token empty")
	}

	var errResponse models.ApiErrorResponse

	// 使用 Resty 发送 POST 请求
	resp, err := t.client.R().
		SetBody(params).
		SetResult(result).
		SetError(&errResponse).
		Post(path)

	if err != nil {
		return err
	}

	// 检查 HTTP 状态码
	if resp.StatusCode() >= 400 {
		if !errResponse.Ok {
			return fmt.Errorf("telegram error: %s", errResponse.Description)
		}
		return fmt.Errorf("HTTP error: %d", resp.StatusCode())
	}

	return nil
}

func (t *Telegram) postMultipart(path string, body *bytes.Buffer, contentType string, result any) error {
	if t.token == "" {
		return fmt.Errorf("bot token empty")
	}

	var errResponse models.ApiErrorResponse

	resp, err := t.client.R().
		SetHeader("Content-Type", contentType).
		SetBody(body.Bytes()).
		SetResult(result).
		SetError(&errResponse).
		Post(path)

	if err != nil {
		return err
	}

	// 检查 HTTP 状态码
	if resp.StatusCode() >= 400 {
		if !errResponse.Ok {
			return fmt.Errorf("telegram error: %s", errResponse.Description)
		}
		return fmt.Errorf("HTTP error: %d", resp.StatusCode())
	}

	return nil
}

// SendMessage 发送消息
func (t *Telegram) SendMessage(params *requests.SendMessage, deleteOptions ...*MessageDeleteOptions) (*models.MessageResponse, error) {
	var apiResponse *models.MessageResponse

	apiPath := "sendMessage"
	if params.Photo != "" {
		apiPath = "sendPhoto"
	} else if params.Video != "" {
		apiPath = "sendVideo"
	} else if params.Audio != "" {
		apiPath = "sendAudio"
	} else if params.Document != "" {
		apiPath = "sendDocument"
	}

	err := t.post(apiPath, params.GetParams(), &apiResponse)
	if err != nil {
		return nil, err
	}

	if len(deleteOptions) > 0 && deleteOptions[0] != nil {
		options := deleteOptions[0]
		messageIds := make([]int64, len(options.MessageIds))
		copy(messageIds, options.MessageIds)
		if options.DeleteReply && params.ReplyParameters != nil {
			messageIds = append(messageIds, params.ReplyParameters.MessageId)
		}
		if apiResponse != nil {
			chatId := apiResponse.Result.Chat.ID
			messageIds = append(messageIds, apiResponse.Result.MessageID)
			go func() {
				time.Sleep(options.Delay)
				_, _ = t.DeleteMessage(&requests.Message{
					ChatId:     chatId,
					MessageIds: &messageIds,
				})
			}()
		}
	}

	return apiResponse, nil
}

// EditMessageText 编辑消息文本
func (t *Telegram) EditMessageText(params *requests.EditMessage) (*models.MessageResponse, error) {
	var apiResponse *models.MessageResponse
	err := t.post("editMessageText", params.GetParams(), &apiResponse)
	if err != nil {
		return nil, err
	}
	return apiResponse, err
}

// DeleteMessage 删除消息
func (t *Telegram) DeleteMessage(params *requests.Message) (*models.ApiResponse, error) {
	var apiResponse *models.ApiResponse
	var err error

	if params.MessageIds != nil {
		err = t.post("deleteMessages", params, &apiResponse)
	} else {
		err = t.post("deleteMessage", params, &apiResponse)
	}

	if err != nil {
		return nil, err
	}

	return apiResponse, err
}

// SendPhoto 图片二进制发布
func (t *Telegram) SendPhoto(params *requests.SendPhoto) (*models.MessageResponse, error) {
	var apiResponse *models.MessageResponse
	body, contentType, err := params.ToMultipart()
	if err != nil {
		return nil, err
	}
	err = t.postMultipart("sendPhoto", body, contentType, &apiResponse)
	if err != nil {
		return nil, err
	}
	return apiResponse, err
}

// SendPoll 发送投票
func (t *Telegram) SendPoll(params *requests.SendPoll) (*models.MessageResponse, error) {
	var apiResponse *models.MessageResponse
	err := t.post("sendPoll", params, &apiResponse)
	if err != nil {
		return nil, err
	}
	return apiResponse, err
}

// SendDice 发送骰子
func (t *Telegram) SendDice(params *requests.SendDice) (*models.MessageResponse, error) {
	var apiResponse *models.MessageResponse
	err := t.post("sendDice", params, &apiResponse)
	if err != nil {
		return nil, err
	}
	return apiResponse, err
}

// AnswerCallbackQuery 回答回调查询
func (t *Telegram) AnswerCallbackQuery(params *requests.AnswerCallbackQuery) (*models.ApiResponse, error) {
	var apiResponse *models.ApiResponse
	err := t.post("answerCallbackQuery", params.GetParams(), &apiResponse)
	if err != nil {
		return nil, err
	}
	return apiResponse, err
}

// GetChat 获取群组信息
func (t *Telegram) GetChat(params *requests.Chat) (*models.ChatResponse, error) {
	var apiResponse *models.ChatResponse
	err := t.post("getChat", params, &apiResponse)
	if err != nil {
		return nil, err
	}
	return apiResponse, err
}

// LeaveChat 离开群组
func (t *Telegram) LeaveChat(params *requests.Chat) (*models.ApiResponse, error) {
	var apiResponse *models.ApiResponse
	err := t.post("leaveChat", params, &apiResponse)
	if err != nil {
		return nil, err
	}
	return apiResponse, err
}

// SetChatPhoto 设置群组头像
func (t *Telegram) SetChatPhoto(params *requests.ChatPhoto) (*models.ApiResponse, error) {
	var apiResponse *models.ApiResponse
	body, contentType, err := params.ToMultipart()
	if err != nil {
		return nil, err
	}
	err = t.postMultipart("setChatPhoto", body, contentType, &apiResponse)
	if err != nil {
		return nil, err
	}
	return apiResponse, err
}

// DeleteChatPhoto 删除群组头像
func (t *Telegram) DeleteChatPhoto(params *requests.Chat) (*models.ApiResponse, error) {
	var apiResponse *models.ApiResponse
	err := t.post("deleteChatPhoto", params, &apiResponse)
	if err != nil {
		return nil, err
	}
	return apiResponse, err
}

// SetChatTitle 设置群组标题
func (t *Telegram) SetChatTitle(params *requests.ChatTitle) (*models.ApiResponse, error) {
	var apiResponse *models.ApiResponse
	err := t.post("setChatTitle", params, &apiResponse)
	if err != nil {
		return nil, err
	}
	return apiResponse, err
}

// SetChatDescription 设置群组描述
func (t *Telegram) SetChatDescription(params *requests.ChatDescription) (*models.ApiResponse, error) {
	var apiResponse *models.ApiResponse
	err := t.post("setChatDescription", params, &apiResponse)
	if err != nil {
		return nil, err
	}
	return apiResponse, err
}

// PinChatMessage 置顶消息
func (t *Telegram) PinChatMessage(params *requests.Message) (*models.ApiResponse, error) {
	var apiResponse *models.ApiResponse
	err := t.post("pinChatMessage", params, &apiResponse)
	if err != nil {
		return nil, err
	}
	return apiResponse, err
}

// UnPinChatMessage 取消置顶消息
func (t *Telegram) UnPinChatMessage(params *requests.Message) (*models.ApiResponse, error) {
	var apiResponse *models.ApiResponse
	err := t.post("unpinChatMessage", params, &apiResponse)
	if err != nil {
		return nil, err
	}
	return apiResponse, err
}

// UnPinAllChatMessage 取消所有置顶消息
func (t *Telegram) UnPinAllChatMessage(params *requests.Chat) (*models.ApiResponse, error) {
	var apiResponse *models.ApiResponse
	err := t.post("unpinAllChatMessage", params, &apiResponse)
	if err != nil {
		return nil, err
	}
	return apiResponse, err
}

// GetChatMember 获取群组成员信息
func (t *Telegram) GetChatMember(params *requests.Member) (*models.ChatMemberResponse, error) {
	var apiResponse *models.ChatMemberResponse
	err := t.post("getChatMember", params, &apiResponse)
	if err != nil {
		return nil, err
	}
	return apiResponse, err
}

// RestrictChatMember 限制群组成员
func (t *Telegram) RestrictChatMember(params *requests.Restrict) (*models.ApiResponse, error) {
	var apiResponse *models.ApiResponse
	err := t.post("restrictChatMember", params, &apiResponse)
	if err != nil {
		return nil, err
	}
	return apiResponse, err
}

// BanChatMember 封禁群组成员
func (t *Telegram) BanChatMember(params *requests.Ban) (*models.ApiResponse, error) {
	var apiResponse *models.ApiResponse
	err := t.post("banChatMember", params, &apiResponse)
	if err != nil {
		return nil, err
	}
	return apiResponse, err
}

// UnBanChatMember 解封群组成员
func (t *Telegram) UnBanChatMember(params *requests.UnBan) (*models.ApiResponse, error) {
	var apiResponse *models.ApiResponse
	err := t.post("unbanChatMember", params, &apiResponse)
	if err != nil {
		return nil, err
	}
	return apiResponse, err
}

// PromoteChatMember 提升群组成员权限
func (t *Telegram) PromoteChatMember(params *requests.Promote) (*models.ApiResponse, error) {
	var apiResponse *models.ApiResponse
	err := t.post("promoteChatMember", params, &apiResponse)
	if err != nil {
		return nil, err
	}
	return apiResponse, err
}

// SetChatPermissions 设置群组权限
func (t *Telegram) SetChatPermissions(params *requests.SetChatPermissions) (*models.ApiResponse, error) {
	var apiResponse *models.ApiResponse
	err := t.post("setChatPermissions", params, &apiResponse)
	if err != nil {
		return nil, err
	}
	return apiResponse, err
}

// SetChatAdministratorCustomTitle 设置管理员自定义头衔
func (t *Telegram) SetChatAdministratorCustomTitle(params *requests.CustomTitle) (*models.ApiResponse, error) {
	var apiResponse *models.ApiResponse
	err := t.post("setChatAdministratorCustomTitle", params, &apiResponse)
	if err != nil {
		return nil, err
	}
	return apiResponse, err
}

// SetChatMemberTag 设置群组成员标签
func (t *Telegram) SetChatMemberTag(params *requests.Tag) (*models.ApiResponse, error) {
	var apiResponse *models.ApiResponse
	err := t.post("setChatMemberTag", params, &apiResponse)
	if err != nil {
		return nil, err
	}
	return apiResponse, err
}

// GetChatAdministrators 获取群组管理员列表
func (t *Telegram) GetChatAdministrators(params *requests.Chat) (*models.ChatMembersResponse, error) {
	var apiResponse *models.ChatMembersResponse
	err := t.post("getChatAdministrators", params, &apiResponse)
	if err != nil {
		return nil, err
	}
	return apiResponse, err
}

// GetChatMemberCount 获取群组成员数量
func (t *Telegram) GetChatMemberCount(params *requests.Chat) (*models.ChatMemberCountResponse, error) {
	var apiResponse *models.ChatMemberCountResponse
	err := t.post("getChatMemberCount", params, &apiResponse)
	if err != nil {
		return nil, err
	}
	return apiResponse, err
}

// CreateChatInviteLink 创建群组邀请链接
func (t *Telegram) CreateChatInviteLink(params *requests.CreateChatInviteLink) (*models.ChatInviteLinkResponse, error) {
	var apiResponse *models.ChatInviteLinkResponse
	err := t.post("createChatInviteLink", params, &apiResponse)
	if err != nil {
		return nil, err
	}
	return apiResponse, err
}

// ApproveChatJoinRequest 批准入群请求
func (t *Telegram) ApproveChatJoinRequest(params *requests.Member) (*models.ApiResponse, error) {
	var apiResponse *models.ApiResponse
	err := t.post("approveChatJoinRequest", params, &apiResponse)
	if err != nil {
		return nil, err
	}
	return apiResponse, err
}

// DeclineChatJoinRequest 拒绝入群请求
func (t *Telegram) DeclineChatJoinRequest(params *requests.Member) (*models.ApiResponse, error) {
	var apiResponse *models.ApiResponse
	err := t.post("declineChatJoinRequest", params, &apiResponse)
	if err != nil {
		return nil, err
	}
	return apiResponse, err
}

// GetFile 获取文件信息
func (t *Telegram) GetFile(params *requests.File) (*models.FileResponse, error) {
	var apiResponse *models.FileResponse
	err := t.post("getFile", params, &apiResponse)
	if err != nil {
		return nil, err
	}
	return apiResponse, err
}

// Download 获取文件信息
func (t *Telegram) Download(params *requests.FilePath) ([]byte, error) {
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	if strings.Contains(t.client.BaseURL, "127.0.0.1") {
		data, err := os.ReadFile(params.FilePath)
		if err == nil {
			return data, nil
		}
	}

	baseUrl := strings.Replace(t.client.BaseURL, "/bot", "/file/bot", 1)
	fileUrl := fmt.Sprintf("%s/%s", baseUrl, params.FilePath)
	resp, err := client.Get(fileUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to download image: %w", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("download failed with status: %s, url: %s", resp.Status, fileUrl)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read image data: %w", err)
	}

	return data, nil
}

// Config Telegram API 配置
type Config struct {
	BaseApi   string        `json:"base_api"`   // API基础地址
	Timeout   time.Duration `json:"timeout"`    // 超时时间
	UserAgent string        `json:"user_agent"` // User-Agent
	Debug     bool          `json:"debug"`      // 是否开启调试模式
}

// 默认配置
var defaultConfig = &Config{
	BaseApi:   "https://api.telegram.org/bot",
	Timeout:   30,
	UserAgent: "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/146.0.0.0 Safari/537.36 Edg/146.0.0.0",
	Debug:     false,
}

// NewTelegram 多机器人实例模式创建
func NewTelegram(botId int64, token string, config ...*Config) *Telegram {
	// 先尝试读取
	mu.RLock()
	if i, exists := instances[botId]; exists {
		mu.RUnlock()
		return i
	}
	mu.RUnlock()

	// 不存在则创建（需要写锁）
	mu.Lock()
	defer mu.Unlock()

	// 双重检查，防止在等待锁的过程中被其他协程创建
	if i, exists := instances[botId]; exists {
		return i
	}

	// 使用传入的配置或默认配置
	cfg := defaultConfig
	if len(config) > 0 && config[0] != nil {
		cfg = config[0]
	}

	// 创建 Resty 客户端（带钩子）
	restyClient := newRestyClient(cfg)
	restyClient.SetBaseURL(fmt.Sprintf("%s%s", cfg.GetBaseApi(), token))

	// 创建新实例
	i := &Telegram{
		client: restyClient,
		debug:  cfg.Debug,
		botId:  botId,
		token:  token,
	}

	// 存入实例池
	instances[botId] = i

	return i
}

// GetTelegram 从实例池中获取指定机器人实例
func GetTelegram(botID int64) *Telegram {
	mu.RLock()
	defer mu.RUnlock()

	instance, exists := instances[botID]
	if !exists {
		return nil
	}
	return instance
}

// GetBaseApi 获取API基础地址
func (params *Config) GetBaseApi() string {
	if params.BaseApi == "" {
		return defaultConfig.BaseApi
	}
	return params.BaseApi
}

// GetTimeout 获取超时时间
func (params *Config) GetTimeout() time.Duration {
	return params.Timeout * time.Second
}

// GetUserAgent 获取User-Agent
func (params *Config) GetUserAgent() string {
	if params.UserAgent == "" {
		return defaultConfig.UserAgent
	}
	return params.UserAgent
}

// NewTelegramApi 单实例模式
func NewTelegramApi(token string, config ...*Config) {
	once.Do(func() {
		// 使用传入的配置或默认配置
		cfg := defaultConfig
		if len(config) > 0 && config[0] != nil {
			cfg = config[0]
		}

		restyClient := newRestyClient(cfg)
		restyClient.SetBaseURL(fmt.Sprintf("%s%s", cfg.GetBaseApi(), token))

		botInfo := strings.Split(token, ":")
		botId, _ := strconv.ParseInt(botInfo[0], 10, 64)

		instance = &Telegram{
			client: restyClient,
			debug:  cfg.Debug,
			botId:  botId,
			token:  token,
		}
	})
}

// GetTelegramApi 获取单实例
func GetTelegramApi() *Telegram {
	if instance == nil {
		panic("telegram instance is nil")
	}

	return instance
}
