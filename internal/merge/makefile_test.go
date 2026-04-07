package merge

import (
	"bytes"
	"testing"
)

func TestMakefile(t *testing.T) {
	dst := []byte(`
.PHONY: build run

build:
	go build -o bin/app main.go

custom:
	echo "custom target"
`)

	src := []byte(`
.PHONY: build run test

build:
	go build -o newbin main.go

run:
	go run main.go

test:
	go test ./...
`)

	result, err := Makefile(dst, src)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !bytes.Contains(result, []byte("test:\n\tgo test ./...")) {
		t.Errorf("expected result to contain test target, got:\n%s", string(result))
	}
	if !bytes.Contains(result, []byte("run:\n\tgo run main.go")) {
		t.Errorf("expected result to contain run target, got:\n%s", string(result))
	}
	if !bytes.Contains(result, []byte("custom:\n\techo \"custom target\"")) {
		t.Errorf("expected result to preserve custom target, got:\n%s", string(result))
	}
}
