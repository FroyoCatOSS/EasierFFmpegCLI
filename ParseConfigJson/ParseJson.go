package ParseConfigJson

import (
	"encoding/json"
	"fmt"
	"os"
	"runtime"
	"strings"
)

// var configPath := "configs/DownloadConfig.json"
var configPath = "configs/DownloadConfig.json"

// DownloadConfig represents the configuration for GitHub mirror usage.
type DownloadConfig struct {
	UseProxy        bool                `json:"UseProxy"`
	ProxyURL        string              `json:"ProxyURL"`
	UseGithubMirror bool                `json:"UseGithubMirror"`
	GithubMirrors   []string            `json:"GithubMirrors"`
	License         []string            `json:"License"`
	DownloadURLS    PlatformDownloadMap `json:"DownloadURLS"`
}

type PlatformDownloadMap map[string][]map[string]string

// type PlatformLinksMap map[string][]map[string]string

// ParseConfig reads and parses the JSON configuration file from the given path.
// Returns the parsed config and an error if any.
func ParseConfig(configPath string) (*DownloadConfig, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file '%s': %w", configPath, err)
	}

	var config DownloadConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse JSON from config file '%s': %w", configPath, err)
	}

	//验证配置合法性
	if config.UseGithubMirror && len(config.GithubMirrors) == 0 {
		return nil, fmt.Errorf("useGithubMirror is true, but githubMirrors list is empty")
	}
	for platform, items := range config.DownloadURLS {
		for i, item := range items {
			for name, url := range item {
				if strings.TrimSpace(url) != url {
					// 修改原 map
					delete(item, name)
					item[name] = strings.TrimSpace(url)
				}
			}
			config.DownloadURLS[platform][i] = item
		}
	}
	// fmt.Printf("Successfully parsed config: %+v\n", config)
	return &config, nil
}

// GetGitHubMirrors reads the config and returns the mirror URLs if enabled.
// Returns nil slice if not enabled or no mirrors are configured.
func GetGitHubMirrors(configPath string) ([]string, error) {
	config, err := ParseConfig(configPath)
	if err != nil {
		return nil, err
	}

	if config.UseGithubMirror && len(config.GithubMirrors) > 0 {
		return config.GithubMirrors, nil
	}

	return nil, nil
}

func GetRunningPlatform() string {
	switch runtime.GOOS {
	case "windows":
		return "windows"
	case "linux":
		return "linux"
		/*
			case "darwin":
				return "Mac"
			default:
				return "Unknown"
			}
		*/
	default:
		return runtime.GOOS
	}
}
func matchesLicense(name string, allowedLicenses []string) bool {
	if len(allowedLicenses) == 0 {
		return false // 没有允许的许可证，则不匹配
	}
	name = strings.ToUpper(name) // 统一转大写比较
	for _, lic := range allowedLicenses {
		if strings.Contains(name, strings.ToUpper(lic)) {
			return true
		}
	}
	return false
}
func GetFFDownloadUrls(configPath string) ([]string, error) {
	config, err := ParseConfig(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	currentPlatform := GetRunningPlatform()
	/*
		if currentPlatform == "Unknown" {
			return nil, fmt.Errorf("unsupported platform: %s", runtime.GOOS)
		}
	*/
	// 获取当前平台的下载项
	downloadItems, ok := config.DownloadURLS[currentPlatform]
	if !ok || len(downloadItems) == 0 {
		return nil, fmt.Errorf("no download URLs found for platform: %s", currentPlatform)
	}

	var filteredURLs []string

	// 遍历当前平台的所有下载项
	for _, item := range downloadItems {
		for name, url := range item {
			trimmedURL := strings.TrimSpace(url)
			if trimmedURL == "" {
				continue
			}

			// 检查 name 是否匹配任意一个 License 类型
			if matchesLicense(name, config.License) {
				filteredURLs = append(filteredURLs, trimmedURL)
				fmt.Printf("  Matched %s for license %v: %s\n", name, config.License, trimmedURL)
			}
		}
	}

	if len(filteredURLs) == 0 {
		return nil, fmt.Errorf("no matching download URLs found for licenses %v on platform %s", config.License, currentPlatform)
	}

	return filteredURLs, nil
}
