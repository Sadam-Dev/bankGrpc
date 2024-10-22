package utils

import (
	"crypto/sha256"
	"encoding/hex"
)

func GenerateHash(text string) string {
	hash := sha256.New()
	hash.Write([]byte(text))
	hashBytes := hash.Sum(nil)
	return hex.EncodeToString(hashBytes[:])
}
