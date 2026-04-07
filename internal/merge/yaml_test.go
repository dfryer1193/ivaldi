package merge

import (
	"bytes"
	"testing"
)

func TestYAML(t *testing.T) {
	dst := []byte(`
version: "2"
issues:
  max-same-issues: 10
linters:
  enable:
    - gofmt
`)

	src := []byte(`
version: "2"
issues:
  max-same-issues: 50
  exclude-use-default: false
linters:
  enable:
    - gofmt
    - govet
  disable-all: true
`)

	// Our deep merge should:
	// - Preserve max-same-issues: 10
	// - Add exclude-use-default: false
	// - For arrays, our naive merge currently merges map keys but not elements.
	// Let's see what the current merge logic outputs. We expect it to add missing keys at mapping level.
	// Array merging is not deeply handled by our simple naive logic; it just treats them as scalar or skips.
	// Actually, wait, our YAML merge logic:
	// "Existing keys in dst are preserved, missing keys are added from src."
	// Let's test the map merging logic.

	result, err := YAML(dst, src)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expectedParts := []string{
		"max-same-issues: 10",
		"exclude-use-default: false",
		"disable-all: true",
	}

	for _, part := range expectedParts {
		if !bytes.Contains(result, []byte(part)) {
			t.Errorf("expected result to contain %q, got:\n%s", part, string(result))
		}
	}
}
