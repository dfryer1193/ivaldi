package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadSaveConfig(t *testing.T) {
	tempDir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", tempDir) // Linux
	t.Setenv("AppData", tempDir)         // Windows
	t.Setenv("HOME", tempDir)            // macOS fallback

	// Ensure UserConfigDir returns a path inside tempDir
	configDir, err := os.UserConfigDir()
	if err != nil {
		t.Fatalf("failed to get user config dir: %v", err)
	}

	// Test Loading non-existent config
	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("unexpected error loading non-existent config: %v", err)
	}
	if cfg.ModulePrefix != "" {
		t.Errorf("expected empty module prefix, got %q", cfg.ModulePrefix)
	}

	// Test Saving config
	cfg.ModulePrefix = "github.com/test/user"
	if err := SaveConfig(cfg); err != nil {
		t.Fatalf("failed to save config: %v", err)
	}

	// Verify file exists
	expectedPath := filepath.Join(configDir, "ivaldi", "config.yaml")
	if _, err := os.Stat(expectedPath); os.IsNotExist(err) {
		t.Errorf("expected config file to be created at %s", expectedPath)
	}

	// Test Loading existing config
	cfg2, err := LoadConfig()
	if err != nil {
		t.Fatalf("failed to load config: %v", err)
	}
	if cfg2.ModulePrefix != "github.com/test/user" {
		t.Errorf("expected module prefix 'github.com/test/user', got %q", cfg2.ModulePrefix)
	}
}
