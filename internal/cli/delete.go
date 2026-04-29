package cli

import (
	"errors"
	"fmt"
	"log"

	"github.com/temaelkin/qwertify/internal/utils"
	"github.com/temaelkin/qwertify/internal/vault"
)

func Delete(url string) {
	utils.ClearScreen()

	s, err := vault.Load()
	if err != nil {
		log.Fatalf("Failed to load safe file: %v", err)
	}

	_, ok := s.Entries[url]
	if !ok {
		fmt.Printf("Entry with URL %q not found.\n", url)
		return
	} else {
		delete(s.Entries, url)
	}

	err = s.SaveOptimistic()
	if err != nil {
		if errors.Is(err, vault.ErrConcurrentModification) {
			fmt.Println("Failed to save: safe was modified by another process.")
			return
		}
		log.Fatalf("Failed to write safe file: %v", err)
	}

	fmt.Printf("Entry %s deleted successfully!", url)
}
