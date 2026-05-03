package cli

import (
	"errors"
	"fmt"
	"log"

	"github.com/atotto/clipboard"
	"github.com/temaelkin/qwertify/internal/crypto"
	"github.com/temaelkin/qwertify/internal/utils"
	"github.com/temaelkin/qwertify/internal/vault"
)

func Get(url string) {
	utils.ClearScreen()

	v, err := vault.Load()
	if err != nil {
		log.Fatalf("Failed to load vault file: %v", err)
	}

	entry, ok := v.Entries[url]
	if !ok {
		fmt.Printf("Entry with URL %q not found.\n", url)
		return
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

	utils.PrintEntry(url, entry)

	authData := vault.FormAuthData(url, entry.Email, entry.Username)

	password, err := entry.Unlock(key, authData)
	if err != nil {
		fmt.Println("Cannot retrieve password. Check input or data integrity.")
		return
	}
	defer crypto.Wipe(password)

	// TODO: different clipboard package
	err = clipboard.WriteAll(string(password))
	if err != nil {
		log.Fatalf("Failed to copy password to clipboard: %v", err)
	}

	fmt.Println("Password copied to clipboard.")
}
