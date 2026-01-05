package cmd

import (
	"fmt"
	"os"

	"github.com/soliton-go/tools/core"
	"github.com/spf13/cobra"
)

var handlerTopicFlag string

var eventHandlerCmd = &cobra.Command{
	Use:   "event-handler [domain] [event]",
	Short: "Generate a domain event handler",
	Long: `Generate an event handler and register it into the module.

Examples:
  soliton-gen event-handler user UserCreated
  soliton-gen event-handler order OrderPaid --topic "order.paid"`,
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		domain := args[0]
		eventName := args[1]

		cfg := core.EventHandlerConfig{
			Domain:    domain,
			EventName: eventName,
			Topic:     handlerTopicFlag,
			Force:     forceFlag,
		}

		if err := core.ValidateEventHandlerConfig(cfg); err != nil {
			fmt.Printf("❌ 配置错误: %v\n", err)
			os.Exit(1)
		}

		result, err := core.GenerateEventHandler(cfg)
		if err != nil {
			fmt.Printf("❌ 生成失败: %v\n", err)
			os.Exit(1)
		}

		printGenerationResult(result)
	},
}

func init() {
	rootCmd.AddCommand(eventHandlerCmd)
	eventHandlerCmd.Flags().StringVar(&handlerTopicFlag, "topic", "", "Event topic (e.g., 'user.created')")
	eventHandlerCmd.Flags().BoolVar(&forceFlag, "force", false, "Force overwrite existing files")
}
