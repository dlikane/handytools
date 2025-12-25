package rename

import (
	"strings"

	"handytools/pkg/common"

	"github.com/spf13/cobra"
)

type RenameConfig struct {
	OutputName string
	InputFiles []string
	Apply      bool
	SortBy     string // created, modified, name, random
	Order      string // asc, desc
}

var (
	logger = common.GetLogger()
	config RenameConfig
)

var Cmd = &cobra.Command{
	Use:   "rename",
	Short: "Rename files with a given base name and counter",
	Long:  "Renames files using the provided output name followed by a counter (e.g., base_0001, base_0002).",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			logger.Error("No files provided. Use wildcard or file list.")
			return
		}
		if config.OutputName == "" {
			logger.Error("Output name (-o) is required.")
			return
		}

		// Normalize and validate flags
		config.SortBy = strings.ToLower(strings.TrimSpace(config.SortBy))
		if config.SortBy == "" {
			config.SortBy = "created"
		}
		switch config.SortBy {
		case "created", "modified", "name", "random":
			// ok
		default:
			logger.Errorf("Invalid --sort value: %s (use 'created', 'modified', 'name', or 'random')", config.SortBy)
			return
		}

		config.Order = strings.ToLower(strings.TrimSpace(config.Order))
		if config.Order == "" {
			config.Order = "asc"
		}
		switch config.Order {
		case "asc", "desc":
			// ok
		default:
			logger.Errorf("Invalid --order value: %s (use 'asc' or 'desc')", config.Order)
			return
		}

		config.InputFiles = common.ExpandWildcards(args)
		logger.Infof("Running rename with config: %+v\n", config)
		renameFiles(config)
	},
}

func init() {
	Cmd.Flags().StringVarP(&config.OutputName, "output", "o", "", "Base name for renaming files")
	Cmd.Flags().BoolVarP(&config.Apply, "apply", "a", false, "Apply changes (default is dry-run)")
	Cmd.Flags().StringVarP(&config.SortBy, "sort", "s", "created", "Sort input files by: created | modified | name | random")
	Cmd.Flags().StringVarP(&config.Order, "order", "r", "asc", "Sort order: asc | desc")
}
