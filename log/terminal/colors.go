package terminal

import (
	"fmt"
)

// Color codes are valid for dark color scheme.
// For bright scheme usually all colors are the same, but white and black are swapped.
const (
	ForegroundColorWhite         = 30
	ForegroundColorRed           = 31
	ForegroundColorGreen         = 32
	ForegroundColorYellow        = 33
	ForegroundColorBlue          = 34
	ForegroundColorMagenta       = 35
	ForegroundColorCyan          = 36
	ForegroundColorGray          = 37
	ForegroundColorDarkGray      = 90
	ForegroundColorBrightRed     = 91
	ForegroundColorBrightGreen   = 92
	ForegroundColorBrightYellow  = 93
	ForegroundColorBrightBlue    = 94
	ForegroundColorBrightMagenta = 95
	ForegroundColorBrightCyan    = 96
	ForegroundColorBlack         = 97

	BackgroundColorWhite         = 40
	BackgroundColorRed           = 41
	BackgroundColorGreen         = 42
	BackgroundColorYellow        = 43
	BackgroundColorBlue          = 44
	BackgroundColorMagenta       = 45
	BackgroundColorCyan          = 46
	BackgroundColorGray          = 47
	BackgroundColorDarkGray      = 100
	BackgroundColorBrightRed     = 101
	BackgroundColorBrightGreen   = 102
	BackgroundColorBrightYellow  = 103
	BackgroundColorBrightBlue    = 104
	BackgroundColorBrightMagenta = 105
	BackgroundColorBrightCyan    = 106
	BackgroundColorBlack         = 107
)

var colorNameByCode = map[int]string{
	ForegroundColorWhite:         "foreground white",
	ForegroundColorRed:           "foreground red",
	ForegroundColorGreen:         "foreground green",
	ForegroundColorYellow:        "foreground yellow",
	ForegroundColorBlue:          "foreground blue",
	ForegroundColorMagenta:       "foreground magenta",
	ForegroundColorCyan:          "foreground cyan",
	ForegroundColorGray:          "foreground gray",
	ForegroundColorDarkGray:      "foreground dark gray",
	ForegroundColorBrightRed:     "foreground bright red",
	ForegroundColorBrightGreen:   "foreground bright green",
	ForegroundColorBrightYellow:  "foreground bright yellow",
	ForegroundColorBrightBlue:    "foreground bright blue",
	ForegroundColorBrightMagenta: "foreground bright magenta",
	ForegroundColorBrightCyan:    "foreground bright cyan",
	ForegroundColorBlack:         "foreground black",

	BackgroundColorWhite:         "background white",
	BackgroundColorRed:           "background red",
	BackgroundColorGreen:         "background green",
	BackgroundColorYellow:        "background yellow",
	BackgroundColorBlue:          "background blue",
	BackgroundColorMagenta:       "background magenta",
	BackgroundColorCyan:          "background cyan",
	BackgroundColorGray:          "background gray",
	BackgroundColorDarkGray:      "background dark gray",
	BackgroundColorBrightRed:     "background bright red",
	BackgroundColorBrightGreen:   "background bright green",
	BackgroundColorBrightYellow:  "background bright yellow",
	BackgroundColorBrightBlue:    "background bright blue",
	BackgroundColorBrightMagenta: "background bright magenta",
	BackgroundColorBrightCyan:    "background bright cyan",
	BackgroundColorBlack:         "background black",
}

func IsValidColor(color int) bool {
	_, ok := colorNameByCode[color]
	return ok
}

func GetColorName(color int) string {
	if name, ok := colorNameByCode[color]; ok {
		return name
	}

	panic(fmt.Errorf("invalid color: %d", color))
}
