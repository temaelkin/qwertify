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

	entry, ok := s.Entries[url]
	if !ok {
		log.Fatalf("Entry with URL %s not found", url)
	} else {
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

		utils.PrintEntry(url, entry)

		associatedData, err := vault.FormAD(url, entry.Email, entry.Username)
		if err != nil {
			log.Fatalf("Error forming associated data: %v", err)
		}

		pwd, err := entry.Unlock(mainKey, associatedData)
		if err != nil {
			log.Fatal("Error unlocking entry:", err)
		}
		defer crypto.Wipe(pwd)

		// unsafe!
		err = clipboard.WriteAll(string(pwd))
		if err != nil {
			log.Fatal("Error copying password to clipboard:", err)
		}

		fmt.Println("Password copied to clipboard.")
	}
}
