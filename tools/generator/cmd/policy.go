package cmd

import (
	"fmt"
	"os"

	"github.com/soliton-go/tools/core"
	"github.com/spf13/cobra"
)

var policyTargetFlag string

var policyCmd = &cobra.Command{
	Use:   "policy [domain] [name]",
	Short: "Generate a policy in domain layer",
	Long: `Generate a domain policy (Policy).

Examples:
  soliton-gen policy user PasswordPolicy --target User
  soliton-gen policy order RefundPolicy --target Order`,
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		domain := args[0]
		name := args[1]

		cfg := core.PolicyConfig{
			Domain: domain,
			Name:   name,
			Target: policyTargetFlag,
			Force:  forceFlag,
		}

		if err := core.ValidatePolicyConfig(cfg); err != nil {
			fmt.Printf("❌ 配置错误: %v\n", err)
			os.Exit(1)
		}

		result, err := core.GeneratePolicy(cfg)
		if err != nil {
			fmt.Printf("❌ 生成失败: %v\n", err)
			os.Exit(1)
		}

		printGenerationResult(result)
	},
}

func init() {
	rootCmd.AddCommand(policyCmd)
	policyCmd.Flags().StringVar(&policyTargetFlag, "target", "", "Target type (e.g., User)")
	policyCmd.Flags().BoolVar(&forceFlag, "force", false, "Force overwrite existing files")
}
