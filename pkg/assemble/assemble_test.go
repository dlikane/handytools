package assemble_test

import (
	"image/color"
	"os"
	"path/filepath"
	"testing"

	"github.com/disintegration/imaging"

	. "handytools/pkg/assemble"
)

func createTempImageFile(t *testing.T, dir string, name string, w, h int) string {
	t.Helper()
	img := imaging.New(w, h, color.Black)
	p := filepath.Join(dir, name)
	if err := imaging.Save(img, p); err != nil {
		t.Fatalf("failed to save temp image: %v", err)
	}
	return p
}

func TestAssembleImagesWithMax_NoValidImages(t *testing.T) {
	tmp := t.TempDir()
	outPrefix := filepath.Join(tmp, "out")

	err := AssembleImagesWithMax([]string{
		filepath.Join(tmp, "does_not_exist_1.jpg"),
		filepath.Join(tmp, "does_not_exist_2.jpg"),
	}, outPrefix, true)

	if err == nil {
		t.Fatalf("expected error for no valid images, got nil")
	}
}

func TestAssembleImagesWithMax_SinglePage_CreatesSingleFile(t *testing.T) {
	tmp := t.TempDir()
	img1 := createTempImageFile(t, tmp, "a1.jpg", 400, 300)
	img2 := createTempImageFile(t, tmp, "a2.jpg", 400, 300)

	outPrefix := filepath.Join(tmp, "assembled")

	if err := AssembleImagesWithMax([]string{img1, img2}, outPrefix, true); err != nil {
		t.Fatalf("AssembleImagesWithMax failed: %v", err)
	}

	// Expect exactly one file: assembled.jpg (no _01, _02... suffixes)
	want := outPrefix + ".jpg"
	if _, err := os.Stat(want); err != nil {
		t.Fatalf("expected output file %s to exist, got stat error: %v", want, err)
	}

	// Ensure no suffixed pages exist
	matches, _ := filepath.Glob(outPrefix + "_*.jpg")
	if len(matches) != 0 {
		t.Fatalf("expected no suffixed pages, found: %v", matches)
	}
}

func TestAssembleImagesWithMax_SkipsBadImage_StillSucceeds(t *testing.T) {
	tmp := t.TempDir()

	// One bad path and one valid image
	bad := filepath.Join(tmp, "missing.jpg")
	good := createTempImageFile(t, tmp, "good.jpg", 500, 350)

	outPrefix := filepath.Join(tmp, "mixout")

	if err := AssembleImagesWithMax([]string{bad, good}, outPrefix, true); err != nil {
		t.Fatalf("AssembleImagesWithMax failed with mixed inputs: %v", err)
	}

	// At least one output should exist (single-page or multi-page)
	// Check for base or suffixed outputs.
	base := outPrefix + ".jpg"
	_, baseErr := os.Stat(base)

	suffixed, _ := filepath.Glob(outPrefix + "_*.jpg")

	if baseErr != nil && len(suffixed) == 0 {
		t.Fatalf("expected at least one output image, got none")
	}
}
