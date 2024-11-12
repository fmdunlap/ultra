package log

// OutputFormat is a type representing the output format of a formatter.
//
// It can be one of the following:
//   - JSON
//   - Text
//
// TODO: Add more output formats [YAML, XML, etc.]
type OutputFormat string

const (
    OutputFormatJSON OutputFormat = "json"
    OutputFormatText OutputFormat = "text"
)

// LogLineArgs are the arguments that are passed to the FormatLogLine function of a LogLineFormatter, and further to the
// FieldFormatter function of a Field. Args are any format-level contextual information that may be needed to format a
// log field or log line.
type LogLineArgs struct {
    Level        Level
    Tag          string
    OutputFormat OutputFormat
}

// FormatResult is a struct that contains the formatted log line and any errors that may have occurred.
type FormatResult struct {
    bytes []byte
    err   error
}

// LogLineFormatter is an interface that defines a formatter for a log line. Implement this interface to create a
// custom formatter for your log lines if you need a specific format, or want to use ultralogger for a datatype that
// isn't built-in.
type LogLineFormatter interface {
    // FormatLogLine formats the log line using the provided data and returns a FormatResult which contains the
    // formatted log line and any errors that may have occurred.
    FormatLogLine(args LogLineArgs, data []any) FormatResult
}

// FormatterOption is a function that takes a LogLineFormatter and returns a new LogLineFormatter that has an option
// applied to it. This is useful for creating custom formatters that have additional options.
type FormatterOption func(f LogLineFormatter) LogLineFormatter

func NewFormatter(outputFormat OutputFormat, fields []Field, opts ...FormatterOption) (LogLineFormatter, error) {
    var f LogLineFormatter

    fieldFormatters := make(map[string]FieldFormatter)
    for _, field := range fields {
        fieldFormatter, err := field.NewFieldFormatter()
        if err != nil {
            return nil, &ErrorFieldFormatterInit{field: field, err: err}
        }
        fieldFormatters[field.Name()] = fieldFormatter
    }

    switch outputFormat {
    case OutputFormatJSON:
        f = &jsonFormatter{Fields: fields, FieldFormatters: fieldFormatters}
    case OutputFormatText:
        f = &textFormatter{Fields: fields, FieldFormatters: fieldFormatters}
    default:
        return nil, &ErrorInvalidOutput{outputFormat: outputFormat}
    }

    for _, opt := range opts {
        f = opt(f)
    }

    return f, nil
}

// WithDefaultColorization enables colorization for the formatter with the default colors.
//
// The default colors are ANSI 3-bit colors, and are compatible with most/virtually all terminals.
// See https://en.wikipedia.org/wiki/ANSI_escape_code#3-bit_and_4-bit for more information.
func WithDefaultColorization() FormatterOption {
    return func(f LogLineFormatter) LogLineFormatter {
        return NewColorizedFormatter(f, nil)
    }
}

// WithColorization enables colorization for the formatter with the provided colors.
//
// colors is a map of level to color. If a level is not present in the map, the default color for that level will be
// used.
func WithColorization(colors map[Level]Color) FormatterOption {
    return func(f LogLineFormatter) LogLineFormatter {
        return NewColorizedFormatter(f, colors)
    }
}
