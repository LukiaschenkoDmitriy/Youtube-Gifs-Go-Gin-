package utils

import (
	"crypto/rand"
	"encoding/hex"
)

func RandomHex(n int) (string, error) {
	bytes := make([]byte, n)

	if _, err := rand.Read(bytes); err != nil {
		return "Random Read Error", err
	}

	return hex.EncodeToString(bytes), nil
}
