package pinterest

import (
	"handytools/pkg/common"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

type Config struct {
	BoardURL string
	Output   string
	Layout   string
}

var config Config
var logger = common.GetLogger()

var Cmd = &cobra.Command{
	Use:   "pinterest",
	Short: "Download pins from a Pinterest board and assemble into image(s)",
	Long:  "Downloads images from a Pinterest board and assembles them into collages with specified layout: compact, medium, large, fit.",
	Run: func(cmd *cobra.Command, args []string) {
		if config.BoardURL == "" {
			logger.Error("Pinterest board URL is required")
			return
		}

		logger.Infof("Fetching board: %s", config.BoardURL)

		tempDir := filepath.Join(os.TempDir(), "pinterest_pins")
		os.MkdirAll(tempDir, 0755)

		imagePaths, err := DownloadPins(config.BoardURL, tempDir)
		if err != nil {
			logger.WithError(err).Error("Failed to download pins")
			return
		}

		if len(imagePaths) == 0 {
			logger.Warn("No images downloaded")
			return
		}

		logger.Infof("Downloaded %d pins", len(imagePaths))

		output := strings.TrimSuffix(config.Output, ".jpg")
		err = AssembleImages(imagePaths, config.Layout, output)
		if err != nil {
			logger.WithError(err).Error("Failed to assemble collage")
		}
	},
}

func init() {
	Cmd.Flags().StringVarP(&config.BoardURL, "url", "u", "", "Pinterest board URL (e.g., https://pin.it/7kvMjAV3t)")
	Cmd.Flags().StringVarP(&config.Output, "output", "o", "pinterest.jpg", "Output image path prefix (e.g., output/pin.jpg → output/pin_01.jpg etc)")
	Cmd.Flags().StringVarP(&config.Layout, "layout", "l", "fit", "Layout mode: compact (6x), medium (3x), large (1x), fit (auto-fit into 1080x1920)")
}
