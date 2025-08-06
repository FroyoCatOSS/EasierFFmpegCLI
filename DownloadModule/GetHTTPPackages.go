package DownloadModule

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func DownloadFFPackage(FinalFFDownloadURL string) error {
	// 创建 HTTP 客户端
	fmt.Printf("Downloading from: %s\n", FinalFFDownloadURL)
	HttpClient := &http.Client{}
	request, err := http.NewRequest("GET", FinalFFDownloadURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %w", err)
	}
	request.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")

	response, err := HttpClient.Do(request)
	if err != nil {
		return fmt.Errorf("failed to download FFmpeg package: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download FFmpeg package: %s", response.Status)
	}
	urlPath := response.Request.URL.Path
	fileName := filepath.Base(urlPath)
	fileName = strings.TrimSuffix(fileName, " ")
	if fileName == "" || fileName == "." {
		fileName = "ffmpeg-download"
	}
	downloadDir := "downloads"
	if err := os.MkdirAll(downloadDir, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create download directory '%s': %w", downloadDir, err)
	}
	filePath := filepath.Join(downloadDir, fileName)
	outFile, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file '%s': %w", filePath, err)
	}
	defer outFile.Close()
	fmt.Printf("Saving to: %s\n", fileName)
	_, err = io.Copy(outFile, response.Body)
	if err != nil {
		return fmt.Errorf("failed to download FFmpeg package: %w", err)
	}
	err = outFile.Sync()
	if err != nil {
		return fmt.Errorf("failed to sync file '%s': %w", filePath, err)
	}
	fmt.Printf("Download completed: %s\n", filePath)
	return nil
}
