package main

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
)

func gui(links []string) {
	app := app.NewWithID("com.yourname.unityassetpipeline")
	window := app.NewWindow("Auto Unity Asset Pipeline")

	labelChooseInput := widget.NewLabel("Choose input file:")
	entryChooseInput := widget.NewEntry()
	buttonChooseInput := widget.NewButton("   Choose file   ", func() {
		dialog := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err != nil {
				dialog.ShowError(err, window)
				return
			}
			if reader == nil {
				return
			}
			defer reader.Close()

			selectedFile := reader.URI().Path()
			entryChooseInput.SetText(selectedFile)
		}, window)

		dialog.SetFilter(storage.NewExtensionFileFilter([]string{".xlsx"}))
		dialog.Show()
	})
	rowChooseInput := container.NewBorder(nil, nil, nil, buttonChooseInput, entryChooseInput)

	labelChooseOutput := widget.NewLabel("Choose output folder:")
	entryChooseOutput := widget.NewEntry()
	buttonChooseOutput := widget.NewButton("Choose folder", func() {
		dialog.ShowFolderOpen(func(lu fyne.ListableURI, err error) {
			selectedFolder := lu.Path()
			entryChooseOutput.SetText(selectedFolder)
		}, window)
	})
	rowChooseOutput := container.NewBorder(nil, nil, nil, buttonChooseOutput, entryChooseOutput)

	labelDownloadOption := widget.NewLabel("Keep downloaded files:")
	checkboxDownloadOption := widget.NewCheck("", func(checked bool) {
		if checked {
			fmt.Println("Đánh dấu giữ file")
		} else {
			fmt.Println("Không giữ file")
		}
	})
	rowDownloadOption := container.NewHBox(
		labelDownloadOption,
		layout.NewSpacer(),
		checkboxDownloadOption,
	)

	labelDecryptOption := widget.NewLabel("Decrypt:")
	checkboxDecryptOption := widget.NewCheck("", func(checked bool) {
		if checked {
			fmt.Println("Đánh dấu giữ file")
		} else {
			fmt.Println("Không giữ file")
		}
	})
	rowDecryptOption := container.NewHBox(
		labelDecryptOption,
		layout.NewSpacer(),
		checkboxDecryptOption,
	)
	entryDecryptFile := widget.NewEntry()
	chooseDecryptFile := widget.NewButton("   Choose file   ", func() {
		dialog.ShowFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err != nil {
				dialog.ShowError(err, window)
				return
			}
			if reader == nil {
				return
			}
			defer reader.Close()

			selectedFile := reader.URI().Path()
			entryDecryptFile.SetText(selectedFile)
		}, window)
	})
	rowDecryptFile := container.NewBorder(nil, nil, nil, chooseDecryptFile, entryDecryptFile)

	startButton := widget.NewButton("Start", func() {
		go processWithBatching(links)
	})

	window.SetContent(container.NewVBox(
		labelChooseInput,
		rowChooseInput,
		labelChooseOutput,
		rowChooseOutput,
		rowDownloadOption,
		rowDecryptOption,
		rowDecryptFile,
		layout.NewSpacer(),
		startButton,
	))

	window.Resize(fyne.NewSize(800, 500))
	window.ShowAndRun()
}
