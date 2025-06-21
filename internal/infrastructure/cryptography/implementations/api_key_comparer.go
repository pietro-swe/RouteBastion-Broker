package implementations

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
)

type APIKeyComparer struct {}

func (c *APIKeyComparer) Compare(secretKey []byte, hash string, plain string) (bool, error) {
    data, err := base64.URLEncoding.DecodeString(hash)
    if err != nil {
        return false, fmt.Errorf("error decoding base64: %w", err)
    }

    block, err := aes.NewCipher(secretKey)
    if err != nil {
        return false, err
    }

    aesGCM, err := cipher.NewGCM(block)
    if err != nil {
        return false, err
    }

    nonceSize := aesGCM.NonceSize()
    if len(data) < nonceSize {
        return false, fmt.Errorf("ciphertext too short")
    }

    nonce, ciphertext := data[:nonceSize], data[nonceSize:]

    plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
    if err != nil {
        return false, fmt.Errorf("error decrypting: %w", err)
    }

    return string(plaintext) == plain, nil
}
