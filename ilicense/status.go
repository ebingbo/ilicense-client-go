package ilicense

// LicenseStatus represents the current runtime status of a loaded license.
type LicenseStatus string

const (
	LicenseStatusValid        LicenseStatus = "valid"
	LicenseStatusExpired      LicenseStatus = "expired"
	LicenseStatusNotActivated LicenseStatus = "not_activated"
)
