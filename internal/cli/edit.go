package cli

import (
	"fmt"
	"log"

	"github.com/temaelkin/qwertify/internal/crypto"
	"github.com/temaelkin/qwertify/internal/utils"
	"github.com/temaelkin/qwertify/internal/vault"
)

func Edit(url string) {
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

	var oldEntry vault.Entry
	var index int
	found := false

	for i, e := range entries {
		if e.URL == url {
			oldEntry = e
			index = i
			found = true
			break
		}
	}

	if !found {
		log.Fatal("Entry with URL not found:", url)
	}

	fmt.Printf("Hint: To keep an old value leave a field empty\n\n")

	email, err := utils.GetWithDefault("Enter new email: ", oldEntry.Email, false)
	if err != nil {
		log.Fatal("Input error:", err)
	}

	username, err := utils.GetWithDefault("Enter username: ", oldEntry.Username, false)
	if err != nil {
		log.Fatal("Input error:", err)
	}

	password, err := utils.GetWithDefault("Enter password: ", oldEntry.Password, true)
	if err != nil {
		log.Fatal("Input error:", err)
	}

	newEntry := vault.Entry{
		URL:      url,
		Password: password,
		Email:    email,
		Username: username,
	}

	utils.ClearScreen()
	utils.PrintEntry(newEntry, true)

	entries[index] = newEntry

	err = s.Lock(inputPassword, entries)
	if err != nil {
		log.Fatalf("Error locking safe: %v", err)
	}

	err = s.SaveOptimistic()
	if err != nil {
		log.Fatalf("Error saving safe: %v", err)
	}

	fmt.Println("Entry updated successfully!")
}
