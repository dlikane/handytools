package gallery

import (
	"bufio"
	"handytools/pkg/assemble"
	"handytools/pkg/common"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

type Config struct {
	BoardURL  string
	Output    string
	Directory string
	InputList string
}

var config Config
var logger = common.GetLogger()

var Cmd = &cobra.Command{
	Use:   "gallery",
	Short: "Assemble a gallery image from a list of images or Pinterest board",
	Long:  "Creates a single or multi-page image gallery from local image paths, a directory, a Pinterest board, or a file containing image paths.",
	Run: func(cmd *cobra.Command, args []string) {
		var imagePaths []string

		switch {
		case config.BoardURL != "":
			logger.Infof("Fetching board: %s", config.BoardURL)
			tempDir := filepath.Join(os.TempDir(), "pinterest_pins")
			os.MkdirAll(tempDir, 0755)
			var err error
			imagePaths, err = DownloadPins(config.BoardURL, tempDir)
			if err != nil {
				logger.WithError(err).Error("Failed to download pins")
				return
			}
		case config.Directory != "":
			logger.Infof("Loading images from directory: %s", config.Directory)
			entries, err := os.ReadDir(config.Directory)
			if err != nil {
				logger.WithError(err).Error("Failed to read directory")
				return
			}
			for _, e := range entries {
				if !e.IsDir() && strings.HasSuffix(e.Name(), ".jpg") {
					imagePaths = append(imagePaths, filepath.Join(config.Directory, e.Name()))
				}
			}
		case config.InputList != "":
			file, err := os.Open(config.InputList)
			if err != nil {
				logger.WithError(err).Error("Failed to open file list")
				return
			}
			defer file.Close()
			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				line := strings.TrimSpace(scanner.Text())
				line = strings.Trim(line, `"`)
				if line != "" {
					imagePaths = append(imagePaths, line)
				}
			}
			if err := scanner.Err(); err != nil {
				logger.WithError(err).Error("Error reading input list")
				return
			}
		default:
			if len(args) == 0 {
				logger.Error("No input source provided. Use --pinterest, --directory, --file or pass image paths as arguments.")
				return
			}
			imagePaths = args
		}

		if len(imagePaths) == 0 {
			logger.Warn("No images to process")
			return
		}

		output := strings.TrimSuffix(config.Output, ".jpg")
		err := assemble.AssembleImagesWithMax(imagePaths, output)
		if err != nil {
			logger.WithError(err).Error("Failed to assemble gallery")
		}
	},
}

func init() {
	Cmd.Flags().StringVarP(&config.BoardURL, "pinterest", "p", "", "Pinterest board URL")
	Cmd.Flags().StringVarP(&config.Output, "output", "o", "gallery.jpg", "Output image path prefix (e.g., out/gallery_01.jpg)")
	Cmd.Flags().StringVarP(&config.Directory, "directory", "d", "", "Directory to read .jpg files from")
	Cmd.Flags().StringVarP(&config.InputList, "file", "f", "", "Text file with list of image paths")
}
