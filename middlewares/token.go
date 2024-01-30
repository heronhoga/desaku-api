package middlewares

import (
	"crypto/rand"
	"encoding/hex"
)

func GenerateSessionToken() (string, error) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	token := hex.EncodeToString(b)
	return token, nil
}
