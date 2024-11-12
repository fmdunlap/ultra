package log

import (
    "fmt"
    "maps"
    "net/http"
    "strconv"
    "strings"
    "time"
)

var defaultDateTimeFormat = "2006-01-02 15:04:05"

// Field is an individual piece of data that is added to a log line. It can be a simple value, or an object.
//
// Most interaction with Fields should be done through the [NewObjectField] function.
type Field interface {
    // NewFieldFormatter returns a FieldFormatter, and an error if an error occurs while creating the FieldFormatter.
    NewFieldFormatter() (FieldFormatter, error)
    // Name returns the name of the field. This is used to identify the field in the formatter.
    Name() string
    // Settings returns the options for the field.
    Settings() FieldSettings
}

type FieldSettings struct {
    HideKey     bool
    AlwaysMatch bool
}

// FieldFormatter is a function that formats a field. It takes a LogLineArgs and the data to be formatted, and returns
// a FieldResult.
type FieldFormatter func(
    args LogLineArgs,
    data any,
) (any, error)

// TODO: Consider adding positioning control to fields. I.e. AlwaysLast or AlwaysFirst. That way we could have a field,
//  like the message field, that always comes last, and a field like the level or current time field that always comes
//  first.

// ObjectField is a field that provides a formatter for a struct of type T.
type ObjectField[T any] struct {
    // name is the name of the field.
    name string
    // options are the options for the field.
    options FieldSettings
    // format is the formatter for the field.
    format FieldFormatter
}

// NewFieldFormatter returns the FieldFormatter for the ObjectField. Typically, the FieldFormatter is computed when we
// create a new ObjectField with [NewObjectField], but a custom implementation can create the FieldFormatter at any
// time.
func (f ObjectField[T]) NewFieldFormatter() (FieldFormatter, error) {
    return f.format, nil
}

// Name returns the name of the field.
func (f ObjectField[T]) Name() string {
    return f.name
}

// FieldSettings returns the options for the field.
func (f ObjectField[T]) Settings() FieldSettings {
    return f.options
}

// ObjectFieldFormatter is a function that formats a struct of type T and returns the formatted data. Note that this
// does not (presently) return a FieldResult, but it may in the future.
type ObjectFieldFormatter[T any] func(
    args LogLineArgs,
    data T,
) (any, error)

// NewObjectField returns a new ObjectField with the specified name and formatter. If the name is empty, an error is
// returned. If the formatter is nil, an error is returned.
//
// The formatter is a function that takes a LogLineArgs and a T, and returns a FieldResult. The FieldResult contains
// the name of the field and the formatted data that the logger will use to create the log line.
func NewObjectField[T any](name string, formatter ObjectFieldFormatter[T], opts ...FieldOption) (ObjectField[T], error) {
    if name == "" {
        return ObjectField[T]{}, ErrorEmptyFieldName
    }
    if formatter == nil {
        return ObjectField[T]{}, ErrorNilFormatter
    }

    objectField := ObjectField[T]{
        name: name,
    }

    for _, opt := range opts {
        if err := opt(&objectField.options); err != nil {
            return ObjectField[T]{}, err
        }
    }

    objectField.format = func(args LogLineArgs, data any) (any, error) {
        if _, ok := data.(T); !ok {
            return nil, &ErrorInvalidFieldDataType{
                field: name,
            }
        }

        return formatter(args, data.(T))
    }

    return objectField, nil
}

type FieldOption func(f *FieldSettings) error

func WithHideKey(hideKey bool) FieldOption {
    return func(s *FieldSettings) error {
        s.HideKey = hideKey
        return nil
    }
}

func WithAlwaysMatch(formatWithoutData bool) FieldOption {
    return func(s *FieldSettings) error {
        s.AlwaysMatch = formatWithoutData
        return nil
    }
}

type LineArgsField struct {
    name   string
    format FieldFormatter
}

type LineArgsFormatter func(args LogLineArgs) (any, error)

func NewLineArgsField(name string, formatter LineArgsFormatter) (Field, error) {
    return &LineArgsField{
        name: name,
        format: func(args LogLineArgs, _ any) (any, error) {
            return formatter(args)
        },
    }, nil
}

func (f *LineArgsField) Name() string {
    return f.name
}

func (f *LineArgsField) Settings() FieldSettings {
    return FieldSettings{
        HideKey:     true,
        AlwaysMatch: true,
    }
}

func (f *LineArgsField) NewFieldFormatter() (FieldFormatter, error) {
    return f.format, nil
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
        func(args LogLineArgs, data string) (any, error) {
            return data, nil
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
        func(args LogLineArgs, data bool) (any, error) {
            if args.OutputFormat != OutputFormatText {
                return data, nil
            }
            if data {
                return "true", nil
            }
            return "false", nil
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
        func(args LogLineArgs, data time.Time) (any, error) {
            if args.OutputFormat == OutputFormatText {
                return data.Format(format), nil
            }
            return data, nil
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
        func(args LogLineArgs, data int) (any, error) {
            if args.OutputFormat == OutputFormatText {
                return strconv.Itoa(data), nil
            }
            return data, nil
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
        func(args LogLineArgs, data float64) (any, error) {
            if args.OutputFormat == OutputFormatText {
                return strconv.FormatFloat(data, 'f', -1, 64), nil
            }
            return data, nil
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
        func(args LogLineArgs, data time.Duration) (any, error) {
            if args.OutputFormat == OutputFormatText {
                return data.String(), nil
            }
            return data, nil
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
        func(args LogLineArgs, data error) (any, error) {
            if args.OutputFormat == OutputFormatText {
                return data.Error(), nil
            }
            return data, nil
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
        func(args LogLineArgs, data []T) (any, error) {
            res := make([]any, len(data))
            var err error
            for i, v := range data {
                res[i], err = formatter(args, v)
                if err != nil {
                    return nil, err
                }
            }

            if args.OutputFormat == OutputFormatText {
                if len(res) == 0 {
                    return "", nil
                }
                stringRes := make([]string, len(res))
                for i, v := range res {
                    stringRes[i] = fmt.Sprintf("%v", v)
                }
                return fmt.Sprintf("[%s]", strings.Join(stringRes, ", ")), nil
            }

            return res, err
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
        func(args LogLineArgs, data map[K]V) (any, error) {
            res := make(map[any]any)
            for k, v := range data {
                key, err := keyFormatter(args, k)
                if err != nil {
                    return nil, err
                }
                value, err := valueFormatter(args, v)
                if err != nil {
                    return nil, err
                }
                res[key] = value
            }

            // At least for JSON (the only currently non-text output format), we need to return a map[string]any.
            // Otherwise, the JSON formatter will try to marshal the map[any]any into JSON, which will fail.
            if args.OutputFormat != OutputFormatText {
                validMap := make(map[string]any)
                for k, v := range res {
                    validMap[fmt.Sprintf("%v", k)] = v
                }
                return validMap, nil
            }

            return res, nil
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
func NewCurrentTimeField(settings *CurrentTimeFieldSettings) Field {
    if settings == nil {
        settings = &CurrentTimeFieldSettings{}
    }
    settings.MergeDefault()

    currentTimeField, err := NewLineArgsField(
        settings.Name,
        func(args LogLineArgs) (any, error) {
            now := time.Now()

            // This would be better if we could inject a fake clock into the field formatter. As is we're wasting a
            // compare operation here.
            if settings.fakeNow != nil {
                now = *settings.fakeNow
            }

            switch args.OutputFormat {
            case OutputFormatJSON:
                return now, nil
            case OutputFormatText:
                return now.Format(settings.Format), nil
            }

            return nil, nil
        },
    )

    if err != nil {
        printSkippingFieldErr(settings.Name, err)
        return nil
    }

    return currentTimeField
}

func NewDefaultCurrentTimeField() Field {
    return NewCurrentTimeField(nil)
}

type CurrentTimeFieldSettings struct {
    // Name is the name of the field.
    Name string
    // Format is the format to use for the current time field.
    Format string

    // for testing
    fakeNow *time.Time
}

var defaultCurrentTimeFieldSettings = CurrentTimeFieldSettings{
    Name:   "currentTime",
    Format: defaultDateTimeFormat,
}

func (s *CurrentTimeFieldSettings) MergeDefault() {
    if s.Name == "" {
        s.Name = defaultCurrentTimeFieldSettings.Name
    }
    if s.Format == "" {
        s.Format = defaultCurrentTimeFieldSettings.Format
    }
}

// TODO: May want different behavior when serializing to non-text output formats. Currently we're returning the string
//  value of the Level. Do we want to keep the brackets? Or maybe we want to output the integer value of the level?
//  Maybe we just want to make the whole thing configurable? ¯\_(ツ)_/¯

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
func NewLevelField(settings *LevelFieldSettings) Field {
    if settings == nil {
        settings = &LevelFieldSettings{}
    }
    settings.MergeDefault()

    textLevelStrings := make(map[any]string)

    // Merge guarantees that there will always be a level string for each level.
    for _, lvl := range AllLevels() {
        textLevelStrings[lvl] = settings.Bracket.Wrap(settings.StringsForLevels[lvl])
    }

    levelField, err := NewLineArgsField(
        settings.Name,
        func(args LogLineArgs) (any, error) {
            if args.OutputFormat == OutputFormatText {
                return textLevelStrings[args.Level], nil
            }
            return settings.StringsForLevels[args.Level], nil
        },
    )

    if err != nil {
        printSkippingFieldErr(settings.Name, err)
        return nil
    }

    return levelField
}

func NewDefaultLevelField() Field {
    return NewLevelField(nil)
}

var defaultLevelStrings = map[Level]string{
    Debug: Debug.String(),
    Info:  Info.String(),
    Warn:  Warn.String(),
    Error: Error.String(),
    Panic: Panic.String(),
}

type LevelFieldSettings struct {
    Name             string
    Bracket          Bracket
    StringsForLevels map[Level]string
}

var defaultLevelFieldSettings = LevelFieldSettings{
    Name:             "level",
    Bracket:          Brackets.Angle,
    StringsForLevels: maps.Clone(defaultLevelStrings),
}

func (s *LevelFieldSettings) MergeDefault() {
    if s.Name == "" {
        s.Name = defaultLevelFieldSettings.Name
    }

    if s.Bracket == nil {
        s.Bracket = defaultLevelFieldSettings.Bracket
    }

    if s.StringsForLevels == nil {
        s.StringsForLevels = defaultLevelFieldSettings.StringsForLevels
    }
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
    msgFieldName := "message"

    msgField, err := NewObjectField[string](
        msgFieldName,
        func(args LogLineArgs, msg string) (any, error) {
            return msg, nil
        },
        WithHideKey(true),
    )

    if err != nil {
        printSkippingFieldErr(msgFieldName, err)
        return nil
    }

    return msgField
}

// NewTagField returns a new Field for the logger tag. The field will format the tag using the provided settings.
// If the logger has no tag, the field will return an empty string.
//
// If the name is empty, an error is returned.
//
// OutputFormats:
//  - OutputFormatText => tag is formatted as a string with the format %v.
//  - OutputFormatJSON => tag is formatted as a tag.
func NewTagField(settings *TagFieldSettings) (Field, error) {
    if settings == nil {
        settings = &TagFieldSettings{}
    }
    settings.MergeDefault()

    tagFmtString := buildTagFormatString(settings.Bracket, settings.PadSettings)

    return NewLineArgsField(
        settings.Name,
        func(args LogLineArgs) (any, error) {
            if args.Tag == "" {
                return "", &ErrorNonFatalFormatterError{settings.Name, ErrorTagFieldActiveButNoTag}
            }

            if args.OutputFormat == OutputFormatText {
                return fmt.Sprintf(tagFmtString, args.Tag), nil
            }
            return args.Tag, nil
        },
    )
}

func NewDefaultTagField() Field {
    f, _ := NewTagField(nil)
    return f
}

func buildTagFormatString(bracket Bracket, padSettings *TagPadSettings) string {
    b := strings.Builder{}

    if padSettings != nil && padSettings.PadChar != "" {
        b.WriteString(strings.Repeat(padSettings.PadChar, padSettings.PrefixPadSize))
    }

    b.WriteString(bracket.Open())
    b.WriteString("%s")
    b.WriteString(bracket.Close())

    if padSettings != nil && padSettings.PadChar != "" {
        b.WriteString(strings.Repeat(padSettings.PadChar, padSettings.SuffixPadSize))
    }

    return b.String()
}

// TagFieldSettings are the settings for the TagField.
type TagFieldSettings struct {
    // Name is the name of the field.
    Name string
    // Bracket is the bracket type to use for the tag field.
    Bracket Bracket
    // PadSettings are the settings for padding the tag field.
    PadSettings *TagPadSettings
}

var defaultTagFieldSettings = TagFieldSettings{
    Name:    "tag",
    Bracket: Brackets.Square,
}

// TagPadSettings are the settings for padding a tag field. If PadChar is empty, it will default to a space.
// Note: for non-text formatters the padding setting may be ignored (it is in the built in JSON formatter).
type TagPadSettings struct {
    // PadChar is the character to use for padding. If empty, it will default to a space.
    PadChar string
    // PrefixPadSize is the number of times PadChar will be added before the tag.
    PrefixPadSize int
    // SuffixPadSize is the number of times PadChar will be added after the tag.
    SuffixPadSize int
}

var defaultTagPadSettings = TagPadSettings{
    PadChar:       " ",
    PrefixPadSize: 0,
    SuffixPadSize: 0,
}

func (s *TagFieldSettings) MergeDefault() {
    if s.Name == "" {
        s.Name = defaultTagFieldSettings.Name
    }
    if s.Bracket == nil {
        s.Bracket = defaultTagFieldSettings.Bracket
    }
    if s.PadSettings == nil {
        s.PadSettings = &TagPadSettings{}
    }
    if s.PadSettings.PadChar == "" {
        s.PadSettings.PadChar = defaultTagPadSettings.PadChar
    }
    if s.PadSettings.PrefixPadSize == 0 {
        s.PadSettings.PrefixPadSize = defaultTagPadSettings.PrefixPadSize
    }
    if s.PadSettings.SuffixPadSize == 0 {
        s.PadSettings.SuffixPadSize = defaultTagPadSettings.SuffixPadSize
    }
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
func NewRequestField(settings *RequestFieldSettings) (Field, error) {
    settings = defaultRequestFieldSettings.Merge(settings)

    return NewObjectField[*http.Request](
        settings.Name,
        func(args LogLineArgs, data *http.Request) (any, error) {
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
                return logEntry.String(settings.TimeFormat), nil
            }
            return logEntry, nil
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
    // Name is the name of the field.
    Name string

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

var defaultRequestFieldSettings = RequestFieldSettings{
    Name:          "request",
    TimeFormat:    defaultDateTimeFormat,
    LogReceivedAt: false,
    LogMethod:     true,
    LogPath:       true,
    LogSourceIP:   false,
}

func (s *RequestFieldSettings) Merge(other *RequestFieldSettings) *RequestFieldSettings {
    if other.Name != "" {
        s.Name = other.Name
    }
    if other.TimeFormat != "" {
        s.TimeFormat = other.TimeFormat
    }
    if other.LogReceivedAt {
        s.LogReceivedAt = other.LogReceivedAt
    }
    if other.LogMethod {
        s.LogMethod = other.LogMethod
    }
    if other.LogPath {
        s.LogPath = other.LogPath
    }
    if other.LogSourceIP {
        s.LogSourceIP = other.LogSourceIP
    }

    return s
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
func NewResponseField(settings *ResponseFieldSettings) (Field, error) {
    settings = defaultResponseFieldSettings.Merge(settings)

    return NewObjectField[*http.Response](
        settings.Name,
        func(args LogLineArgs, data *http.Response) (any, error) {
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
                return logEntry.String(), nil
            }
            return logEntry, nil
        },
    )
}

type ResponseFieldSettings struct {
    // Name is the name of the field.
    Name string
    // LogStatus determines whether to include the http.Response.Status field in the formatted output.
    LogStatus bool
    // LogStatusCode determines whether to include the http.Response.StatusCode field in the formatted output.
    LogStatusCode bool
    // LogPath determines whether to include the associated http.Request.URL.Path field in the formatted output.
    LogPath bool
}

var defaultResponseFieldSettings = ResponseFieldSettings{
    Name:          "response",
    LogStatus:     true,
    LogStatusCode: false,
    LogPath:       true,
}

func (s *ResponseFieldSettings) Merge(other *ResponseFieldSettings) *ResponseFieldSettings {
    if other == nil {
        return s
    }

    if other.Name != "" {
        s.Name = other.Name
    }
    if other.LogStatus {
        s.LogStatus = other.LogStatus
    }
    if other.LogStatusCode {
        s.LogStatusCode = other.LogStatusCode
    }
    if other.LogPath {
        s.LogPath = other.LogPath
    }

    return s
}

type ResponseLogEntry struct {
    StatusCode int
    Status     string
    Path       string
}

func (r *ResponseLogEntry) String() string {
    parts := make([]string, 0)
    if r.StatusCode != 0 {
        parts = append(parts, strconv.Itoa(r.StatusCode))
    }
    if r.Status != "" {
        parts = append(parts, r.Status)
    }
    if r.Path != "" {
        parts = append(parts, r.Path)
    }
    return strings.Join(parts, " ")
}
