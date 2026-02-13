package examples

import (
	"errors"
	"testing"

	"github.com/ebingbo/ilicense-client-go/ilicense"
)

func TestActivateInvalidCode(t *testing.T) {
	cfg := ilicense.DefaultConfig()
	cfg.PublicKey = "invalid-public-key"
	cfg.StoragePath = t.TempDir() + "/license.dat"

	client := ilicense.NewClient(&cfg)
	_, err := client.Activate("")
	if err == nil {
		t.Fatalf("expected activation error, got nil")
	}
}

func TestCheckLicenseWithoutActivation(t *testing.T) {
	cfg := ilicense.DefaultConfig()
	client := ilicense.NewClient(&cfg)

	err := client.CheckLicense()
	if !errors.Is(err, ilicense.ErrLicenseNotFound) {
		t.Fatalf("expected ErrLicenseNotFound, got: %v", err)
	}
}

func TestCheckModuleWithoutActivation(t *testing.T) {
	cfg := ilicense.DefaultConfig()
	client := ilicense.NewClient(&cfg)

	err := client.CheckModule("m-a")
	if !errors.Is(err, ilicense.ErrLicenseNotFound) {
		t.Fatalf("expected ErrLicenseNotFound, got: %v", err)
	}
}
