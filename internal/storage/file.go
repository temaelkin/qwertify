package storage

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/gofrs/flock"
)

const file = ".qwertify/vault.json"

var ErrFileLocked = errors.New("storage file is locked by another process")

func getPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user home directory: %w", err)
	}

	return filepath.Join(home, file), nil
}

func ReadRaw() ([]byte, error) {
	path, err := getPath()
	if err != nil {
		return nil, err
	}

	lock := flock.New(path + ".lock")
	locked, err := lock.TryLock()
	if err != nil {
		return nil, fmt.Errorf("failed to acquire file lock: %w", err)
	}
	if !locked {
		return nil, ErrFileLocked
	}
	defer lock.Unlock()

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read storage file at %q: %w", path, err)
	}

	return data, nil
}

func WriteRaw(data []byte) error {
	path, err := getPath()
	if err != nil {
		return err
	}

	err = os.MkdirAll(filepath.Dir(path), 0700)
	if err != nil {
		return fmt.Errorf("failed to create storage directory: %w", err)
	}

	lock := flock.New(path + ".lock")
	locked, err := lock.TryLock()
	if err != nil {
		return fmt.Errorf("failed to acquire file lock: %w", err)
	}
	if !locked {
		return ErrFileLocked
	}
	defer lock.Unlock()

	tmp := path + ".tmp"
	defer os.Remove(tmp)

	err = os.WriteFile(tmp, data, 0600)
	if err != nil {
		return fmt.Errorf("failed to write temporary file: %w", err)
	}

	err = os.Rename(tmp, path)
	if err != nil {
		return fmt.Errorf("failed to replace storage file: %w", err)
	}

	return nil
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
		return false, fmt.Errorf("failed to check storage file existence: %w", err)
	}

	return true, nil
}
