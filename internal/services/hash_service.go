package services

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"io"

	"github.com/AlexWilliam12/silent-signal/internal/dtos"
)

func GenerateHash(credentials *dtos.UserRequest) (string, error) {
	salt := make([]byte, 16)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return "", err
	}

	hasher := sha256.New()
	hasher.Write([]byte(credentials.Username + credentials.Password + hex.EncodeToString(salt)))
	return hex.EncodeToString(hasher.Sum(nil)), nil
}
