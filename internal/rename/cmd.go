package rename

import (
	"handytools/pkg/common"

	"github.com/spf13/cobra"
)

type RenameConfig struct {
	OutputName string
	InputFiles []string
	Apply      bool
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

		config.InputFiles = common.ExpandWildcards(args)
		logger.Infof("Running rename with config: %+v\n", config)
		renameFiles(config)
	},
}

func init() {
	Cmd.Flags().StringVarP(&config.OutputName, "output", "o", "", "Base name for renaming files")
	Cmd.Flags().BoolVarP(&config.Apply, "apply", "a", false, "Apply changes (default is dry-run)")
}
