package cloud

import "time"

// CloudProfile represents the encrypted blob stored in Supabase
type CloudProfile struct {
	UserID    string    `json:"user_id"`
	Email     string    `json:"email"`
	UpdatedAt time.Time `json:"updated_at"`

	// Encrypted Data (Client-Side Encrypted JSON)
	// Contains: AppConfig, Accounts
	DataBlob string `json:"data_blob"`
	IV       string `json:"iv"` // Initialization Vector for AES-GCM
}

// SyncState tracks local vs remote state
type SyncState struct {
	LastSync time.Time `json:"last_sync"`
	IsDirty  bool      `json:"is_dirty"` // True if local changes exist since last sync
}

// LoginRequest used for Wails
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// RegisterRequest used for Wails
type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
