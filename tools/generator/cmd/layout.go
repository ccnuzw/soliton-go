package cmd

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
	start, err := os.Getwd()
	if err != nil {
		start = "."
	}

	for dir := start; ; dir = filepath.Dir(dir) {
		// Standalone layout: <root>/internal with <root>/go.mod.
		if isDir(filepath.Join(dir, "internal")) && isFile(filepath.Join(dir, "go.mod")) {
			return buildLayout(dir)
		}

		// Monorepo layout: <root>/application/internal with <root>/application/go.mod.
		appDir := filepath.Join(dir, "application")
		if isDir(filepath.Join(appDir, "internal")) && isFile(filepath.Join(appDir, "go.mod")) {
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
	modulePath, err := readModulePath(filepath.Join(moduleDir, "go.mod"))
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

func readModulePath(goModPath string) (string, error) {
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

func isDir(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}

func isFile(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}
