package core

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ProjectLayout describes where to write generated code and which module path to use.
type ProjectLayout struct {
	ModulePath    string
	ModuleDir     string
	InternalDir   string
	DomainDir     string
	AppDir        string
	InfraDir      string
	InterfacesDir string
}

// ResolveProjectLayout finds a project layout by walking up from the current working directory.
func ResolveProjectLayout() (ProjectLayout, error) {
	return ResolveProjectLayoutFrom("")
}

// ResolveProjectLayoutFrom finds a project layout starting from a specified directory.
// If startDir is empty, uses current working directory.
func ResolveProjectLayoutFrom(startDir string) (ProjectLayout, error) {
	if startDir == "" {
		var err error
		startDir, err = os.Getwd()
		if err != nil {
			startDir = "."
		}
	}

	// PRIORITY 1: Check current directory first (standalone layout)
	// This prevents finding parent project when running in a subdirectory
	if IsDir(filepath.Join(startDir, "internal")) && IsFile(filepath.Join(startDir, "go.mod")) {
		return buildLayout(startDir)
	}

	// PRIORITY 2: Check for monorepo layout in current directory
	appDir := filepath.Join(startDir, "application")
	if IsDir(filepath.Join(appDir, "internal")) && IsFile(filepath.Join(appDir, "go.mod")) {
		return buildLayout(appDir)
	}

	// PRIORITY 3: Walk up the directory tree (only if current dir doesn't have a project)
	for dir := filepath.Dir(startDir); ; dir = filepath.Dir(dir) {
		// Standalone layout: <root>/internal with <root>/go.mod.
		if IsDir(filepath.Join(dir, "internal")) && IsFile(filepath.Join(dir, "go.mod")) {
			return buildLayout(dir)
		}

		// Monorepo layout: <root>/application/internal with <root>/application/go.mod.
		appDir := filepath.Join(dir, "application")
		if IsDir(filepath.Join(appDir, "internal")) && IsFile(filepath.Join(appDir, "go.mod")) {
			return buildLayout(appDir)
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
	}

	return ProjectLayout{}, fmt.Errorf("could not locate project layout (missing go.mod or internal directory)")
}

func buildLayout(moduleDir string) (ProjectLayout, error) {
	modulePath, err := ReadModulePath(filepath.Join(moduleDir, "go.mod"))
	if err != nil {
		return ProjectLayout{}, err
	}

	internalDir := filepath.Join(moduleDir, "internal")
	return ProjectLayout{
		ModulePath:    modulePath,
		ModuleDir:     moduleDir,
		InternalDir:   internalDir,
		DomainDir:     filepath.Join(internalDir, "domain"),
		AppDir:        filepath.Join(internalDir, "application"),
		InfraDir:      filepath.Join(internalDir, "infrastructure", "persistence"),
		InterfacesDir: filepath.Join(internalDir, "interfaces", "http"),
	}, nil
}

// ReadModulePath reads the module path from a go.mod file.
func ReadModulePath(goModPath string) (string, error) {
	f, err := os.Open(goModPath)
	if err != nil {
		return "", fmt.Errorf("failed to read go.mod: %w", err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "module ") {
			modulePath := strings.TrimSpace(strings.TrimPrefix(line, "module"))
			if modulePath == "" {
				return "", fmt.Errorf("module path is empty in %s", goModPath)
			}
			return modulePath, nil
		}
	}
	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("failed to read go.mod: %w", err)
	}

	return "", fmt.Errorf("module path not found in %s", goModPath)
}

// IsDir checks if a path is a directory.
func IsDir(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}

// IsFile checks if a path is a file.
func IsFile(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}
