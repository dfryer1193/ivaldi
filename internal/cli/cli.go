//nolint:forbidigo // This file is a cli UI intended to be used for user interaction, so fmt and os are appropriate here.
package cli //nolint:cyclop // Complex logic for project initialization.

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"ivaldi/internal/config"
	"ivaldi/internal/prompt"
	"ivaldi/internal/scaffold"
)

const (
	choiceHTTP        = 0
	choiceWorker      = 1
	choiceCLI         = 2
	choiceInteractive = 3
)

//nolint:gocognit,nestif,funlen // Complex project initialization logic is inherently nested.
func Run(mode string, modulePrefix string) error {
	cfg, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	if modulePrefix != "" {
		cfg.ModulePrefix = modulePrefix
	}

	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	baseDirName := filepath.Base(cwd)

	defaultModulePath := baseDirName
	if cfg.ModulePrefix != "" {
		defaultModulePath = filepath.Join(cfg.ModulePrefix, baseDirName)
	}

	projectCfg := &scaffold.ProjectConfig{
		GoVersion: "1.26",
	}

	p := prompt.New(os.Stdin, os.Stdout)

	if mode == "init" {
		fmt.Println("=== Ivaldi Project Scaffold ===")
		projectCfg.ModulePath = p.String("Module Path", defaultModulePath)

		fmt.Println("\nConfigure binaries (comma separated names). Leave blank to finish.")
		binariesInput := p.String("Binaries", baseDirName)

		var binaries []scaffold.Binary
		for name := range strings.SplitSeq(binariesInput, ",") {
			name = strings.TrimSpace(name)
			if name == "" {
				continue
			}

			opts := []string{
				"HTTP Server (mjolnir + zerolog)",
				"Background Worker (zerolog)",
				"CLI Tool (standard library)",
				"Interactive CLI",
			}

			choice := p.Select(fmt.Sprintf("Select type for binary '%s':", name), opts)

			var binType string
			switch choice {
			case choiceHTTP:
				binType = "http"
			case choiceWorker:
				binType = "worker"
			case choiceCLI:
				binType = "cli"
			case choiceInteractive:
				binType = "interactive"
			}

			binaries = append(binaries, scaffold.Binary{
				Name: name,
				Type: binType,
			})
		}
		projectCfg.Binaries = binaries

		projectCfg.SetupCI = p.Bool("Setup GitHub Actions CI?", false)

		fmt.Println("\nInitializing go.mod...")
		err = scaffold.InitGoMod(projectCfg.ModulePath)
		if err != nil {
			return fmt.Errorf("go mod init failed: %w", err)
		}

		fmt.Println("Creating directories...")
		err = scaffold.CreateDirectories(projectCfg.Binaries)
		if err != nil {
			return err
		}

		fmt.Println("Writing main.go files...")
		err = scaffold.WriteMainFiles(projectCfg)
		if err != nil {
			return err
		}

		fmt.Println("Running go mod tidy to fetch dependencies...")
		err = scaffold.RunGoModTidy()
		if err != nil {
			fmt.Printf("Warning: go mod tidy failed: %v\n", err)
		}
	} else {
		// For update and clobber, we try to detect existing binaries from the Makefile or cmd/ dir
		// We'll just infer binaries from the directories in cmd/
		var entries []os.DirEntry
		entries, err = os.ReadDir("cmd")
		if err == nil {
			for _, entry := range entries {
				if entry.IsDir() {
					projectCfg.Binaries = append(projectCfg.Binaries, scaffold.Binary{
						Name: entry.Name(),
						Type: "cli", // Type doesn't matter for makefile generation except Name
					})
				}
			}
		}
	}

	fmt.Printf("Writing tooling files (mode: %s)...\n", mode)
	err = scaffold.WriteTooling(projectCfg, mode)
	if err != nil {
		return err
	}

	fmt.Println("Done!")
	return nil
}
