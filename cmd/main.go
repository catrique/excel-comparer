package main

import (
	"excel-comparer/internal/service"
	"excel-comparer/internal/ui"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

func main() {
	meuApp := app.New()
	janela := meuApp.NewWindow("Comparador de Planilhas")
	janela.Resize(fyne.NewSize(900, 650))

	if icone, err := fyne.LoadResourceFromPath("Icon.png"); err == nil {
		meuApp.SetIcon(icone)
	}

	ss := service.NewSpreadsheetService()
	cs := service.NewCompareService()
	es := service.NewExportService()

	mainWindow := ui.NewMainWindow(janela, ss, cs, es)
	mainWindow.BuildUI()

	janela.ShowAndRun()
}
