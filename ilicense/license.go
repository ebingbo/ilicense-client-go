package ilicense

import (
	"strings"
	"time"
)

// License is the public license model exposed by the SDK.
type License struct {
	LicenseCode  string    `json:"license_code"`
	CustomerCode string    `json:"customer_code"`
	CustomerName string    `json:"customer_name"`
	ProductCode  string    `json:"product_code"`
	ProductName  string    `json:"product_name"`
	IssuerCode   string    `json:"issuer_code"`
	IssuerName   string    `json:"issuer_name"`
	IssueAt      time.Time `json:"issue_at"`
	ExpireAt     time.Time `json:"expire_at"`
	Modules      string    `json:"modules"`
	MaxInstances int       `json:"max_instances"`

	Valid    bool  `json:"valid"`
	DaysLeft int64 `json:"days_left"`
}

// IsExpired reports whether ExpireAt is before the given time.
func (l *License) IsExpired(now time.Time) bool {
	if l.ExpireAt.IsZero() {
		return false
	}
	return l.ExpireAt.Before(now)
}

// HasModule reports whether Modules contains an exact module token.
func (l *License) HasModule(moduleName string) bool {
	moduleName = strings.TrimSpace(moduleName)
	if moduleName == "" {
		return false
	}
	for _, m := range strings.Split(l.Modules, ",") {
		if strings.TrimSpace(m) == moduleName {
			return true
		}
	}
	return false
}
