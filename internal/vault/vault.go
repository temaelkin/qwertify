package vault

import (
	"crypto/sha256"
	"encoding/json"
	"errors"

	"github.com/temaelkin/qwertify/internal/crypto"
	"github.com/temaelkin/qwertify/internal/storage"
)

type Safe struct {
	User string `json:"user"`
	// why do we save hash as string?
	// TODO: make it []byte
	// and in crypto.go too
	HashedMaster string `json:"hashed_master_key"`
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

func hashData(data []byte) [32]byte {
	return sha256.Sum256(data)
}

func Load() (Safe, error) {
	data, err := storage.ReadRaw()
	if err != nil {
		return Safe{}, err
	}

	var s Safe
	err = json.Unmarshal(data, &s)
	if err != nil {
		return Safe{}, err
	}

	s.OriginalHash = hashData(data)

	return s, nil
}

func Save(s Safe) error {
	data, err := json.MarshalIndent(s, "", " ")
	if err != nil {
		return err
	}

	return storage.WriteRaw(data)
}

func (s *Safe) SaveOptimistic() error {
	current, err := storage.ReadRaw()
	if err != nil {
		return err
	}

	if hashData(current) != s.OriginalHash {
		return errors.New("safe was changed by another process")
	}

	new, err := json.MarshalIndent(s, "", " ")
	if err != nil {
		return err
	}

	s.OriginalHash = hashData(new)

	return storage.WriteRaw(new)
}

func (s *Safe) Authenticate(masterKey []byte) ([]byte, error) {
	if !crypto.VerifyPassword(s.HashedMaster, masterKey) {
		return nil, errors.New("invalid password")
	}
	return crypto.GetMainKey(s.KeySalt, masterKey)
}

func (e *Entry) Unlock(mainKey []byte, associatedData string) ([]byte, error) {
	// TODO: meta
	decryptedPwd, err := crypto.DecryptData(e.EncryptedPassword, mainKey, []byte(associatedData))
	if err != nil {
		return nil, err
	}
	return decryptedPwd, nil
}

func (e *Entry) Lock(pwd []byte, mainKey []byte, associatedData string) error {
	encryptedPwd, err := crypto.EncryptData(pwd, mainKey, []byte(associatedData))
	if err != nil {
		return err
	}
	e.EncryptedPassword = encryptedPwd
	return nil
}
