// Package cmd provides CLI commands for the generator.
// This file re-exports layout functionality from the core package for backward compatibility.
package cmd

import "github.com/soliton-go/tools/core"

// ProjectLayout re-exports core.ProjectLayout for backward compatibility.
type ProjectLayout = core.ProjectLayout

// ResolveProjectLayout re-exports core.ResolveProjectLayout for backward compatibility.
func ResolveProjectLayout() (ProjectLayout, error) {
	return core.ResolveProjectLayout()
}
