package cli

import (
	"bytes"
	"errors"
	"fmt"
	"log"

	"github.com/temaelkin/qwertify/internal/crypto"
	"github.com/temaelkin/qwertify/internal/storage"
	"github.com/temaelkin/qwertify/internal/utils"
	"github.com/temaelkin/qwertify/internal/vault"
)

func Init() {
	exists, err := storage.FileExists()
	if err != nil {
		log.Fatalf("Failed to check if safe file exists: %v", err)
	}
	if exists {
		fmt.Println("Safe already exists.")
		return
	}

	userName, err := utils.GetInput("Enter your name: ")
	if err != nil {
		log.Fatalf("Failed to read input: %v", err)
	}

	inputPassword, err := utils.GetPassword("Enter master password: ")
	if err != nil {
		log.Fatalf("Failed to read password: %v", err)
	}
	defer crypto.Wipe(inputPassword)

	inputPasswordConfirm, err := utils.GetPassword("Confirm master password: ")
	if err != nil {
		log.Fatalf("Failed to read password: %v", err)
	}
	defer crypto.Wipe(inputPasswordConfirm)

	// TODO: subtle.ConstantTimeComparison
	if !bytes.Equal(inputPassword, inputPasswordConfirm) {
		fmt.Println("Passwords do not match! Please try again.")
		return
	}

	saltForKey, err := crypto.GenSalt()
	if err != nil {
		log.Fatalf("Failed to generate salt: %v", err)
	}

	hashedPassword, err := crypto.HashPassword(inputPassword)
	if err != nil {
		log.Fatalf("Failed to hash password: %v", err)
	}

	safe := vault.Safe{
		User:         userName,
		HashedMaster: hashedPassword,
		KeySalt:      saltForKey,
		Entries:      map[string]vault.Entry{},
	}

	err = vault.Save(safe)
	if err != nil {
		if errors.Is(err, storage.ErrFileLocked) {
			log.Fatal("Another program is using the safe. Close it and try again.")
		}
		log.Fatalf("Faild to save safe file: %v", err)
	}

	fmt.Println("Safe created successfully!")
}
