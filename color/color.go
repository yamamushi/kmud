package color

import (
	"fmt"
	"regexp"
	"strings"
)

var ColorRegex = regexp.MustCompile("([@#][0-6]|@@|##)")

type ColorMode int

const (
	ModeLight ColorMode = iota
	ModeDark  ColorMode = iota
	ModeNone  ColorMode = iota
)

type Color string

const (
	Red     Color = "@0"
	Green   Color = "@1"
	Yellow  Color = "@2"
	Blue    Color = "@3"
	Magenta Color = "@4"
	Cyan    Color = "@5"
	White   Color = "@6"

	DarkRed     Color = "#0"
	DarkGreen   Color = "#1"
	DarkYellow  Color = "#2"
	DarkBlue    Color = "#3"
	DarkMagenta Color = "#4"
	DarkCyan    Color = "#5"
	Black       Color = "#6"

	Gray   Color = "@@"
	Normal Color = "##"
)

type colorCode string

const (
	red     colorCode = "\033[01;31m"
	green   colorCode = "\033[01;32m"
	yellow  colorCode = "\033[01;33m"
	blue    colorCode = "\033[01;34m"
	magenta colorCode = "\033[01;35m"
	cyan    colorCode = "\033[01;36m"
	white   colorCode = "\033[01;37m"

	darkRed     colorCode = "\033[22;31m"
	darkGreen   colorCode = "\033[22;32m"
	darkYellow  colorCode = "\033[22;33m"
	darkBlue    colorCode = "\033[22;34m"
	darkMagenta colorCode = "\033[22;35m"
	darkCyan    colorCode = "\033[22;36m"
	black       colorCode = "\033[22;30m"

	gray   colorCode = "\033[22;37m"
	normal colorCode = "\033[0m"
)

func getAnsiCode(mode ColorMode, color Color) string {
	if mode == ModeNone {
		return ""
	}

	var code colorCode
	switch color {
	case Normal:
		code = normal
	case Red:
		code = red
	case Green:
		code = green
	case Yellow:
		code = yellow
	case Blue:
		code = blue
	case Magenta:
		code = magenta
	case Cyan:
		code = cyan
	case White:
		code = white
	case DarkRed:
		code = darkRed
	case DarkGreen:
		code = darkGreen
	case DarkYellow:
		code = darkYellow
	case DarkBlue:
		code = darkBlue
	case DarkMagenta:
		code = darkMagenta
	case DarkCyan:
		code = darkCyan
	case Black:
		code = black
	case Gray:
		code = gray
	}

	if mode == ModeDark {
		if code == white {
			return string(black)
		} else if code == black {
			return string(white)
		} else if strings.Contains(string(code), "01") {
			return strings.Replace(string(code), "01", "22", 1)
		} else {
			return strings.Replace(string(code), "22", "01", 1)
		}
	}

	return string(code)
}

// Wraps the given text in the given color, followed by a color reset
func Colorize(color Color, text string) string {
	return fmt.Sprintf("%s%s%s", string(color), text, string(Normal))
}

var Lookup = map[Color]bool{
	Red:         true,
	Green:       true,
	Yellow:      true,
	Blue:        true,
	Magenta:     true,
	Cyan:        true,
	White:       true,
	DarkRed:     true,
	DarkGreen:   true,
	DarkYellow:  true,
	DarkBlue:    true,
	DarkMagenta: true,
	DarkCyan:    true,
	Black:       true,
	Gray:        true,
	Normal:      true,
}

// Strips MUD color codes and replaces them with ansi color codes
func ProcessColors(text string, cm ColorMode) string {
	replace := func(match string) string {
		found := Lookup[Color(match)]

		if found {
			return getAnsiCode(cm, Color(match))
		}

		return match
	}

	after := ColorRegex.ReplaceAllStringFunc(text, replace)
	return after
}

func StripColors(text string) string {
	return ColorRegex.ReplaceAllString(text, "")
}
