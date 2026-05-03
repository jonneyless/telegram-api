package telegram_api

import (
	"telegram_api/models"
	"telegram_api/requests"
	"time"
)

var telegram *Telegram

type Telegram struct {
	client  *httpClient
	baseApi string
}

type MessageDeleteOptions struct {
	Delay       time.Duration
	DeleteReply bool
	MessageIds  []int64
}

func (t *Telegram) SetToken(token string) *Telegram {
	t.client.setBaseUrl(t.baseApi + token + "/")

	return t
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

	_, err := t.post(apiPath, params.GetParams(), &apiResponse)
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
				time.Sleep(time.Second * options.Delay)
				_, _ = t.DeleteMessage(&requests.Message{
					ChatId:     apiResponse.Result.Chat.ID,
					MessageId:  apiResponse.Result.MessageID,
					MessageIds: messageIds,
				})
			}()
		}

		return apiResponse, nil
	}

	return apiResponse, nil
}

func (t *Telegram) post(path string, params interface{}, result interface{}) (interface{}, error) {
	err := t.client.post(path, params, &result)

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (t *Telegram) get(path string, params map[string]string, result interface{}) (interface{}, error) {
	err := t.client.get(path, params, &result)

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (t *Telegram) DeleteMessage(params *requests.Message) (*models.ApiResponse, error) {
	var apiResponse *models.ApiResponse

	result, err := t.post("deleteMessage", params, &apiResponse)

	return result.(*models.ApiResponse), err
}

func (t *Telegram) SendPoll(params *requests.SendPoll) (*models.MessageResponse, error) {
	var apiResponse *models.MessageResponse

	result, err := t.post("sendPoll", params, &apiResponse)

	return result.(*models.MessageResponse), err
}

func (t *Telegram) SendDice(params *requests.SendDice) (*models.MessageResponse, error) {
	var apiResponse *models.MessageResponse

	result, err := t.post("sendDice", params, &apiResponse)

	return result.(*models.MessageResponse), err
}

func (t *Telegram) AnswerCallbackQuery(params *requests.AnswerCallbackQuery) (*models.ApiResponse, error) {
	var apiResponse *models.ApiResponse

	result, err := t.post("answerCallbackQuery", params.GetParams(), &apiResponse)

	return result.(*models.ApiResponse), err
}

func (t *Telegram) GetChat(params *requests.Chat) (*models.ChatResponse, error) {
	var apiResponse *models.ChatResponse

	result, err := t.post("getChat", params, &apiResponse)

	return result.(*models.ChatResponse), err
}

func (t *Telegram) GetChatMember(params *requests.Member) (*models.ChatMemberResponse, error) {
	var apiResponse *models.ChatMemberResponse

	result, err := t.post("getChatMember", params, &apiResponse)

	return result.(*models.ChatMemberResponse), err
}

func (t *Telegram) RestrictChatMember(params *requests.Restrict) (*models.ApiResponse, error) {
	var apiResponse *models.ApiResponse

	result, err := t.post("restrictChatMember", params, &apiResponse)

	return result.(*models.ApiResponse), err
}

func (t *Telegram) PromoteChatMember(params *requests.Promote) (*models.ApiResponse, error) {
	var apiResponse *models.ApiResponse

	result, err := t.post("promoteChatMember", params, &apiResponse)

	return result.(*models.ApiResponse), err
}

func (t *Telegram) PromoteChatMember2(params map[string]string) (*models.ApiResponse, error) {
	var apiResponse *models.ApiResponse

	result, err := t.get("promoteChatMember", params, &apiResponse)

	return result.(*models.ApiResponse), err
}

func (t *Telegram) GetChatAdministrators(params *requests.Chat) (*models.ChatMembersResponse, error) {
	var apiResponse *models.ChatMembersResponse

	result, err := t.post("getChatAdministrators", params, &apiResponse)

	return result.(*models.ChatMembersResponse), err
}

func (t *Telegram) CreateChatInviteLink(params *requests.CreateChatInviteLink) (*models.ChatInviteLinkResponse, error) {
	var apiResponse *models.ChatInviteLinkResponse

	result, err := t.post("createChatInviteLink", params, &apiResponse)

	return result.(*models.ChatInviteLinkResponse), err
}

func (t *Telegram) ApproveChatJoinRequest(params *requests.Member) (*models.ApiResponse, error) {
	var apiResponse *models.ApiResponse

	result, err := t.post("approveChatJoinRequest", params, &apiResponse)

	return result.(*models.ApiResponse), err
}

func (t *Telegram) DeclineChatJoinRequest(params *requests.Member) (*models.ApiResponse, error) {
	var apiResponse *models.ApiResponse

	result, err := t.post("declineChatJoinRequest", params, &apiResponse)

	return result.(*models.ApiResponse), err
}

func (t *Telegram) GetFile(params *requests.File) (*models.FileResponse, error) {
	var apiResponse *models.FileResponse

	result, err := t.post("getFile", params, &apiResponse)

	return result.(*models.FileResponse), err
}

type TelegramApi struct {
	BaseApi   string        `json:"base_api"`
	Timeout   time.Duration `json:"timeout"`
	UserAgent string        `json:"user_agent"`
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
	}
}

func GetTelegramApi() *Telegram {
	return telegram
}
