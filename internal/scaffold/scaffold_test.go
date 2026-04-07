package scaffold_test

import (
	"bytes"
	"ivaldi/internal/scaffold"
	"os"
	"path/filepath"
	"testing"
)

func TestCreateDirectories(t *testing.T) {
	tempDir := t.TempDir()
	t.Chdir(tempDir)

	binaries := []scaffold.Binary{
		{Name: "api", Type: "http"},
		{Name: "worker", Type: "worker"},
	}

	if err := scaffold.CreateDirectories(binaries); err != nil {
		t.Fatalf("failed to create directories: %v", err)
	}

	expectedDirs := []string{
		"internal",
		filepath.Join("cmd", "api"),
		filepath.Join("cmd", "worker"),
	}

	for _, dir := range expectedDirs {
		stat, err := os.Stat(dir)
		if err != nil {
			t.Errorf("expected directory %s to exist, but got error: %v", dir, err)
			continue
		}
		if !stat.IsDir() {
			t.Errorf("expected %s to be a directory", dir)
		}
	}
}

func TestWriteMainFiles(t *testing.T) {
	tempDir := t.TempDir()
	t.Chdir(tempDir)

	binaries := []scaffold.Binary{
		{Name: "api", Type: "http"},
		{Name: "cli", Type: "cli"},
	}

	if err := scaffold.CreateDirectories(binaries); err != nil {
		t.Fatalf("failed to create directories: %v", err)
	}

	config := &scaffold.ProjectConfig{
		ModulePath: "github.com/test/app",
		Binaries:   binaries,
	}

	if err := scaffold.WriteMainFiles(config); err != nil {
		t.Fatalf("failed to write main files: %v", err)
	}

	apiMain := filepath.Join("cmd", "api", "main.go")
	content, err := os.ReadFile(apiMain)
	if err != nil {
		t.Fatalf("failed to read api main file: %v", err)
	}
	if !bytes.Contains(content, []byte("github.com/dfryer1193/mjolnir/router")) {
		t.Errorf("expected api main to contain mjolnir router import, got:\n%s", string(content))
	}

	cliMain := filepath.Join("cmd", "cli", "main.go")
	content, err = os.ReadFile(cliMain)
	if err != nil {
		t.Fatalf("failed to read cli main file: %v", err)
	}
	if !bytes.Contains(content, []byte("flag.Parse()")) {
		t.Errorf("expected cli main to contain flag parsing, got:\n%s", string(content))
	}
}

func TestWriteTooling(t *testing.T) {
	tempDir := t.TempDir()
	t.Chdir(tempDir)

	config := &scaffold.ProjectConfig{
		ModulePath: "github.com/test/app",
		Binaries:   []scaffold.Binary{{Name: "api", Type: "http"}},
		SetupCI:    true,
		GoVersion:  "1.26",
	}

	// Test init mode
	if err := scaffold.WriteTooling(config, "init"); err != nil {
		t.Fatalf("failed to write tooling in init mode: %v", err)
	}

	if _, err := os.Stat("Makefile"); err != nil {
		t.Errorf("expected Makefile to be created")
	}
	if _, err := os.Stat(".golangci.yml"); err != nil {
		t.Errorf("expected .golangci.yml to be created")
	}
	if _, err := os.Stat(filepath.Join(".github", "workflows", "ci.yml")); err != nil {
		t.Errorf("expected CI workflow to be created")
	}

	// Test clobber mode (should overwrite)
	if err := scaffold.WriteTooling(config, "clobber"); err != nil {
		t.Fatalf("failed to write tooling in clobber mode: %v", err)
	}
}
