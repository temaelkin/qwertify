package vault

import (
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/temaelkin/qwertify/internal/crypto"
	"github.com/temaelkin/qwertify/internal/storage"
)

var ErrConcurrentModification = errors.New("safe was modified by another process")
var ErrInvalidPassword = errors.New("invalid password")

type Safe struct {
	User string `json:"user"`

	HashedMaster []byte `json:"hashed_master_key"`
	KeySalt      []byte `json:"key_salt"`

	Entries map[string]Entry `json:"entries"`

	OriginalHash [32]byte `json:"-"`
}

type Entry struct {
	EncryptedPassword []byte `json:"encrypted_password"`

	Email    string `json:"email"`
	Username string `json:"user_name"`

	Meta string `json:"meta"`
}

type AD struct {
	URL      string `json:"url"`
	Email    string `json:"email"`
	Username string `json:"username"`
}

func FormAD(url string, email string, username string) []byte {
	data, _ := json.Marshal(AD{URL: url, Email: email, Username: username})
	return data
}

func hashData(data []byte) [32]byte {
	return sha256.Sum256(data)
}

func Load() (Safe, error) {
	data, err := storage.ReadRaw()
	if err != nil {
		return Safe{}, fmt.Errorf("failed to load safe: %w", err)
	}

	var s Safe
	err = json.Unmarshal(data, &s)
	if err != nil {
		return Safe{}, fmt.Errorf("failed to parse safe data: invalid or corrupted JSON: %w", err)
	}

	s.OriginalHash = hashData(data)

	return s, nil
}

func Save(s Safe) error {
	data, err := json.MarshalIndent(s, "", " ")
	if err != nil {
		return fmt.Errorf("failed to serialize safe data: %w", err)
	}

	err = storage.WriteRaw(data)
	if err != nil {
		return fmt.Errorf("failed to save safe to storage: %w", err)
	}

	return nil
}

func (s *Safe) SaveOptimistic() error {
	current, err := storage.ReadRaw()
	if err != nil {
		return fmt.Errorf("failed to read current safe state for optimistic save: %w", err)
	}

	if hashData(current) != s.OriginalHash {
		return ErrConcurrentModification
	}

	new, err := json.MarshalIndent(s, "", " ")
	if err != nil {
		return fmt.Errorf("failed to serialize safe for optimistic save: %w", err)
	}

	s.OriginalHash = hashData(new)

	err = storage.WriteRaw(new)
	if err != nil {
		return fmt.Errorf("failed to write optimistic save: %w", err)
	}

	return nil
}

func (s *Safe) Authenticate(masterKey []byte) ([]byte, error) {
	if !crypto.VerifyPassword(s.HashedMaster, masterKey) {
		return nil, ErrInvalidPassword
	}

	mainKey, err := crypto.GetMainKey(s.KeySalt, masterKey)
	if err != nil {
		return nil, fmt.Errorf("failed to derive main key during authentication: %w", err)
	}

	return mainKey, nil
}

func (e *Entry) Unlock(mainKey []byte, associatedData []byte) ([]byte, error) {
	decryptedPwd, err := crypto.DecryptData(e.EncryptedPassword, mainKey, associatedData)
	if err != nil {
		// Do not wrap the error: preserve security semantics of decryption failure
		return nil, err
	}
	return decryptedPwd, nil
}

func (e *Entry) Lock(pwd []byte, mainKey []byte, associatedData []byte) error {
	encryptedPwd, err := crypto.EncryptData(pwd, mainKey, associatedData)
	if err != nil {
		return fmt.Errorf("failed to encrypt password for entry: %w", err)
	}
	e.EncryptedPassword = encryptedPwd
	return nil
}
