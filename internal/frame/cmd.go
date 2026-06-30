package frame

import (
	"handytools/pkg/common"

	"github.com/spf13/cobra"
)

type Config struct {
	InputFiles []string
	OutputDir  string
	FramePct   float64
	Color      string
	Torn       bool
	TornDepth  float64
}

var (
	logger = common.GetLogger()
	config Config
)

var Cmd = &cobra.Command{
	Use:   "frame [files...]",
	Short: "Add a border frame to images",
	Long: `Adds a configurable border frame to one or more images.

Accepts glob patterns (e.g. *.jpg). Output goes to -o directory;
use -o . to overwrite the originals in place.`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if config.OutputDir == "" {
			logger.Error("Output directory (-o) is required. Use '.' to overwrite originals.")
			return
		}
		config.InputFiles = common.ExpandWildcards(args)
		if len(config.InputFiles) == 0 {
			logger.Error("No matching files found.")
			return
		}
		runFrame(config)
	},
}

func init() {
	Cmd.Flags().StringVarP(&config.OutputDir, "output", "o", "", "Output directory ('.' to overwrite originals)")
	Cmd.Flags().Float64Var(&config.FramePct, "frame", 1.0, "Frame border width in % of image width (e.g. 5 = 5% border)")
	Cmd.Flags().StringVar(&config.Color, "color", "white", "Frame color: white, black, cream, or #RRGGBB")
	Cmd.Flags().BoolVar(&config.Torn, "torn", false, "Torn-edge effect on inner frame border")
	Cmd.Flags().Float64Var(&config.TornDepth, "torn-depth", 50.0, "Tear depth as percentage of frame width")
}
