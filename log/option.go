package log

import (
    "io"
    "os"
)

// LoggerOption is a function that takes a Logger and returns a new Logger that has an option applied to it. This is
// useful for creating custom loggers that have additional options.
type LoggerOption func(l *ultraLogger) error

// WithMinLevel sets the minimum log level that will be output.
func WithMinLevel(level Level) LoggerOption {
    return func(l *ultraLogger) error {
        l.minLevel = level
        return nil
    }
}

// WithFields sets the fields for the logger.
func WithFields(writer io.Writer, fields []Field) LoggerOption {
    return func(l *ultraLogger) error {
        if l.destinations == nil {
            l.destinations = map[io.Writer]LogLineFormatter{}
        }
        formatter, err := NewFormatter(OutputFormatText, fields)
        if err != nil {
            return err
        }

        l.destinations[writer] = formatter

        return nil
    }
}

// WithStdoutFormatter sets the formatter to use for stdout.
// Note: This will not overwrite existing, non-stdout destinations, if any.
func WithStdoutFormatter(formatter LogLineFormatter) LoggerOption {
    return func(l *ultraLogger) error {
        if formatter == nil {
            return ErrorNilFormatter
        }
        if l.destinations == nil {
            l.destinations = map[io.Writer]LogLineFormatter{}
        }

        l.destinations[os.Stdout] = formatter
        return nil
    }
}

// WithDestination sets the destination for the logger. If the formatter is nil, the destination will be ignored.
// If the logger already has destinations, this will overwrite them.
func WithDestination(destination io.Writer, formatter LogLineFormatter) LoggerOption {
    return func(l *ultraLogger) error {
        if len(l.destinations) == 0 {
            l.destinations = map[io.Writer]LogLineFormatter{}
        }
        l.destinations[destination] = formatter
        return nil
    }
}

// WithDestinations sets the destinations for the logger. If the formatter is nil, the destination will be ignored.
// If the logger already has destinations, this will overwrite them.
func WithDestinations(destinations map[io.Writer]LogLineFormatter) LoggerOption {
    return func(l *ultraLogger) error {
        l.destinations = destinations
        return nil
    }
}

// WithSilent enables silent mode.
func WithSilent(silent bool) LoggerOption {
    return func(l *ultraLogger) error {
        l.silent = silent
        return nil
    }
}

// WithFallbackEnabled enables fallback to writing to os.Stdout.
func WithFallbackEnabled(fallback bool) LoggerOption {
    return func(l *ultraLogger) error {
        l.fallback = fallback
        return nil
    }
}

// WithPanicOnPanicLevel enables panic on panic level.
func WithPanicOnPanicLevel(panicOnPanicLevel bool) LoggerOption {
    return func(l *ultraLogger) error {
        l.panicOnPanicLevel = panicOnPanicLevel
        return nil
    }
}

// WithDefaultColorizationEnabled enables colorization for the formatter with the default colors.
//
// The default formatter will be used if no formatter has been set for the provided writer.
//
// The default colors are ANSI 3-bit colors, and are compatible with most/virtually all terminals.
// See https://en.wikipedia.org/wiki/ANSI_escape_code#3-bit_and_4-bit for more information.
func WithDefaultColorizationEnabled(writer io.Writer) LoggerOption {
    return func(l *ultraLogger) error {
        if len(l.destinations) == 0 {
            defaultFormatter, _ := NewFormatter(OutputFormatText, defaultFields)
            l.destinations = map[io.Writer]LogLineFormatter{writer: defaultFormatter}
        }

        l.destinations[writer] = NewColorizedFormatter(l.destinations[writer], nil)
        return nil
    }
}

// WithCustomColorization enables colorization for the formatter with the default colors.
//
// The default formatter will be used if no formatter has been set for the provided writer.
//
// The default colors are ANSI 3-bit colors, and are compatible with most/virtually all terminals.
// See https://en.wikipedia.org/wiki/ANSI_escape_code#3-bit_and_4-bit for more information.
func WithCustomColorization(writer io.Writer, colors map[Level]Color) LoggerOption {
    return func(l *ultraLogger) error {
        if l.destinations == nil {
            defaultFormatter, _ := NewFormatter(OutputFormatText, defaultFields)
            l.destinations = map[io.Writer]LogLineFormatter{writer: defaultFormatter}
        }

        l.destinations[writer] = NewColorizedFormatter(l.destinations[writer], colors)
        return nil
    }
}

// WithTag sets the tag for the logger.
func WithTag(tag string) LoggerOption {
    return func(l *ultraLogger) error {
        l.SetTag(tag)
        return nil
    }
}

// WithAsync enables async logging. Default=true.
//
// If async is true, the logger will write logs asynchronously. This is useful when writing to a file or a network
// connection, as it allows the logger to continue writing logs while
func WithAsync(async bool) LoggerOption {
    return func(l *ultraLogger) error {
        l.async = async
        return nil
    }
}
