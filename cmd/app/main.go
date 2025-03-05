package main

import (
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"

	"lr2/internal/app/data"
	"lr2/internal/app/design"
	"lr2/internal/app/parser"
	"lr2/internal/app/repository"
	"lr2/internal/config"
	"lr2/internal/constants"
)

func main() {
	a := app.New()
	w := a.NewWindow(constants.WindowName.String())

	cfg := config.MustLoad()

	d := data.New()
	r := repository.New(filepath.Join(cfg.StorageRoot, cfg.Rules))
	p := parser.New(r, d)

	des := design.MustLoad(cfg, d, p, w)

	w.SetContent(des)
	w.Resize(fyne.Size{Width: 400, Height: 200})
	w.ShowAndRun()
}
