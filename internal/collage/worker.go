package collage

import (
	"handytools/pkg/common"
	"image"
	"image/color"
	"math"
	"os"

	"github.com/disintegration/imaging"
)

func scaleImages(images []image.Image, width, height int) []image.Image {
	for i, img := range images {
		images[i] = imaging.Fill(img, width, height, imaging.Center, imaging.Lanczos)
	}
	return images
}

func createGrid(rows, cols, width, height, frame int) *image.NRGBA {
	return imaging.New(width*cols+(cols-1)*frame, height*rows+(rows-1)*frame, color.NRGBA{0xff, 0xff, 0xff, 0xff})
}

func placeImagesInGrid(grid *image.NRGBA, images []image.Image, rows, cols, width, height, frame int) *image.NRGBA {
	i := 0
	for y := 0; y < rows; y++ {
		for x := 0; x < cols; x++ {
			if i >= len(images) {
				break
			}
			grid = imaging.Paste(grid, images[i], image.Pt(x*(width+frame), y*(height+frame)))
			i++
		}
	}
	return grid
}

func createCollage(cfg Config) {
	logger := common.GetLogger()

	var images []image.Image
	for _, imgPath := range cfg.InputFiles {
		logger.Info("Opening image file: " + imgPath)
		infile, err := os.Open(imgPath)
		if err != nil {
			logger.Error("Failed to open image file: " + imgPath)
			return
		}
		img, _, err := image.Decode(infile)
		infile.Close()
		if err != nil {
			logger.Error("Failed to decode image file: " + imgPath)
			return
		}
		images = append(images, img)
		logger.Info("Image file loaded: " + imgPath)
	}

	if len(images) == 0 {
		logger.Error("No input images provided")
		return
	}

	refWidth := images[0].Bounds().Dx()
	refHeight := images[0].Bounds().Dy()

	maxTotalWidth := 2 * refWidth
	maxTotalHeight := 2 * refHeight

	cellWidth := int(math.Min(float64(refWidth), float64(maxTotalWidth/cfg.Columns)))
	cellHeight := int(math.Min(float64(refHeight), float64(maxTotalHeight/cfg.Rows)))

	if cfg.AspectRatio == "4x5" {
		totalWidth := cellWidth * cfg.Columns
		totalHeight := int(float64(totalWidth) * 5.0 / 4.0)
		cellHeight = totalHeight / cfg.Rows
	}

	images = scaleImages(images, cellWidth, cellHeight)
	logger.Info("Images scaled")

	grid := createGrid(cfg.Rows, cfg.Columns, cellWidth, cellHeight, 0)
	logger.Info("Grid created")

	grid = placeImagesInGrid(grid, images, cfg.Rows, cfg.Columns, cellWidth, cellHeight, 0)
	logger.Info("Images placed in grid")

	logger.Info("Writing output file: " + cfg.OutputFile)
	outfile, err := os.Create(cfg.OutputFile)
	if err != nil {
		logger.Error("Failed to create output file: " + cfg.OutputFile)
		return
	}
	defer outfile.Close()

	err = imaging.Encode(outfile, grid, imaging.JPEG)
	if err != nil {
		logger.Error("Failed to write output file: " + cfg.OutputFile)
		return
	}
	logger.Info("Merged image saved successfully: " + cfg.OutputFile)
}
