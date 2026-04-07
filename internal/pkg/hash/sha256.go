package hash

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base32"
	"encoding/hex"
)

func SHA256Hex(data string) string {
	h := sha256.Sum256([]byte(data))
	return hex.EncodeToString(h[:])
}

func GenerateAPIKey() (fullKey string, prefix string, keyHash string) {
	randomBytes := make([]byte, 32)
	_, _ = rand.Read(randomBytes)

	fullKey = "sk-proj-" + base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(randomBytes)[:48]
	prefix = fullKey[:12]
	keyHash = SHA256Hex(fullKey)

	return fullKey, prefix, keyHash
}
