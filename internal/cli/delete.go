package cli

import (
	"fmt"
	"log"

	"github.com/temaelkin/qwertify/internal/crypto"
	"github.com/temaelkin/qwertify/internal/utils"
	"github.com/temaelkin/qwertify/internal/vault"
)

func Delete(url string) {
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

	if _, ok := entries[url]; !ok {
		log.Fatalf("Entry with URL %s not found", url)
	} else {
		delete(entries, url)
	}

	err = s.Lock(inputPassword, entries)
	if err != nil {
		log.Fatalf("Error locking safe: %v", err)
	}

	err = s.SaveOptimistic()
	if err != nil {
		log.Fatalf("Error saving safe: %v", err)
	}

	fmt.Printf("Entry %s deleted successfully!", url)
}
