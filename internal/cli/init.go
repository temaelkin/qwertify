package cli

import (
	"bytes"
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
		log.Fatalf("Error getting path: %v", err)
	}

	if exists {
		fmt.Println("Safe already exists.")
		return
	}

	userName, err := utils.GetInput("Enter your name: ")
	if err != nil {
		log.Fatal("Input error:", err)
	}

	inputPassword, err := utils.GetPassword("Enter master password: ")
	if err != nil {
		log.Fatal("Error getting password:", err)
	}
	defer crypto.Wipe(inputPassword)

	inputPasswordConfirm, err := utils.GetPassword("Confirm master password: ")
	if err != nil {
		log.Fatal("Error getting password:", err)
	}
	defer crypto.Wipe(inputPasswordConfirm)

	if !bytes.Equal(inputPassword, inputPasswordConfirm) {
		log.Fatal("Passwords do not match!")
	}

	saltForKey, err := crypto.GenSalt()
	if err != nil {
		log.Fatal("Error generating salt:", err)
	}

	hashedPassword, err := crypto.HashPassword(inputPassword)
	if err != nil {
		log.Fatal("Error hashing password:", err)
	}

	safe := vault.Safe{
		User:         userName,
		HashedMaster: hashedPassword,
		KeySalt:      saltForKey,
		Entries:      map[string]vault.Entry{},
	}

	err = vault.Save(safe)
	if err != nil {
		log.Fatalf("Error saving safe: %v", err)
	}

	fmt.Println("Safe created successfully!")
}
