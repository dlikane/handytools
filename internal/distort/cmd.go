package distort

import (
	"path/filepath"
	"strings"

	"handytools/pkg/common"

	"github.com/spf13/cobra"
)

type Config struct {
	InputFile  string
	OutputFile string
	Mode       string
	Intensity  float64
	Seed       int64
}

var (
	logger = common.GetLogger()
	config Config
)

var Cmd = &cobra.Command{
	Use:   "distort [file]",
	Short: "Apply databending distortion to an image",
	Long: `Deliberately corrupts JPEG data to produce glitch artifacts (databending/datamosh).

Modes:
  corrupt  Randomly flips bytes in the JPEG scan data, causing Huffman decode
           cascade errors — the image-file equivalent of I-frame removal.
  shift    Displaces rows of pixels horizontally by random amounts (CRT scan error look).
  melt     Copies chunks of scan data forward, smearing earlier content over later rows.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		config.InputFile = args[0]
		if config.OutputFile == "" {
			ext := filepath.Ext(config.InputFile)
			base := strings.TrimSuffix(config.InputFile, ext)
			config.OutputFile = base + "_distorted" + ext
		}
		runDistort(config)
	},
}

func init() {
	Cmd.Flags().StringVarP(&config.OutputFile, "output", "o", "", "Output file (default: <input>_distorted.jpg)")
	Cmd.Flags().StringVarP(&config.Mode, "mode", "m", "corrupt", "Distortion mode: corrupt | shift | melt")
	Cmd.Flags().Float64Var(&config.Intensity, "intensity", 0.05, "Distortion intensity 0.0–1.0")
	Cmd.Flags().Int64Var(&config.Seed, "seed", 0, "Random seed for reproducibility (0 = random)")
}
