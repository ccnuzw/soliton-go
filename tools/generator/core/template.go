package core

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

func generateTemplateFile(path string, tmpl string, data interface{}, force bool, previewOnly bool) GeneratedFile {
	genFile := GeneratedFile{
		Path: path,
	}

	// Check if file exists
	exists := false
	if _, err := os.Stat(path); err == nil {
		exists = true
		if !force && !previewOnly {
			genFile.Status = FileStatusSkip
			return genFile
		}
	}

	funcMap := template.FuncMap{
		"enumConst": EnumConst,
		"lower":     strings.ToLower,
		"upper":     strings.ToUpper,
	}

	t := template.Must(template.New("file").Funcs(funcMap).Parse(tmpl))
	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		genFile.Status = FileStatusError
		return genFile
	}

	content := buf.String()
	genFile.Content = content

	if previewOnly {
		if exists {
			genFile.Status = FileStatusOverwrite
		} else {
			genFile.Status = FileStatusNew
		}
		return genFile
	}

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
	genFile.Content = ""

	return genFile
}
