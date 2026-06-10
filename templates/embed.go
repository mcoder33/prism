// Package templates ships the tool-neutral workflow texts inside the binary.
package templates

import "embed"

//go:embed conventions.md commands/*.md
var FS embed.FS
