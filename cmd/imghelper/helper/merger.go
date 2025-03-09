package helper

import (
	"fmt"
	"github.com/disintegration/imaging"
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

func MergeImages(cliParams CliParams) error {
	logger := GetLogger()

	var images []image.Image
	for _, imgPath := range cliParams.Images {
		logger.Info("Opening image file: " + imgPath)
		infile, err := os.Open(imgPath)
		if err != nil {
			logger.Error("Failed to open image file: " + imgPath)
			return fmt.Errorf("failed to open file: %w", err)
		}

		img, _, err := image.Decode(infile)
		if err != nil {
			logger.Error("Failed to decode image file: " + imgPath)
			return fmt.Errorf("failed to decode file: %w", err)
		}

		infile.Close()

		images = append(images, img)
		logger.Info("Image file loaded: " + imgPath)
	}

	if len(images) == 0 {
		logger.Error("No input images provided")
		return fmt.Errorf("no images provided")
	}

	width := images[0].Bounds().Dx()
	height := images[0].Bounds().Dy()

	images = scaleImages(images, width, height)
	logger.Info("Images scaled")

	grid := createGrid(cliParams.Rows, cliParams.Cols, width, height, cliParams.Whitespace)
	logger.Info("Grid created")

	grid = placeImagesInGrid(grid, images, cliParams.Rows, cliParams.Cols, width, height, cliParams.Whitespace)
	logger.Info("Images placed in grid")

	// Save the final image
	logger.Info("Writing output file: " + cliParams.Output)
	outfile, err := os.Create(cliParams.Output)
	if err != nil {
		logger.Error("Failed to create output file: " + cliParams.Output)
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer outfile.Close()

	err = imaging.Encode(outfile, grid, imaging.JPEG)
	if err != nil {
		logger.Error("Failed to write output file: " + cliParams.Output)
		return fmt.Errorf("failed to write output file: %w", err)
	}

	logger.Info("Merged image saved successfully: " + cliParams.Output)
	return nil
}
