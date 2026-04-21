package cli

import (
	"fmt"
	"log"

	"github.com/temaelkin/qwertify/internal/utils"
	"github.com/temaelkin/qwertify/internal/vault"
)

func Delete(url string) {
	utils.ClearScreen()

	s, err := vault.Load()
	if err != nil {
		log.Fatalf("Error while loading safe: %v", err)
	}

	_, ok := s.Entries[url]
	if !ok {
		log.Fatalf("Entry with URL %s not found", url)
	} else {
		delete(s.Entries, url)
	}

	err = s.SaveOptimistic()
	if err != nil {
		log.Fatalf("Error saving safe: %v", err)
	}

	fmt.Printf("Entry %s deleted successfully!", url)
}
