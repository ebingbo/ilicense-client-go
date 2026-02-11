package core

import (
	"strings"
	"time"
)

// License mirrors Java License fields.
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

func (l License) IsExpired(now time.Time) bool {
	if l.ExpireAt.IsZero() {
		return false
	}
	return l.ExpireAt.Before(now)
}

func (l License) HasModule(moduleName string) bool {
	return strings.Contains(l.Modules, moduleName)
}
