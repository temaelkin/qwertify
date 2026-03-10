package utils

import (
	"fmt"
	"os"
	"strings"

	"github.com/temaelkin/qwertify/internal/vault"
	"golang.org/x/term"
)

func ClearScreen() {
	if term.IsTerminal(int(os.Stdout.Fd())) {
		fmt.Print("\033c")
	} else {
		fmt.Println(strings.Repeat("\n", 100))
	}
}

func hidePassword(pw string) string {
	var builder strings.Builder

	for i := 0; i < len(pw); i++ {
		builder.WriteString("*")
	}

	return builder.String()
}

func PrintHelp() {
	ClearScreen()

	fmt.Println("=========================")
	fmt.Println()

	fmt.Println(`                              _  ___       `)
	fmt.Println(`                         _   (_)/ __)      `)
	fmt.Println(`  ____ _ _ _  ____  ____| |_  _| |__ _   _ `)
	fmt.Println(` / _  | | | |/ _  )/ ___)  _)| |  __) | | |`)
	fmt.Println(`| | | | | | ( (/ /| |   | |__| | |  | |_| |`)
	fmt.Println(` \_|| |\____|\____)_|    \___)_|_|   \__  |`)
	fmt.Println(`    |_|                             (____/ `)

	fmt.Println()
	fmt.Println("qwertify")
	fmt.Println()
	fmt.Println("Secure offline CLI password manager.")
	fmt.Println("Works and stores only on your machine.")
	fmt.Println()
	fmt.Println("temaelkin, 2026, v.1.0.0")
	fmt.Println()
	fmt.Println("=========================")

	fmt.Println()
	fmt.Println("How to use:")
	fmt.Println("qwfy <command> <optional>")
	fmt.Println()
	fmt.Println("1. add <url> - add new entry with a URL")
	fmt.Println("2. get <url> - get exiting entry")
	fmt.Println("3. edit <url> - edit existing entry")
	fmt.Println("4. all - get all entries")
	fmt.Println("5. help - you are here!")
	fmt.Println()
}

func PrintEntry(e vault.Entry, hide bool) {
	fmt.Println("=========================")
	fmt.Println(e.URL)
	fmt.Println("=========================")
	fmt.Println("Email:    ", e.Email)
	fmt.Println("Username: ", e.Username)
	if hide {
		fmt.Println("Password: ", hidePassword(e.Password))
	} else {
		fmt.Println("Password: ", e.Password)
	}
	fmt.Println()
}
