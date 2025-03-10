package main

import (
	"fmt"
	"os"

	"handytools/internal/collage"
	"handytools/internal/optimise"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "img",
	Short: "A handy tool for image processing",
	Long:  "img allows you to create collages and optimise images.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Use 'img --help' to see available commands.")
	},
}

func main() {
	// Add subcommands
	rootCmd.AddCommand(collage.Cmd)
	rootCmd.AddCommand(optimise.Cmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
