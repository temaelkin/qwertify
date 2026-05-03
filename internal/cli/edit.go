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

	var old vault.Entry
	var new vault.Entry

	entry, ok := v.Entries[url]
	if !ok {
		fmt.Printf("Entry with URL %q not found.\n", url)
		return
	}
	old = entry

	utils.PrintEntry(url, old)

	fmt.Printf("Hint: To keep an old value leave a field empty\n\n")

	email, err := utils.GetInput("Enter new email: ")
	if err != nil {
		log.Fatalf("Failed to read input: %v", err)
	}
	if email == "" {
		new.Email = old.Email
	} else {
		new.Email = email
	}

	username, err := utils.GetInput("Enter new username: ")
	if err != nil {
		log.Fatalf("Failed to read input: %v", err)
	}
	if username == "" {
		new.Username = old.Username
	} else {
		new.Username = username
	}

	password, err := utils.GetPassword("Enter password: ")
	if err != nil {
		log.Fatalf("Failed to read password: %v", err)
	}
	defer crypto.Wipe(password)

	if len(password) == 0 {
		// Keeping the old password
		oldAuthData := vault.FormAuthData(url, old.Email, old.Username)

		oldPassword, err := old.Unlock(key, oldAuthData)
		if err != nil {
			fmt.Println("Cannot retrieve password. Check input or data integrity.")
			return
		}
		defer crypto.Wipe(oldPassword)

		newAuthData := vault.FormAuthData(url, new.Email, new.Username)

		err = new.Lock(oldPassword, key, newAuthData)
		if err != nil {
			log.Fatalf("Failed to lock entry: %v", err)
		}
	} else {
		// Using new password
		authData := vault.FormAuthData(url, new.Email, new.Username)

		err = new.Lock(password, key, authData)
		if err != nil {
			log.Fatalf("Failed to lock entry: %v", err)
		}
	}

	utils.PrintEntry(url, new)

	v.Entries[url] = new

	err = v.SaveOptimistic()
	if err != nil {
		if errors.Is(err, vault.ErrConcurrentModification) {
			fmt.Println("Failed to save: vault was modified by another process.")
			return
		}
		log.Fatalf("Failed to write vault file: %v", err)
	}

	fmt.Println("Entry updated successfully!")
}
