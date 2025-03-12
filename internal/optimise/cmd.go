package optimise

import (
	"github.com/spf13/cobra"
	"handytools/pkg/common"
)

type Config struct {
	InputFiles []string
	Apply      bool
}

var (
	logger = common.GetLogger()
	config Config
)

var Cmd = &cobra.Command{
	Use:   "optimise",
	Short: "Optimise image file size",
	Long:  "Optimises images by reducing file size while maintaining quality.",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			logger.Error("No images provided. Use wildcard or file list.")
			return
		}
		config.InputFiles = common.ExpandWildcards(args)
		logger.Infof("Running optimise with config: %+v\n", config)
		optimiseImages(config)
	},
}

func init() {
	Cmd.Flags().BoolVarP(&config.Apply, "apply", "a", false, "Apply changes (default is dry-run)")
}
