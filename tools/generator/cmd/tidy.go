package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

var tidyCmd = &cobra.Command{
	Use:   "tidy",
	Short: "Run go mod tidy to update dependencies",
	Long: `Run go mod tidy to download missing dependencies and remove unused ones.

This command is equivalent to running:
  GOWORK=off go mod tidy`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("📦 更新依赖...")

		// Run go mod tidy with GOWORK=off
		tidyCmd := exec.Command("go", "mod", "tidy")
		tidyCmd.Env = append(os.Environ(), "GOWORK=off")
		tidyCmd.Stdout = os.Stdout
		tidyCmd.Stderr = os.Stderr

		if err := tidyCmd.Run(); err != nil {
			fmt.Printf("❌ go mod tidy 失败: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("✅ 依赖更新完成")
	},
}

func init() {
	rootCmd.AddCommand(tidyCmd)
}
