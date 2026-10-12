package helpers

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

// GenerateRandomHexToken generates a secure random hex-encoded token of given byte length.
func GenerateRandomHexToken(byteLength int) (string, error) {
	bytes := make([]byte, byteLength)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// GenerateApiKey generates a standard API key with 'pay2u_live_' prefix.
func GenerateApiKey() (string, error) {
	token, err := GenerateRandomHexToken(24)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("pay2u_live_%s", token), nil
}

// GenerateWebhookSecret generates a standard webhook secret with 'whsec_' prefix.
func GenerateWebhookSecret() (string, error) {
	token, err := GenerateRandomHexToken(24)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("whsec_%s", token), nil
}
