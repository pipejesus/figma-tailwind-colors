package basecolors

import (
	"bufio"
	"bytes"
	"regexp"
	"slices"
	"strings"
)

type ColorMap map[string]map[string]string

var regExpIsComment *regexp.Regexp
var regExpIsLineWithBgColor *regexp.Regexp
var regExpExtractColor *regexp.Regexp
var regExpColorNameSlashWeight *regexp.Regexp

func init() {
	regExpIsComment = regexp.MustCompile("(?i)/\\*\\s.*\\s\\*/")
	regExpIsLineWithBgColor = regexp.MustCompile("(?i)background:")
	regExpExtractColor = regexp.MustCompile("(?i)background:(.*#(?:[A-Fa-f0-9]{3}){1,2}\\b);")
	regExpColorNameSlashWeight = regexp.MustCompile("(?i)/\\*\\s(.*)\\s\\*/")
}

// Process reads the exported Figma CSS and tries to return the colors.
// The criteria for a valid color is that it's a css background property
// with hex color preceded by a comment formatted as "color name/color weight"
func Process(r *bytes.Reader) ColorMap {
	var lineIsCandidate bool
	var currentKey string

	colors := make(ColorMap)
	scanner := bufio.NewScanner(r)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		line := scanner.Text()

		if line == "" {
			continue
		}

		if lineIsCandidate {
			ok := processCandidateLine(line, currentKey, &colors)
			lineIsCandidate = false
			currentKey = ""

			if !ok {
				continue
			}
		}

		if !isComment(line) {
			continue
		}

		lineIsCandidate = true
		currentKey = line
	}

	if len(colors) > 0 {
		colors = addDefaultWeight(colors)
	}

	return colors
}

func processCandidateLine(line string, currentKey string, colors *ColorMap) bool {
	nameSlashHue := extractColorNameSlashWeight(currentKey)
	colorParts := strings.Split(nameSlashHue, "/")

	if !isLineWithBgColor(line) || len(colorParts) != 2 {
		return false
	}

	colorName := colorParts[0]
	colorWeight := colorParts[1]

	if colorName == "" || colorWeight == "" {
		return false
	}

	if _, ok := (*colors)[colorName]; !ok {
		(*colors)[colorName] = make(map[string]string)
	}

	(*colors)[colorName][colorWeight] = extractColor(line)

	return true
}

func addDefaultWeight(colors ColorMap) ColorMap {
	for colorName, colorValues := range colors {
		if len(colorValues) > 0 {
			collectedHues := make([]string, 0)

			for hueName, _ := range colorValues {
				collectedHues = append(collectedHues, hueName)
			}

			slices.Sort(collectedHues)

			colors[colorName]["DEFAULT"] = collectedHues[len(collectedHues)-1]
		}
	}

	return colors
}

func isComment(line string) bool {
	return regExpIsComment.MatchString(line)
}

func isLineWithBgColor(line string) bool {
	return regExpIsLineWithBgColor.MatchString(line)
}

func extractColor(line string) string {
	result := regExpExtractColor.FindStringSubmatch(line)

	if len(result) < 2 {
		return ""
	}

	return strings.TrimSpace(result[1])
}

func extractColorNameSlashWeight(line string) string {
	result := regExpColorNameSlashWeight.FindStringSubmatch(line)

	if len(result) < 2 {
		return ""
	}

	return strings.Replace(strings.TrimSpace(result[1]), " ", "", -1)
}
