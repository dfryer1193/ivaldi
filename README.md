# Ivaldi

Ivaldi is an opinionated, interactive CLI tool for scaffolding and maintaining Go projects. It sets up standard project structures, configures linters and CI, and provides tailored templates for different types of Go binaries.

## Overview

When you start a new Go project, there is often a lot of boilerplate involved. Ivaldi handles this by generating a standard project layout, an interactive Makefile, GitHub Actions for CI, and a strict `golangci-lint` configuration based on [maratori's golden config](https://github.com/maratori/golangci-lint-config).

It also allows you to define multiple binaries within the same project and select tailored templates for each, including HTTP servers (using `mjolnir`), background workers, and CLI tools.

## Installation

To install Ivaldi, simply use `go install`:

```bash
go install github.com/your-username/ivaldi/cmd/ivaldi@latest
```

Ensure your `$(go env GOPATH)/bin` directory is in your system's `PATH`.

## Configuration

Ivaldi can read a default module path prefix from a configuration file so you don't have to type it out every time you initialize a new project.

Create a configuration file at `~/.config/ivaldi/config.yaml` (or the equivalent user config directory on your OS):

```yaml
module_prefix: "github.com/your-username"
```

You can also override this on the fly using the `-p` or `--module-prefix` flag:

```bash
ivaldi -p github.com/workorg init
```

## Commands

Ivaldi provides three primary subcommands to manage the lifecycle of your project's tooling:

### `init`

The `init` command is used to bootstrap a brand new project. It runs interactively and will:

1. Prompt for your module path (defaulting to your configured prefix + the current directory name).
2. Ask for a comma-separated list of binaries you plan to build.
3. Prompt you to select a template for each binary:
   * **HTTP Server:** Sets up a `chi` router using `mjolnir`, configures `zerolog`, and scaffolds `/health`, `/live`, and `/ready` endpoints.
   * **Background Worker:** Configures `zerolog` and sets up a standard context-based worker loop with graceful shutdown.
   * **CLI Tool:** A basic template utilizing the standard library `flag` package.
   * **Interactive CLI:** A template demonstrating an interactive `bufio.Scanner` loop.
4. Optionally set up GitHub Actions CI.
5. Initialize the `go.mod`, generate the boilerplate code in `cmd/`, and write out the `Makefile` and `.golangci.yml`.

```bash
ivaldi init
```

### `update`

The `update` command safely injects new configurations into an existing project.

* **Makefile:** It parses your existing Makefile and appends any missing standard targets (like `build`, `run`, `install`) without modifying your custom targets.
* **.golangci.yml:** It performs a deep merge of the latest golden config into your existing lint configuration, preserving any custom overrides or rules you have added.

```bash
ivaldi update
```

### `clobber`

The `clobber` command forcefully replaces the `Makefile`, `.golangci.yml`, and `.github/workflows/ci.yml` with the latest embedded templates from Ivaldi.

This is a destructive command for the tooling files, but it strictly ignores your application code in `cmd/` and `internal/` to ensure your business logic is never overwritten.

```bash
ivaldi clobber
```

## Project Structure

Ivaldi enforces a standard Go project layout:

```
.
├── cmd/
│   ├── api/
│   │   └── main.go       # Entrypoint for the 'api' binary
│   └── worker/
│       └── main.go       # Entrypoint for the 'worker' binary
├── internal/             # Private application code
├── .github/
│   └── workflows/
│       └── ci.yml        # GitHub Actions configuration
├── .golangci.yml         # Strict linter configuration
├── go.mod
├── go.sum
└── Makefile              # Generated build and run targets
```
