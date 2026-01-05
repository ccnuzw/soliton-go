package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/soliton-go/tools/core"
	"github.com/spf13/cobra"
)

var fieldsFlag string
var forceFlag bool
var tableNameFlag string
var routeBaseFlag string
var wireFlag bool
var softDeleteFlag bool

var domainCmd = &cobra.Command{
	Use:   "domain [name]",
	Short: "Generate a complete domain module ready for production",
	Long: `Generate a complete domain entity including:
  - Entity with ID and aggregate root
  - Repository interface and GORM implementation
  - Domain events (Created, Updated, Deleted)
  - Application layer (Commands, Queries, DTOs)
  - HTTP Handler with CRUD endpoints
  - Fx dependency injection module
  - Database migration support

Examples:
  soliton-gen domain User
  soliton-gen domain User --fields "username,email,status:enum(active|inactive)"
  soliton-gen domain User --fields "username:string:用户名,email::邮箱" --wire
  soliton-gen domain User --fields "..." --force  # Overwrite existing files`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		if forceFlag {
			fmt.Printf("🚀 Generating domain: %s (force mode)\n\n", name)
		} else {
			fmt.Printf("🚀 Generating domain: %s\n\n", name)
		}

		// Parse fields from CLI format to core.FieldConfig
		fields := parseFieldsToConfig(fieldsFlag)

		// Create domain configuration
		cfg := core.DomainConfig{
			Name:       name,
			Fields:     fields,
			TableName:  tableNameFlag,
			RouteBase:  routeBaseFlag,
			SoftDelete: softDeleteFlag,
			Wire:       wireFlag,
			Force:      forceFlag,
		}

		// Validate configuration
		if err := core.ValidateDomainConfig(cfg); err != nil {
			fmt.Printf("❌ 配置错误: %v\n", err)
			os.Exit(1)
		}

		// Generate domain using core package
		result, err := core.GenerateDomain(cfg)
		if err != nil {
			fmt.Printf("❌ 生成失败: %v\n", err)
			os.Exit(1)
		}

		// Print result
		printGenerationResult(result)
	},
}

func init() {
	rootCmd.AddCommand(domainCmd)
	domainCmd.Flags().StringVarP(&fieldsFlag, "fields", "f", "", "Comma-separated list of fields (e.g., 'name,email,status:enum(active|inactive)')")
	domainCmd.Flags().BoolVar(&forceFlag, "force", false, "Force overwrite existing files")
	domainCmd.Flags().StringVar(&tableNameFlag, "table", "", "Override database table name")
	domainCmd.Flags().StringVar(&routeBaseFlag, "route", "", "Override route base path (e.g., users)")
	domainCmd.Flags().BoolVar(&wireFlag, "wire", false, "Auto-wire module into main.go (requires init template structure)")
	domainCmd.Flags().BoolVar(&softDeleteFlag, "soft-delete", false, "Enable soft delete (adds deleted_at field)")

	// Add subcommands
	domainCmd.AddCommand(domainListCmd)
	domainCmd.AddCommand(domainDeleteCmd)
}

// domainListCmd lists all domains
var domainListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all domain models in the project",
	Run: func(cmd *cobra.Command, args []string) {
		domains, err := core.ListDomains()
		if err != nil {
			fmt.Printf("❌ 错误: %v\n", err)
			os.Exit(1)
		}

		if len(domains) == 0 {
			fmt.Println("未检测到领域模型")
			return
		}

		fmt.Printf("已检测到 %d 个领域模型：\n\n", len(domains))
		for _, d := range domains {
			fmt.Printf("  • %s\n", d)
		}
	},
}

// domainDeleteCmd deletes a domain
var domainDeleteCmd = &cobra.Command{
	Use:   "delete [name]",
	Short: "Delete a domain model and all related files",
	Long: `Delete a domain model and all related files including:
  - internal/domain/<name>/
  - internal/application/<name>/
  - internal/infrastructure/persistence/<name>_repo.go
  - internal/interfaces/http/<name>_handler.go
  - main.go injections`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		domainName := args[0]

		// Confirm deletion unless --force is set
		if !forceFlag {
			fmt.Printf("⚠️  警告: 即将删除领域 '%s' 及其所有相关文件\n", domainName)
			fmt.Print("确认删除? (y/N): ")
			reader := bufio.NewReader(os.Stdin)
			response, _ := reader.ReadString('\n')
			response = strings.TrimSpace(strings.ToLower(response))
			if response != "y" && response != "yes" {
				fmt.Println("已取消删除")
				return
			}
		}

		result := core.DeleteDomain(domainName)
		if result.Success {
			fmt.Printf("✅ %s\n\n", result.Message)
			fmt.Println("已删除:")
			for _, item := range result.DeletedItems {
				fmt.Printf("  • %s\n", item)
			}
		} else {
			fmt.Printf("❌ %s\n", result.Message)
			if len(result.Errors) > 0 {
				fmt.Println("错误:")
				for _, e := range result.Errors {
					fmt.Printf("  • %s\n", e)
				}
			}
			os.Exit(1)
		}
	},
}

// parseFieldsToConfig parses CLI field format to core.FieldConfig slice
// Format: name:type:comment (type and comment are optional)
// Examples: "username", "price:int64", "email::邮箱", "status:enum(active|inactive):状态"
func parseFieldsToConfig(fieldsStr string) []core.FieldConfig {
	if fieldsStr == "" {
		return nil
	}

	var fields []core.FieldConfig
	parts := strings.Split(fieldsStr, ",")

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		field := parseFieldDefinition(part)
		fields = append(fields, field)
	}

	return fields
}

// parseFieldDefinition parses a single field definition
// Format: name:type:comment (comment is optional)
// Examples: "user_id:uuid:用户ID", "order_no::订单号", "name:string"
func parseFieldDefinition(def string) core.FieldConfig {
	var fieldName, fieldType, comment string

	// Check if it contains enum(...)
	if enumIdx := strings.Index(def, "enum("); enumIdx != -1 {
		// Find the closing parenthesis
		closeIdx := strings.LastIndex(def, ")")
		if closeIdx > enumIdx {
			// Extract enum part
			beforeEnum := def[:enumIdx]
			enumPart := def[enumIdx : closeIdx+1]
			afterEnum := ""
			if closeIdx+1 < len(def) {
				afterEnum = def[closeIdx+1:]
			}

			// Parse name from before enum (removing trailing colon if any)
			fieldName = strings.TrimSuffix(beforeEnum, ":")
			fieldType = "enum"

			// Extract enum values
			enumContent := enumPart[5 : len(enumPart)-1] // Remove "enum(" and ")"
			enumValues := parseEnumValues(enumContent)

			// Parse comment from after enum (after colon if any)
			if strings.HasPrefix(afterEnum, ":") {
				comment = strings.TrimPrefix(afterEnum, ":")
			}

			return core.FieldConfig{
				Name:       fieldName,
				Type:       fieldType,
				Comment:    comment,
				EnumValues: enumValues,
			}
		}
	}

	// Regular parsing: name:type:comment
	colonParts := strings.SplitN(def, ":", 3)
	fieldName = colonParts[0]
	if len(colonParts) >= 2 {
		fieldType = colonParts[1]
	}
	if len(colonParts) >= 3 {
		comment = colonParts[2]
	}
	if fieldType == "" {
		fieldType = "string"
	}

	return core.FieldConfig{
		Name:    fieldName,
		Type:    fieldType,
		Comment: comment,
	}
}

// parseEnumValues parses enum values from a pipe-separated string
func parseEnumValues(raw string) []string {
	parts := strings.Split(raw, "|")
	values := make([]string, 0, len(parts))
	for _, part := range parts {
		value := strings.TrimSpace(part)
		if value == "" {
			continue
		}
		values = append(values, value)
	}
	return values
}

// printGenerationResult prints the generation result in CLI-friendly format
func printGenerationResult(result *core.GenerationResult) {
	if !result.Success {
		fmt.Println("❌ 生成失败")
		for _, e := range result.Errors {
			fmt.Printf("  • %s\n", e)
		}
		os.Exit(1)
	}

	fmt.Println("✅ 生成成功！\n")
	fmt.Println("已生成文件:")
	for _, f := range result.Files {
		status := ""
		switch f.Status {
		case core.FileStatusNew:
			status = "[NEW]"
		case core.FileStatusOverwrite:
			status = "[OVERWRITE]"
		case core.FileStatusSkip:
			status = "[SKIP]"
		case core.FileStatusError:
			status = "[ERROR]"
		}
		fmt.Printf("  %s %s\n", status, f.Path)
	}

	if result.Message != "" {
		fmt.Printf("\n%s\n", result.Message)
	}
}
