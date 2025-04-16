package pinterest

import (
	"context"
	"fmt"
	"handytools/pkg/common"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/chromedp/chromedp"
)

func DownloadPins(boardURL, outputDir string) ([]string, error) {
	logger := common.GetLogger()
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()
	ctx, cancel = context.WithTimeout(ctx, 45*time.Second)
	defer cancel()

	var imageLinks []string

	logger.Info("Opening board in headless Chrome...")
	err := chromedp.Run(ctx,
		chromedp.Navigate(boardURL),
		chromedp.Sleep(3*time.Second),

		// Scroll 6 times (stable amount)
		chromedp.ActionFunc(func(ctx context.Context) error {
			for i := 0; i < 6; i++ {
				if err := chromedp.Run(ctx,
					chromedp.Evaluate(`window.scrollBy(0, 1500)`, nil),
					chromedp.Sleep(1*time.Second),
				); err != nil {
					return err
				}
			}
			return nil
		}),

		// Use only <img> based sources
		chromedp.Evaluate(`
			Array.from(document.querySelectorAll("img")).map(img => {
				if (img.src && img.src.includes("pinimg.com")) return img.src;
				if (img.dataset?.src?.includes("pinimg.com")) return img.dataset.src;
				if (img.srcset && img.srcset.includes("pinimg.com")) {
					return img.srcset.split(",").pop().trim().split(" ")[0];
				}
				return null;
			}).filter(Boolean)
		`, &imageLinks),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch pins: %w", err)
	}

	seen := map[string]bool{}
	var imagePaths []string
	for _, link := range imageLinks {
		if seen[link] {
			continue
		}
		seen[link] = true
	}
	logger.Infof("Found %d unique image links", len(seen))

	i := 1
	for link := range seen {
		path := filepath.Join(outputDir, fmt.Sprintf("pin_%03d.jpg", i))
		if err := downloadImage(link, path); err != nil {
			logger.WithError(err).Warnf("Failed to download: %s", link)
			continue
		}
		imagePaths = append(imagePaths, path)
		logger.Infof("Saved: %s", path)
		i++
	}

	return imagePaths, nil
}

func downloadImage(url, path string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(path)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}
