package updater

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCheckForUpdates_Available(t *testing.T) {
	// Mock GitHub API
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/repos/test/repo/releases/latest" {
			t.Errorf("Expected path /repos/test/repo/releases/latest, got %s", r.URL.Path)
		}

		response := GitHubRelease{
			TagName: "v1.3.0",
			Body:    "New features",
			HtmlUrl: "http://github.com/test/repo/releases/v1.3.0",
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer ts.Close()

	svc := &Service{
		CurrentVersion: "v1.2.0",
		RepoOwner:      "test",
		RepoName:       "repo",
		APIURL:         ts.URL,
	}

	info := svc.CheckForUpdates()

	if !info.Available {
		t.Errorf("Expected update available, got false")
	}
	if info.LatestVersion != "v1.3.0" {
		t.Errorf("Expected v1.3.0, got %s", info.LatestVersion)
	}
	if info.DownloadURL != "http://github.com/test/repo/releases/v1.3.0" {
		t.Errorf("Expected correct download URL, got %s", info.DownloadURL)
	}
}

func TestCheckForUpdates_NoUpdate(t *testing.T) {
	// Mock GitHub API
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := GitHubRelease{
			TagName: "v1.2.0",
			Body:    "Fixes",
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer ts.Close()

	svc := &Service{
		CurrentVersion: "v1.2.0",
		RepoOwner:      "test",
		RepoName:       "repo",
		APIURL:         ts.URL,
	}

	info := svc.CheckForUpdates()

	if info.Available {
		t.Errorf("Expected no update available, got true")
	}
}
