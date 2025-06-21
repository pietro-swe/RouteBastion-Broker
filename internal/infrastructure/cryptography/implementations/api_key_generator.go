/*
Package implementations provide concrete implementations for cryptography utilities
*/
package implementations

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"

	uuid "github.com/satori/go.uuid"
)

type APIKeyGenerator struct {}

func NewAPIKeyGenerator() *APIKeyGenerator {
	return &APIKeyGenerator{}
}

func (u *APIKeyGenerator) Generate(secretKey []byte, plain string) (string, error) {
	_, err := uuid.FromString(plain)
	if err != nil {
		return "", fmt.Errorf("invalid uuid provided: %w", err)
	}

	block, err := aes.NewCipher(secretKey)
	if err != nil {
		return "", fmt.Errorf("error creating cipher: %w", err)
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("error creating GCM: %w", err)
	}

	nonce := make([]byte, aesGCM.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("error generating nonce: %w", err)
	}

	ciphertext := aesGCM.Seal(nil, nonce, []byte(plain), nil)

	encryptedWithNonce := append(nonce, ciphertext...)
	encoded := base64.URLEncoding.EncodeToString(encryptedWithNonce)

	return encoded, nil
}
