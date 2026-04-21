package cli

import (
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
		log.Fatalf("Error while loading safe: %v", err)
	}

	inputPassword, err := utils.GetPassword("Enter master password: ")
	if err != nil {
		log.Fatal("Error getting password:", err)
	}
	defer crypto.Wipe(inputPassword)

	mainKey, err := s.Authenticate(inputPassword)
	if err != nil {
		log.Fatal("Error authenticating:", err)
	}
	defer crypto.Wipe(mainKey)

	email, err := utils.GetInput("Enter email: ")
	if err != nil {
		log.Fatal("Input error:", err)
	}

	username, err := utils.GetInput("Enter username: ")
	if err != nil {
		log.Fatal("Input error:", err)
	}

	password, err := utils.GetPassword("Enter password: ")
	if err != nil {
		log.Fatal("Input error:", err)
	}
	defer crypto.Wipe(password)

	associatedData, err := vault.FormAD(url, email, username)
	if err != nil {
		log.Fatalf("Error forming associated data: %v", err)
	}

	newEntry := vault.Entry{
		EncryptedPassword: nil,
		Email:             email,
		Username:          username,
		Meta:              "",
	}

	err = newEntry.Lock(password, mainKey, associatedData)
	if err != nil {
		log.Fatalf("Error locking entry: %v", err)
	}

	utils.PrintEntry(url, newEntry)

	s.Entries[url] = newEntry

	err = vault.Save(s)
	if err != nil {
		log.Fatalf("Error saving safe: %v", err)
	}

	fmt.Println("Entry added successfully!")
}
