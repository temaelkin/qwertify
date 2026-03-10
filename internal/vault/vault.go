package vault

import (
	"crypto/sha256"
	"encoding/json"
	"errors"

	"github.com/temaelkin/qwertify/internal/crypto"
	"github.com/temaelkin/qwertify/internal/storage"
)

type Safe struct {
	User          string `json:"user"`
	HashedMaster  string `json:"hashed_master_key"`
	KeySalt       []byte `json:"key_salt"`
	EncryptedData []byte `json:"encrypted_data"`

	OriginalHash [32]byte `json:"-"`
}

type Entry struct {
	URL      string `json:"url"`
	Password string `json:"password"`
	Email    string `json:"email"`
	Username string `json:"user_name"`
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

func (s *Safe) Unlock(password []byte) ([]Entry, error) {
	if !crypto.VerifyPassword(s.HashedMaster, password) {
		return nil, errors.New("invalid password")
	}

	mainKey, err := crypto.GetMainKey(s.KeySalt, password)
	if err != nil {
		return nil, err
	}
	defer crypto.Wipe(mainKey)

	decryptedData, err := crypto.DecryptData(s.EncryptedData, mainKey)
	if err != nil {
		return nil, err
	}
	defer crypto.Wipe(decryptedData)

	var entries []Entry
	err = json.Unmarshal(decryptedData, &entries)
	if err != nil {
		return nil, err
	}

	return entries, nil
}

func (s *Safe) Lock(password []byte, entries []Entry) error {
	newData, err := json.Marshal(entries)
	if err != nil {
		return err
	}

	mainKey, err := crypto.GetMainKey(s.KeySalt, password)
	if err != nil {
		return err
	}
	defer crypto.Wipe(mainKey)

	encryptedNewData, err := crypto.EncryptData(newData, mainKey)
	if err != nil {
		return err
	}

	s.EncryptedData = encryptedNewData
	return nil
}
