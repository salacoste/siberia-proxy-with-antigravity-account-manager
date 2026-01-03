package updater

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// UpdateInfo holds information about the latest available release
type UpdateInfo struct {
	Available      bool   `json:"available"`
	CurrentVersion string `json:"current_version"`
	LatestVersion  string `json:"latest_version"`
	DownloadURL    string `json:"download_url"`
	ReleaseNotes   string `json:"release_notes"`
	Error          string `json:"error,omitempty"`
}

// Service handles update checks
type Service struct {
	CurrentVersion string
	RepoOwner      string
	RepoName       string
}

// NewService creates a new updater service
func NewService(currentVersion string) *Service {
	// Default to local repo for now
	return &Service{
		CurrentVersion: currentVersion,
		RepoOwner:      "salacoste",
		RepoName:       "siberia",
	}
}

// GitHubRelease represents the minimal structure of the GitHub Release API response
type GitHubRelease struct {
	TagName string `json:"tag_name"`
	Body    string `json:"body"`
	HtmlUrl string `json:"html_url"`
}

// CheckForUpdates queries the GitHub API for the latest release
func (s *Service) CheckForUpdates() UpdateInfo {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", s.RepoOwner, s.RepoName)

	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return UpdateInfo{Error: "Failed to create request: " + err.Error()}
	}

	resp, err := client.Do(req)
	if err != nil {
		return UpdateInfo{Error: "Network error: " + err.Error()}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return UpdateInfo{Error: fmt.Sprintf("GitHub API returned status: %d", resp.StatusCode)}
	}

	var release GitHubRelease
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return UpdateInfo{Error: "Failed to parse release info"}
	}

	latest := strings.TrimPrefix(release.TagName, "v")
	current := strings.TrimPrefix(s.CurrentVersion, "v")

	// Simple string comparison for MVP.
	// In production, use "github.com/Masterminds/semver/v3"
	// But since this is a desktop app, let's keep dependencies low if possible or just do simple check.
	// Actually, just != check is enough to say "Update Available" if we assume strictly increasing tags.
	updateAvailable := latest != current

	return UpdateInfo{
		Available:      updateAvailable,
		CurrentVersion: s.CurrentVersion,
		LatestVersion:  release.TagName,
		DownloadURL:    release.HtmlUrl, // Point to release page for manual download (MVP Safety)
		ReleaseNotes:   release.Body,
	}
}
