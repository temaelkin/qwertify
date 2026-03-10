package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"io"

	"golang.org/x/crypto/bcrypt"
	"golang.org/x/crypto/scrypt"
)

const saltSize = 32
const scryptNParam = 65536

func DecryptData(encryptedData []byte, mainKey []byte) ([]byte, error) {
	block, err := aes.NewCipher(mainKey)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(encryptedData) < nonceSize {
		return nil, errors.New("malformed ciphertext")
	}

	nonce, ciphertext := encryptedData[:nonceSize], encryptedData[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}

func EncryptData(plaintext []byte, mainKey []byte) ([]byte, error) {
	block, err := aes.NewCipher(mainKey)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)

	return ciphertext, nil
}

func GenSalt() ([]byte, error) {
	b := make([]byte, saltSize)

	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func HashPassword(password []byte) ([]byte, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	return hashedPassword, nil
}

func VerifyPassword(hashedPassword string, password []byte) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), password)

	return err == nil
}

func GetMainKey(salt, password []byte) ([]byte, error) {
	mainKey, err := scrypt.Key(password, salt, scryptNParam, 8, 1, 32)
	if err != nil {
		return nil, err
	}

	return mainKey, nil
}

func Wipe(data []byte) {
	for i := range data {
		data[i] = 0
	}
}
