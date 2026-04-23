package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"fmt"
	"io"

	"golang.org/x/crypto/bcrypt"
	"golang.org/x/crypto/scrypt"
)

const saltSize = 32
const scryptNParam = 65536

var ErrMalformedCiphertext = errors.New("malformed ciphertext: insufficient data for nonce")

func DecryptData(encryptedData []byte, mainKey []byte, associatedData []byte) ([]byte, error) {
	block, err := aes.NewCipher(mainKey)
	if err != nil {
		return nil, fmt.Errorf("invalid key for AES cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize GCM mode: %w", err)
	}

	nonceSize := gcm.NonceSize()
	if len(encryptedData) < nonceSize {
		return nil, ErrMalformedCiphertext
	}

	nonce, ciphertext := encryptedData[:nonceSize], encryptedData[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, associatedData)
	if err != nil {
		// Do not wrap the error: preserve security semantics of decryption failure
		return nil, err
	}

	return plaintext, nil
}

func EncryptData(plaintext []byte, mainKey []byte, associatedData []byte) ([]byte, error) {
	block, err := aes.NewCipher(mainKey)
	if err != nil {
		return nil, fmt.Errorf("invalid key for AES cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize GCM mode: %w", err)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, fmt.Errorf("failed to read random nonce: %w", err)
	}

	ciphertext := gcm.Seal(nonce, nonce, plaintext, associatedData)

	return ciphertext, nil
}

func GenSalt() ([]byte, error) {
	b := make([]byte, saltSize)

	_, err := rand.Read(b)
	if err != nil {
		return nil, fmt.Errorf("failed to generate random salt: %w", err)
	}

	return b, nil
}

func HashPassword(password []byte) ([]byte, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	return hashedPassword, nil
}

func VerifyPassword(hashedPassword []byte, password []byte) bool {
	err := bcrypt.CompareHashAndPassword(hashedPassword, password)
	return err == nil
}

func GetMainKey(salt, password []byte) ([]byte, error) {
	mainKey, err := scrypt.Key(password, salt, scryptNParam, 8, 1, 32)
	if err != nil {
		return nil, fmt.Errorf("failed to derive key with scrypt: %w", err)
	}

	return mainKey, nil
}

func Wipe(data []byte) {
	for i := range data {
		data[i] = 0
	}
}
