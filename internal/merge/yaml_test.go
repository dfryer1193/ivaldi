package merge_test

import (
	"bytes"
	"ivaldi/internal/merge"
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

	result, err := merge.YAML(dst, src)
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
