package updater

import (
	"encoding/json"
	"fmt"
	"net/http"
	"runtime"
	"strings"
	"time"

	"github.com/hashicorp/go-version"
)

const (
	RepoOwner = "salacoste"
	RepoName  = "siberia-proxy-with-antigravity-account-manager"
)

type UpdateInfo struct {
	Available   bool   `json:"available"`
	Version     string `json:"version"`
	ReleaseURL  string `json:"release_url"`
	DownloadURL string `json:"download_url"`
	Description string `json:"description"`
}

type Service struct {
	currentVersion string
}

func NewService(currentVersion string) *Service {
	return &Service{
		currentVersion: currentVersion,
	}
}

type githubRelease struct {
	TagName string `json:"tag_name"`
	HtmlURL string `json:"html_url"`
	Body    string `json:"body"`
	Assets  []struct {
		Name               string `json:"name"`
		BrowserDownloadURL string `json:"browser_download_url"`
	} `json:"assets"`
}

func (s *Service) CheckForUpdates() (*UpdateInfo, error) {
	client := &http.Client{Timeout: 10 * time.Second}
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", RepoOwner, RepoName)

	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to check updates: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("github api returned status: %d", resp.StatusCode)
	}

	var release githubRelease
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	// Compare Versions
	vCurrent, err := version.NewVersion(s.currentVersion)
	if err != nil {
		return nil, fmt.Errorf("invalid current version: %v", err)
	}
	vRemote, err := version.NewVersion(release.TagName)
	if err != nil {
		return nil, fmt.Errorf("invalid remote version: %v", err)
	}

	info := &UpdateInfo{
		Available:   vRemote.GreaterThan(vCurrent),
		Version:     release.TagName,
		ReleaseURL:  release.HtmlURL,
		Description: release.Body,
	}

	// Find asset for current OS
	os := runtime.GOOS
	// Mapping: darwin -> macos, windows -> windows, linux -> linux
	searchStr := ""
	switch os {
	case "darwin":
		searchStr = "macos"
	case "windows":
		searchStr = "windows"
	case "linux":
		searchStr = "linux"
	}

	for _, asset := range release.Assets {
		// extremely simple matching logic
		// if searchStr is in name
		if searchStr != "" && len(asset.BrowserDownloadURL) > 0 {
			// Check if filename contains our OS key
			// Note: This is a bit naive but works for our naming convention "siberia-linux-..."
			if strings.Contains(asset.Name, searchStr) {
				info.DownloadURL = asset.BrowserDownloadURL
				break
			}
		}
	}
	// Fallback if no asset match, use release URL
	if info.DownloadURL == "" {
		info.DownloadURL = release.HtmlURL
	}

	return info, nil
}
