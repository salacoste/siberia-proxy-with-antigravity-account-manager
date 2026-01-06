package updater

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/Masterminds/semver/v3"
)

type UpdateInfo struct {
	Available      bool   `json:"available"`
	CurrentVersion string `json:"current_version"`
	LatestVersion  string `json:"latest_version"`
	DownloadURL    string `json:"download_url"`
	ReleaseNotes   string `json:"release_notes"`
	Error          string `json:"error,omitempty"`
}

type UpdateService struct {
	currentVersion string
	githubRepo     string // format: "owner/repo"
	apiURL         string // Base URL for GitHub API (can be overridden for tests)
	client         *http.Client
}

func NewUpdateService(currentVersion, githubRepo string) *UpdateService {
	// Normalize version
	v := strings.TrimPrefix(currentVersion, "v")

	return &UpdateService{
		currentVersion: v,
		githubRepo:     githubRepo,
		apiURL:         "https://api.github.com",
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// CheckForUpdates queries GitHub for the latest release
func (s *UpdateService) CheckForUpdates() (*UpdateInfo, error) {
	url := fmt.Sprintf("%s/repos/%s/releases/latest", s.apiURL, s.githubRepo)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	// Important: Add User-Agent for GitHub API
	req.Header.Set("User-Agent", "Siberia-Proxy-Updater")

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to check for updates: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("github api returned status: %d", resp.StatusCode)
	}

	var release struct {
		TagName string `json:"tag_name"`
		Body    string `json:"body"`
		HTMLURL string `json:"html_url"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return nil, fmt.Errorf("failed to parse github response: %w", err)
	}

	latestVerStr := strings.TrimPrefix(release.TagName, "v")

	currentV, err := semver.NewVersion(s.currentVersion)
	if err != nil {
		// If current version is dev/invalid, assume update is available if tag exists
		return &UpdateInfo{
			Available:      true,
			CurrentVersion: s.currentVersion,
			LatestVersion:  release.TagName,
			DownloadURL:    release.HTMLURL,
			ReleaseNotes:   release.Body,
			Error:          "Current version invalid (dev build?)",
		}, nil
	}

	latestV, err := semver.NewVersion(latestVerStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse remote version: %w", err)
	}

	if latestV.GreaterThan(currentV) {
		return &UpdateInfo{
			Available:      true,
			CurrentVersion: s.currentVersion,
			LatestVersion:  release.TagName,
			DownloadURL:    release.HTMLURL,
			ReleaseNotes:   release.Body,
		}, nil
	}

	return &UpdateInfo{
		Available:      false,
		CurrentVersion: s.currentVersion,
		LatestVersion:  release.TagName,
	}, nil
}
