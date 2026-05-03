package models

import (
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"time"
)

// File struct
type File struct {
	FileId       string `json:"file_id"`
	FileUniqueId string `json:"file_unique_id"`
	FileSize     int    `json:"file_size"`
	FilePath     string `json:"file_path"`
}

func (file *File) GetFileName() string {
	return fmt.Sprintf("%s%s", file.FileUniqueId, file.GetExt())
}

func (file *File) GetExt() string {
	return filepath.Ext(file.FilePath)
}

// Download 下载文件
func (file *File) Download(token string) ([]byte, error) {
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	fileUrl := fmt.Sprintf("https://api.telegram.org/file/bot%s/%s", token, file.FilePath)
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

	// 4. 读取数据
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read image data: %w", err)
	}

	return data, nil
}
