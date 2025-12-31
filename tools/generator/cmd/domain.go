package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/spf13/cobra"
)

var domainCmd = &cobra.Command{
	Use:   "domain [name]",
	Short: "Generate a new domain entity",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		fmt.Printf("Generating domain entity: %s\n", name)

		generateDomain(name)
	},
}

func init() {
	rootCmd.AddCommand(domainCmd)
}

func generateDomain(name string) {
	// wd, _ := os.Getwd()
	// Assuming we are running from project root, or we need to find it.
	// For simplicity in this scaffold, we assume running from project root or tools/generator
	// Let's assume absolute path for prototype safety or relative to "application"

	// Try to find application directory
	baseDir := "../../application/internal/domain"
	if _, err := os.Stat(baseDir); os.IsNotExist(err) {
		// Fallback if running from root
		baseDir = "application/internal/domain"
	}

	targetDir := filepath.Join(baseDir, name)
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		fmt.Printf("Error creating directory %s: %v\n", targetDir, err)
		return
	}

	// 1. Generate Entity
	entityFile := filepath.Join(targetDir, name+".go")
	generateFile(entityFile, entityTemplate, map[string]string{
		"PackageName": name,
		"EntityName":  name,
	})

	// 2. Generate Repository
	repoFile := filepath.Join(targetDir, "repository.go")
	generateFile(repoFile, repoTemplate, map[string]string{
		"PackageName": name,
		"EntityName":  name,
	})

	// 3. Generate Mapper (SQL Support)
	mapperFile := filepath.Join(targetDir, "mapper.go")
	generateFile(mapperFile, mapperTemplate, map[string]string{
		"PackageName": name,
		"EntityName":  name,
	})

	// 4. Generate Infrastructure Implementation (Persistence)
	// Try to find infrastructure/persistence directory
	infraDir := "../../application/internal/infrastructure/persistence"
	if _, err := os.Stat(infraDir); os.IsNotExist(err) {
		infraDir = "application/internal/infrastructure/persistence"
	}
	// Create if not exists
	_ = os.MkdirAll(infraDir, 0755)

	// Repo Impl
	repoImplFile := filepath.Join(infraDir, name+"_repo.go")
	generateFile(repoImplFile, repoImplTemplate, map[string]string{
		"PackageName": name,
		"EntityName":  name,
	})

	// Mapper Impl
	mapperImplFile := filepath.Join(infraDir, name+"_mapper.go")
	generateFile(mapperImplFile, mapperImplTemplate, map[string]string{
		"PackageName": name,
		"EntityName":  name,
	})
}

// generateFile creates a file from template if it doesn't strictly exist (Lock mechanism)
func generateFile(path string, tmpl string, data interface{}) {
	if _, err := os.Stat(path); err == nil {
		fmt.Printf("[LOCK] Skipping %s: file already exists\n", path)
		return
	}

	f, err := os.Create(path)
	if err != nil {
		fmt.Printf("Error creating file %s: %v\n", path, err)
		return
	}
	defer f.Close()

	t := template.Must(template.New("file").Parse(tmpl))
	if err := t.Execute(f, data); err != nil {
		fmt.Printf("Error executing template for %s: %v\n", path, err)
	} else {
		fmt.Printf("[CREATED] %s\n", path)
	}
}

const entityTemplate = `package {{.PackageName}}

import "github.com/soliton-go/framework/ddd"

type {{.EntityName}}ID string

func (id {{.EntityName}}ID) String() string {
	return string(id)
}

type {{.EntityName}} struct {
	ddd.BaseAggregateRoot
	ID {{.EntityName}}ID ` + "`gorm:\"primaryKey\"`" + `
}

func New{{.EntityName}}(id string) *{{.EntityName}} {
	return &{{.EntityName}}{
		ID: {{.EntityName}}ID(id),
	}
}

func (e *{{.EntityName}}) GetID() ddd.ID {
	return e.ID
}
`

const repoTemplate = `package {{.PackageName}}

import (
	"github.com/soliton-go/framework/orm"
)

type {{.EntityName}}Repository interface {
	orm.Repository[*{{.EntityName}}, {{.EntityName}}ID]
}
`

const mapperTemplate = `package {{.PackageName}}

import (
	"github.com/soliton-go/framework/orm"
)

type {{.EntityName}}Mapper interface {
	orm.SQLMapper[{{.EntityName}}]
}
`

const repoImplTemplate = `package persistence

import (
	"github.com/soliton-go/application/internal/domain/{{.PackageName}}"
	"github.com/soliton-go/framework/orm"
	"gorm.io/gorm"
)

type {{.EntityName}}RepoImpl struct {
	*orm.GormRepository[*{{.PackageName}}.{{.EntityName}}, {{.PackageName}}.{{.EntityName}}ID]
}

func New{{.EntityName}}Repository(db *gorm.DB) {{.PackageName}}.{{.EntityName}}Repository {
	return &{{.EntityName}}RepoImpl{
		GormRepository: orm.NewGormRepository[*{{.PackageName}}.{{.EntityName}}, {{.PackageName}}.{{.EntityName}}ID](db),
	}
}
`

const mapperImplTemplate = `package persistence

import (
	"github.com/soliton-go/application/internal/domain/{{.PackageName}}"
	"github.com/soliton-go/framework/orm"
	"gorm.io/gorm"
)

type {{.EntityName}}MapperImpl struct {
	*orm.GormMapper[{{.PackageName}}.{{.EntityName}}]
}

func New{{.EntityName}}Mapper(db *gorm.DB) {{.PackageName}}.{{.EntityName}}Mapper {
	return &{{.EntityName}}MapperImpl{
		GormMapper: orm.NewGormMapper[{{.PackageName}}.{{.EntityName}}](db),
	}
}
`
