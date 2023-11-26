package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/widget"
	"golang.design/x/clipboard"
	"os"
	"totailwind/internal/basecolors"
)

var ctx context.Context
var ch <-chan []byte
var items basecolors.ColorMap
var preview *widget.Label

func init() {
	err := clipboard.Init()
	if err != nil {
		fmt.Println("Error initializing clipboard:", err)
		os.Exit(0)
	}

	ctx = context.TODO()
	ch = clipboard.Watch(ctx, clipboard.FmtText)
}

func watchClipboard() {
	for data := range ch {
		items = basecolors.Process(bytes.NewReader(data))

		if len(items) < 1 {
			fmt.Println("Nie znaleziono kolorów w skopiowanym CSS. Spróbuj ponownie.")
			continue
		}

		itemsJson, _ := json.MarshalIndent(items, "", "\t")
		fmt.Println("Kolory zostały skopiowane do schowka!")
		preview.SetText(string(itemsJson))
		clipboard.Write(clipboard.FmtText, itemsJson)
		break
	}
}

func main() {
	fmt.Println("Czekam na skopiowanie CSS do schowka... (Ctrl+C aby wyjść)")

	a := app.New()
	w := a.NewWindow("FigTail siempre")
	w.Resize(fyne.NewSize(640, 480))
	w.CenterOnScreen()
	preview = widget.NewLabel("empty!")
	w.SetContent(preview)
	go watchClipboard()

	w.ShowAndRun()
}
