package utils

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"syscall"

	"github.com/temaelkin/qwertify/internal/crypto"
	"golang.org/x/term"
)

func GetInput(prompt string) (string, error) {
	fmt.Print(prompt)

	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(input), nil
}

func GetWithDefault(prompt, defaultValue string, hide bool) (string, error) {
	var input string

	if hide {
		fmt.Printf("Old value: %s\n", hidePassword(defaultValue))

		pw, err := GetPassword(prompt)
		if err != nil {
			return "", err
		}
		defer crypto.Wipe(pw)

		input = string(pw)
	} else {
		fmt.Printf("Old value: %s\n", defaultValue)

		val, err := GetInput(prompt)
		if err != nil {
			return "", err
		}

		input = val
	}

	if input == "" {
		return defaultValue, nil
	}

	return input, nil
}

func GetPassword(prompt string) ([]byte, error) {
	fmt.Print(prompt)

	bytepw, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return nil, err
	}
	fmt.Println()

	return bytepw, nil
}
