package migration

// DiscoveredAccount represents an account found in an external IDE configuration.
type DiscoveredAccount struct {
	Name        string `json:"name"`         // e.g. "AndroidStudio2023.1"
	Path        string `json:"path"`         // Config Path
	RawToken    string `json:"raw_token"`    // Hidden from user in listing if needed, but required for import
	MaskedToken string `json:"masked_token"` // For display
	Source      string `json:"source"`       // "JetBrains", "AndroidStudio"
}

// LegacyStatus represents the status of legacy data.
type LegacyStatus struct {
	Found bool `json:"found"`
	Count int  `json:"count"`
}

// LegacyConfig represents the structure of the old agent's config file.
type LegacyConfig struct {
	Accounts []LegacyAccount `json:"accounts"`
}

// LegacyAccount represents an account in the legacy config.
type LegacyAccount struct {
	Email        string `json:"email"`
	RefreshToken string `json:"refresh_token"`
}
