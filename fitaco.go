package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fitaco/internal/basecolors"
	"fitaco/internal/helpers"
	"fmt"
	"golang.design/x/clipboard"
	"os"
)

var ctx context.Context
var ch <-chan []byte
var messagesChannel chan string
var items basecolors.ColorMap
var outputBytes []byte
var outputJs string
var statusMsg string
var status bool

const msgGoToFigma = "Czekam aż skopiujesz styl Piotra z Figmy do schowka"
const msgCopySuccess = "Wynik skopiowany do schowka!"
const msgColorsDetected = "Kolory zostały rozpoznane, gotowe!"
const msgColorsNotDetected = "Styl CSS nie został rozpoznany :("

func init() {
	err := clipboard.Init()
	if err != nil {
		fmt.Println("Nie mogę zaninicjować schowka:", err)
		os.Exit(0)
	}

	ctx = context.TODO()
	ch = clipboard.Watch(ctx, clipboard.FmtText)
	messagesChannel = make(chan string)
}

func main() {
	helpers.DisplayLogoNewProject()
	helpers.Pretty("@todo", msgGoToFigma)
	go watchClipboard()

	for message := range messagesChannel {
		helpers.Pretty("@info", message)
	}

	clipboard.Write(clipboard.FmtText, outputBytes)
}

func watchClipboard() {
	defer close(messagesChannel)

	for data := range ch {
		items = basecolors.Process(bytes.NewReader(data))

		if len(items) < 1 {
			messagesChannel <- msgColorsNotDetected

			continue
		}

		messagesChannel <- msgColorsDetected
		break
	}

	outputBytes, _ = json.MarshalIndent(items, "", "\t")

	messagesChannel <- msgCopySuccess
}
