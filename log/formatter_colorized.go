package log

import (
    "maps"
)

var defaultLevelColors = map[Level]Color{
    Debug: Colors.Green,
    Info:  Colors.White,
    Warn:  Colors.Yellow,
    Error: Colors.Red,
    Panic: Colors.Magenta,
}

// ColorizedFormatter colorizes the bytes of the base formatter using the provided colors.
type ColorizedFormatter struct {
    BaseFormatter LogLineFormatter
    LevelColors   map[Level]Color
}

// FormatLogLine formats the log line using the provided data and returns a FormatResult which contains the formatted
// log line and any errors that may have occurred.
func (f *ColorizedFormatter) FormatLogLine(args LogLineArgs, data []any) FormatResult {
    res := f.BaseFormatter.FormatLogLine(args, data)
    if res.err != nil {
        return res
    }

    color, ok := f.LevelColors[args.Level]
    if !ok {
        return FormatResult{res.bytes, &ErrorMissingLevelColor{level: args.Level}}
    }

    return FormatResult{color.Colorize(res.bytes), nil}
}

// NewColorizedFormatter returns a new ColorizedFormatter that formats the provided base formatter with the provided
// colors.
func NewColorizedFormatter(baseFormatter LogLineFormatter, levelColors map[Level]Color) *ColorizedFormatter {
    if levelColors == nil {
        levelColors = make(map[Level]Color)
        maps.Copy(levelColors, defaultLevelColors)
    }

    return &ColorizedFormatter{
        BaseFormatter: baseFormatter,
        LevelColors:   levelColors,
    }
}
