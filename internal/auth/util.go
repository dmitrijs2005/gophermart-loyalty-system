package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"

	"golang.org/x/crypto/pbkdf2"
)

func GenerateSalt(size int) ([]byte, error) {
	salt := make([]byte, size)
	_, err := rand.Read(salt)
	if err != nil {
		return nil, err
	}
	return salt, nil
}

func HashPassword(password string, salt []byte) string {
	iter := 100_000
	keyLen := 32

	hash := pbkdf2.Key([]byte(password), salt, iter, keyLen, sha256.New)
	return base64.StdEncoding.EncodeToString(hash)
}
