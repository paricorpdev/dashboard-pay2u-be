package keys

import (
	"crypto/rand"
	"encoding/base64"
)

func GenerateDashboardMSKey() (string, error) {
	value := make([]byte, 32)
	if _, err := rand.Read(value); err != nil {
		return "", err
	}
	return "dba_" + base64.RawURLEncoding.EncodeToString(value), nil
}
