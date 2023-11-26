package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"golang.design/x/clipboard"
	"os"
	"totailwind/internal/basecolors"
)

var ctx context.Context
var ch <-chan []byte
var items basecolors.ColorMap
var preview *widget.Label
var button *widget.Button

const msgGoToFigma = "Idź do Figmy i skopiuj styl CSS do schowka"
const msgCopyResultToClipboard = "Skopiuj wynik do schowka"
const msgColorsDetected = "Kolorwy zostały rozpoznane, gotowe!"
const msgColorsNotDetected = "Styl CSS nie został rozpoznany :("

func init() {
	err := clipboard.Init()
	if err != nil {
		fmt.Println("Nie mogę zaninicjować schowka:", err)
		os.Exit(0)
	}

	ctx = context.TODO()
	ch = clipboard.Watch(ctx, clipboard.FmtText)

}

func main() {
	a := app.New()
	a.Settings().Theme().Size("20")
	w := a.NewWindow("Figma Tailwind Colors")
	w.CenterOnScreen()

	preview = widget.NewLabel(msgGoToFigma)
	preview.Alignment = fyne.TextAlignCenter

	button = widget.NewButton(msgCopyResultToClipboard, func() {
		itemsJson, _ := json.MarshalIndent(items, "", "\t")
		clipboard.Write(clipboard.FmtText, itemsJson)
		button.Disable()
		preview.SetText(msgGoToFigma)

		ctx.Done()
		ctx = context.TODO()
		ch = clipboard.Watch(ctx, clipboard.FmtText)

		go watchClipboard()
	})
	button.Disable()

	w.SetContent(container.NewVBox(
		preview,
		button,
	))

	go watchClipboard()
	w.ShowAndRun()
}

func watchClipboard() {
	for data := range ch {
		items = basecolors.Process(bytes.NewReader(data))

		if len(items) < 1 {
			preview.SetText(msgColorsNotDetected)
			continue
		}

		preview.SetText(msgColorsDetected)
		button.Enable()

		break
	}
}
