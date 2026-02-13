package main

import (
	"fmt"
	"log"
	"os"

	"github.com/ebingbo/ilicense-client-go/ilicense"
)

func main() {
	publicKey := os.Getenv("ILICENSE_PUBLIC_KEY")
	activationCode := os.Getenv("ILICENSE_ACTIVATION_CODE")
	if publicKey == "" || activationCode == "" {
		log.Fatal("set ILICENSE_PUBLIC_KEY and ILICENSE_ACTIVATION_CODE")
	}

	cfg := ilicense.DefaultConfig()
	cfg.PublicKey = publicKey
	cfg.ValidateOnStartup = true
	cfg.AllowStartWhenExpired = false

	client := ilicense.NewClient(&cfg)
	if err := client.Init(); err != nil {
		log.Fatalf("init failed: %v", err)
	}

	lic, err := client.Activate(activationCode)
	if err != nil {
		log.Fatalf("activate failed: %v", err)
	}

	fmt.Printf("activated customer=%s product=%s expireAt=%s\n", lic.CustomerName, lic.ProductName, lic.ExpireAt.Format("2006-01-02"))
}
