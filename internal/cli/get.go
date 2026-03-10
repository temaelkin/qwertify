package cli

import (
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

	found := false

	for _, e := range entries {
		if e.URL == url {
			utils.PrintEntry(e, true)
			err := clipboard.WriteAll(e.Password)
			if err != nil {
				log.Fatal("Error copying password to clipboard:", err)
			}

			fmt.Println("Password copied to clipboard.")

			found = true
			break
		}
	}

	if !found {
		log.Fatal("Entry not found")
	}
}
