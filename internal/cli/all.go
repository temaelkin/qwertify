package cli

import (
	"fmt"
	"log"

	"github.com/temaelkin/qwertify/internal/utils"
	"github.com/temaelkin/qwertify/internal/vault"
)

func All() {
	utils.ClearScreen()

	s, err := vault.Load()
	if err != nil {
		log.Fatalf("Failed to load safe file: %v", err)
	}

	for url := range s.Entries {
		fmt.Println(url)
		fmt.Println()
	}
}
