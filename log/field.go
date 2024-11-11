package log

import (
    "fmt"
    "net/http"
    "strconv"
    "strings"
    "time"
)

// Field is an individual piece of data that is added to a log line. It can be a simple value, or an object.
//
// Most interaction with Fields should be done through the [NewObjectField] function.
type Field interface {
    // NewFieldFormatter returns a FieldFormatter, and an error if an error occurs while creating the FieldFormatter.
    NewFieldFormatter() (FieldFormatter, error)
}

// FieldResult is the result of formatting a field. It is used by the logger to annotate entries with metadata, such as
// the name of the field.
type FieldResult struct {
    Name string
    Data any
}

// FieldFormatter is a function that formats a field. It takes a LogLineArgs and the data to be formatted, and returns
// a FieldResult.
type FieldFormatter func(
    args LogLineArgs,
    data any,
) (FieldResult, error)

// ObjectField is a field that provides a formatter for a struct of type T.
type ObjectField[T any] struct {
    format FieldFormatter
}

// ObjectFieldFormatter is a function that formats a struct of type T and returns the formatted data. Note that this
// does not (presently) return a FieldResult, but it may in the future.
type ObjectFieldFormatter[T any] func(
    args LogLineArgs,
    data T,
) any

// NewFieldFormatter returns the FieldFormatter for the ObjectField. Typically, the FieldFormatter is computed when we
// create a new ObjectField with [NewObjectField], but a custom implementation can create the FieldFormatter at any
// time.
func (f ObjectField[T]) NewFieldFormatter() (FieldFormatter, error) {
    return f.format, nil
}

// NewObjectField returns a new ObjectField with the specified name and formatter. If the name is empty, an error is
// returned. If the formatter is nil, an error is returned.
//
// The formatter is a function that takes a LogLineArgs and a T, and returns a FieldResult. The FieldResult contains
// the name of the field and the formatted data that the logger will use to create the log line.
func NewObjectField[T any](name string, formatter ObjectFieldFormatter[T]) (ObjectField[T], error) {
    if name == "" {
        return ObjectField[T]{}, ErrorEmptyFieldName
    }
    if formatter == nil {
        return ObjectField[T]{}, ErrorNilFormatter
    }
    return ObjectField[T]{
        format: func(args LogLineArgs, data any) (FieldResult, error) {
            result := FieldResult{
                Name: name,
            }

            _, ok := data.(T)
            if !ok {
                return result, &ErrorInvalidFieldDataType{
                    field: name,
                }
            }
            result.Data = formatter(args, data.(T))

            return result, nil
        },
    }, nil
}

// NewStringField returns a new Field that formats a string into a string. The field will format the string using the
// String() method of the string.
//
// If the name is empty, an error is returned.
//
// Output Formats:
//  - All OutputFormats => remains unchanged.
func NewStringField(name string) (Field, error) {
    return NewObjectField[string](
        name,
        func(args LogLineArgs, data string) any {
            return data
        },
    )
}

// NewBoolField returns a new Field that formats a bool into a string. The field will format the bool using the
// Format() method of the bool.
//
// If the name is empty, an error is returned.
//
// OutputFormats:
//  - OutputFormatText => bool is formatted as a string with the format %v.
//  - OutputFormatJSON => bool is formatted as a bool.
func NewBoolField(name string) (Field, error) {
    return NewObjectField[bool](
        name,
        func(args LogLineArgs, data bool) any {
            if args.OutputFormat == OutputFormatText {
                if data {
                    return "true"
                }
                return "false"
            }
            return data
        },
    )
}

// NewTimeField returns a new Field that formats a time.Time into a string. The field will format the time using the
// Format() method of the time.Time.
//
// If the name is empty or the format is empty, an error is returned.
//
// OutputFormats:
//  - OutputFormatText => time.Time is formatted as a string with the format provided in the format argument.
//  - OutputFormatJSON => time.Time is formatted as a time.Time.
func NewTimeField(name, format string) (Field, error) {
    return NewObjectField[time.Time](
        name,
        func(args LogLineArgs, data time.Time) any {
            if args.OutputFormat == OutputFormatText {
                return data.Format(format)
            }
            return data
        },
    )
}

// NewIntField returns a new Field that formats an int.
//
// If the name is empty, an error is returned.
//
// OutputFormats:
//  - OutputFormatText => int is formatted as a string with the format %d using strconv.Itoa().
//  - OutputFormatJSON => int is formatted as a int.
func NewIntField(name string) (Field, error) {
    return NewObjectField[int](
        name,
        func(args LogLineArgs, data int) any {
            if args.OutputFormat == OutputFormatText {
                return strconv.Itoa(data)
            }
            return data
        },
    )
}

// NewFloatField returns a new Field that formats a float64.
//
// If the name is empty, an error is returned.
//
// OutputFormats:
//  - OutputFormatText => float64 is formatted as a string with the format '%f'.
//  - OutputFormatJSON => float64 is formatted as a float64.
func NewFloatField(name string) (Field, error) {
    return NewObjectField[float64](
        name,
        func(args LogLineArgs, data float64) any {
            if args.OutputFormat == OutputFormatText {
                return strconv.FormatFloat(data, 'f', -1, 64)
            }
            return data
        },
    )
}

// NewDurationField returns a new Field that formats a time.Duration.
//
// If the name is empty, an error is returned.
//
// OutputFormats:
//  - OutputFormatText => time.Duration is formatted as a string with the format %s.
//  - OutputFormatJSON => time.Duration is formatted as a time.Duration.
func NewDurationField(name string) (Field, error) {
    return NewObjectField[time.Duration](
        name,
        func(args LogLineArgs, data time.Duration) any {
            if args.OutputFormat == OutputFormatText {
                return data.String()
            }
            return data
        },
    )
}

// NewErrorField returns a new Field that formats an error into a string. The field will format the error using the
// Error() method of the error.
//
// If the name is empty, an error is returned.
//
// OutputFormats:
//  - OutputFormatText => error is formatted as a string with the format %v.
//  - OutputFormatJSON => error is formatted as a error.
func NewErrorField(name string) (Field, error) {
    return NewObjectField[error](
        name,
        func(args LogLineArgs, data error) any {
            if args.OutputFormat == OutputFormatText {
                return data.Error()
            }
            return data
        },
    )
}

// NewArrayField returns a new Field that formats a slice of type T into a slice of any. The field will format each
// element of the slice using the provided formatter.
//
// If the name is empty or the formatter is nil, an error is returned.
//
// OutputFormats:
//  - OutputFormatText => slice is formatted into a string with square brackets and comma separated elements. Each
//    element is formatted using the formatter. If the slice is empty, an empty string is returned. If the slice has
//    only one element, the element is returned in brackets.
//  - OutputFormatJSON => slice is formatted as a slice.
func NewArrayField[T any](name string, formatter ObjectFieldFormatter[T]) (Field, error) {
    if name == "" {
        return ObjectField[[]T]{}, ErrorEmptyFieldName
    }
    return NewObjectField[[]T](
        name,
        func(args LogLineArgs, data []T) any {
            res := make([]any, len(data))
            for i, v := range data {
                res[i] = formatter(args, v)
            }

            if args.OutputFormat == OutputFormatText {
                if len(res) == 0 {
                    return ""
                }
                stringRes := make([]string, len(res))
                for i, v := range res {
                    stringRes[i] = fmt.Sprintf("%v", v)
                }
                return fmt.Sprintf("[%s]", strings.Join(stringRes, ", "))
            }

            return res
        },
    )
}

// NewMapField returns a new Field that formats a map of type K and V into a map of K and V. The field will format each
// key and value of the map using the provided formatters.
//
// If the name is empty or the formatters are nil, an error is returned.
//
// OutputFormats:
//  - OutputFormatText => map is formatted into a string with curly brackets and comma separated key-value pairs. Each
//    key-value pair is formatted using the keyFormatter and valueFormatter. If the map is empty, an empty string is
//    returned. If the map has only one key-value pair, the key-value pair is returned in brackets.
//  - OutputFormatJSON => map is formatted as a map.
func NewMapField[K comparable, V any](name string, keyFormatter ObjectFieldFormatter[K], valueFormatter ObjectFieldFormatter[V]) (Field, error) {
    if name == "" {
        return ObjectField[map[K]V]{}, ErrorEmptyFieldName
    }
    if keyFormatter == nil {
        return ObjectField[map[K]V]{}, ErrorNilFormatter
    }
    if valueFormatter == nil {
        return ObjectField[map[K]V]{}, ErrorNilFormatter
    }
    return NewObjectField[map[K]V](
        name,
        func(args LogLineArgs, data map[K]V) any {
            res := make(map[any]any)
            for k, v := range data {
                res[keyFormatter(args, k)] = valueFormatter(args, v)
            }

            if args.OutputFormat != OutputFormatText {
                validMap := make(map[string]any)
                for k, v := range res {
                    validMap[fmt.Sprintf("%v", k)] = v
                }
                return validMap
            }

            return res
        },
    )

}

// NewCurrentTimeField returns a new Field that formats the current time into a string. The field will format the time
// using the provided format string.
//
// If the name is empty or the format is empty, an error is returned.
//
// OutputFormats:
//  - OutputFormatText => time is formatted as a string with the format provided in the format argument.
//  - OutputFormatJSON => time is formatted as a time.Time.
func NewCurrentTimeField(name, format string) (Field, error) {
    if name == "" {
        return &currentTimeField{}, ErrorEmptyFieldName
    }

    ctf := &currentTimeField{
        name:      name,
        fmtString: format,
        clock:     &realClock{},
    }

    return ctf, nil
}

type currentTimeField struct {
    name      string
    fmtString string
    clock     clock
}

func (f *currentTimeField) NewFieldFormatter() (FieldFormatter, error) {
    return f.format, nil
}

func (f *currentTimeField) format(args LogLineArgs, _ any) (FieldResult, error) {
    result := FieldResult{
        Name: f.name,
    }

    now := f.clock.Now()

    switch args.OutputFormat {
    case OutputFormatJSON:
        result.Data = now
    case OutputFormatText:
        result.Data = now.Format(f.fmtString)
    }

    return result, nil
}

// NewLevelField returns a new Field that formats a level into a string. The field will format the level using the
// String() method of the level.
//
// name: "level"
//
// If the bracket type is empty, the default bracket type is used.
//
// OutputFormats:
//  - OutputFormatText => level is formatted as a string with the format %v and wrapped in the bracket type.
//  - OutputFormatJSON => level is formatted as a level. Not wrapped in the bracket type.
//
// TODO: May want different behavior when serializing to non-text output formats. Currently we're returning the string
//  value of the Level. Do we want to keep the brackets? Or maybe we want to output the integer value of the level?
//  Maybe we just want to make the whole thing configurable? ¯\_(ツ)_/¯
func NewLevelField(bracket Bracket) Field {
    return &levelField{
        bracket: bracket,
    }
}

type levelField struct {
    bracket      Bracket
    levelStrings map[Level]string
}

func (f *levelField) NewFieldFormatter() (FieldFormatter, error) {
    if f.levelStrings == nil {
        f.levelStrings = make(map[Level]string)

        for _, lvl := range AllLevels() {
            f.levelStrings[lvl] = f.bracket.Wrap(lvl.String())
        }
    }

    return f.format, nil
}

func (f *levelField) format(args LogLineArgs, _ any) (FieldResult, error) {
    if args.OutputFormat == OutputFormatText {
        return FieldResult{
            Name: "level",
            Data: f.levelStrings[args.Level],
        }, nil
    }

    return FieldResult{
        Name: "level",
        Data: args.Level.String(),
    }, nil
}

// NewMessageField returns a new Field that formats a message into a string. The field will format the message using the
// String() method of the message.
//
// name: "message"
//
// OutputFormats:
//  - OutputFormatText => message is formatted as a string with the format %v.
//  - OutputFormatJSON => message is formatted as a message.
func NewMessageField() Field {
    return &fieldMessage{}
}

type fieldMessage struct{}

func (f *fieldMessage) NewFieldFormatter() (FieldFormatter, error) {
    return f.format, nil
}

func (f *fieldMessage) format(_ LogLineArgs, message any) (FieldResult, error) {
    result := FieldResult{
        Name: "message",
    }

    switch message.(type) {
    case string:
        result.Data = message.(string)
    case fmt.Stringer:
        result.Data = message.(fmt.Stringer).String()
    default:
        return result, &ErrorInvalidFieldDataType{
            field: "message",
        }
    }

    return result, nil
}

// NewRequestField returns a new Field that formats an http.Request into a string. The field will format the request
// using the provided settings [RequestFieldSettings].
//
// If the name is empty or the settings are nil, an error is returned.
//
// OutputFormats:
//  - OutputFormatText => request is formatted as a string. Http request fields are included based on the settings
//    [RequestFieldSettings]. Included fields are returned as a space separated string with key=value elements. Returns
//    an empty string if [RequestFieldSettings] has no true fields.
//  - OutputFormatJSON => [RequestLogEntry].
func NewRequestField(name string, settings RequestFieldSettings) (Field, error) {
    return NewObjectField[*http.Request](
        name,
        func(args LogLineArgs, data *http.Request) any {
            logEntry := RequestLogEntry{}

            if settings.LogReceivedAt {
                logEntry.ReceivedAt = time.Now()
            }

            if settings.LogSourceIP {
                logEntry.SourceIP = data.RemoteAddr
            }

            if settings.LogMethod {
                logEntry.Method = data.Method
            }

            if settings.LogPath {
                logEntry.Path = data.URL.Path
            }

            if args.OutputFormat == OutputFormatText {
                return logEntry.String(settings.TimeFormat)
            }
            return logEntry
        },
    )
}

// RequestFieldSettings is a struct that contains settings for the RequestField.
//
// The settings are used to determine which fields of the http.Request struct to include in the formatted output, as
// well as the format to use for the fields.
//
// If the time format is empty, the default time format is used.
type RequestFieldSettings struct {
    // TimeFormat is the format to use for the ReceivedAt field.
    TimeFormat string

    // LogReceivedAt determines whether to include the ReceivedAt field in the formatted output.
    LogReceivedAt bool
    // LogMethod determines whether to include the Method field in the formatted output.
    LogMethod bool
    // LogPath determines whether to include the Path field in the formatted output.
    LogPath bool
    // LogSourceIP determines whether to include the SourceIP field in the formatted output.
    LogSourceIP bool
}

// RequestLogEntry is a struct that represents a formatted http.Request.
type RequestLogEntry struct {
    ReceivedAt time.Time
    Method     string
    Path       string
    SourceIP   string
}

func (r *RequestLogEntry) String(timeFmt string) string {
    parts := []string{}
    if !r.ReceivedAt.IsZero() {
        parts = append(parts, r.ReceivedAt.Format(timeFmt))
    }
    if r.Method != "" {
        parts = append(parts, r.Method)
    }
    if r.Path != "" {
        parts = append(parts, r.Path)
    }
    if r.SourceIP != "" {
        parts = append(parts, r.SourceIP)
    }
    return strings.Join(parts, " ")
}

// NewResponseField returns a new Field that formats an http.Response into a string. The field will format the response
// using the provided settings [ResponseFieldSettings].
//
// An error is returned if the name is empty or the settings are nil.
//
// OutputFormats:
//  - OutputFormatText => response is formatted as a string. http.Response fields are included based on the settings
//    [ResponseFieldSettings]. Included fields are returned as a space separated string with key=value elements. Returns
//    an empty string if [RequestFieldSettings] has no true fields.
//  - OutputFormatJSON => [ResponseLogEntry].
func NewResponseField(name string, settings ResponseFieldSettings) (Field, error) {
    return NewObjectField[*http.Response](
        name,
        func(args LogLineArgs, data *http.Response) any {
            logEntry := ResponseLogEntry{}

            if settings.LogStatus {
                logEntry.Status = data.Status
            }

            if settings.LogStatusCode {
                logEntry.StatusCode = data.StatusCode
            }

            if settings.LogPath {
                logEntry.Path = data.Request.URL.Path
            }

            if args.OutputFormat == OutputFormatText {
                return logEntry.String()
            }
            return logEntry
        },
    )
}

type ResponseFieldSettings struct {
    // LogStatus determines whether to include the http.Response.Status field in the formatted output.
    LogStatus bool
    // LogStatusCode determines whether to include the http.Response.StatusCode field in the formatted output.
    LogStatusCode bool
    // LogPath determines whether to include the associated http.Request.URL.Path field in the formatted output.
    LogPath bool
}

type ResponseLogEntry struct {
    StatusCode int
    Status     string
    Path       string
}

func (r *ResponseLogEntry) String() string {
    parts := []string{}
    if r.StatusCode != 0 {
        parts = append(parts, strconv.Itoa(r.StatusCode))
    }
    return strings.Join(parts, " ")
}
