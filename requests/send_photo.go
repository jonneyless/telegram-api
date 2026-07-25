package requests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"

	"github.com/jonneyless/telegram-api/models"
)

type SendPhoto struct {
	ChatId              int64               `json:"chat_id"`
	Photo               []byte              `json:"photo"`
	Caption             *string             `json:"caption,omitempty"`
	ParseMode           *string             `json:"parse_mode,omitempty"`
	ReplyMarkup         *models.ReplyMarkup `json:"reply_markup,omitempty"`
	Buttons             [][]map[string]any  `json:"buttons"`
	ButtonType          string              `json:"button_type"`
	DisableNotification bool                `json:"disable_notification"`
	LinkPreviewOptions  *LinkPreviewOptions `json:"link_preview_options,omitempty"`
	ReplyParameters     *ReplyParameters    `json:"reply_parameters,omitempty"`
}

func (s *SendPhoto) ToMultipart() (*bytes.Buffer, string, error) {
	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)

	part, err := writer.CreateFormFile("photo", "qrcode.png")
	if err != nil {
		return nil, "", fmt.Errorf("创建 photo 字段失败: %w", err)
	}
	if _, err := io.Copy(part, bytes.NewReader(s.Photo)); err != nil {
		return nil, "", fmt.Errorf("写入 photo 数据失败: %w", err)
	}

	if err := writer.WriteField("chat_id", fmt.Sprintf("%d", s.ChatId)); err != nil {
		return nil, "", fmt.Errorf("写入 chat_id 失败: %w", err)
	}

	if s.Caption != nil && *s.Caption != "" {
		if err := writer.WriteField("caption", *s.Caption); err != nil {
			return nil, "", fmt.Errorf("写入 caption 失败: %w", err)
		}
	}

	parseMode := "html"
	if s.ParseMode != nil && *s.ParseMode != "" {
		parseMode = *s.ParseMode
	}
	if err := writer.WriteField("parse_mode", parseMode); err != nil {
		return nil, "", fmt.Errorf("写入 parse_mode 失败: %w", err)
	}

	if s.DisableNotification {
		if err := writer.WriteField("disable_notification", "true"); err != nil {
			return nil, "", fmt.Errorf("写入 disable_notification 失败: %w", err)
		}
	}

	if s.ReplyMarkup != nil {
		replyMarkupJSON, err := json.Marshal(s.ReplyMarkup)
		if err != nil {
			return nil, "", fmt.Errorf("序列化 reply_markup 失败: %w", err)
		}
		if err := writer.WriteField("reply_markup", string(replyMarkupJSON)); err != nil {
			return nil, "", fmt.Errorf("写入 reply_markup 失败: %w", err)
		}
	} else if len(s.Buttons) > 0 {
		if err := s.writeButtonsAsReplyMarkup(writer); err != nil {
			return nil, "", err
		}
	}

	if s.LinkPreviewOptions != nil {
		linkPreviewJSON, err := json.Marshal(s.LinkPreviewOptions)
		if err != nil {
			return nil, "", fmt.Errorf("序列化 link_preview_options 失败: %w", err)
		}
		if err := writer.WriteField("link_preview_options", string(linkPreviewJSON)); err != nil {
			return nil, "", fmt.Errorf("写入 link_preview_options 失败: %w", err)
		}
	}

	if s.ReplyParameters != nil {
		replyParamsJSON, err := json.Marshal(s.ReplyParameters)
		if err != nil {
			return nil, "", fmt.Errorf("序列化 reply_parameters 失败: %w", err)
		}
		if err := writer.WriteField("reply_parameters", string(replyParamsJSON)); err != nil {
			return nil, "", fmt.Errorf("写入 reply_parameters 失败: %w", err)
		}
	}

	if err := writer.Close(); err != nil {
		return nil, "", fmt.Errorf("关闭 writer 失败: %w", err)
	}

	return &requestBody, writer.FormDataContentType(), nil
}

func (s *SendPhoto) writeButtonsAsReplyMarkup(writer *multipart.Writer) error {
	var keyboard [][]models.InlineKeyboardButton

	for _, row := range s.Buttons {
		var buttonRow []models.InlineKeyboardButton
		for _, btn := range row {
			button := models.InlineKeyboardButton{}

			if text, ok := btn["text"].(string); ok {
				button.Text = text
			}
			if callbackData, ok := btn["callback_data"].(string); ok {
				button.CallbackData = &callbackData
			}
			if url, ok := btn["url"].(string); ok {
				button.URL = &url
			}

			buttonRow = append(buttonRow, button)
		}
		keyboard = append(keyboard, buttonRow)
	}

	replyMarkup := models.ReplyMarkup{
		InlineKeyboard: keyboard,
	}

	replyMarkupJSON, err := json.Marshal(replyMarkup)
	if err != nil {
		return fmt.Errorf("序列化 buttons 到 reply_markup 失败: %w", err)
	}

	if err := writer.WriteField("reply_markup", string(replyMarkupJSON)); err != nil {
		return fmt.Errorf("写入 reply_markup 失败: %w", err)
	}

	return nil
}
