package design

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	"lr2/internal/app/charts"
	"lr2/internal/app/data"
	"lr2/internal/app/parser"
	"lr2/internal/config"
	"lr2/internal/constants"
	"lr2/internal/utils"
)

var serverRunning = false

func startServer(cfg *config.Config) {
	if !serverRunning {
		serverRunning = true

		go func() {
			http.Handle("/", http.FileServer(http.Dir(cfg.StorageRoot)))
			err := http.ListenAndServe(":"+cfg.ServerPort, nil)
			if err != nil {
				log.Printf(constants.StartServerError.String(), err)
			}

			serverRunning = false
		}()
	}

	go utils.OpenBrowser(fmt.Sprintf("http://localhost:%s/%s", cfg.ServerPort, cfg.Chart))
}

func validateEntry(s string) error {
	if s == "" {
		return errors.New(constants.ValueEmptyError.String())
	}
	if _, err := strconv.ParseFloat(s, 64); err != nil {
		return errors.New(constants.IncorrectValueError.String())
	}
	return nil
}

func MustLoad(cfg *config.Config, d *data.Data, p *parser.Parser, w fyne.Window) *fyne.Container {
	gEntry := widget.NewEntry()
	tEntry := widget.NewEntry()

	gEntry.Validator = validateEntry
	tEntry.Validator = validateEntry

	resultLabel := widget.NewLabel("")

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: constants.GEntryLabel.String(), Widget: gEntry},
			{Text: constants.TEntryLabel.String(), Widget: tEntry},
		},
	}

	form.OnSubmit = func() {
		gValue, _ := strconv.ParseFloat(gEntry.Text, 64)
		tValue, _ := strconv.ParseFloat(tEntry.Text, 64)

		result, err := p.Parse(gValue, tValue)
		if err != nil {
			dialog.ShowError(err, w)
		}
		resultLabel.SetText(*result)
	}

	chartFile := filepath.Join(cfg.StorageRoot, cfg.Chart)
	_, err := os.Stat(chartFile)
	if errors.Is(err, os.ErrNotExist) {
		c := charts.New(chartFile, d)
		c.Generate()
	}

	chartsBtn := widget.NewButton(constants.ShowCharts.String(), func() {
		startServer(cfg)
	})

	content := container.NewVBox(form, chartsBtn, resultLabel)
	return content
}
