package ilicense

import (
	"errors"
	"testing"
	"time"
)

func TestNewClientCopiesConfig(t *testing.T) {
	cfg := DefaultConfig()
	cfg.Enabled = true

	client := NewClient(&cfg)
	cfg.Enabled = false

	if !client.config.Enabled {
		t.Fatalf("expected client config to be independent copy")
	}
}

func TestCheckLicenseStatusNotActivated(t *testing.T) {
	client := NewClient(nil)

	status, err := client.CheckLicenseStatus()
	if status != LicenseStatusNotActivated {
		t.Fatalf("expected status %q, got %q", LicenseStatusNotActivated, status)
	}
	if !errors.Is(err, ErrLicenseNotFound) {
		t.Fatalf("expected ErrLicenseNotFound, got %v", err)
	}
}

func TestCheckLicenseStatusExpired(t *testing.T) {
	client := NewClient(nil)
	client.setCurrentLicense(&License{ExpireAt: time.Now().Add(-time.Minute)})

	status, err := client.CheckLicenseStatus()
	if status != LicenseStatusExpired {
		t.Fatalf("expected status %q, got %q", LicenseStatusExpired, status)
	}
	if !errors.Is(err, ErrLicenseExpired) {
		t.Fatalf("expected ErrLicenseExpired, got %v", err)
	}
}

func TestCheckLicenseStatusValid(t *testing.T) {
	client := NewClient(nil)
	client.setCurrentLicense(&License{ExpireAt: time.Now().Add(time.Hour)})

	status, err := client.CheckLicenseStatus()
	if status != LicenseStatusValid {
		t.Fatalf("expected status %q, got %q", LicenseStatusValid, status)
	}
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
}

func TestCheckModuleUnauthorizedError(t *testing.T) {
	client := NewClient(nil)
	client.setCurrentLicense(&License{ExpireAt: time.Now().Add(time.Hour), Modules: "m-a,m-b"})

	err := client.CheckModule("m-c")
	if !errors.Is(err, ErrModuleUnauthorized) {
		t.Fatalf("expected ErrModuleUnauthorized, got %v", err)
	}

	var moduleErr *ModuleUnauthorizedError
	if !errors.As(err, &moduleErr) {
		t.Fatalf("expected ModuleUnauthorizedError, got %T", err)
	}
	if moduleErr.Module != "m-c" {
		t.Fatalf("expected module m-c, got %s", moduleErr.Module)
	}
}
