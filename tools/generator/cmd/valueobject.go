package cmd

import (
	"fmt"
	"os"

	"github.com/soliton-go/tools/core"
	"github.com/spf13/cobra"
)

var valueObjectFieldsFlag string

var valueObjectCmd = &cobra.Command{
	Use:   "valueobject [domain] [name]",
	Short: "Generate a value object in domain layer",
	Long: `Generate a value object for a specific domain.

Examples:
  soliton-gen valueobject user EmailAddress --fields "value:string"
  soliton-gen valueobject order Money --fields "amount:decimal,currency:string"`,
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		domain := args[0]
		name := args[1]

		fields := core.ParseFields(valueObjectFieldsFlag)
		cfg := core.ValueObjectConfig{
			Domain: domain,
			Name:   name,
			Fields: fields,
			Force:  forceFlag,
		}

		if err := core.ValidateValueObjectConfig(cfg); err != nil {
			fmt.Printf("❌ 配置错误: %v\n", err)
			os.Exit(1)
		}

		result, err := core.GenerateValueObject(cfg)
		if err != nil {
			fmt.Printf("❌ 生成失败: %v\n", err)
			os.Exit(1)
		}

		printGenerationResult(result)
	},
}

func init() {
	rootCmd.AddCommand(valueObjectCmd)
	valueObjectCmd.Flags().StringVarP(&valueObjectFieldsFlag, "fields", "f", "", "Comma-separated list of fields (e.g., 'amount:decimal,currency:string')")
	valueObjectCmd.Flags().BoolVar(&forceFlag, "force", false, "Force overwrite existing files")
}
