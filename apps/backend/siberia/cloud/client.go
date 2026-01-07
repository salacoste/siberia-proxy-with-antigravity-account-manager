package cloud

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type SupabaseClient struct {
	BaseURL string
	APIKey  string
	Client  *http.Client
}

func NewSupabaseClient(url, key string) *SupabaseClient {
	return &SupabaseClient{
		BaseURL: url,
		APIKey:  key,
		Client:  &http.Client{Timeout: 10 * time.Second},
	}
}

type AuthResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	User         struct {
		ID    string `json:"id"`
		Email string `json:"email"`
	} `json:"user"`
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

type ProfileData struct {
	UserID      string `json:"user_id"`
	Email       string `json:"email"`
	DataBlob    string `json:"data_blob"`
	IV          string `json:"iv"`
	UpdatedAt   string `json:"updated_at"`
	AccessToken string `json:"-"` // Helper
}

func (s *SupabaseClient) SignIn(email, password string) (*AuthResponse, error) {
	endpoint := fmt.Sprintf("%s/auth/v1/token?grant_type=password", s.BaseURL)
	payload := map[string]string{
		"email":    email,
		"password": password,
	}
	// Sign In returns body we need to parse. doAuthRequest helper logic needs to be inside doRequest or similar.
	// Actually doRequest returns *http.Response now. We need to parse it here.
	resp, err := s.doRequest("POST", endpoint, payload, "")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bytes, _ := io.ReadAll(resp.Body)
	var authResp AuthResponse
	if err := json.Unmarshal(bytes, &authResp); err != nil {
		return nil, err
	}
	return &authResp, nil
}

func (s *SupabaseClient) SignUp(email, password string) (*AuthResponse, error) {
	endpoint := fmt.Sprintf("%s/auth/v1/signup", s.BaseURL)
	payload := map[string]string{
		"email":    email,
		"password": password,
	}
	resp, err := s.doRequest("POST", endpoint, payload, "")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bytes, _ := io.ReadAll(resp.Body)
	var authResp AuthResponse
	if err := json.Unmarshal(bytes, &authResp); err != nil {
		return nil, err
	}
	return &authResp, nil
}

// DB Methods

func (s *SupabaseClient) UpsertProfile(userID string, data ProfileData) error {
	endpoint := fmt.Sprintf("%s/rest/v1/profiles", s.BaseURL)
	// Upsert logic usually requires Prefer: resolution=merge-duplicates header

	// Payload matching table schema
	payload := map[string]interface{}{
		"user_id":    userID,
		"email":      data.Email,
		"data_blob":  data.DataBlob,
		"iv":         data.IV,
		"updated_at": time.Now(),
	}

	_, err := s.doRequest("POST", endpoint, payload, data.AccessToken, map[string]string{
		"Prefer": "resolution=merge-duplicates",
	})
	return err
}

func (s *SupabaseClient) GetProfile(userID, accessToken string) (*ProfileData, error) {
	endpoint := fmt.Sprintf("%s/rest/v1/profiles?user_id=eq.%s&select=*", s.BaseURL, userID)

	// Response is a list
	resp, err := s.doAuthRequest("GET", endpoint, nil, accessToken)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bytes, _ := io.ReadAll(resp.Body)
	var profiles []ProfileData
	if err := json.Unmarshal(bytes, &profiles); err != nil {
		return nil, err
	}

	if len(profiles) == 0 {
		return nil, fmt.Errorf("profile not found")
	}
	return &profiles[0], nil
}

// Refactored to support custom headers
func (s *SupabaseClient) doRequest(method, url string, body interface{}, authToken string, headers ...map[string]string) (*http.Response, error) {
	var bodyReader io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		bodyReader = bytes.NewBuffer(jsonBody)
	}

	req, err := http.NewRequest(method, url, bodyReader)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("apikey", s.APIKey)
	if authToken != "" {
		req.Header.Set("Authorization", "Bearer "+authToken)
	}

	for _, h := range headers {
		for k, v := range h {
			req.Header.Set(k, v)
		}
	}

	resp, err := s.Client.Do(req)
	if err != nil {
		return nil, err
	}

	// We read body for error checking in high level methods, but here we return raw response for GET/Reading?
	// Actually original doRequest read the body. Let's fix that inconsistency.
	// For Upsert (POST), we assume 200/201/204.
	if resp.StatusCode >= 400 {
		defer resp.Body.Close()
		b, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("api error %d: %s", resp.StatusCode, string(b))
	}

	return resp, nil
}

// doAuthRequest wrapper
func (s *SupabaseClient) doAuthRequest(method, url string, body interface{}, authToken string) (*http.Response, error) {
	return s.doRequest(method, url, body, authToken)
}
