package log

import "fmt"

type ColorAnsiBackground = []byte

// BackgroundColors are the default backgrounds supported by Ultralogger. All of these colors are the 3-bit ANSI colors
// supported by *most* terminals They can be used in a ColorizedFormatter to colorize log lines by level.
//
// Note: These colors are not interchangeable with normal Color Ansi colors. They are only used for background colors
// because Colorization must be applied in a specific order.
//
// See https://en.wikipedia.org/wiki/ANSI_escape_code#3-bit_and_4-bit for more info on ANSI colors.
var BackgroundColors = struct {
    Black   ColorAnsiBackground
    Red     ColorAnsiBackground
    Green   ColorAnsiBackground
    Yellow  ColorAnsiBackground
    Blue    ColorAnsiBackground
    Magenta ColorAnsiBackground
    Cyan    ColorAnsiBackground
    White   ColorAnsiBackground
}{
    Black:   ColorAnsiBackground("40"),
    Red:     ColorAnsiBackground("41"),
    Green:   ColorAnsiBackground("42"),
    Yellow:  ColorAnsiBackground("43"),
    Blue:    ColorAnsiBackground("44"),
    Magenta: ColorAnsiBackground("45"),
    Cyan:    ColorAnsiBackground("46"),
    White:   ColorAnsiBackground("47"),
}

// BackgroundRGB returns a ColorAnsiBackground that represents an RGB background color.
func BackgroundRGB(r, g, b int) ColorAnsiBackground {
    return ColorAnsiBackground(fmt.Sprintf("48;2;%d;%d;%d", r, g, b))
}
