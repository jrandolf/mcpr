package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "mcpr",
	Short: "MCP Registry - Manage MCP servers across clients",
	Long: `mcpr is a CLI tool for managing Model Context Protocol (MCP) servers.

It allows you to:
  - Add MCP server configurations
  - Install servers to various MCP clients (Claude Desktop, Claude Code, Cursor, Windsurf)
  - Manage your MCP server configurations in a central location`,
}

// Execute runs the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(removeCmd)
	rootCmd.AddCommand(clientCmd)
	rootCmd.AddCommand(listCmd)
}
