package cli

import (
	"fmt"
	"log"

	"github.com/temaelkin/qwertify/internal/utils"
	"github.com/temaelkin/qwertify/internal/vault"
)

func All() {
	utils.ClearScreen()

	v, err := vault.Load()
	if err != nil {
		log.Fatalf("Failed to load vault file: %v", err)
	}

	for url := range v.Entries {
		fmt.Println(url)
		fmt.Println()
	}
}
