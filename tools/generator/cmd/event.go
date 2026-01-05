package cmd

import (
	"fmt"
	"os"

	"github.com/soliton-go/tools/core"
	"github.com/spf13/cobra"
)

var eventFieldsFlag string
var eventTopicFlag string

var eventCmd = &cobra.Command{
	Use:   "event [domain] [name]",
	Short: "Generate a domain event",
	Long: `Generate a custom domain event with registration.

Examples:
  soliton-gen event user UserActivated --fields "user_id:uuid"
  soliton-gen event order OrderPaid --fields "order_id:uuid,amount:decimal" --topic "order.paid"`,
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		domain := args[0]
		name := args[1]

		fields := core.ParseFieldsAllowReserved(eventFieldsFlag)
		cfg := core.EventConfig{
			Domain: domain,
			Name:   name,
			Topic:  eventTopicFlag,
			Fields: fields,
			Force:  forceFlag,
		}

		if err := core.ValidateEventConfig(cfg); err != nil {
			fmt.Printf("❌ 配置错误: %v\n", err)
			os.Exit(1)
		}

		result, err := core.GenerateEvent(cfg)
		if err != nil {
			fmt.Printf("❌ 生成失败: %v\n", err)
			os.Exit(1)
		}

		printGenerationResult(result)
	},
}

func init() {
	rootCmd.AddCommand(eventCmd)
	eventCmd.Flags().StringVarP(&eventFieldsFlag, "fields", "f", "", "Comma-separated list of fields (e.g., 'user_id:uuid,amount:decimal')")
	eventCmd.Flags().StringVar(&eventTopicFlag, "topic", "", "Event topic (e.g., 'user.activated')")
	eventCmd.Flags().BoolVar(&forceFlag, "force", false, "Force overwrite existing files")
}
