package telegram_api

import (
	"fmt"
	"sync"
	"time"

	"github.com/bytedance/sonic"
	"github.com/jonneyless/telegram-api/models"
	"github.com/jonneyless/telegram-api/requests"
)

// 单例模式实例
var legacyTelegram *Telegram

// 多实例模式存储
var (
	instances = make(map[int64]*Telegram)
	mu        sync.RWMutex
)

// Telegram Telegram API 客户端
type Telegram struct {
	client  *httpClient // HTTP客户端
	baseApi string      // API基础地址
	debug   bool        // 是否开启调试模式
	token   string      // Bot Token
	botID   int64       // Bot ID
}

// MessageDeleteOptions 消息删除选项
type MessageDeleteOptions struct {
	Delay       time.Duration // 延迟删除时间
	DeleteReply bool          // 是否删除回复消息
	MessageIds  []int64       // 要删除的消息ID列表
}

// SetToken 设置Bot Token
func (t *Telegram) SetToken(token string) *Telegram {
	t.token = token
	t.client.setBaseUrl(t.baseApi + token + "/")
	return t
}

// GetBotID 获取Bot ID
func (t *Telegram) GetBotID() int64 {
	return t.botID
}

// GetToken 获取Bot Token
func (t *Telegram) GetToken() string {
	return t.token
}

// post 发送POST请求
func (t *Telegram) post(path string, params interface{}, result interface{}) error {
	if t.debug {
		logger.Debug(fmt.Sprintf("Path: %v", path))
		paramsJson, _ := sonic.MarshalIndent(params, "", "  ")
		logger.Debug(fmt.Sprintf("Params: %s", string(paramsJson)))
	}

	err := t.client.post(path, params, result)

	if t.debug {
		if result != nil {
			resultJson, _ := sonic.MarshalIndent(result, "", "  ")
			logger.Debug(fmt.Sprintf("Result: %s", string(resultJson)))
		} else {
			logger.Debug("Result: nil")
		}
	}

	if err != nil {
		return err
	}
	return nil
}

// get 发送GET请求
func (t *Telegram) get(path string, params map[string]string, result interface{}) (interface{}, error) {
	err := t.client.get(path, params, result)
	if err != nil {
		return nil, err
	}
	return result, nil
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
		messageIds := options.MessageIds
		if options.DeleteReply && params.ReplyParameters != nil {
			messageIds = append(messageIds, params.ReplyParameters.MessageId)
		}
		if apiResponse != nil {
			go func() {
				messageIds = append(messageIds, apiResponse.Result.MessageID)
				time.Sleep(time.Second * options.Delay)
				_, _ = t.DeleteMessage(&requests.Message{
					ChatId:     apiResponse.Result.Chat.ID,
					MessageIds: &messageIds,
				})
			}()
		}
		return apiResponse, nil
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

// GetChat 获取聊天信息
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
	err := t.post("setChatPhoto", params, &apiResponse)
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

// TelegramApi Telegram API 配置
type TelegramApi struct {
	BaseApi   string        `json:"base_api"`   // API基础地址
	Timeout   time.Duration `json:"timeout"`    // 超时时间
	UserAgent string        `json:"user_agent"` // User-Agent
	Debug     bool          `json:"debug"`      // 是否开启调试模式
}

// 默认配置
var defaultConfig = &TelegramApi{
	BaseApi:   "https://api.telegram.org/bot",
	Timeout:   30,
	UserAgent: "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/146.0.0.0 Safari/537.36 Edg/146.0.0.0",
	Debug:     false,
}

// GetOrCreateTelegram 获取或创建Telegram客户端实例（多实例模式）
func GetOrCreateTelegram(botID int64, token string, config ...*TelegramApi) *Telegram {
	// 先尝试读取
	mu.RLock()
	if instance, exists := instances[botID]; exists {
		mu.RUnlock()
		return instance
	}
	mu.RUnlock()

	// 不存在则创建（需要写锁）
	mu.Lock()
	defer mu.Unlock()

	// 双重检查，防止在等待锁的过程中被其他协程创建
	if instance, exists := instances[botID]; exists {
		return instance
	}

	// 使用传入的配置或默认配置
	cfg := defaultConfig
	if len(config) > 0 && config[0] != nil {
		cfg = config[0]
	}

	// 创建新实例
	instance := &Telegram{
		client: newHTTPClient(
			"",
			withTimeout(cfg.GetTimeout()),
			withHeaders(map[string]string{
				"User-Agent": cfg.GetUserAgent(),
			}),
		),
		baseApi: cfg.GetBaseApi(),
		debug:   cfg.Debug,
		token:   token,
		botID:   botID,
	}

	// 设置token
	instance.SetToken(token)

	// 存入实例池
	instances[botID] = instance

	return instance
}

// GetTelegram 获取Telegram客户端实例（多实例模式）
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
func (params *TelegramApi) GetBaseApi() string {
	if params.BaseApi == "" {
		params.BaseApi = "https://api.telegram.org/bot"
	}
	return params.BaseApi
}

// GetTimeout 获取超时时间
func (params *TelegramApi) GetTimeout() time.Duration {
	return params.Timeout * time.Second
}

// GetUserAgent 获取User-Agent
func (params *TelegramApi) GetUserAgent() string {
	if params.UserAgent == "" {
		params.UserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/146.0.0.0 Safari/537.36 Edg/146.0.0.0"
	}
	return params.UserAgent
}

// NewTelegramApi 创建Telegram API客户端（单实例模式）
func NewTelegramApi(config ...*TelegramApi) {
	// 使用传入的配置或默认配置
	cfg := defaultConfig
	if len(config) > 0 && config[0] != nil {
		cfg = config[0]
	}

	legacyTelegram = &Telegram{
		client: newHTTPClient(
			"",
			withTimeout(cfg.GetTimeout()),
			withHeaders(map[string]string{
				"User-Agent": cfg.GetUserAgent(),
			}),
		),
		baseApi: cfg.GetBaseApi(),
		debug:   cfg.Debug,
	}
}

// GetTelegramApi 获取Telegram API客户端（单实例模式）
func GetTelegramApi() *Telegram {
	return legacyTelegram
}
