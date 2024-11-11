package log

import "fmt"

var ansiReset = []byte("\033[0m")

var ansiCSInit = []byte("\033[")
var ansiCSEnd = byte('m')
var ansiCSSeparator = byte(';')

// TODO: 256 color (maybe)

// Colors are the default colors supported by Ultralogger. All of these colors are the 3-bit ANSI colors supported by
// *most* terminals. They can be used in a ColorizedFormatter to colorize log lines by level.
//
// See https://en.wikipedia.org/wiki/ANSI_escape_code#3-bit_and_4-bit for more info on ANSI colors.
var Colors = struct {
    Black   ColorAnsi
    Red     ColorAnsi
    Green   ColorAnsi
    Yellow  ColorAnsi
    Blue    ColorAnsi
    Magenta ColorAnsi
    Cyan    ColorAnsi
    White   ColorAnsi
    Default ColorAnsi
}{
    Black:   ColorAnsi{Code: []byte("30")},
    Red:     ColorAnsi{Code: []byte("31")},
    Green:   ColorAnsi{Code: []byte("32")},
    Yellow:  ColorAnsi{Code: []byte("33")},
    Blue:    ColorAnsi{Code: []byte("34")},
    Magenta: ColorAnsi{Code: []byte("35")},
    Cyan:    ColorAnsi{Code: []byte("36")},
    White:   ColorAnsi{Code: []byte("37")},
    Default: ColorAnsi{Code: []byte("39")},
}

// ColorAnsi is a type that represents an ANSI color. It can be used to colorize arbitrary []byte content.
type ColorAnsi struct {
    // Code is the ANSI code that represents the foreground color of a string of text. Colors are typically multi-byte
    // codes, depending on the color. For example, the code for the Red color is "31", but an RGB color is a specially
    // formatted code that looks like "38;2;{R};{G};{B}" where {R}, {G}, and {B} are the 0-255 values for the red,
    // green, and blue components of the color, respectively.
    Code []byte

    // Background is the ANSI code that represents the background color. It is applied to the []byte content after
    // the foreground color is applied.
    Background ColorAnsiBackground

    // Settings are the ANSI Settings that are applied to the color. For example, Bold, Dim, Italic, Underline,
    // SlowBlink, and Strikethrough are Settings that can be applied to a color.
    Settings []AnsiSetting
}

// ColorAnsiRGB returns a ColorAnsi that represents an RGB color.
func ColorAnsiRGB(r, g, b int) ColorAnsi {
    return ColorAnsi{
        Code:     []byte(fmt.Sprintf("38;2;%d;%d;%d", r, g, b)),
        Settings: []AnsiSetting{},
    }
}

// SetBackground returns a new ColorAnsi with the specified background color.
func (ac ColorAnsi) SetBackground(background ColorAnsiBackground) ColorAnsi {
    return ColorAnsi{
        Code:       ac.Code,
        Settings:   ac.Settings,
        Background: background,
    }
}

// Bold returns a new ColorAnsi with the Bold setting applied.
func (ac ColorAnsi) Bold() ColorAnsi {
    return ColorAnsi{
        Code:       ac.Code,
        Settings:   append(ac.Settings, ColorSettings.Bold),
        Background: ac.Background,
    }
}

// Dim returns a new ColorAnsi with the Dim setting applied.
func (ac ColorAnsi) Dim() ColorAnsi {
    return ColorAnsi{
        Code:       ac.Code,
        Settings:   append(ac.Settings, ColorSettings.Dim),
        Background: ac.Background,
    }
}

// Italic returns a new ColorAnsi with the Italic setting applied.
func (ac ColorAnsi) Italic() ColorAnsi {
    return ColorAnsi{
        Code:       ac.Code,
        Settings:   append(ac.Settings, ColorSettings.Italic),
        Background: ac.Background,
    }
}

// Underline returns a new ColorAnsi with the Underline setting applied.
func (ac ColorAnsi) Underline() ColorAnsi {
    return ColorAnsi{
        Code:       ac.Code,
        Settings:   append(ac.Settings, ColorSettings.Underline),
        Background: ac.Background,
    }
}

// SlowBlink returns a new ColorAnsi with the SlowBlink setting applied.
func (ac ColorAnsi) SlowBlink() ColorAnsi {
    return ColorAnsi{
        Code:       ac.Code,
        Settings:   append(ac.Settings, ColorSettings.Blink),
        Background: ac.Background,
    }
}

// Colorize returns the []byte content with the ANSI color applied.
//
// If the content is empty, an empty []byte is returned.
//
// A colorized byte array will not have the same length as the original byte array; it will be longer. This is because
// the colorized byte array will have the ANSI escape codes added to it, and the length of the escape codes will be
// different than the length of the original byte array.
//
// Colorization is always applied in the following order: ControlSequenceInitializer, Settings, Background, Code,
// AnsiEnd, CONTENT, AnsiResetSequence. Each section of the colorization is separated by the ansiCSSeparator byte
// (almost always a semicolon). Effectively, we're prefixing the content with the ANSI escape codes, and then
// resetting the ANSI escape codes after the content.
func (ac ColorAnsi) Colorize(content []byte) []byte {
    if len(content) == 0 {
        return content
    }

    buf := make([]byte, ac.totalBufferLength(content))
    cursor := 0

    copy(buf, ansiCSInit)
    cursor += len(ansiCSInit)

    for _, setting := range ac.Settings {
        copy(buf[cursor:], setting)
        cursor += len(setting)
        buf[cursor] = ansiCSSeparator
        cursor++
    }

    if len(ac.Background) > 0 {
        copy(buf[cursor:], ac.Background)
        cursor += len(ac.Background)
        buf[cursor] = ansiCSSeparator
        cursor++
    }

    copy(buf[cursor:], ac.Code)
    cursor += len(ac.Code)
    buf[cursor] = ansiCSEnd
    cursor++

    copy(buf[cursor:], content)
    cursor += len(content)

    copy(buf[cursor:], ansiReset)
    cursor += len(ansiReset)

    return buf
}

func (ac ColorAnsi) totalBufferLength(content []byte) int {
    settingsLength := 0
    for _, setting := range ac.Settings {
        settingsLength += len(setting) + 1
    }
    backgroundLength := 0
    if ac.Background != nil {
        backgroundLength = len(ac.Background) + 1
    }

    return len(ansiCSInit) + settingsLength + backgroundLength + len(ac.Code) + 1 + len(content) + len(ansiReset)
}
