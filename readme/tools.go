//go:build tools

// Package for Go build and development tools that are excluded from the main package.
// See https://github.com/golang/go/wiki/Modules#how-can-i-track-tool-dependencies-for-a-module
// for more information.

package tools

import (
	// Linting (make lint).
	_ "github.com/golangci/golangci-lint/cmd/golangci-lint"
	// Test coverage reporting (make coverage).
	_ "github.com/boumenot/gocover-cobertura"
	// Markdown documentation generator (go generate ./...).
	_ "github.com/princjef/gomarkdoc/cmd/gomarkdoc"
	// Code vulnerability checks (make check-vuln).
	_ "golang.org/x/vuln/cmd/govulncheck"
	// Formatting and linting (make gofumpt).
	_ "mvdan.cc/gofumpt"
)
