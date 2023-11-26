package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"golang.design/x/clipboard"
	"os"
	"regexp"
	"slices"
	"strings"
)

func init() {
	err := clipboard.Init()
	if err != nil {
		fmt.Println("Error initializing clipboard:", err)
		os.Exit(0)
	}
}

func main() {
	ctx := context.TODO()
	ch := clipboard.Watch(ctx, clipboard.FmtText)
	fmt.Println("Czekam na skopiowanie CSS do schowka... (Ctrl+C aby wyjść)")

	for data := range ch {
		r := bytes.NewReader(data)
		items := initialScan(r)
		if len(items) < 1 {
			fmt.Println("Nie znaleziono kolorów w skopiowanym CSS. Spróbuj ponownie lub CTRL+C aby wyjść.")
			continue
		}

		itemsJson, _ := json.MarshalIndent(items, "", "\t")
		clipboard.Write(clipboard.FmtText, itemsJson)
		fmt.Println("Kolory zostały skopiowane do schowka.")
		break
	}
}

func initialScan(r *bytes.Reader) map[string]map[string]string {
	var lineIsCandidate bool
	var currentKey string

	initialItems := make(map[string]map[string]string)
	scanner := bufio.NewScanner(r)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		line := scanner.Text()

		if line == "" {
			continue
		}

		if lineIsCandidate {
			nameSlashHue := extractNameSlashHue(currentKey)
			colorParts := strings.Split(nameSlashHue, "/")
			lineIsCandidate = false
			currentKey = ""

			if !isLineWithBgColor(line) || len(colorParts) != 2 {
				continue
			}

			if colorParts[0] == "" || colorParts[1] == "" {
				continue
			}

			if _, ok := initialItems[colorParts[0]]; !ok {
				initialItems[colorParts[0]] = make(map[string]string)
			}

			initialItems[colorParts[0]][colorParts[1]] = extractColor(line)
		}

		if !isComment(line) {
			continue
		}

		lineIsCandidate = true
		currentKey = line
	}

	if len(initialItems) > 0 {
		for colorName, colorValues := range initialItems {
			if len(colorValues) > 0 {
				collectedHues := make([]string, 0)

				for hueName, _ := range colorValues {
					collectedHues = append(collectedHues, hueName)
				}

				slices.Sort(collectedHues)

				initialItems[colorName]["DEFAULT"] = collectedHues[len(collectedHues)-1]
			}
		}
	}

	return initialItems
}

func isComment(line string) bool {
	re := regexp.MustCompile("(?i)/\\*\\s.*\\s\\*/")
	return re.MatchString(line)
}

func isLineWithBgColor(line string) bool {
	re := regexp.MustCompile("(?i)background:")
	return re.MatchString(line)
}

func extractColor(line string) string {
	re := regexp.MustCompile("(?i)background:(.*#(?:[A-Fa-f0-9]{3}){1,2}\\b);")
	result := re.FindStringSubmatch(line)
	if len(result) < 2 {
		return ""
	}
	return strings.TrimSpace(result[1])
}

func extractNameSlashHue(line string) string {
	re := regexp.MustCompile("(?i)/\\*\\s(.*)\\s\\*/")
	result := re.FindStringSubmatch(line)
	if len(result) < 2 {
		return ""
	}
	return strings.Replace(strings.TrimSpace(result[1]), "", "", -1)
}
