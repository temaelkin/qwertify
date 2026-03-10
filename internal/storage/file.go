package storage

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/gofrs/flock"
)

const file = ".qwertify/safe.json"

func getPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	path := filepath.Join(home, file)

	return path, nil
}

func ReadRaw() ([]byte, error) {
	path, err := getPath()
	if err != nil {
		return nil, err
	}

	lock := flock.New(path + ".lock")

	locked, err := lock.TryLock()
	if err != nil {
		return nil, err
	}

	if !locked {
		return nil, errors.New("file is locked by another process")
	}

	defer lock.Unlock()

	return os.ReadFile(path)
}

func WriteRaw(data []byte) error {
	path, err := getPath()
	if err != nil {
		return err
	}

	err = os.MkdirAll(filepath.Dir(path), 0700)
	if err != nil {
		return err
	}

	lock := flock.New(path + ".lock")

	locked, err := lock.TryLock()
	if err != nil {
		return err
	}

	if !locked {
		return errors.New("file is locked by another process")
	}

	defer lock.Unlock()

	tmp := path + ".tmp"
	defer os.Remove(tmp)

	err = os.WriteFile(tmp, data, 0600)
	if err != nil {
		return err
	}

	return os.Rename(tmp, path)
}

func FileExists() (bool, error) {
	path, err := getPath()
	if err != nil {
		return false, err
	}

	_, err = os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}
