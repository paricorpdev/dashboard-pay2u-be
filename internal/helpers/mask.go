package helpers

import "strings"

// MaskApiKey masks an API key showing only prefix and suffix.
// e.g. "pg_localoka-v2_8772a2b2c30522f8ecb340bd2e49bbc9" -> "pg_localoka-v2_****bbc9"
func MaskApiKey(apiKey string) string {
	if len(apiKey) <= 12 {
		return "****"
	}

	lastUnderscore := strings.LastIndex(apiKey, "_")
	var prefix string
	if lastUnderscore != -1 && lastUnderscore < len(apiKey)-5 {
		prefix = apiKey[:lastUnderscore+1]
	} else if len(apiKey) > 10 {
		prefix = apiKey[:10]
	} else {
		prefix = apiKey[:4]
	}

	suffix := apiKey[len(apiKey)-4:]
	return prefix + "****" + suffix
}

// MaskWebhookSecret masks a webhook secret showing only prefix and suffix.
// e.g. "whsec_947e8ba5418dadf69bc019120110c88695ce180b2d0789f8" -> "whsec_****89f8"
func MaskWebhookSecret(secret string) string {
	if len(secret) <= 10 {
		return "****"
	}

	prefix := "whsec_"
	if !strings.HasPrefix(secret, "whsec_") {
		if len(secret) > 6 {
			prefix = secret[:6]
		} else {
			prefix = secret[:3]
		}
	}

	suffix := secret[len(secret)-4:]
	return prefix + "****" + suffix
}
