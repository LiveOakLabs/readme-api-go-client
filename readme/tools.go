//go:build tools

// Package for Go build and development tools that are excluded from the main package.
// See https://github.com/golang/go/wiki/Modules#how-can-i-track-tool-dependencies-for-a-module
// for more information.

package tools

import (
	_ "github.com/boumenot/gocover-cobertura"               // Test coverage reporting (make coverage).
	_ "github.com/golangci/golangci-lint/cmd/golangci-lint" // Linting (make lint).
	_ "github.com/princjef/gomarkdoc/cmd/gomarkdoc"         // Markdown documentation generator (go generate ./...).
	_ "github.com/segmentio/golines"                        // Long line fixer.
	_ "golang.org/x/vuln/cmd/govulncheck"                   // Code vulnerability checks (make check-vuln).
	_ "mvdan.cc/gofumpt"                                    // Formatting and linting (make gofumpt).
)
