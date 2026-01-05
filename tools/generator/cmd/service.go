package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/soliton-go/tools/core"
	"github.com/spf13/cobra"
)

var serviceMethods string

var serviceCmd = &cobra.Command{
	Use:   "service [name]",
	Short: "Generate an application service",
	Long: `Generate an application service for cross-domain business logic.
Services orchestrate multiple domains to implement complex use cases.

Examples:
  soliton-gen service OrderService
  soliton-gen service OrderService --methods "CreateOrder,CancelOrder,GetUserOrders"
  soliton-gen service PaymentService --methods "ProcessPayment,RefundPayment"`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		fmt.Printf("🚀 Generating service: %s\n\n", name)

		// Smart detection: show service type info
		detection, err := core.DetectServiceType(name)
		if err == nil {
			if detection.DomainExists {
				fmt.Printf("📌 %s\n", detection.Message)
				fmt.Printf("   类型: 领域服务 (Domain Service)\n")
			} else {
				fmt.Printf("📌 %s\n", detection.Message)
				fmt.Printf("   类型: 跨领域服务 (Cross-domain Service)\n")
			}
			fmt.Printf("   目标: %s/service.go\n", detection.TargetDir)
			if detection.ShouldReuseDTO {
				fmt.Println("   ✓ 复用现有 DTO")
			}
			fmt.Println()
		}

		// Parse methods from comma-separated string
		var methods []string
		if serviceMethods != "" {
			for _, m := range strings.Split(serviceMethods, ",") {
				m = strings.TrimSpace(m)
				if m != "" {
					methods = append(methods, m)
				}
			}
		}

		// Create service configuration
		cfg := core.ServiceConfig{
			Name:    name,
			Methods: methods,
			Force:   forceFlag,
		}

		// Validate configuration
		if err := core.ValidateServiceConfig(cfg); err != nil {
			fmt.Printf("❌ 配置错误: %v\n", err)
			os.Exit(1)
		}

		// Generate service using core package
		result, err := core.GenerateService(cfg)
		if err != nil {
			fmt.Printf("❌ 生成失败: %v\n", err)
			os.Exit(1)
		}

		// Print result
		printGenerationResult(result)
	},
}

func init() {
	rootCmd.AddCommand(serviceCmd)
	serviceCmd.Flags().StringVar(&serviceMethods, "methods", "", "Comma-separated list of service methods (e.g., 'CreateOrder,CancelOrder')")
	serviceCmd.Flags().BoolVar(&forceFlag, "force", false, "Force overwrite existing files")

	// Add subcommands
	serviceCmd.AddCommand(serviceListCmd)
	serviceCmd.AddCommand(serviceDeleteCmd)
}

// serviceListCmd lists all services
var serviceListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all application services in the project",
	Run: func(cmd *cobra.Command, args []string) {
		services, err := core.ListServices()
		if err != nil {
			fmt.Printf("❌ 错误: %v\n", err)
			os.Exit(1)
		}

		if len(services) == 0 {
			fmt.Println("未检测到应用服务")
			return
		}

		fmt.Printf("已检测到 %d 个应用服务：\n\n", len(services))
		for _, s := range services {
			fmt.Printf("  • %s\n", s)
		}
	},
}

// serviceDeleteCmd deletes a service
var serviceDeleteCmd = &cobra.Command{
	Use:   "delete [name]",
	Short: "Delete an application service and related files",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		serviceName := args[0]

		// Confirm deletion unless --force is set
		if !forceFlag {
			fmt.Printf("⚠️  警告: 即将删除服务 '%s' 及其相关文件\n", serviceName)
			fmt.Print("确认删除? (y/N): ")
			reader := bufio.NewReader(os.Stdin)
			response, _ := reader.ReadString('\n')
			response = strings.TrimSpace(strings.ToLower(response))
			if response != "y" && response != "yes" {
				fmt.Println("已取消删除")
				return
			}
		}

		result := core.DeleteService(serviceName)
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
