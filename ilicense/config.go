package ilicense

import "os"

// Logger defines optional logging hook for SDK runtime messages.
type Logger interface {
	Printf(format string, v ...any)
	Println(v ...any)
}

// Config mirrors license properties from the Java SDK.
type Config struct {
	Enabled               bool   `json:"enabled"`
	PublicKey             string `json:"public_key"`
	StoragePath           string `json:"storage_path"`
	ValidateOnStartup     bool   `json:"validate_on_startup"`
	AllowStartWhenExpired bool   `json:"allow_start_when_expired"`
	Logger                Logger `json:"-"`
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
	home, err := os.UserHomeDir()
	if err != nil {
		return ".license/license.dat"
	}
	if home == "" {
		return ".license/license.dat"
	}
	return home + "/.license/license.dat"
}
