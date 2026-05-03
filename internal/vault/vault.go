package vault

import (
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/temaelkin/qwertify/internal/crypto"
	"github.com/temaelkin/qwertify/internal/storage"
)

var (
	ErrConcurrentModification = errors.New("vault was modified by another process")
	ErrInvalidPassword        = errors.New("invalid password")
)

type Vault struct {
	User string `json:"user"`

	HashedMaster   []byte `json:"hashed_master_key"`
	DerivationSalt []byte `json:"key_salt"`

	Entries map[string]Entry `json:"entries"`

	StateHash [32]byte `json:"-"`
}

type Entry struct {
	EncryptedPassword []byte `json:"encrypted_password"`

	Email    string `json:"email"`
	Username string `json:"username"`

	Meta string `json:"meta"`
}

type AuthData struct {
	URL      string `json:"url"`
	Email    string `json:"email"`
	Username string `json:"username"`
}

func FormAuthData(url string, email string, username string) []byte {
	data, _ := json.Marshal(AuthData{URL: url, Email: email, Username: username})
	return data
}

func hashData(data []byte) [32]byte {
	return sha256.Sum256(data)
}

func Load() (Vault, error) {
	data, err := storage.ReadRaw()
	if err != nil {
		return Vault{}, fmt.Errorf("failed to load vault: %w", err)
	}

	var v Vault
	err = json.Unmarshal(data, &v)
	if err != nil {
		return Vault{}, fmt.Errorf("failed to parse vault data: invalid or corrupted JSON: %w", err)
	}

	v.StateHash = hashData(data)

	return v, nil
}

func Save(v Vault) error {
	data, err := json.MarshalIndent(v, "", " ")
	if err != nil {
		return fmt.Errorf("failed to serialize vault data: %w", err)
	}

	err = storage.WriteRaw(data)
	if err != nil {
		return fmt.Errorf("failed to save vault to storage: %w", err)
	}

	return nil
}

func (v *Vault) SaveOptimistic() error {
	current, err := storage.ReadRaw()
	if err != nil {
		return fmt.Errorf("failed to read current vault state for optimistic save: %w", err)
	}

	if hashData(current) != v.StateHash {
		return ErrConcurrentModification
	}

	new, err := json.MarshalIndent(v, "", " ")
	if err != nil {
		return fmt.Errorf("failed to serialize vault for optimistic save: %w", err)
	}

	v.StateHash = hashData(new)

	err = storage.WriteRaw(new)
	if err != nil {
		return fmt.Errorf("failed to write optimistic save: %w", err)
	}

	return nil
}

func (v *Vault) Authenticate(master []byte) ([]byte, error) {
	if !crypto.VerifyPassword(v.HashedMaster, master) {
		return nil, ErrInvalidPassword
	}

	encryptionKey, err := crypto.GetEncryptionKey(v.DerivationSalt, master)
	if err != nil {
		return nil, fmt.Errorf("failed to derive encryption key during authentication: %w", err)
	}

	return encryptionKey, nil
}

func (e *Entry) Unlock(encryptionKey []byte, authData []byte) ([]byte, error) {
	decryptedPassword, err := crypto.DecryptData(e.EncryptedPassword, encryptionKey, authData)
	if err != nil {
		// Do not wrap the error: preserve security semantics of decryption failure
		return nil, err
	}
	return decryptedPassword, nil
}

func (e *Entry) Lock(password []byte, encryptionKey []byte, authData []byte) error {
	encryptedPassword, err := crypto.EncryptData(password, encryptionKey, authData)
	if err != nil {
		return fmt.Errorf("failed to encrypt password for entry: %w", err)
	}
	e.EncryptedPassword = encryptedPassword
	return nil
}
