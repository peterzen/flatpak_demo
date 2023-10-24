package main

import (
	_ "embed"
	"io"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"fyne.io/systray"
)

//go:embed Icon.png
var iconData []byte

func main() {
	a := app.NewWithID("io.fyne.flatpak_demo")
	a.SetIcon(&fyne.StaticResource{StaticName: "Icon.png", StaticContent: iconData})

	w := a.NewWindow("Flatpak Demo")

	markdown := &widget.Entry{MultiLine: true, Wrapping: fyne.TextWrapWord}
	preview := &widget.RichText{Wrapping: fyne.TextWrapWord}
	markdown.OnChanged = preview.ParseMarkdown

	open := &widget.Button{Text: "Open file", Icon: theme.ContentAddIcon(), OnTapped: func() {
		files := dialog.NewFileOpen(func(file fyne.URIReadCloser, err error) {
			if err != nil {
				log.Println(err)
				return
			} else if file == nil {
				return
			}

			text, err := io.ReadAll(file)
			if err != nil {
				log.Println(err)
				return
			}

			markdown.SetText(string(text))
			a.SendNotification(&fyne.Notification{Title: "Opened a file", Content: file.URI().Name() + " was opened correctly."})
		}, w)

		files.SetFilter(storage.NewExtensionFileFilter([]string{".md"}))
		files.Show()
	}}

	w.SetContent(
		container.NewBorder(container.NewHBox(open), nil, nil, nil,
			container.NewHSplit(markdown, preview),
		),
	)

	systray.Run(func() {
		systray.SetTitle("Flatpak Demo")
		systray.SetTooltip("Demo application")

		mOpen := systray.AddMenuItem("Open", "Open main window")
		openC := make(chan struct{}) // no buffer. Is ignored if window is currently open
		go func() {
			for range mOpen.ClickedCh {
				select {
				case openC <- struct{}{}:
					log.Printf("Received window reopen request from system tray")
				default:
					log.Printf("Ignored a window open request from the system tray")
				}
			}
		}()

	}, nil)

	w.Resize(fyne.NewSize(600, 400))
	w.ShowAndRun()
}
