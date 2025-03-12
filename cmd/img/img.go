package main

import (
	"fmt"
	"os"

	"handytools/internal/batchrename"
	"handytools/internal/collage"
	"handytools/internal/optimise"
	"handytools/internal/rename"

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
	rootCmd.AddCommand(batchrename.Cmd)
	rootCmd.AddCommand(collage.Cmd)
	rootCmd.AddCommand(optimise.Cmd)
	rootCmd.AddCommand(rename.Cmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
