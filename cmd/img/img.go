package main

import (
	"fmt"
	"os"

	"handytools/internal/batchrename"
	"handytools/internal/collage"
	"handytools/internal/grab"
	"handytools/internal/imgui"
	"handytools/internal/optimise"
	"handytools/internal/rename"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "img",
	Short: "A handy tool for image processing",
	Long:  "img allows you to create collages, optimise images, and explore files.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Use 'img --help' to see available commands.")
	},
}

func init() {
	rootCmd.AddCommand(batchrename.Cmd)
	rootCmd.AddCommand(collage.Cmd)
	rootCmd.AddCommand(grab.Cmd)
	rootCmd.AddCommand(optimise.Cmd)
	rootCmd.AddCommand(rename.Cmd)

	rootCmd.AddCommand(&cobra.Command{
		Use:   "exploreui",
		Short: "Launch GUI file explorer",
		Run: func(cmd *cobra.Command, args []string) {
			imgui.Run()
		},
	})
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
