package telegram_api

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/bytedance/sonic"
	"github.com/jonneyless/telegram-api/models"
	"github.com/jonneyless/telegram-api/requests"
)

var telegram *Telegram

type Telegram struct {
	client  *httpClient
	baseApi string
	botId   int64
	debug   bool
}

type MessageDeleteOptions struct {
	Delay       time.Duration
	DeleteReply bool
	MessageIds  []int64
}

func (t *Telegram) IsNewBot(botId int64) bool {
	return t.botId != botId
}

func (t *Telegram) SetToken(token string) *Telegram {
	parts := strings.Split(token, ":")
	botId, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return t
	}

	t.botId = botId
	t.client.setBaseUrl(t.baseApi + token + "/")

	return t
}

func (t *Telegram) post(path string, params interface{}, result interface{}) error {
	err := t.client.post(path, params, result)
	if t.debug {
		logger.Debug(fmt.Sprintf("Path: %v", path))
		paramsJson, _ := sonic.MarshalIndent(params, "", "  ")
		logger.Debug(fmt.Sprintf("Params: %s", string(paramsJson)))
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

func (t *Telegram) get(path string, params map[string]string, result interface{}) (interface{}, error) {
	err := t.client.get(path, params, result)

	if err != nil {
		return nil, err
	}

	return result, nil
}

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

func (t *Telegram) EditMessageText(params *requests.EditMessage) (*models.MessageResponse, error) {
	var apiResponse *models.MessageResponse

	err := t.post("editMessageText", params.GetParams(), &apiResponse)
	if err != nil {
		return nil, err
	}

	return apiResponse, err
}

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

func (t *Telegram) SendPoll(params *requests.SendPoll) (*models.MessageResponse, error) {
	var apiResponse *models.MessageResponse

	err := t.post("sendPoll", params, &apiResponse)
	if err != nil {
		return nil, err
	}

	return apiResponse, err
}

func (t *Telegram) SendDice(params *requests.SendDice) (*models.MessageResponse, error) {
	var apiResponse *models.MessageResponse

	err := t.post("sendDice", params, &apiResponse)
	if err != nil {
		return nil, err
	}

	return apiResponse, err
}

func (t *Telegram) AnswerCallbackQuery(params *requests.AnswerCallbackQuery) (*models.ApiResponse, error) {
	var apiResponse *models.ApiResponse

	err := t.post("answerCallbackQuery", params.GetParams(), &apiResponse)
	if err != nil {
		return nil, err
	}

	return apiResponse, err
}

func (t *Telegram) GetChat(params *requests.Chat) (*models.ChatResponse, error) {
	var apiResponse *models.ChatResponse

	err := t.post("getChat", params, &apiResponse)
	if err != nil {
		return nil, err
	}

	return apiResponse, err
}

func (t *Telegram) LeaveChat(params *requests.Chat) (*models.ApiResponse, error) {
	var apiResponse *models.ApiResponse

	err := t.post("leaveChat", params, &apiResponse)
	if err != nil {
		return nil, err
	}

	return apiResponse, err
}

func (t *Telegram) SetChatPhoto(params *requests.ChatPhoto) (*models.ApiResponse, error) {
	var apiResponse *models.ApiResponse

	err := t.post("setChatPhoto", params, &apiResponse)
	if err != nil {
		return nil, err
	}

	return apiResponse, err
}

func (t *Telegram) DeleteChatPhoto(params *requests.Chat) (*models.ApiResponse, error) {
	var apiResponse *models.ApiResponse

	err := t.post("deleteChatPhoto", params, &apiResponse)
	if err != nil {
		return nil, err
	}

	return apiResponse, err
}

func (t *Telegram) SetChatTitle(params *requests.ChatTitle) (*models.ApiResponse, error) {
	var apiResponse *models.ApiResponse

	err := t.post("setChatTitle", params, &apiResponse)
	if err != nil {
		return nil, err
	}

	return apiResponse, err
}

func (t *Telegram) SetChatDescription(params *requests.ChatDescription) (*models.ApiResponse, error) {
	var apiResponse *models.ApiResponse

	err := t.post("setChatDescription", params, &apiResponse)
	if err != nil {
		return nil, err
	}

	return apiResponse, err
}

func (t *Telegram) PinChatMessage(params *requests.Message) (*models.ApiResponse, error) {
	var apiResponse *models.ApiResponse

	err := t.post("pinChatMessage", params, &apiResponse)
	if err != nil {
		return nil, err
	}

	return apiResponse, err
}

func (t *Telegram) UnPinChatMessage(params *requests.Message) (*models.ApiResponse, error) {
	var apiResponse *models.ApiResponse

	err := t.post("unpinChatMessage", params, &apiResponse)
	if err != nil {
		return nil, err
	}

	return apiResponse, err
}

func (t *Telegram) UnPinAllChatMessage(params *requests.Chat) (*models.ApiResponse, error) {
	var apiResponse *models.ApiResponse

	err := t.post("unpinAllChatMessage", params, &apiResponse)
	if err != nil {
		return nil, err
	}

	return apiResponse, err
}

func (t *Telegram) GetChatMember(params *requests.Member) (*models.ChatMemberResponse, error) {
	var apiResponse *models.ChatMemberResponse

	err := t.post("getChatMember", params, &apiResponse)
	if err != nil {
		return nil, err
	}

	return apiResponse, err
}

func (t *Telegram) RestrictChatMember(params *requests.Restrict) (*models.ApiResponse, error) {
	var apiResponse *models.ApiResponse

	err := t.post("restrictChatMember", params, &apiResponse)
	if err != nil {
		return nil, err
	}

	return apiResponse, err
}

func (t *Telegram) BanChatMember(params *requests.Ban) (*models.ApiResponse, error) {
	var apiResponse *models.ApiResponse

	err := t.post("banChatMember", params, &apiResponse)
	if err != nil {
		return nil, err
	}

	return apiResponse, err
}

func (t *Telegram) UnBanChatMember(params *requests.UnBan) (*models.ApiResponse, error) {
	var apiResponse *models.ApiResponse

	err := t.post("unbanChatMember", params, &apiResponse)
	if err != nil {
		return nil, err
	}

	return apiResponse, err
}

func (t *Telegram) PromoteChatMember(params *requests.Promote) (*models.ApiResponse, error) {
	var apiResponse *models.ApiResponse

	err := t.post("promoteChatMember", params, &apiResponse)
	if err != nil {
		return nil, err
	}

	return apiResponse, err
}

func (t *Telegram) SetChatPermissions(params *requests.SetChatPermissions) (*models.ApiResponse, error) {
	var apiResponse *models.ApiResponse

	err := t.post("setChatPermissions", params, &apiResponse)
	if err != nil {
		return nil, err
	}

	return apiResponse, err
}

func (t *Telegram) SetChatAdministratorCustomTitle(params *requests.CustomTitle) (*models.ApiResponse, error) {
	var apiResponse *models.ApiResponse

	err := t.post("setChatAdministratorCustomTitle", params, &apiResponse)
	if err != nil {
		return nil, err
	}

	return apiResponse, err
}

func (t *Telegram) SetChatMemberTag(params *requests.Tag) (*models.ApiResponse, error) {
	var apiResponse *models.ApiResponse

	err := t.post("setChatMemberTag", params, &apiResponse)
	if err != nil {
		return nil, err
	}

	return apiResponse, err
}

func (t *Telegram) GetChatAdministrators(params *requests.Chat) (*models.ChatMembersResponse, error) {
	var apiResponse *models.ChatMembersResponse

	err := t.post("getChatAdministrators", params, &apiResponse)
	if err != nil {
		return nil, err
	}

	return apiResponse, err
}

func (t *Telegram) GetChatMemberCount(params *requests.Chat) (*models.ChatMemberCountResponse, error) {
	var apiResponse *models.ChatMemberCountResponse

	err := t.post("getChatMemberCount", params, &apiResponse)
	if err != nil {
		return nil, err
	}

	return apiResponse, err
}

func (t *Telegram) CreateChatInviteLink(params *requests.CreateChatInviteLink) (*models.ChatInviteLinkResponse, error) {
	var apiResponse *models.ChatInviteLinkResponse

	err := t.post("createChatInviteLink", params, &apiResponse)
	if err != nil {
		return nil, err
	}

	return apiResponse, err
}

func (t *Telegram) ApproveChatJoinRequest(params *requests.Member) (*models.ApiResponse, error) {
	var apiResponse *models.ApiResponse

	err := t.post("approveChatJoinRequest", params, &apiResponse)
	if err != nil {
		return nil, err
	}

	return apiResponse, err
}

func (t *Telegram) DeclineChatJoinRequest(params *requests.Member) (*models.ApiResponse, error) {
	var apiResponse *models.ApiResponse

	err := t.post("declineChatJoinRequest", params, &apiResponse)
	if err != nil {
		return nil, err
	}

	return apiResponse, err
}

func (t *Telegram) GetFile(params *requests.File) (*models.FileResponse, error) {
	var apiResponse *models.FileResponse

	err := t.post("getFile", params, &apiResponse)
	if err != nil {
		return nil, err
	}

	return apiResponse, err
}

type TelegramApi struct {
	BaseApi   string        `json:"base_api"`
	Timeout   time.Duration `json:"timeout"`
	UserAgent string        `json:"user_agent"`
	Debug     bool          `json:"debug"`
}

func (params *TelegramApi) GetBaseApi() string {
	if params.BaseApi == "" {
		params.BaseApi = "https://api.telegram.org/bot"
	}

	return params.BaseApi
}

func (params *TelegramApi) GetTimeout() time.Duration {
	return params.Timeout * time.Second
}

func (params *TelegramApi) GetUserAgent() string {
	if params.UserAgent == "" {
		params.UserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/146.0.0.0 Safari/537.36 Edg/146.0.0.0"
	}

	return params.UserAgent
}

func NewTelegramApi(params *TelegramApi) {
	telegram = &Telegram{
		client: newHTTPClient(
			"",
			withTimeout(params.GetTimeout()),
			withHeaders(map[string]string{
				"User-Agent": params.GetUserAgent(),
			}),
		),
		baseApi: params.GetBaseApi(),
		debug:   params.Debug,
	}
}

func GetTelegramApi() *Telegram {
	return telegram
}
