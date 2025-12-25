package assemble

import (
	"image"
	"image/color"
	"testing"

	"github.com/disintegration/imaging"
)

func genImages(n, w, h int, c color.NRGBA) []image.Image {
	images := make([]image.Image, 0, n)
	for i := 0; i < n; i++ {
		images = append(images, imaging.New(w, h, c))
	}
	return images
}

func TestBuildContinuousCanvas_SinglePage_WhenFitOnePageTrue(t *testing.T) {
	t.Parallel()

	// Many images that would normally overflow multiple pages at defaultRowH.
	images := genImages(24, 1000, 1000, color.NRGBA{R: 10, G: 20, B: 30, A: 255})

	canvas, pageBreaks, err := buildContinuousCanvas(images, true)
	if err != nil {
		t.Fatalf("buildContinuousCanvas returned error: %v", err)
	}

	if canvas == nil {
		t.Fatalf("expected non-nil canvas")
	}

	// When fitOnePage=true, rowHeight should be reduced until it fits one page (or min row height).
	if len(pageBreaks) > 1 {
		t.Fatalf("expected <= 1 page when fitOnePage=true, got %d", len(pageBreaks))
	}

	if gotW := canvas.Bounds().Dx(); gotW != maxWidth {
		t.Fatalf("expected canvas width %d, got %d", maxWidth, gotW)
	}

	if gotH := canvas.Bounds().Dy(); gotH <= 0 || gotH > maxHeight {
		t.Fatalf("expected canvas height in (0, %d], got %d", maxHeight, gotH)
	}
}

func TestBuildContinuousCanvas_MultiPage_WhenFitOnePageFalse(t *testing.T) {
	t.Parallel()

	// Force multiple pages by using a large number of images.
	images := genImages(48, 1200, 800, color.NRGBA{R: 200, G: 50, B: 50, A: 255})

	canvas, pageBreaks, err := buildContinuousCanvas(images, false)
	if err != nil {
		t.Fatalf("buildContinuousCanvas returned error: %v", err)
	}

	if canvas == nil {
		t.Fatalf("expected non-nil canvas")
	}

	// With fitOnePage=false, pagination is allowed; expect at least 2 pages for a large set.
	if len(pageBreaks) < 2 {
		t.Fatalf("expected multiple pages, got %d", len(pageBreaks))
	}

	if gotW := canvas.Bounds().Dx(); gotW != maxWidth {
		t.Fatalf("expected canvas width %d, got %d", maxWidth, gotW)
	}

	if gotH := canvas.Bounds().Dy(); gotH <= maxHeight {
		// Not strictly required to be > maxHeight, but with a large set it should be taller than a single page.
		t.Fatalf("expected canvas height to exceed single-page height %d, got %d", maxHeight, gotH)
	}
}

func TestBuildContinuousCanvas_BasicPlacementDimensions(t *testing.T) {
	t.Parallel()

	images := []image.Image{
		imaging.New(800, 1200, color.NRGBA{R: 0, G: 0, B: 0, A: 255}),
		imaging.New(1600, 1200, color.NRGBA{R: 0, G: 0, B: 0, A: 255}),
		imaging.New(1200, 800, color.NRGBA{R: 0, G: 0, B: 0, A: 255}),
	}

	canvas, pageBreaks, err := buildContinuousCanvas(images, false)
	if err != nil {
		t.Fatalf("buildContinuousCanvas returned error: %v", err)
	}

	if canvas == nil {
		t.Fatalf("expected non-nil canvas")
	}

	if len(pageBreaks) == 0 {
		t.Fatalf("expected at least one page break (end of first page), got 0")
	}

	if canvas.Bounds().Dx() != maxWidth {
		t.Fatalf("expected canvas width %d, got %d", maxWidth, canvas.Bounds().Dx())
	}

	if canvas.Bounds().Dy() <= 0 {
		t.Fatalf("expected positive canvas height, got %d", canvas.Bounds().Dy())
	}
}
