package scaffold

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"text/template"
	"time"

	"ivaldi/internal/merge"
	"ivaldi/internal/templates"
)

type Binary struct {
	Name string
	Type string // "http", "worker", "cli", "interactive"
}

type ProjectConfig struct {
	ModulePath string
	Binaries   []Binary
	SetupCI    bool
	GoVersion  string
}

const (
	modeClobber    = "clobber"
	commandTimeout = 30 * time.Second
)

func InitGoMod(modulePath string) error {
	if _, err := os.Stat("go.mod"); err == nil {
		//nolint:forbidigo // This is a CLI tool, so using fmt for user feedback is appropriate
		fmt.Println("go.mod already exists, skipping init")
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), commandTimeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, "go", "mod", "init", modulePath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func CreateDirectories(binaries []Binary) error {
	dirs := []string{"internal"}
	for _, bin := range binaries {
		dirs = append(dirs, filepath.Join("cmd", bin.Name))
	}

	for _, dir := range dirs {
		//nolint:gosec // We are creating directories in the current working directory, so this is safe
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}
	return nil
}

func WriteMainFiles(config *ProjectConfig) error {
	for _, bin := range binariesToTemplateData(config) {
		path := filepath.Join("cmd", bin.Name, "main.go")
		if _, err := os.Stat(path); err == nil {
			//nolint:forbidigo // This is a CLI tool, so using fmt for user feedback is appropriate
			fmt.Printf("File %s already exists, skipping...\n", path)
			continue
		}

		var tmplContent string
		switch bin.Type {
		case "http":
			tmplContent = templates.MainHTTP
		case "worker":
			tmplContent = templates.MainWorker
		case "cli":
			tmplContent = templates.MainCLI
		case "interactive":
			tmplContent = templates.MainInteractive
		default:
			tmplContent = templates.MainCLI
		}

		tmpl, err := template.New("main").Parse(tmplContent)
		if err != nil {
			return err
		}

		var buf bytes.Buffer
		if err = tmpl.Execute(&buf, map[string]string{
			"ModulePath": config.ModulePath,
			"BinaryName": bin.Name,
		}); err != nil {
			return err
		}

		//nolint:gosec // We are creating files in the current working directory, so this is safe
		if err = os.WriteFile(path, buf.Bytes(), 0644); err != nil {
			return err
		}
	}
	return nil
}

func generateTemplate(name, content string, data any) ([]byte, error) {
	tmpl, err := template.New(name).Parse(content)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	if err = tmpl.Execute(&buf, data); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func WriteTooling(config *ProjectConfig, mode string) error {
	// mode: "init", "update", "clobber"

	// 1. Makefile
	makefileContent, err := generateTemplate("makefile", templates.Makefile, config)
	if err != nil {
		return err
	}
	if err = writeOrMerge("Makefile", makefileContent, mode, merge.Makefile); err != nil {
		return err
	}

	// 2. golangci.yml
	if err = writeOrMerge(".golangci.yml", []byte(templates.GolangCI), mode, merge.YAML); err != nil {
		return err
	}

	// 3. CI
	if config.SetupCI || mode == modeClobber {
		var ciContent []byte
		ciContent, err = generateTemplate("ci", templates.CI, config)
		if err != nil {
			return err
		}
		//nolint:gosec // We are creating directories in the current working directory, so this is safe
		if err = os.MkdirAll(filepath.Join(".github", "workflows"), 0755); err != nil {
			return err
		}

		// For CI, we generally don't merge, just replace if clobber, or skip if init and exists
		ciPath := filepath.Join(".github", "workflows", "ci.yml")
		if mode == modeClobber {
			//nolint:gosec // We are creating files in the current working directory, so this is safe
			_ = os.WriteFile(ciPath, ciContent, 0644)
		} else if _, err = os.Stat(ciPath); os.IsNotExist(err) {
			//nolint:gosec // We are creating files in the current working directory, so this is safe
			_ = os.WriteFile(ciPath, ciContent, 0644)
		}
	}

	return nil
}

func writeOrMerge(path string, srcContent []byte, mode string, mergeFunc func(dst, src []byte) ([]byte, error)) error {
	if mode == modeClobber {
		//nolint:gosec // We are creating files in the current working directory, so this is safe
		return os.WriteFile(path, srcContent, 0644)
	}

	existing, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			//nolint:gosec // We are creating files in the current working directory, so this is safe
			return os.WriteFile(path, srcContent, 0644)
		}
		return err
	}

	if mode == "update" {
		var merged []byte
		merged, err = mergeFunc(existing, srcContent)
		if err != nil {
			return err
		}
		//nolint:gosec // We are creating files in the current working directory, so this is safe
		return os.WriteFile(path, merged, 0644)
	}

	// init mode: skip if exists
	//nolint:forbidigo // This is a CLI tool, so using fmt for user feedback is appropriate
	fmt.Printf("File %s already exists, skipping...\n", path)
	return nil
}

func binariesToTemplateData(config *ProjectConfig) []Binary {
	return config.Binaries
}

func RunGoModTidy() error {
	ctx, cancel := context.WithTimeout(context.Background(), commandTimeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, "go", "mod", "tidy")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
