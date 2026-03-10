package cli

import (
	"log"

	"github.com/temaelkin/qwertify/internal/utils"
)

func Handle(args []string) {
	if len(args) < 2 {
		log.Fatal("Usage: qwfy <command>")
	}

	cmd := args[1]

	switch cmd {
	case "init":
		Init()
	case "add":
		if len(args) < 3 {
			log.Fatal("Usage: qwfy add <url>")
		}
		url := args[2]
		Add(url)
	case "get":
		if len(args) < 3 {
			log.Fatal("Usage: qwfy get <url>")
		}
		url := args[2]
		Get(url)
	case "edit":
		if len(args) < 3 {
			log.Fatal("Usage: qwfy edit <url>")
		}
		url := args[2]
		Edit(url)
	case "all":
		All()
	case "help":
		utils.PrintHelp()
	default:
		log.Fatalf("Unknown command: %s", cmd)
	}
}
