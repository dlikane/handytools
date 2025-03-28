package optimise

import (
	"github.com/spf13/cobra"
	"handytools/pkg/common"
)

type Config struct {
	InputFiles []string
	Profile    string
	Apply      bool
	Stat       bool
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
		if config.Apply && config.Stat {
			logger.Error("Cannot use --apply and --stat together")
			return
		}

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
	Cmd.Flags().BoolVarP(&config.Stat, "stat", "s", false, "Collect image profile statistics (mutually exclusive with --apply)")
	Cmd.Flags().StringVarP(&config.Profile, "profile", "p", "insta", "Profile size (x-small, small, med, large, x-large, insta)")
}
