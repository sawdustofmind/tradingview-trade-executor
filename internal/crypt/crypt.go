package crypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
)

type Crypter struct {
	salt []byte
}

func NewCrypter(salt []byte) *Crypter {
	return &Crypter{
		salt: salt,
	}
}

func (c *Crypter) Decrypt(encrypted string) (string, error) {
	ciphertext := []byte(encrypted)

	block, err := aes.NewCipher(c.salt)
	if err != nil {
		return "", err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// the nonce is prepended to the cipher text so we need to make sure it is still there and length matches up
	nonceSize := aesgcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return "", err
	}

	// now we split the nonce from the ciptertext
	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]

	plaintext, err := aesgcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

func (c *Crypter) Encrypt(plaintext string) (string, error) {
	bplaintext := []byte(plaintext)

	block, err := aes.NewCipher(c.salt)
	if err != nil {
		return "", err
	}

	// Never use more than 2^32 random nonces with a given key because of the risk of a repeat.
	nonce := make([]byte, 12)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// encrypt an prepend the nonce to the ciphertext before returning it
	ciphertext := aesgcm.Seal(nonce, nonce, bplaintext, nil)

	return string(ciphertext), nil
}
