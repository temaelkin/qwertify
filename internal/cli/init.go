package cli

import (
	"crypto/subtle"
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
		log.Fatalf("Failed to check if vault file exists: %v", err)
	}
	if exists {
		fmt.Println("Vault already exists.")
		return
	}

	user, err := utils.GetInput("Enter your name: ")
	if err != nil {
		log.Fatalf("Failed to read input: %v", err)
	}

	master, err := utils.GetPassword("Enter master password: ")
	if err != nil {
		log.Fatalf("Failed to read password: %v", err)
	}
	defer crypto.Wipe(master)

	masterConfirm, err := utils.GetPassword("Confirm master password: ")
	if err != nil {
		log.Fatalf("Failed to read password: %v", err)
	}
	defer crypto.Wipe(masterConfirm)

	if subtle.ConstantTimeCompare(master, masterConfirm) != 1 {
		fmt.Println("Passwords do not match! Please try again.")
		return
	}

	salt, err := crypto.GenSalt()
	if err != nil {
		log.Fatalf("Failed to generate salt: %v", err)
	}

	hashedMaster, err := crypto.HashPassword(master)
	if err != nil {
		log.Fatalf("Failed to hash password: %v", err)
	}

	v := vault.Vault{
		User:           user,
		HashedMaster:   hashedMaster,
		DerivationSalt: salt,
		Entries:        map[string]vault.Entry{},
	}

	err = vault.Save(v)
	if err != nil {
		if errors.Is(err, storage.ErrFileLocked) {
			log.Fatal("Another program is using the vault. Close it and try again.")
		}
		log.Fatalf("Faild to save v file: %v", err)
	}

	fmt.Println("Vault created successfully!")
}
