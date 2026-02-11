package ilicense

import (
	"ilicense-client-go/core"
)

// Config mirrors license properties from the Java SDK.
type Config struct {
	Enabled               bool   `json:"enabled"`
	PublicKey             string `json:"public_key"`
	StoragePath           string `json:"storage_path"`
	ValidateOnStartup     bool   `json:"validate_on_startup"`
	AllowStartWhenExpired bool   `json:"allow_start_when_expired"`
}

// DefaultConfig returns the Java-equivalent defaults.
func DefaultConfig() Config {
	return Config{
		Enabled:               true,
		StoragePath:           defaultStoragePath(),
		ValidateOnStartup:     false,
		AllowStartWhenExpired: true,
	}
}

func defaultStoragePath() string {
	home := core.UserHomeDir()
	if home == "" {
		return ".license/license.dat"
	}
	return home + "/.license/license.dat"
}
