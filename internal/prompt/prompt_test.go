package prompt_test

import (
	"bytes"
	"ivaldi/internal/prompt"
	"strings"
	"testing"
)

func TestString(t *testing.T) {
	tests := []struct {
		name         string
		message      string
		defaultValue string
		input        string
		expectedOut  string
		expectedAns  string
	}{
		{
			name:         "empty input uses default",
			message:      "Module",
			defaultValue: "github.com/foo",
			input:        "\n",
			expectedOut:  "Module [github.com/foo]: ",
			expectedAns:  "github.com/foo",
		},
		{
			name:         "user input overrides default",
			message:      "Module",
			defaultValue: "github.com/foo",
			input:        "github.com/bar\n",
			expectedOut:  "Module [github.com/foo]: ",
			expectedAns:  "github.com/bar",
		},
		{
			name:         "no default value",
			message:      "Name",
			defaultValue: "",
			input:        "app\n",
			expectedOut:  "Name: ",
			expectedAns:  "app",
		},
		{
			name:         "spaces around input",
			message:      "Name",
			defaultValue: "",
			input:        "   app   \n",
			expectedOut:  "Name: ",
			expectedAns:  "app",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			in := strings.NewReader(tt.input)
			var out bytes.Buffer

			p := prompt.New(in, &out)
			ans := p.String(tt.message, tt.defaultValue)

			if ans != tt.expectedAns {
				t.Errorf("expected answer %q, got %q", tt.expectedAns, ans)
			}
			if out.String() != tt.expectedOut {
				t.Errorf("expected output %q, got %q", tt.expectedOut, out.String())
			}
		})
	}
}

func TestBool(t *testing.T) {
	tests := []struct {
		name         string
		message      string
		defaultValue bool
		input        string
		expectedOut  string
		expectedAns  bool
	}{
		{
			name:         "empty input uses default false",
			message:      "Enable",
			defaultValue: false,
			input:        "\n",
			expectedOut:  "Enable [y/N]: ",
			expectedAns:  false,
		},
		{
			name:         "empty input uses default true",
			message:      "Enable",
			defaultValue: true,
			input:        "\n",
			expectedOut:  "Enable [Y/n]: ",
			expectedAns:  true,
		},
		{
			name:         "y input returns true",
			message:      "Enable",
			defaultValue: false,
			input:        "y\n",
			expectedOut:  "Enable [y/N]: ",
			expectedAns:  true,
		},
		{
			name:         "yes input returns true",
			message:      "Enable",
			defaultValue: false,
			input:        "YES\n",
			expectedOut:  "Enable [y/N]: ",
			expectedAns:  true,
		},
		{
			name:         "n input returns false",
			message:      "Enable",
			defaultValue: true,
			input:        "n\n",
			expectedOut:  "Enable [Y/n]: ",
			expectedAns:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			in := strings.NewReader(tt.input)
			var out bytes.Buffer

			p := prompt.New(in, &out)
			ans := p.Bool(tt.message, tt.defaultValue)

			if ans != tt.expectedAns {
				t.Errorf("expected answer %v, got %v", tt.expectedAns, ans)
			}
			if out.String() != tt.expectedOut {
				t.Errorf("expected output %q, got %q", tt.expectedOut, out.String())
			}
		})
	}
}

func TestSelect(t *testing.T) {
	tests := []struct {
		name        string
		message     string
		options     []string
		input       string
		expectedAns int
	}{
		{
			name:        "valid choice 1",
			message:     "Choose",
			options:     []string{"A", "B", "C"},
			input:       "1\n",
			expectedAns: 0,
		},
		{
			name:        "valid choice 3",
			message:     "Choose",
			options:     []string{"A", "B", "C"},
			input:       "3\n",
			expectedAns: 2,
		},
		{
			name:        "invalid then valid",
			message:     "Choose",
			options:     []string{"A", "B", "C"},
			input:       "0\n4\nfoo\n2\n",
			expectedAns: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			in := strings.NewReader(tt.input)
			var out bytes.Buffer

			p := prompt.New(in, &out)
			ans := p.Select(tt.message, tt.options)

			if ans != tt.expectedAns {
				t.Errorf("expected answer %v, got %v", tt.expectedAns, ans)
			}
		})
	}
}
