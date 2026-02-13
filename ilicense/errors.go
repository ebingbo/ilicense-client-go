package ilicense

import (
	"errors"

	core "github.com/ebingbo/ilicense-client-go/internal/core"
)

var (
	// ErrLicenseNotFound means the system has not been activated.
	ErrLicenseNotFound = errors.New("system not activated")
	// ErrLicenseExpired means the currently loaded license is expired.
	ErrLicenseExpired = errors.New("license expired")
	// ErrSignatureInvalid means activation code signature verification failed.
	ErrSignatureInvalid = core.ErrSignatureInvalid
)

// LicenseError wraps a low-level error with context.
type LicenseError struct {
	Msg string
	Err error
}

func (e *LicenseError) Error() string {
	if e.Err == nil {
		return e.Msg
	}
	return e.Msg + ": " + e.Err.Error()
}

func (e *LicenseError) Unwrap() error { return e.Err }
