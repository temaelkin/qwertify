package cli

import (
	"log"

	"github.com/temaelkin/qwertify/internal/crypto"
	"github.com/temaelkin/qwertify/internal/utils"
	"github.com/temaelkin/qwertify/internal/vault"
)

func All() {
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

	for _, e := range entries {
		utils.PrintEntry(e, true)
	}
}
