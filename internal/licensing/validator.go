package licensing

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"strings"
	"time"
)

var ErrSignatureInvalid = errors.New("signature verification failed")

func Validate(publicKey string, activationCode string) (*License, error) {
	cleaned := strings.TrimSpace(strings.Map(func(r rune) rune {
		switch r {
		case ' ', '\n', '\r', '\t':
			return -1
		default:
			return r
		}
	}, activationCode))

	decoded, err := decodeBase64URL(cleaned)
	if err != nil {
		return nil, fmt.Errorf("license validation failed: %w", err)
	}

	if len(decoded) < 8 {
		return nil, errors.New("license validation failed: activation payload too short")
	}

	dataLen := int(binary.BigEndian.Uint32(decoded[:4]))
	if dataLen < 0 || 4+dataLen+4 > len(decoded) {
		return nil, errors.New("license validation failed: invalid data length")
	}
	dataBytes := decoded[4 : 4+dataLen]

	sigLenOffset := 4 + dataLen
	sigLen := int(binary.BigEndian.Uint32(decoded[sigLenOffset : sigLenOffset+4]))
	if sigLen < 0 || sigLenOffset+4+sigLen > len(decoded) {
		return nil, errors.New("license validation failed: invalid signature length")
	}
	signatureBytes := decoded[sigLenOffset+4 : sigLenOffset+4+sigLen]

	pubKey, err := loadPublicKey(publicKey)
	if err != nil {
		return nil, fmt.Errorf("license validation failed: %w", err)
	}

	if err := verifySignature(dataBytes, signatureBytes, pubKey); err != nil {
		return nil, err
	}
	info, err := parseLicenseData(dataBytes)
	if err != nil {
		return nil, fmt.Errorf("license validation failed: %w", err)
	}

	now := time.Now()
	info.Valid = !info.IsExpired(now)
	if !info.ExpireAt.IsZero() {
		info.DaysLeft = int64(info.ExpireAt.Sub(now).Hours() / 24)
	}

	return info, nil
}

func decodeBase64URL(s string) ([]byte, error) {
	if s == "" {
		return nil, errors.New("empty activation code")
	}
	if data, err := base64.RawURLEncoding.DecodeString(s); err == nil {
		return data, nil
	}
	return base64.URLEncoding.DecodeString(s)
}

func parseLicenseData(data []byte) (*License, error) {
	var license License
	if err := json.Unmarshal(data, &license); err != nil {
		return nil, err
	}
	return &license, nil
}

func loadPublicKey(publicKeyStr string) (*rsa.PublicKey, error) {
	if strings.TrimSpace(publicKeyStr) == "" {
		return nil, errors.New("public key is empty")
	}

	cleaned := strings.TrimSpace(publicKeyStr)
	if strings.Contains(cleaned, "BEGIN PUBLIC KEY") {
		block, _ := pem.Decode([]byte(cleaned))
		if block == nil {
			return nil, errors.New("invalid PEM public key")
		}
		key, err := x509.ParsePKIXPublicKey(block.Bytes)
		if err != nil {
			return nil, err
		}
		pub, ok := key.(*rsa.PublicKey)
		if !ok {
			return nil, errors.New("public key is not RSA")
		}
		return pub, nil
	}

	cleaned = strings.ReplaceAll(cleaned, "-----BEGIN PUBLIC KEY-----", "")
	cleaned = strings.ReplaceAll(cleaned, "-----END PUBLIC KEY-----", "")
	cleaned = strings.ReplaceAll(cleaned, "\n", "")
	cleaned = strings.ReplaceAll(cleaned, "\r", "")
	cleaned = strings.ReplaceAll(cleaned, " ", "")
	keyBytes, err := base64.StdEncoding.DecodeString(cleaned)
	if err != nil {
		return nil, err
	}
	key, err := x509.ParsePKIXPublicKey(keyBytes)
	if err != nil {
		return nil, err
	}
	pub, ok := key.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("public key is not RSA")
	}
	return pub, nil
}

func verifySignature(data, signature []byte, publicKey *rsa.PublicKey) error {
	hash := sha256.Sum256(data)
	if err := rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, hash[:], signature); err != nil {
		return ErrSignatureInvalid
	}
	return nil
}
