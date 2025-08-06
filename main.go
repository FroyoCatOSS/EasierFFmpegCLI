package main

import (
	"EasierFFmpegCLI/DownloadModule"
	"EasierFFmpegCLI/ParseConfigJson"
	"fmt"
	"os"
	"strings"
)

// isGitHubURL 判断是否为 GitHub 下载链接
func isGitHubURL(url string) bool {
	url = strings.TrimSpace(url)
	return strings.Contains(url, "github.com")
}

func main() {
	fmt.Println("Easier FFmpeg CLI is Now Running...")

	configPath := "configs/DownloadConfig.json"

	// 检查配置文件是否存在
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		fmt.Printf("Error: Config file '%s' not found.\n", configPath)
		return
	}

	// 获取 GitHub 镜像列表
	mirrors, err := ParseConfigJson.GetGitHubMirrors(configPath)
	if err != nil {
		fmt.Println("Warning: Failed to load GitHub mirrors:", err)
	} else if len(mirrors) > 0 {
		fmt.Println("GitHub Mirrors Enabled:", mirrors)
	} else {
		fmt.Println("GitHub Mirrors are configured but not enabled or empty.")
	}

	// 获取当前平台和许可证匹配的原始下载链接
	rawURLs, err := ParseConfigJson.GetFFDownloadUrls(configPath)
	if err != nil {
		fmt.Println("Error getting FFmpeg download URLs:", err)
		return
	}

	fmt.Println("\nRaw Download URLs:")
	for i, url := range rawURLs {
		host := "Other"
		if isGitHubURL(url) {
			host = "GitHub"
		}
		fmt.Printf("   %d. [%s] %s\n", i+1, host, url)
	}

	// 决定最终使用的链接（应用镜像或直连）
	var finalURLs []string
	useMirror := len(mirrors) > 0

	if useMirror {
		fmt.Println("\nApplying GitHub Mirror...")
		// 使用 DownloadModule 中的函数批量处理
		finalURLs = DownloadModule.ApplyMirrorsToURLS(rawURLs, mirrors)

		fmt.Println("   Mirror Host:", mirrors[0])
		for i, original := range rawURLs {
			if isGitHubURL(original) {
				fmt.Printf("%s\n%s\n", original, finalURLs[i])
			} else {
				fmt.Printf("(Direct) %s\n", original)
			}
		}
	} else {
		finalURLs = rawURLs
		fmt.Println("\nMirror is disabled. Using direct download URLs.")
	}

	// 输出最终下载计划
	fmt.Println("\nFinal Download Plan:")
	for _, url := range finalURLs {
		fmt.Printf("%s\n", url)
	}
	fmt.Println("\nDownload preparation complete. Ready to start downloading.")
	// 下载 FFmpeg 包
	for _, url := range finalURLs {
		err := DownloadModule.DownloadFFPackage(url)
		if err != nil {
			fmt.Printf("Error downloading '%s': %v\n", url, err)
		}
	}
}
