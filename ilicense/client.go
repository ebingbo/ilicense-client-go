package ilicense

import (
	"errors"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	"ilicense-client-go/core"
)

type Client struct {
	config     *Config
	mu         sync.RWMutex
	licensePtr *core.License
}

func NewClient(config *Config) *Client {
	return &Client{
		config: config,
	}
}

func (m *Client) Init() error {
	if !m.config.Enabled {
		log.Printf("license validation disabled")
		return nil
	}
	if m.config.ValidateOnStartup {
		return m.performStartupValidation()
	}
	return nil
}

func (m *Client) performStartupValidation() error {
	if err := m.loadLicenseFromFile(); err != nil {
		log.Printf("license initialization failed: %v", err)
		if !m.config.AllowStartWhenExpired {
			return err
		}
		return nil
	}

	license := m.getCurrentLicense()
	if license == nil {
		return m.handleNoLicense()
	}
	if license.IsExpired(time.Now()) {
		return m.handleExpiredLicense()
	}
	m.handleValidLicense(*license)
	return nil
}

func (m *Client) handleNoLicense() error {
	log.Println("system not activated - please upload a license activation code to activate")
	if !m.config.AllowStartWhenExpired {
		return ErrLicenseNotFound
	}
	return nil
}

func (m *Client) handleExpiredLicense() error {
	if !m.config.AllowStartWhenExpired {
		return ErrLicenseExpired
	}
	return nil
}

func (m *Client) handleValidLicense(license core.License) {
	log.Printf("license validation successful - customer: %s, product: %s, expiry: %s, days left: %d",
		license.CustomerName,
		license.ProductName,
		license.ExpireAt.Format("2006-01-02 15:04:04"),
		license.DaysLeft,
	)
}

// CheckLicenseStatus matches Java scheduled check behavior.
func (m *Client) CheckLicenseStatus() {
	license := m.getCurrentLicense()
	if license == nil {
		log.Println("skipping check: not activated")
		return
	}

	if license.IsExpired(time.Now()) {
		log.Println("periodic check: license expired")
		return
	}
}

// Activate validates activation code and persists it.
func (m *Client) Activate(activationCode string) (*core.License, error) {
	log.Println("starting license activation")
	license, err := core.Validate(m.config.PublicKey, activationCode)
	if err != nil {
		return nil, err
	}

	if license.IsExpired(time.Now()) {
		return nil, ErrLicenseExpired
	}

	if err := m.saveLicenseToFile(activationCode); err != nil {
		return nil, err
	}

	m.setCurrentLicense(license)
	log.Printf("license activated successfully: %s", license.CustomerName)
	return license, nil
}

func (m *Client) GetCurrentLicense() *core.License {
	return m.getCurrentLicense()
}

func (m *Client) IsValid() bool {
	license := m.getCurrentLicense()
	return license != nil && !license.IsExpired(time.Now())
}

func (m *Client) HasModule(moduleName string) bool {
	license := m.getCurrentLicense()
	return license != nil && license.HasModule(moduleName)
}

func (m *Client) CheckLicense() error {
	license := m.getCurrentLicense()
	if license == nil {
		return ErrLicenseNotFound
	}
	if license.IsExpired(time.Now()) {
		return ErrLicenseExpired
	}
	return nil
}

func (m *Client) CheckModule(moduleName string) error {
	if err := m.CheckLicense(); err != nil {
		return err
	}
	license := m.getCurrentLicense()
	if license == nil || !license.HasModule(moduleName) {
		return &LicenseError{Msg: "unauthorized module: " + moduleName}
	}
	return nil
}

func (m *Client) loadLicenseFromFile() error {
	path := m.config.StoragePath
	if path == "" {
		return nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			log.Printf("license file does not exist: %s", path)
			return nil
		}
		return &LicenseError{Msg: "failed to load license file", Err: err}
	}

	license, err := core.Validate(m.config.PublicKey, string(data))
	if err != nil {
		return err
	}
	m.setCurrentLicense(license)
	log.Println("license loaded successfully from file")
	return nil
}

func (m *Client) saveLicenseToFile(activationCode string) error {
	path := m.config.StoragePath
	if path == "" {
		return &LicenseError{Msg: "failed to save license", Err: errors.New("storage path is empty")}
	}
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return &LicenseError{Msg: "failed to save license", Err: err}
	}
	if err := os.WriteFile(path, []byte(activationCode), 0o600); err != nil {
		return &LicenseError{Msg: "failed to save license", Err: err}
	}
	log.Printf("license saved: %s", path)
	return nil
}

func (m *Client) getCurrentLicense() *core.License {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if m.licensePtr == nil {
		return nil
	}
	clone := *m.licensePtr
	return &clone
}

func (m *Client) setCurrentLicense(license *core.License) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.licensePtr = license
}

func truncate(str string, max int) string {
	if len(str) <= max {
		return str
	}
	return str[:max-3] + "..."
}
