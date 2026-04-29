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

	s, err := vault.Load()
	if err != nil {
		log.Fatalf("Failed to load safe file: %v", err)
	}

	entry, ok := s.Entries[url]
	if !ok {
		fmt.Printf("Entry with URL %q not found.\n", url)
		return
	} else {
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

		utils.PrintEntry(url, entry)

		associatedData := vault.FormAD(url, entry.Email, entry.Username)

		pwd, err := entry.Unlock(mainKey, associatedData)
		if err != nil {
			fmt.Println("Cannot retrieve password. Check input or data integrity.")
			return
		}
		defer crypto.Wipe(pwd)

		// TODO: different clipboard package
		err = clipboard.WriteAll(string(pwd))
		if err != nil {
			log.Fatalf("Failed to copy password to clipboard: %v", err)
		}

		fmt.Println("Password copied to clipboard.")
	}
}
