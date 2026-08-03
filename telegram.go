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

	// 使用 Resty 发送 POST 请求
	resp, err := t.client.R().
		SetBody(params).
		SetResult(&result).
		SetError(&result).
		Post(path)

	if err != nil {
		return err
	}

	// 检查 HTTP 状态码
	if resp.StatusCode() >= 400 {
		return fmt.Errorf("HTTP error: %d", resp.StatusCode())
	}

	return nil
}

func (t *Telegram) postMultipart(path string, body *bytes.Buffer, contentType string, result any) error {
	if t.token == "" {
		return fmt.Errorf("bot token empty")
	}

	resp, err := t.client.R().
		SetHeader("Content-Type", contentType).
		SetBody(body.Bytes()).
		SetResult(&result).
		Post(path)

	if err != nil {
		return err
	}

	// 检查 HTTP 状态码
	if resp.StatusCode() >= 400 {
		return fmt.Errorf("HTTP error: %d", resp.StatusCode())
	}

	return nil
}

// SendMessage 发送消息
func (t *Telegram) SendMessage(params *requests.SendMessage, deleteOptions ...*MessageDeleteOptions) (*models.Message, error) {
	var apiResponse *models.Response[models.Message]

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

	if !apiResponse.Ok {
		return nil, fmt.Errorf("telegram error: %s", apiResponse.Description)
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
					MessageIds: messageIds,
				})
			}()
		}
	}

	return apiResponse.Result, nil
}

// EditMessageText 编辑消息文本
func (t *Telegram) EditMessageText(params *requests.EditMessageText) (*models.Message, error) {
	var apiResponse *models.Response[models.Message]
	err := t.post("editMessageText", params.GetParams(), &apiResponse)
	if err != nil {
		return nil, err
	}
	if !apiResponse.Ok {
		return nil, fmt.Errorf("telegram error: %s", apiResponse.Description)
	}
	return apiResponse.Result, nil
}

// EditMessageCaption 编辑消息文本
func (t *Telegram) EditMessageCaption(params *requests.EditMessageCaption) (*models.Message, error) {
	var apiResponse *models.Response[models.Message]
	err := t.post("editMessageText", params, &apiResponse)
	if err != nil {
		return nil, err
	}
	if !apiResponse.Ok {
		return nil, fmt.Errorf("telegram error: %s", apiResponse.Description)
	}
	return apiResponse.Result, nil
}

// EditMessageMedia 编辑媒体先遨嬉
func (t *Telegram) EditMessageMedia(params *requests.EditMessageMedia) (*models.Message, error) {
	var apiResponse *models.Response[models.Message]
	err := t.post("editMessageMedia", params, &apiResponse)
	if err != nil {
		return nil, err
	}
	if !apiResponse.Ok {
		return nil, fmt.Errorf("telegram error: %s", apiResponse.Description)
	}
	return apiResponse.Result, nil
}

// DeleteMessage 删除消息
func (t *Telegram) DeleteMessage(params *requests.Message) (bool, error) {
	var apiResponse *models.Response[bool]
	var err error

	if params.MessageIds != nil {
		err = t.post("deleteMessages", params, &apiResponse)
	} else {
		err = t.post("deleteMessage", params, &apiResponse)
	}

	if err != nil {
		return false, err
	}

	if !apiResponse.Ok {
		return false, fmt.Errorf("telegram error: %s", apiResponse.Description)
	}

	return apiResponse.Ok, nil
}

// SendPhoto 图片二进制发布
func (t *Telegram) SendPhoto(params *requests.SendPhoto) (*models.Message, error) {
	var apiResponse *models.Response[models.Message]
	body, contentType, err := params.ToMultipart()
	if err != nil {
		return nil, err
	}
	err = t.postMultipart("sendPhoto", body, contentType, &apiResponse)
	if err != nil {
		return nil, err
	}

	if !apiResponse.Ok {
		return nil, fmt.Errorf("telegram error: %s", apiResponse.Description)
	}

	return apiResponse.Result, nil
}

// SendPoll 发送投票
func (t *Telegram) SendPoll(params *requests.SendPoll) (*models.Message, error) {
	var apiResponse *models.Response[models.Message]
	err := t.post("sendPoll", params, &apiResponse)
	if err != nil {
		return nil, err
	}
	if !apiResponse.Ok {
		return nil, fmt.Errorf("telegram error: %s", apiResponse.Description)
	}
	return apiResponse.Result, nil
}

// SendDice 发送骰子
func (t *Telegram) SendDice(params *requests.SendDice) (*models.Message, error) {
	var apiResponse *models.Response[models.Message]
	err := t.post("sendDice", params, &apiResponse)
	if err != nil {
		return nil, err
	}
	if !apiResponse.Ok {
		return nil, fmt.Errorf("telegram error: %s", apiResponse.Description)
	}
	return apiResponse.Result, nil
}

// AnswerCallbackQuery 回答回调查询
func (t *Telegram) AnswerCallbackQuery(params *requests.AnswerCallbackQuery) (bool, error) {
	var apiResponse *models.Response[bool]
	err := t.post("answerCallbackQuery", params, &apiResponse)
	if err != nil {
		return false, err
	}
	if !apiResponse.Ok {
		return false, fmt.Errorf("telegram error: %s", apiResponse.Description)
	}
	return apiResponse.Ok, nil
}

// GetChat 获取群组信息
func (t *Telegram) GetChat(params *requests.Chat) (*models.ChatFull, error) {
	var apiResponse *models.Response[models.ChatFull]
	err := t.post("getChat", params, &apiResponse)
	if err != nil {
		return nil, err
	}
	if !apiResponse.Ok {
		return nil, fmt.Errorf("telegram error: %s", apiResponse.Description)
	}
	return apiResponse.Result, nil
}

// LeaveChat 离开群组
func (t *Telegram) LeaveChat(params *requests.Chat) (bool, error) {
	var apiResponse *models.Response[bool]
	err := t.post("leaveChat", params, &apiResponse)
	if err != nil {
		return false, err
	}
	if !apiResponse.Ok {
		return false, fmt.Errorf("telegram error: %s", apiResponse.Description)
	}
	return apiResponse.Ok, nil
}

// SetChatPhoto 设置群组头像
func (t *Telegram) SetChatPhoto(params *requests.ChatPhoto) (bool, error) {
	var apiResponse *models.Response[bool]
	body, contentType, err := params.ToMultipart()
	if err != nil {
		return false, err
	}
	err = t.postMultipart("setChatPhoto", body, contentType, &apiResponse)
	if err != nil {
		return false, err
	}
	if !apiResponse.Ok {
		return false, fmt.Errorf("telegram error: %s", apiResponse.Description)
	}
	return apiResponse.Ok, nil
}

// DeleteChatPhoto 删除群组头像
func (t *Telegram) DeleteChatPhoto(params *requests.Chat) (bool, error) {
	var apiResponse *models.Response[bool]
	err := t.post("deleteChatPhoto", params, &apiResponse)
	if err != nil {
		return false, err
	}
	if !apiResponse.Ok {
		return false, fmt.Errorf("telegram error: %s", apiResponse.Description)
	}
	return apiResponse.Ok, nil
}

// SetChatTitle 设置群组标题
func (t *Telegram) SetChatTitle(params *requests.ChatTitle) (bool, error) {
	var apiResponse *models.Response[bool]
	err := t.post("setChatTitle", params, &apiResponse)
	if err != nil {
		return false, err
	}
	if !apiResponse.Ok {
		return false, fmt.Errorf("telegram error: %s", apiResponse.Description)
	}
	return apiResponse.Ok, nil
}

// SetChatDescription 设置群组描述
func (t *Telegram) SetChatDescription(params *requests.ChatDescription) (bool, error) {
	var apiResponse *models.Response[bool]
	err := t.post("setChatDescription", params, &apiResponse)
	if err != nil {
		return false, err
	}
	if !apiResponse.Ok {
		return false, fmt.Errorf("telegram error: %s", apiResponse.Description)
	}
	return apiResponse.Ok, nil
}

// PinChatMessage 置顶消息
func (t *Telegram) PinChatMessage(params *requests.Message) (bool, error) {
	var apiResponse *models.Response[bool]
	err := t.post("pinChatMessage", params, &apiResponse)
	if err != nil {
		return false, err
	}
	if !apiResponse.Ok {
		return false, fmt.Errorf("telegram error: %s", apiResponse.Description)
	}
	return apiResponse.Ok, nil
}

// UnPinChatMessage 取消置顶消息
func (t *Telegram) UnPinChatMessage(params *requests.Message) (bool, error) {
	var apiResponse *models.Response[bool]
	err := t.post("unpinChatMessage", params, &apiResponse)
	if err != nil {
		return false, err
	}
	if !apiResponse.Ok {
		return false, fmt.Errorf("telegram error: %s", apiResponse.Description)
	}
	return apiResponse.Ok, nil
}

// UnPinAllChatMessage 取消所有置顶消息
func (t *Telegram) UnPinAllChatMessage(params *requests.Chat) (bool, error) {
	var apiResponse *models.Response[bool]
	err := t.post("unpinAllChatMessage", params, &apiResponse)
	if err != nil {
		return false, err
	}
	if !apiResponse.Ok {
		return false, fmt.Errorf("telegram error: %s", apiResponse.Description)
	}
	return apiResponse.Ok, nil
}

// GetChatMember 获取群组成员信息
func (t *Telegram) GetChatMember(params *requests.Member) (*models.Member, error) {
	var apiResponse *models.Response[models.Member]
	err := t.post("getChatMember", params, &apiResponse)
	if err != nil {
		return nil, err
	}
	if !apiResponse.Ok {
		return nil, fmt.Errorf("telegram error: %s", apiResponse.Description)
	}
	return apiResponse.Result, nil
}

// RestrictChatMember 限制群组成员
func (t *Telegram) RestrictChatMember(params *requests.Restrict) (bool, error) {
	var apiResponse *models.Response[bool]
	err := t.post("restrictChatMember", params, &apiResponse)
	if err != nil {
		return false, err
	}
	if !apiResponse.Ok {
		return false, fmt.Errorf("telegram error: %s", apiResponse.Description)
	}
	return apiResponse.Ok, nil
}

// BanChatMember 封禁群组成员
func (t *Telegram) BanChatMember(params *requests.Ban) (bool, error) {
	var apiResponse *models.Response[bool]
	err := t.post("banChatMember", params, &apiResponse)
	if err != nil {
		return false, err
	}
	if !apiResponse.Ok {
		return false, fmt.Errorf("telegram error: %s", apiResponse.Description)
	}
	return apiResponse.Ok, nil
}

// UnBanChatMember 解封群组成员
func (t *Telegram) UnBanChatMember(params *requests.UnBan) (bool, error) {
	var apiResponse *models.Response[bool]
	err := t.post("unbanChatMember", params, &apiResponse)
	if err != nil {
		return false, err
	}
	if !apiResponse.Ok {
		return false, fmt.Errorf("telegram error: %s", apiResponse.Description)
	}
	return apiResponse.Ok, nil
}

// PromoteChatMember 提升群组成员权限
func (t *Telegram) PromoteChatMember(params *requests.Promote) (bool, error) {
	var apiResponse *models.Response[bool]
	err := t.post("promoteChatMember", params, &apiResponse)
	if err != nil {
		return false, err
	}
	if !apiResponse.Ok {
		return false, fmt.Errorf("telegram error: %s", apiResponse.Description)
	}
	return apiResponse.Ok, nil
}

// SetChatPermissions 设置群组权限
func (t *Telegram) SetChatPermissions(params *requests.SetChatPermissions) (bool, error) {
	var apiResponse *models.Response[bool]
	err := t.post("setChatPermissions", params, &apiResponse)
	if err != nil {
		return false, err
	}
	if !apiResponse.Ok {
		return false, fmt.Errorf("telegram error: %s", apiResponse.Description)
	}
	return apiResponse.Ok, nil
}

// SetChatAdministratorCustomTitle 设置管理员自定义头衔
func (t *Telegram) SetChatAdministratorCustomTitle(params *requests.CustomTitle) (bool, error) {
	var apiResponse *models.Response[bool]
	err := t.post("setChatAdministratorCustomTitle", params, &apiResponse)
	if err != nil {
		return false, err
	}
	if !apiResponse.Ok {
		return false, fmt.Errorf("telegram error: %s", apiResponse.Description)
	}
	return apiResponse.Ok, nil
}

// SetChatMemberTag 设置群组成员标签
func (t *Telegram) SetChatMemberTag(params *requests.Tag) (bool, error) {
	var apiResponse *models.Response[bool]
	err := t.post("setChatMemberTag", params, &apiResponse)
	if err != nil {
		return false, err
	}
	if !apiResponse.Ok {
		return false, fmt.Errorf("telegram error: %s", apiResponse.Description)
	}
	return apiResponse.Ok, nil
}

// GetChatAdministrators 获取群组管理员列表
func (t *Telegram) GetChatAdministrators(params *requests.Chat) (*[]models.Member, error) {
	var apiResponse *models.Response[[]models.Member]
	err := t.post("getChatAdministrators", params, &apiResponse)
	if err != nil {
		return nil, err
	}
	if !apiResponse.Ok {
		return nil, fmt.Errorf("telegram error: %s", apiResponse.Description)
	}
	return apiResponse.Result, nil
}

// GetChatMemberCount 获取群组成员数量
func (t *Telegram) GetChatMemberCount(params *requests.Chat) (*int64, error) {
	var apiResponse *models.Response[int64]
	err := t.post("getChatMemberCount", params, &apiResponse)
	if err != nil {
		return nil, err
	}
	if !apiResponse.Ok {
		return nil, fmt.Errorf("telegram error: %s", apiResponse.Description)
	}
	return apiResponse.Result, nil
}

// CreateChatInviteLink 创建群组邀请链接
func (t *Telegram) CreateChatInviteLink(params *requests.CreateChatInviteLink) (*models.InviteLink, error) {
	var apiResponse *models.Response[models.InviteLink]
	err := t.post("createChatInviteLink", params, &apiResponse)
	if err != nil {
		return nil, err
	}
	if !apiResponse.Ok {
		return nil, fmt.Errorf("telegram error: %s", apiResponse.Description)
	}
	return apiResponse.Result, nil
}

// EditChatInviteLink 编辑群组邀请链接
func (t *Telegram) EditChatInviteLink(params *requests.CreateChatInviteLink) (*models.InviteLink, error) {
	var apiResponse *models.Response[models.InviteLink]
	err := t.post("editChatInviteLink", params, &apiResponse)
	if err != nil {
		return nil, err
	}
	if !apiResponse.Ok {
		return nil, fmt.Errorf("telegram error: %s", apiResponse.Description)
	}
	return apiResponse.Result, nil
}

// RevokeChatInviteLink 撤销群组邀请链接
func (t *Telegram) RevokeChatInviteLink(params *requests.RevokeChatInviteLink) (*models.InviteLink, error) {
	var apiResponse *models.Response[models.InviteLink]
	err := t.post("createChatInviteLink", params, &apiResponse)
	if err != nil {
		return nil, err
	}
	if !apiResponse.Ok {
		return nil, fmt.Errorf("telegram error: %s", apiResponse.Description)
	}
	return apiResponse.Result, nil
}

// ApproveChatJoinRequest 批准入群请求
func (t *Telegram) ApproveChatJoinRequest(params *requests.Member) (bool, error) {
	var apiResponse *models.Response[bool]
	err := t.post("approveChatJoinRequest", params, &apiResponse)
	if err != nil {
		return false, err
	}
	if !apiResponse.Ok {
		return false, fmt.Errorf("telegram error: %s", apiResponse.Description)
	}
	return apiResponse.Ok, nil
}

// DeclineChatJoinRequest 拒绝入群请求
func (t *Telegram) DeclineChatJoinRequest(params *requests.Member) (bool, error) {
	var apiResponse *models.Response[bool]
	err := t.post("declineChatJoinRequest", params, &apiResponse)
	if err != nil {
		return false, err
	}
	if !apiResponse.Ok {
		return false, fmt.Errorf("telegram error: %s", apiResponse.Description)
	}
	return apiResponse.Ok, nil
}

// GetFile 获取文件信息
func (t *Telegram) GetFile(params *requests.File) (*models.File, error) {
	var apiResponse *models.Response[models.File]
	err := t.post("getFile", params, &apiResponse)
	if err != nil {
		return nil, err
	}
	if !apiResponse.Ok {
		return nil, fmt.Errorf("telegram error: %s", apiResponse.Description)
	}
	return apiResponse.Result, nil
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
func NewTelegramApi(token string, config ...*Config) *Telegram {
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

	return instance
}

// GetTelegramApi 获取单实例
func GetTelegramApi() *Telegram {
	if instance == nil {
		panic("telegram instance is nil")
	}

	return instance
}
