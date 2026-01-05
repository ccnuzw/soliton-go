package cmd

import (
	"fmt"
	"os"

	"github.com/soliton-go/tools/core"
	"github.com/spf13/cobra"
)

var moduleName string
var frameworkVersion string
var frameworkReplace string

var initCmd = &cobra.Command{
	Use:   "init [project-name]",
	Short: "Initialize a new Soliton-Go project",
	Long: `Initialize a new Soliton-Go project with complete directory structure:
  - cmd/main.go (Entry point with Fx setup)
  - configs/config.yaml (Configuration)
  - internal/ (DDD layer structure)
  - go.mod (Dependencies)

Examples:
  soliton-gen init my-project
  soliton-gen init my-project --module github.com/myorg/my-project`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		projectName := args[0]
		if moduleName == "" {
			moduleName = core.GetDefaultModuleName(projectName)
		}
		fmt.Printf("🚀 Initializing project: %s\n", projectName)
		fmt.Printf("   Module: %s\n\n", moduleName)

		// Create project configuration
		cfg := core.ProjectConfig{
			Name:             projectName,
			ModuleName:       moduleName,
			FrameworkVersion: frameworkVersion,
			FrameworkReplace: frameworkReplace,
		}

		// Validate configuration
		if err := core.ValidateProjectConfig(cfg); err != nil {
			fmt.Printf("❌ 配置错误: %v\n", err)
			os.Exit(1)
		}

		// Initialize project using core package
		result, err := core.InitProject(cfg)
		if err != nil {
			fmt.Printf("❌ 初始化失败: %v\n", err)
			os.Exit(1)
		}

		// Print result
		printGenerationResult(result)

		fmt.Println("\n📦 下一步:")
		fmt.Printf("   cd %s\n", projectName)
		fmt.Println("   go mod tidy")
		fmt.Println("   go run ./cmd/main.go")
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
	initCmd.Flags().StringVarP(&moduleName, "module", "m", "", "Go module name (default: github.com/soliton-go/<project-name>)")
	initCmd.Flags().StringVar(&frameworkVersion, "framework-version", "", "Framework version (default: auto)")
	initCmd.Flags().StringVar(&frameworkReplace, "framework-replace", "", "Replace github.com/soliton-go/framework with a local path")
}
