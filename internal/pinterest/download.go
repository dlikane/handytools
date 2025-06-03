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
	ctx, cancel = context.WithTimeout(ctx, 120*time.Second)
	defer cancel()

	var imageLinks []string

	logger.Info("Opening board in headless Chrome...")

	err := chromedp.Run(ctx,
		chromedp.Navigate(boardURL),
		chromedp.Sleep(3*time.Second),

		chromedp.ActionFunc(func(ctx context.Context) error {
			var allLinks []string
			for i := 0; i < 40; i++ {
				logger.Infof("Scrolling step %d/40", i+1)
				var newLinks []string
				err := chromedp.Run(ctx,
					chromedp.Evaluate(`
						Array.from(document.querySelectorAll('div[data-grid-item="true"] img')).filter(img => {
							let parent = img.closest('div[data-test-id="related-interests-multi-column-module"]');
							let boardParent = img.closest('div[data-test-id="board-feed"]');
							return !parent && boardParent;
						}).map(img => img.src).filter(src => src.includes("pinimg.com"));
					`, &newLinks),
				)
				if err != nil {
					return err
				}
				allLinks = append(allLinks, newLinks...)
				if len(allLinks) > 0 && len(newLinks) == 0 {
					logger.Infof("No new images, break")
					break
				}

				if err := chromedp.Run(ctx,
					chromedp.Evaluate(`window.scrollBy(0, 500)`, nil),
					chromedp.Sleep(500*time.Millisecond),
				); err != nil {
					return err
				}
			}
			// remove duplicates
			seen := map[string]bool{}
			imageLinks = nil
			for _, link := range allLinks {
				if !seen[link] {
					seen[link] = true
					imageLinks = append(imageLinks, link)
				}
			}
			return nil
		}),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch pins: %w", err)
	}

	logger.Infof("Found %d unique image links", len(imageLinks))

	var imagePaths []string
	i := 1
	for _, link := range imageLinks {
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
