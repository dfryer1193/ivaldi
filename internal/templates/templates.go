package templates

import (
	_ "embed"
)

//go:embed makefile.tmpl
var Makefile string

//go:embed ci.tmpl
var CI string

//go:embed golangci.yml
var GolangCI string

//go:embed main_http.tmpl
var MainHTTP string

//go:embed main_worker.tmpl
var MainWorker string

//go:embed main_cli.tmpl
var MainCLI string

//go:embed main_interactive.tmpl
var MainInteractive string
