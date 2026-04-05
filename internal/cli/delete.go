package cli

import (
	"fmt"
	"log"
	"slices"

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

	var index int
	found := false

	for i, e := range entries {
		if e.URL == url {
			index = i
			found = true
			break
		}
	}

	if !found {
		log.Fatal("Entry with URL not found:", url)
	}

	utils.ClearScreen()

	entries = slices.Delete(entries, index, index+1)

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
