package cli

import (
	"errors"
	"fmt"
	"log"

	"github.com/temaelkin/qwertify/internal/crypto"
	"github.com/temaelkin/qwertify/internal/utils"
	"github.com/temaelkin/qwertify/internal/vault"
)

func Add(url string) {
	utils.ClearScreen()

	v, err := vault.Load()
	if err != nil {
		log.Fatalf("Failed to load vault file: %v", err)
	}

	master, err := utils.GetPassword("Enter master password: ")
	if err != nil {
		log.Fatalf("Failed to read password: %v", err)
	}
	defer crypto.Wipe(master)

	key, err := v.Authenticate(master)
	if err != nil {
		if errors.Is(err, vault.ErrInvalidPassword) {
			fmt.Println("Invalid password. Please try again.")
			return
		}
		log.Fatalf("Failed to authenticate master password: %v", err)
	}
	defer crypto.Wipe(key)

	email, err := utils.GetInput("Enter email: ")
	if err != nil {
		log.Fatalf("Failed to read input: %v", err)
	}

	username, err := utils.GetInput("Enter username: ")
	if err != nil {
		log.Fatalf("Failed to read input: %v", err)
	}

	password, err := utils.GetPassword("Enter password: ")
	if err != nil {
		log.Fatalf("Failed to read password: %v", err)
	}
	defer crypto.Wipe(password)

	authData := vault.FormAuthData(url, email, username)

	entry := vault.Entry{
		EncryptedPassword: nil,
		Email:             email,
		Username:          username,
		Meta:              "",
	}

	err = entry.Lock(password, key, authData)
	if err != nil {
		log.Fatalf("Failed to lock entry: %v", err)
	}

	utils.PrintEntry(url, entry)

	v.Entries[url] = entry

	err = vault.Save(v)
	if err != nil {
		log.Fatalf("Failed to write vault file: %v", err)
	}

	fmt.Println("Entry added successfully!")
}
