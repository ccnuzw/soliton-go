package cmd

import (
	"fmt"

	"github.com/soliton-go/tools/server"
	"github.com/spf13/cobra"
)

var servePort int
var serveHost string

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the Web GUI server",
	Long: `Start a web server that provides a GUI for the code generator.

The web interface allows you to:
  - Initialize new projects with a wizard
  - Generate domain modules with a visual field editor
  - Generate application services
  - Preview generated code before writing

Examples:
  soliton-gen serve                    # Start on default port 3000
  soliton-gen serve --port 8080        # Start on custom port
  soliton-gen serve --host 0.0.0.0     # Listen on all interfaces`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("🚀 Starting Soliton-Gen Web GUI\n")
		fmt.Printf("   URL: http://%s:%d\n\n", serveHost, servePort)
		server.Start(serveHost, servePort)
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
	serveCmd.Flags().IntVarP(&servePort, "port", "p", 3000, "Server port")
	serveCmd.Flags().StringVar(&serveHost, "host", "127.0.0.1", "Server host")
}
