package cli

import (
	"fmt"
	"log"

	"github.com/temaelkin/qwertify/internal/crypto"
	"github.com/temaelkin/qwertify/internal/utils"
	"github.com/temaelkin/qwertify/internal/vault"
)

func Edit(url string) {
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

	mainKey, err := s.Authenticate(inputPassword)
	if err != nil {
		log.Fatal("Error authenticating:", err)
	}
	defer crypto.Wipe(mainKey)

	var oldEntry vault.Entry
	var newEntry vault.Entry

	entry, ok := s.Entries[url]
	if !ok {
		log.Fatalf("Entry with URL %s not found", url)
	} else {
		oldEntry = entry
	}

	utils.PrintEntry(url, oldEntry)

	fmt.Printf("Hint: To keep an old value leave a field empty\n\n")

	email, err := utils.GetInput("Enter new email: ")
	if err != nil {
		log.Fatal("Input error:", err)
	}
	if email == "" {
		newEntry.Email = oldEntry.Email
	} else {
		newEntry.Email = email
	}

	username, err := utils.GetInput("Enter new username: ")
	if err != nil {
		log.Fatal("Input error:", err)
	}
	if username == "" {
		newEntry.Username = oldEntry.Username
	} else {
		newEntry.Username = username
	}

	password, err := utils.GetPassword("Enter password: ")
	if err != nil {
		log.Fatal("Input error:", err)
	}
	defer crypto.Wipe(password)

	if len(password) == 0 {
		oldAssociatedData, err := vault.FormAD(url, oldEntry.Email, oldEntry.Username)
		if err != nil {
			log.Fatalf("Error forming associated data: %v", err)
		}

		oldPwd, err := oldEntry.Unlock(mainKey, oldAssociatedData)
		if err != nil {
			log.Fatalf("Error unlocking entry: %v", err)
		}
		defer crypto.Wipe(oldPwd)

		newAssociatedData, err := vault.FormAD(url, newEntry.Email, newEntry.Username)
		if err != nil {
			log.Fatalf("Error forming associated data: %v", err)
		}

		err = newEntry.Lock(oldPwd, mainKey, newAssociatedData)
		if err != nil {
			log.Fatalf("Error locking entry: %v", err)
		}
	} else {
		associatedData, err := vault.FormAD(url, newEntry.Email, newEntry.Username)
		if err != nil {
			log.Fatalf("Error forming associated data: %v", err)
		}

		err = newEntry.Lock(password, mainKey, associatedData)
		if err != nil {
			log.Fatalf("Error locking entry: %v", err)
		}
	}

	utils.PrintEntry(url, newEntry)

	s.Entries[url] = newEntry

	err = s.SaveOptimistic()
	if err != nil {
		log.Fatalf("Error saving safe: %v", err)
	}

	fmt.Println("Entry updated successfully!")
}
