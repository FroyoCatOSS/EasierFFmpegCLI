package parseconfigjson

import (
	"encoding/json"
	"fmt"
	"os"
)

type DownloadConfig struct {
	UseGithubMirror bool     `json:"UseGithubMirror"`
	GithubMirrors   []string `json:"GithubMirrors"`
}

func ParseJson(JsonConfigPath string) {
	var githubDownloadConfig DownloadConfig
	githubDownloadConfigFile, err := os.ReadFile(JsonConfigPath)
	if err != nil {
		fmt.Println("Error reading config file:", err)
		return
	}

	err = json.Unmarshal(githubDownloadConfigFile, &githubDownloadConfig)
	if err != nil {
		fmt.Println("Error parsing config file:", err)
		return
	}

	fmt.Println("Parsed config successfully:", githubDownloadConfig)
}
