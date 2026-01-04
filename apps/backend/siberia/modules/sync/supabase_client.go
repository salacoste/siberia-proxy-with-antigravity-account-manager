package sync

import (
	"context"
	"fmt"
	"time"

	"github.com/nedpals/supabase-go"
)

// CloudProfile matches the Supabase 'profiles' table structure
type CloudProfile struct {
	UserID    string    `json:"user_id"`
	DataBlob  string    `json:"data_blob"` // Encrypted hex string
	UpdatedAt time.Time `json:"updated_at"`
}

// SupabaseClient wraps the official client for our specific use case
type SupabaseClient struct {
	client *supabase.Client
	userID string
}

// NewSupabaseClient creates a client using credentials
func NewSupabaseClient(url, key string) *SupabaseClient {
	return &SupabaseClient{
		client: supabase.CreateClient(url, key),
		userID: "", // Starts unauthenticated
	}
}

// Push uploads the encrypted blob to Supabase
// For MVP, we use UPSERT (Insert or Update) based on UserID
func (s *SupabaseClient) Push(encryptedData string) error {
	if s.userID == "" {
		return fmt.Errorf("not logged in")
	}
	profile := CloudProfile{
		UserID:    s.userID,
		DataBlob:  encryptedData,
		UpdatedAt: time.Now(),
	}

	// Supabase Go client "DB" provides a query builder
	// Upsert into "profiles" table
	var results []CloudProfile
	err := s.client.DB.From("profiles").Upsert(profile).Execute(&results)
	if err != nil {
		return fmt.Errorf("supabase push error: %w", err)
	}

	return nil
}

// Pull downloads the encrypted blob from Supabase
func (s *SupabaseClient) Pull() (string, error) {
	if s.userID == "" {
		return "", fmt.Errorf("not logged in")
	}
	var results []CloudProfile
	// Select * from profiles where user_id = s.userID
	err := s.client.DB.From("profiles").Select("*").Eq("user_id", s.userID).Execute(&results)
	if err != nil {
		return "", fmt.Errorf("supabase pull error: %w", err)
	}

	if len(results) == 0 {
		return "", fmt.Errorf("no profile found for user %s", s.userID)
	}

	return results[0].DataBlob, nil
}

// SignUp registers a new user
func (s *SupabaseClient) SignUp(email, password string) error {
	ctx := context.Background()
	user, err := s.client.Auth.SignUp(ctx, supabase.UserCredentials{
		Email:    email,
		Password: password,
	})
	if err != nil {
		return err
	}
	s.userID = user.ID
	return nil
}

// SignIn logs in an existing user
func (s *SupabaseClient) SignIn(email, password string) error {
	ctx := context.Background()
	user, err := s.client.Auth.SignIn(ctx, supabase.UserCredentials{
		Email:    email,
		Password: password,
	})
	if err != nil {
		return err
	}
	s.userID = user.User.ID
	return nil
}

// GetUser returns the current user ID or empty
func (s *SupabaseClient) GetUser() string {
	return s.userID
}
