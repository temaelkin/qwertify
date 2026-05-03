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

	v, err := vault.Load()
	if err != nil {
		log.Fatalf("Failed to load vault file: %v", err)
	}

	_, ok := v.Entries[url]
	if !ok {
		fmt.Printf("Entry with URL %q not found.\n", url)
		return
	}

	delete(v.Entries, url)

	err = v.SaveOptimistic()
	if err != nil {
		if errors.Is(err, vault.ErrConcurrentModification) {
			fmt.Println("Failed to save: vault was modified by another process.")
			return
		}
		log.Fatalf("Failed to write vault file: %v", err)
	}

	fmt.Printf("Entry %s deleted successfully!", url)
}
