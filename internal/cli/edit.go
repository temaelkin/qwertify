package cli

import (
	"errors"
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

	var oldEntry vault.Entry
	var newEntry vault.Entry

	entry, ok := s.Entries[url]
	if !ok {
		fmt.Printf("Entry with URL %q not found.\n", url)
		return
	}
	oldEntry = entry

	utils.PrintEntry(url, oldEntry)

	fmt.Printf("Hint: To keep an old value leave a field empty\n\n")

	email, err := utils.GetInput("Enter new email: ")
	if err != nil {
		log.Fatalf("Failed to read input: %v", err)
	}
	if email == "" {
		newEntry.Email = oldEntry.Email
	} else {
		newEntry.Email = email
	}

	username, err := utils.GetInput("Enter new username: ")
	if err != nil {
		log.Fatalf("Failed to read input: %v", err)
	}
	if username == "" {
		newEntry.Username = oldEntry.Username
	} else {
		newEntry.Username = username
	}

	password, err := utils.GetPassword("Enter password: ")
	if err != nil {
		log.Fatalf("Failed to read password: %v", err)
	}
	defer crypto.Wipe(password)

	if len(password) == 0 {
		// Keep the old password
		oldAssociatedData := vault.FormAD(url, oldEntry.Email, oldEntry.Username)

		oldPwd, err := oldEntry.Unlock(mainKey, oldAssociatedData)
		if err != nil {
			fmt.Println("Cannot retrieve password. Check input or data integrity.")
			return
		}
		defer crypto.Wipe(oldPwd)

		newAssociatedData := vault.FormAD(url, newEntry.Email, newEntry.Username)

		err = newEntry.Lock(oldPwd, mainKey, newAssociatedData)
		if err != nil {
			log.Fatalf("Failed to lock entry: %v", err)
		}
	} else {
		// Use new password
		associatedData := vault.FormAD(url, newEntry.Email, newEntry.Username)

		err = newEntry.Lock(password, mainKey, associatedData)
		if err != nil {
			log.Fatalf("Failed to lock entry: %v", err)
		}
	}

	utils.PrintEntry(url, newEntry)

	s.Entries[url] = newEntry

	err = s.SaveOptimistic()
	if err != nil {
		if errors.Is(err, vault.ErrConcurrentModification) {
			fmt.Println("Failed to save: safe was modified by another process.")
			return
		}
		log.Fatalf("Failed to write safe file: %v", err)
	}

	fmt.Println("Entry updated successfully!")
}
