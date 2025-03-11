package utils

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

func GenerateRandomString(length int) (string, error) {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("GenerateRandomString: %w", err)
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}
