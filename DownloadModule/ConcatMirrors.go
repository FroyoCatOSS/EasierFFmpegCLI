package DownloadModule

import (
	"strings"
)

func ApplyGitHubMirror(FFDownloadURL string, GithubMirror string) string {
	FFDownloadURL = strings.TrimSpace(FFDownloadURL)
	if !strings.HasPrefix(GithubMirror, "http://") && !strings.HasPrefix(GithubMirror, "https://") {
		GithubMirror = "https://" + GithubMirror
	}
	return GithubMirror + "/https://" + strings.TrimPrefix(FFDownloadURL, "https://")
}

func ApplyMirrorsToURLS(FFDownloadURLs []string, GithubMirrors []string) []string {
	var mirroredURLS []string
	for _, url := range FFDownloadURLs {
		url = strings.TrimSpace(url)
		if strings.Contains(url, "github.com") {
			mirroredURL := ApplyGitHubMirror(url, GithubMirrors[0])
			mirroredURLS = append(mirroredURLS, mirroredURL)
		} else {
			mirroredURLS = append(mirroredURLS, url)
		}
	}
	return mirroredURLS
}
