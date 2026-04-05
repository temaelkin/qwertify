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

	if entry, ok := entries[url]; !ok {
		log.Fatalf("Entry with URL %s not found", url)
	} else {
		utils.PrintEntry(url, entry, true)
		err := clipboard.WriteAll(entry.Password)
		if err != nil {
			log.Fatal("Error copying password to clipboard:", err)
		}

		fmt.Println("Password copied to clipboard.")
	}
}
