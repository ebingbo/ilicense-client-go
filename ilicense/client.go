package ilicense

import (
	"errors"
	"os"
	"path/filepath"
	"sync"
	"time"

	core "github.com/ebingbo/ilicense-client-go/internal/core"
)

type Client struct {
	config     *Config
	mu         sync.RWMutex
	licensePtr *License
}

// NewClient creates a license client with the provided config.
// If config is nil, DefaultConfig() is applied.
func NewClient(config *Config) *Client {
	if config == nil {
		c := DefaultConfig()
		config = &c
	}
	return &Client{
		config: config,
	}
}

// Init performs startup checks based on config flags.
func (m *Client) Init() error {
	if !m.config.Enabled {
		m.logf("license validation disabled")
		return nil
	}
	if m.config.ValidateOnStartup {
		return m.performStartupValidation()
	}
	return nil
}

func (m *Client) performStartupValidation() error {
	if err := m.loadLicenseFromFile(); err != nil {
		m.logf("license initialization failed: %v", err)
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
	m.logln("system not activated - please upload a license activation code to activate")
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

func (m *Client) handleValidLicense(license License) {
	m.logf("license validation successful - customer: %s, product: %s, expiry: %s, days left: %d",
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
		m.logln("skipping check: not activated")
		return
	}

	if license.IsExpired(time.Now()) {
		m.logln("periodic check: license expired")
		return
	}
}

// Activate validates activation code and persists it.
func (m *Client) Activate(activationCode string) (*License, error) {
	m.logln("starting license activation")
	raw, err := core.Validate(m.config.PublicKey, activationCode)
	if err != nil {
		return nil, err
	}
	license := fromCoreLicense(raw)

	if license.IsExpired(time.Now()) {
		return nil, ErrLicenseExpired
	}

	if err := m.saveLicenseToFile(activationCode); err != nil {
		return nil, err
	}

	m.setCurrentLicense(license)
	m.logf("license activated successfully: %s", license.CustomerName)
	return license, nil
}

// GetCurrentLicense returns a snapshot of the currently loaded license.
func (m *Client) GetCurrentLicense() *License {
	return m.getCurrentLicense()
}

// IsValid reports whether a non-expired license is currently loaded.
func (m *Client) IsValid() bool {
	license := m.getCurrentLicense()
	return license != nil && !license.IsExpired(time.Now())
}

// HasModule reports whether the loaded license grants the given module.
func (m *Client) HasModule(moduleName string) bool {
	license := m.getCurrentLicense()
	return license != nil && license.HasModule(moduleName)
}

// CheckLicense validates that a non-expired license is loaded.
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

// CheckModule validates both license validity and module authorization.
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
			m.logf("license file does not exist: %s", path)
			return nil
		}
		return &LicenseError{Msg: "failed to load license file", Err: err}
	}

	raw, err := core.Validate(m.config.PublicKey, string(data))
	if err != nil {
		return err
	}
	license := fromCoreLicense(raw)
	m.setCurrentLicense(license)
	m.logln("license loaded successfully from file")
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
	m.logf("license saved: %s", path)
	return nil
}

func (m *Client) getCurrentLicense() *License {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if m.licensePtr == nil {
		return nil
	}
	clone := *m.licensePtr
	return &clone
}

func (m *Client) setCurrentLicense(license *License) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.licensePtr = license
}

func fromCoreLicense(in *core.License) *License {
	if in == nil {
		return nil
	}
	return &License{
		LicenseCode:  in.LicenseCode,
		CustomerCode: in.CustomerCode,
		CustomerName: in.CustomerName,
		ProductCode:  in.ProductCode,
		ProductName:  in.ProductName,
		IssuerCode:   in.IssuerCode,
		IssuerName:   in.IssuerName,
		IssueAt:      in.IssueAt,
		ExpireAt:     in.ExpireAt,
		Modules:      in.Modules,
		MaxInstances: in.MaxInstances,
		Valid:        in.Valid,
		DaysLeft:     in.DaysLeft,
	}
}

func (m *Client) logf(format string, v ...any) {
	if m.config != nil && m.config.Logger != nil {
		m.config.Logger.Printf(format, v...)
	}
}

func (m *Client) logln(v ...any) {
	if m.config != nil && m.config.Logger != nil {
		m.config.Logger.Println(v...)
	}
}
