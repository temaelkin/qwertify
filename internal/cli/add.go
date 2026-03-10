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

	entries, err := s.Unlock(inputPassword)
	if err != nil {
		log.Fatal("Error unlocking safe:", err)
	}

	utils.ClearScreen()

	email, err := utils.GetInput("Enter email: ")
	if err != nil {
		log.Fatal("Input error:", err)
	}

	username, err := utils.GetInput("Enter username: ")
	if err != nil {
		log.Fatal("Input error:", err)
	}

	pwBytes, err := utils.GetPassword("Enter password: ")
	if err != nil {
		log.Fatal("Input error:", err)
	}
	defer crypto.Wipe(pwBytes)

	password := string(pwBytes)

	newEntry := vault.Entry{
		URL:      url,
		Password: password,
		Email:    email,
		Username: username,
	}

	utils.ClearScreen()
	utils.PrintEntry(newEntry, true)

	entries = append(entries, newEntry)

	err = s.Lock(inputPassword, entries)
	if err != nil {
		log.Fatalf("Error locking safe: %v", err)
	}

	err = vault.Save(s)
	if err != nil {
		log.Fatalf("Error saving safe: %v", err)
	}

	fmt.Println("Entry added successfully!")
}
