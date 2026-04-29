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

	s, err := vault.Load()
	if err != nil {
		log.Fatalf("Failed to load safe file: %v", err)
	}

	inputPassword, err := utils.GetPassword("Enter master password: ")
	if err != nil {
		log.Fatalf("Failed to read password: %v", err)
	}
	defer crypto.Wipe(inputPassword)

	mainKey, err := s.Authenticate(inputPassword)
	if err != nil {
		if errors.Is(err, vault.ErrInvalidPassword) {
			fmt.Println("Invalid password. Please try again.")
			return
		}
		log.Fatalf("Failed to authenticate master password: %v", err)
	}
	defer crypto.Wipe(mainKey)

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

	associatedData := vault.FormAD(url, email, username)

	newEntry := vault.Entry{
		EncryptedPassword: nil,
		Email:             email,
		Username:          username,
		Meta:              "",
	}

	err = newEntry.Lock(password, mainKey, associatedData)
	if err != nil {
		log.Fatalf("Failed to lock entry: %v", err)
	}

	utils.PrintEntry(url, newEntry)

	s.Entries[url] = newEntry

	err = vault.Save(s)
	if err != nil {
		log.Fatalf("Failed to write safe file: %v", err)
	}

	fmt.Println("Entry added successfully!")
}
