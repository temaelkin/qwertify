package main

import (
	"os"

	"github.com/temaelkin/qwertify/internal/cli"
)

func main() {
	cli.Handle(os.Args)
}
