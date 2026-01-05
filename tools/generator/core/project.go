package core

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"text/template"
)

// InitProject initializes a new project with the given configuration.
func InitProject(cfg ProjectConfig) (*GenerationResult, error) {
	return initProjectInternal(cfg, false)
}

// PreviewInitProject previews what files would be created without actually creating them.
func PreviewInitProject(cfg ProjectConfig) (*GenerationResult, error) {
	return initProjectInternal(cfg, true)
}

func initProjectInternal(cfg ProjectConfig, previewOnly bool) (*GenerationResult, error) {
	result := &GenerationResult{
		Success: true,
		Files:   []GeneratedFile{},
	}

	projectName := cfg.Name
	moduleName := cfg.ModuleName
	if moduleName == "" {
		moduleName = "github.com/soliton-go/" + projectName
	}

	// Determine project root path FIRST
	// Create project at ../../<project-name> (parallel to tools directory)
	cwd, _ := os.Getwd()
	var projectRoot string
	var projectParentDir string

	fmt.Printf("   DEBUG cwd: %s\n", cwd)
	fmt.Printf("   DEBUG pattern: %s\n", filepath.Join("tools", "generator"))
	fmt.Printf("   DEBUG contains: %v\n", strings.Contains(cwd, filepath.Join("tools", "generator")))

	if strings.Contains(cwd, filepath.Join("tools", "generator")) {
		// Running from tools/generator, create at ../../<project>
		projectRoot = filepath.Join("..", "..", projectName)
		projectParentDir = filepath.Join(cwd, "..", "..")
	} else {
		// Running from elsewhere
		projectRoot = projectName
		projectParentDir = cwd
	}

	// Now detect framework based on where project will be created
	frameworkVersion := cfg.FrameworkVersion
	frameworkReplace := cfg.FrameworkReplace

	if frameworkReplace == "" {
		frameworkPath := filepath.Join(projectParentDir, "framework")
		if info, err := os.Stat(frameworkPath); err == nil && info.IsDir() {
			frameworkReplace = filepath.ToSlash(filepath.Join("..", "framework"))
		}
	}

	if frameworkVersion == "" {
		if frameworkReplace != "" {
			frameworkVersion = "v0.0.0-00010101000000-000000000000"
		} else {
			frameworkVersion = "v0.1.0"
		}
	}

	// Detect Go version
	goVersion := "1.23" // Default safe fallback

	// 1. Try to read from go.work (Project Preference)
	// We look for go.work in current dir or parent dirs, but InitProject usually runs at root or uses logical root.
	// For simplicity, we check the probable workspace root: projectParentDir
	goWorkPath := filepath.Join(projectParentDir, "go.work")
	if goWorkData, err := os.ReadFile(goWorkPath); err == nil {
		lines := strings.Split(string(goWorkData), "\n")
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if strings.HasPrefix(line, "go ") {
				v := strings.TrimSpace(strings.TrimPrefix(line, "go "))
				if v != "" {
					goVersion = v
					break
				}
			}
		}
	} else {
		// 2. Try runtime version (Environment Preference)
		runtimeVer := runtime.Version()
		runtimeVer = strings.TrimPrefix(runtimeVer, "go")
		if parts := strings.Split(runtimeVer, "."); len(parts) >= 2 {
			// Extract major.minor (e.g. 1.23)
			goVersion = parts[0] + "." + parts[1]
		}
	}

	// Ensure we don't output "devel" or unstable versions inadvertently
	if strings.Contains(goVersion, "devel") {
		goVersion = "1.22"
	}

	data := ProjectData{
		ProjectName:      projectName,
		ModuleName:       moduleName,
		FrameworkVersion: frameworkVersion,
		FrameworkReplace: frameworkReplace,
		GoVersion:        goVersion,
		FxVersion:        "v1.22.0",
	}

	// projectRoot was already determined above

	// Create project root directory
	if !previewOnly {
		if err := os.MkdirAll(projectRoot, 0755); err != nil {
			return nil, fmt.Errorf("failed to create project directory: %w", err)
		}
	}

	// Create directory structure
	dirs := []string{
		"cmd",
		"configs",
		"internal/domain",
		"internal/application",
		"internal/infrastructure/persistence",
		"internal/interfaces/http",
	}

	if !previewOnly {
		for _, dir := range dirs {
			path := filepath.Join(projectRoot, dir)
			if err := os.MkdirAll(path, 0755); err != nil {
				result.Errors = append(result.Errors, fmt.Sprintf("failed to create directory %s: %v", dir, err))
			}
		}
	}

	// Files to generate
	files := []struct {
		path     string
		template string
	}{
		{"go.mod", GoModTemplate},
		{"cmd/main.go", MainTemplate},
		{"configs/config.yaml", ConfigTemplate},
		{"configs/config.example.yaml", ConfigExampleTemplate},
		{"internal/interfaces/http/response.go", ResponseTemplate},
		{".gitignore", GitignoreTemplate},
		{"README.md", ReadmeTemplate},
		{"Makefile", MakefileTemplate},
	}

	for _, f := range files {
		fullPath := filepath.Join(projectRoot, f.path)
		genFile := generateProjectFile(fullPath, f.template, data, previewOnly)
		result.Files = append(result.Files, genFile)
		if genFile.Status == FileStatusError {
			result.Errors = append(result.Errors, fmt.Sprintf("%s: generation failed", f.path))
		}
	}

	if len(result.Errors) > 0 {
		result.Success = false
	}

	// Include project path in message
	absPath, _ := filepath.Abs(projectRoot)
	result.Message = fmt.Sprintf("Project %s initialized successfully at %s", projectName, absPath)
	return result, nil
}

func generateProjectFile(path string, tmpl string, data ProjectData, previewOnly bool) GeneratedFile {
	genFile := GeneratedFile{
		Path: path,
	}

	// Check if file exists
	exists := false
	if _, err := os.Stat(path); err == nil {
		exists = true
		genFile.Status = FileStatusSkip
		if previewOnly {
			return genFile
		}
		return genFile
	}

	// Render template
	t := template.Must(template.New("file").Parse(tmpl))
	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		genFile.Status = FileStatusError
		return genFile
	}

	content := buf.String()
	genFile.Content = content

	if previewOnly {
		genFile.Status = FileStatusNew
		return genFile
	}

	// Create file
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		genFile.Status = FileStatusError
		return genFile
	}

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		genFile.Status = FileStatusError
		return genFile
	}

	if exists {
		genFile.Status = FileStatusOverwrite
	} else {
		genFile.Status = FileStatusNew
	}
	genFile.Content = "" // Don't include content in non-preview mode

	return genFile
}

// GetDefaultModuleName returns the default module name for a project.
func GetDefaultModuleName(projectName string) string {
	return "github.com/soliton-go/" + projectName
}

// ValidateProjectConfig validates the project configuration.
func ValidateProjectConfig(cfg ProjectConfig) error {
	if cfg.Name == "" {
		return fmt.Errorf("project name is required")
	}
	if strings.ContainsAny(cfg.Name, " \t\n/\\") {
		return fmt.Errorf("project name contains invalid characters")
	}
	return nil
}
