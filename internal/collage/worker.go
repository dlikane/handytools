package collage

import (
	"github.com/disintegration/imaging"
	"handytools/pkg/common"
	"image"
	"image/color"
	"os"
)

// scaleImages scales the provided images to match the dimensions of the first image
func scaleImages(images []image.Image, width, height int) []image.Image {
	for i, img := range images {
		// Resize the image to width x height and crop it if necessary.
		images[i] = imaging.Fill(img, width, height, imaging.Center, imaging.Lanczos)
	}
	return images
}

// createGrid creates a blank image to serve as the grid for the combined image, filling it
// with a white color and accounting for the frame dimensions
func createGrid(rows, cols, width, height, frame int) *image.NRGBA {
	grid := imaging.New(width*cols+(cols-1)*frame, height*rows+(rows-1)*frame, color.NRGBA{0xff, 0xff, 0xff, 0xff})
	return grid
}

// placeImagesInGrid places each image at the correct position in the grid
func placeImagesInGrid(grid *image.NRGBA, images []image.Image, rows, cols, width, height, frame int) *image.NRGBA {
	i := 0
	for y := 0; y < rows; y++ {
		for x := 0; x < cols; x++ {
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
		if err != nil {
			logger.Error("Failed to decode image file: " + imgPath)
			return
		}
		infile.Close()

		images = append(images, img)
		logger.Info("Image file loaded: " + imgPath)
	}

	if len(images) == 0 {
		logger.Error("No input images provided")
		return
	}

	width := images[0].Bounds().Dx()
	height := images[0].Bounds().Dy()

	images = scaleImages(images, width, height)
	logger.Info("Images scaled")

	grid := createGrid(cfg.Rows, cfg.Columns, width, height, 0 /*cfg.Whitespace*/)
	logger.Info("Grid created")

	grid = placeImagesInGrid(grid, images, cfg.Rows, cfg.Columns, width, height, 0 /*cfg.Whitespace*/)
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

	return
}
