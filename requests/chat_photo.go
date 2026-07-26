package requests

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
)

type ChatPhoto struct {
	ChatId any    `json:"chat_id"`
	Photo  []byte `json:"photo"`
}

func (s *ChatPhoto) ToMultipart() (*bytes.Buffer, string, error) {
	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)

	part, err := writer.CreateFormFile("photo", "logo.png")
	if err != nil {
		return nil, "", fmt.Errorf("创建 photo 字段失败: %w", err)
	}
	if _, err := io.Copy(part, bytes.NewReader(s.Photo)); err != nil {
		return nil, "", fmt.Errorf("写入 photo 数据失败: %w", err)
	}

	if err := writer.WriteField("chat_id", fmt.Sprintf("%d", s.ChatId)); err != nil {
		return nil, "", fmt.Errorf("写入 chat_id 失败: %w", err)
	}

	if err := writer.Close(); err != nil {
		return nil, "", fmt.Errorf("关闭 writer 失败: %w", err)
	}

	return &requestBody, writer.FormDataContentType(), nil
}
