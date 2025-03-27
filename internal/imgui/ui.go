package imgui

import (
	"fmt"
	"image"
	"os"
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func Run() {
	a := app.New()
	w := a.NewWindow("Image Explorer")
	w.Resize(fyne.NewSize(1000, 700))

	root := "." // starting directory

	tree := widget.NewTree(
		func(uid string) []string {
			files, _ := os.ReadDir(uid)
			var children []string
			for _, f := range files {
				children = append(children, filepath.Join(uid, f.Name()))
			}
			return children
		},
		func(uid string) bool {
			fi, err := os.Stat(uid)
			return err == nil && fi.IsDir()
		},
		func(branch bool) fyne.CanvasObject {
			return widget.NewLabel("template")
		},
		func(uid string, branch bool, obj fyne.CanvasObject) {
			obj.(*widget.Label).SetText(filepath.Base(uid))
		},
	)

	tree.Root = root

	var preview fyne.CanvasObject = canvas.NewText("Select a JPG to preview", nil)

	tree.OnSelected = func(path string) {
		if isJPG(path) {
			img := loadImage(path)
			if img != nil {
				imageObj := canvas.NewImageFromImage(img)
				imageObj.FillMode = canvas.ImageFillContain
				preview = imageObj
				content := container.NewHSplit(tree, preview)
				content.Offset = 0.3
				w.SetContent(content)
			}
		}
	}

	content := container.NewHSplit(tree, preview)
	content.Offset = 0.3
	w.SetContent(content)
	w.ShowAndRun()
}

func isJPG(path string) bool {
	ext := filepath.Ext(path)
	return ext == ".jpg" || ext == ".jpeg"
}

func loadImage(path string) image.Image {
	file, err := os.Open(path)
	if err != nil {
		fmt.Println("open error:", err)
		return nil
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		fmt.Println("decode error:", err)
		return nil
	}
	return img
}
