package migration

// LegacyAccount represents the structure of an account in the Python/Bash agent config
type LegacyAccount struct {
	Email        string `json:"email"`
	RefreshToken string `json:"refresh_token"`
	// We ignore other fields like access_token as they are likely expired
}

// LegacyConfig represents the root JSON object of the legacy config file
type LegacyConfig struct {
	Accounts []LegacyAccount `json:"accounts"`
}

type LegacyStatus struct {
	Found bool `json:"found"`
	Count int  `json:"count"`
}
