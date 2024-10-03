package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/sshaparenko/quartz/internal/git"
)

func init() {
	rootCmd.AddCommand(track)
}

var rootCmd = &cobra.Command{
	Use:     "quartz",
	Long:    "Quartz is a synchronization tool for your Obsidian vaults",
	Args:    cobra.MinimumNArgs(1),
	Version: "1.0.0-beta.1",
	CompletionOptions: cobra.CompletionOptions{
		DisableDefaultCmd: true,
	},
}

var track = &cobra.Command{
	Use:     "track",
	Aliases: []string{"t"},
	Short:   "Tracks obsidian vault",
	Long:    "Tracks obsidian vault via specified path",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := git.Init(args[0]); err != nil {
			fmt.Printf("Failed to initialize Quartz: %s\n", err.Error())
			os.Exit(1)
		}
	},
}

func Execute() error {
	if err := rootCmd.Execute(); err != nil {
		return fmt.Errorf("in cmd.Execute: %w", err)
	}
	return nil
}
