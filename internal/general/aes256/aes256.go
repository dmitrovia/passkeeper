package aes256

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
)

func generateRandom(size int) ([]byte, error) {
	random := make([]byte, size)

	_, err := rand.Read(random)
	if err != nil {
		return nil, fmt.Errorf("generateRandom->Read: %w", err)
	}

	return random, nil
}

func Encrypt(data *[]byte, key *[]byte) (*[]byte, error) {
	aesblock, err := aes.NewCipher(*key)
	if err != nil {
		return nil, fmt.Errorf("Encrypt->NewCipher: %w", err)
	}

	aesgcm, err := cipher.NewGCM(aesblock)
	if err != nil {
		return nil, fmt.Errorf("Encrypt->NewGCM: %w", err)
	}

	nonce, err := generateRandom(aesgcm.NonceSize())
	if err != nil {
		return nil, fmt.Errorf("Encrypt->generateRandom: %w", err)
	}

	est := aesgcm.Seal(nonce, nonce, *data, nil)

	return &est, nil
}

func Decrypt(data *[]byte, key *[]byte) (*[]byte, error) {
	aesblock, err := aes.NewCipher(*key)
	if err != nil {
		return nil, fmt.Errorf("Decrypt->NewCipher: %w", err)
	}

	aesgcm, err := cipher.NewGCM(aesblock)
	if err != nil {
		return nil, fmt.Errorf("Decrypt->NewGCM: %w", err)
	}

	nonceSize := aesgcm.NonceSize()

	nonce, ciphertext := (*data)[:nonceSize],
		(*data)[nonceSize:]

	dec, err := aesgcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("Decrypt->Open: %w", err)
	}

	return &dec, nil
}
