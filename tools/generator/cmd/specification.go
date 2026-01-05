package cmd

import (
	"fmt"
	"os"

	"github.com/soliton-go/tools/core"
	"github.com/spf13/cobra"
)

var specTargetFlag string

var specificationCmd = &cobra.Command{
	Use:   "spec [domain] [name]",
	Short: "Generate a specification in domain layer",
	Long: `Generate a domain specification (Specification).

Examples:
  soliton-gen spec user ActiveUserSpec --target User
  soliton-gen spec order PaidOrderSpec --target Order`,
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		domain := args[0]
		name := args[1]

		cfg := core.SpecificationConfig{
			Domain: domain,
			Name:   name,
			Target: specTargetFlag,
			Force:  forceFlag,
		}

		if err := core.ValidateSpecificationConfig(cfg); err != nil {
			fmt.Printf("❌ 配置错误: %v\n", err)
			os.Exit(1)
		}

		result, err := core.GenerateSpecification(cfg)
		if err != nil {
			fmt.Printf("❌ 生成失败: %v\n", err)
			os.Exit(1)
		}

		printGenerationResult(result)
	},
}

func init() {
	rootCmd.AddCommand(specificationCmd)
	specificationCmd.Flags().StringVar(&specTargetFlag, "target", "", "Target type (e.g., User)")
	specificationCmd.Flags().BoolVar(&forceFlag, "force", false, "Force overwrite existing files")
}
