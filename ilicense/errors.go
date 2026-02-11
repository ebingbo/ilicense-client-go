package ilicense

import "errors"

var (
	ErrLicenseNotFound  = errors.New("system not activated")
	ErrLicenseExpired   = errors.New("license expired")
	ErrSignatureInvalid = errors.New("signature verification failed")
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
