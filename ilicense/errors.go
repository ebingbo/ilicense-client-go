package ilicense

import (
	"errors"

	licensing "github.com/ebingbo/ilicense-client-go/internal/licensing"
)

var (
	// ErrLicenseNotFound means the system has not been activated.
	ErrLicenseNotFound = errors.New("system not activated")
	// ErrLicenseExpired means the currently loaded license is expired.
	ErrLicenseExpired = errors.New("license expired")
	// ErrModuleUnauthorized means current license does not grant a module.
	ErrModuleUnauthorized = errors.New("unauthorized module")
	// ErrSignatureInvalid means activation code signature verification failed.
	ErrSignatureInvalid = licensing.ErrSignatureInvalid
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

// ModuleUnauthorizedError contains the unauthorized module name.
type ModuleUnauthorizedError struct {
	Module string
}

func (e *ModuleUnauthorizedError) Error() string {
	if e.Module == "" {
		return ErrModuleUnauthorized.Error()
	}
	return ErrModuleUnauthorized.Error() + ": " + e.Module
}

func (e *ModuleUnauthorizedError) Unwrap() error { return ErrModuleUnauthorized }
