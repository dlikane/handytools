package main

import (
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

const (
	inputDir     = "frames"
	outputDir    = "processed"
	fps          = 25.0
	grayscaleDur = 3
	invertDur    = 2
	numWorkers   = 8
)

func loadBeats(path string) ([]float64, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var beats []float64
	err = json.Unmarshal(data, &beats)
	return beats, err
}

func isInEffectWindow(frameTime float64, beats []float64) string {
	for _, beat := range beats {
		delta := frameTime - beat
		if delta >= 0 && delta < float64(grayscaleDur)/fps {
			return "grayscale"
		}
		if delta >= float64(grayscaleDur)/fps && delta < float64(grayscaleDur+invertDur)/fps {
			return "invert"
		}
	}
	return ""
}

func applyGrayscale(img image.Image) image.Image {
	bounds := img.Bounds()
	gray := image.NewRGBA(bounds)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			l := uint8((r*299 + g*587 + b*114 + 500) / 1000 >> 8)
			gray.Set(x, y, color.RGBA{l, l, l, uint8(a >> 8)})
		}
	}
	return gray
}

func applyInvert(img image.Image) image.Image {
	bounds := img.Bounds()
	inv := image.NewRGBA(bounds)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			inv.Set(x, y, color.RGBA{
				255 - uint8(r>>8),
				255 - uint8(g>>8),
				255 - uint8(b>>8),
				uint8(a >> 8),
			})
		}
	}
	return inv
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, in)
	return err
}

func processFrame(file os.FileInfo, index int, beats []float64) {
	framePath := filepath.Join(inputDir, file.Name())
	outPath := filepath.Join(outputDir, file.Name())

	frameTime := float64(index) / fps
	effect := isInEffectWindow(frameTime, beats)

	if effect == "" {
		if err := copyFile(framePath, outPath); err != nil {
			log.Printf("Failed to copy %s: %v", file.Name(), err)
		} else {
			fmt.Printf("Copied frame %d: no effect\n", index+1)
		}
		return
	}

	f, err := os.Open(framePath)
	if err != nil {
		log.Printf("Failed to open %s: %v", framePath, err)
		return
	}
	img, err := jpeg.Decode(f)
	f.Close()
	if err != nil {
		log.Printf("Failed to decode %s: %v", framePath, err)
		return
	}

	switch effect {
	case "grayscale":
		img = applyGrayscale(img)
	case "invert":
		img = applyInvert(img)
	}

	outFile, err := os.Create(outPath)
	if err != nil {
		log.Printf("Failed to create output %s: %v", outPath, err)
		return
	}
	jpeg.Encode(outFile, img, &jpeg.Options{Quality: 95})
	outFile.Close()

	fmt.Printf("Processed frame %d: %s\n", index+1, effect)
}

func main() {
	beats, err := loadBeats("beats.json")
	if err != nil {
		log.Fatalf("Failed to load beats: %v", err)
	}

	files, err := ioutil.ReadDir(inputDir)
	if err != nil {
		log.Fatalf("Failed to read input frames: %v", err)
	}

	os.MkdirAll(outputDir, 0755)

	jobs := make(chan struct {
		file  os.FileInfo
		index int
	}, numWorkers)

	var wg sync.WaitGroup

	for w := 0; w < numWorkers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for job := range jobs {
				processFrame(job.file, job.index, beats)
			}
		}()
	}

	for i, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".jpg") {
			continue
		}
		jobs <- struct {
			file  os.FileInfo
			index int
		}{file, i}
	}
	close(jobs)
	wg.Wait()
}
