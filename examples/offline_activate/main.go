package main

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/ebingbo/ilicense-client-go/ilicense"
)

func main() {
	cfg := ilicense.DefaultConfig()
	cfg.PublicKey = os.Getenv("ILICENSE_PUBLIC_KEY")
	cfg.ValidateOnStartup = true
	cfg.AllowStartWhenExpired = true

	if cfg.PublicKey == "" {
		log.Fatal("set ILICENSE_PUBLIC_KEY")
	}

	client := ilicense.NewClient(&cfg)
	if err := client.Init(); err != nil {
		log.Fatalf("init failed: %v", err)
	}

	if err := client.CheckModule("m-a"); err == nil {
		fmt.Println("module m-a already granted")
		return
	} else if !errors.Is(err, ilicense.ErrLicenseNotFound) && !errors.Is(err, ilicense.ErrLicenseExpired) {
		log.Fatalf("unexpected check error: %v", err)
	}

	code := os.Getenv("ILICENSE_ACTIVATION_CODE")
	if code == "" {
		log.Fatal("license missing or expired; set ILICENSE_ACTIVATION_CODE to activate")
	}

	if _, err := client.Activate(code); err != nil {
		log.Fatalf("activate failed: %v", err)
	}

	if err := client.CheckModule("m-a"); err != nil {
		log.Fatalf("module m-a still unauthorized: %v", err)
	}
	fmt.Println("module m-a authorized")
}
