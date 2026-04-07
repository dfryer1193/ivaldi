package prompt

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var scanner = bufio.NewScanner(os.Stdin)

// String prompts for a string with a default value
func String(message, defaultValue string) string {
	if defaultValue != "" {
		fmt.Printf("%s [%s]: ", message, defaultValue)
	} else {
		fmt.Printf("%s: ", message)
	}

	if !scanner.Scan() {
		return defaultValue
	}

	input := strings.TrimSpace(scanner.Text())
	if input == "" {
		return defaultValue
	}
	return input
}

// Bool prompts for a yes/no answer
func Bool(message string, defaultValue bool) bool {
	defaultStr := "y/N"
	if defaultValue {
		defaultStr = "Y/n"
	}

	fmt.Printf("%s [%s]: ", message, defaultStr)

	if !scanner.Scan() {
		return defaultValue
	}

	input := strings.ToLower(strings.TrimSpace(scanner.Text()))
	if input == "" {
		return defaultValue
	}

	return input == "y" || input == "yes"
}

// Select prompts the user to select an option from a list
func Select(message string, options []string) int {
	fmt.Println(message)
	for i, opt := range options {
		fmt.Printf("%d. %s\n", i+1, opt)
	}

	for {
		fmt.Print("Select an option: ")
		if !scanner.Scan() {
			return 0
		}

		input := strings.TrimSpace(scanner.Text())
		choice, err := strconv.Atoi(input)
		if err == nil && choice >= 1 && choice <= len(options) {
			return choice - 1
		}
		fmt.Println("Invalid selection. Please try again.")
	}
}
