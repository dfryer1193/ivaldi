package prompt

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

// Prompter handles interactive user input
type Prompter struct {
	in      io.Reader
	out     io.Writer
	scanner *bufio.Scanner
}

// New creates a new Prompter with the given input and output streams
func New(in io.Reader, out io.Writer) *Prompter {
	return &Prompter{
		in:      in,
		out:     out,
		scanner: bufio.NewScanner(in),
	}
}

// String prompts for a string with a default value
func (p *Prompter) String(message, defaultValue string) string {
	if defaultValue != "" {
		fmt.Fprintf(p.out, "%s [%s]: ", message, defaultValue)
	} else {
		fmt.Fprintf(p.out, "%s: ", message)
	}

	if !p.scanner.Scan() {
		return defaultValue
	}

	input := strings.TrimSpace(p.scanner.Text())
	if input == "" {
		return defaultValue
	}
	return input
}

// Bool prompts for a yes/no answer
func (p *Prompter) Bool(message string, defaultValue bool) bool {
	defaultStr := "y/N"
	if defaultValue {
		defaultStr = "Y/n"
	}

	fmt.Fprintf(p.out, "%s [%s]: ", message, defaultStr)

	if !p.scanner.Scan() {
		return defaultValue
	}

	input := strings.ToLower(strings.TrimSpace(p.scanner.Text()))
	if input == "" {
		return defaultValue
	}

	return input == "y" || input == "yes"
}

// Select prompts the user to select an option from a list
func (p *Prompter) Select(message string, options []string) int {
	fmt.Fprintln(p.out, message)
	for i, opt := range options {
		fmt.Fprintf(p.out, "%d. %s\n", i+1, opt)
	}

	for {
		fmt.Fprint(p.out, "Select an option: ")
		if !p.scanner.Scan() {
			return 0
		}

		input := strings.TrimSpace(p.scanner.Text())
		choice, err := strconv.Atoi(input)
		if err == nil && choice >= 1 && choice <= len(options) {
			return choice - 1
		}
		fmt.Fprintln(p.out, "Invalid selection. Please try again.")
	}
}
