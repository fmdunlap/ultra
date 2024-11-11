package log

import (
    "context"
    "errors"
    "fmt"
    "io"
    "os"
    "time"
)

// Logger defines the interface for a structured ultraLogger in Go.
//
// This interface is useful for either creating your own logger or for using an existing logger, and preventing changes
// to the loggers formatting Settings.
type Logger interface {
    // Log logs at the specified level without formatting.
    Log(level Level, data any)

    // Debug logs a debug-level message.
    Debug(data any)

    // Info logs an info-level message.
    Info(data any)

    // Warn logs a warning-level message.
    Warn(data any)

    // Error logs an error-level message.
    Error(data any)

    // Panic logs a panic-level message and then panics.
    Panic(data any)

    // SetMinLevel sets the minimum logging level that will be output.
    SetMinLevel(level Level)

    // SetTag sets the tag for the logger.
    SetTag(tag string)

    Silence(enable bool)
}

const loglineTimeout = time.Millisecond * 250

var defaultDateTimeFormat = "2006-01-02 15:04:05"
var defaultLevelBracket = Brackets.Angle

var defaultCurrentTimeField, _ = NewCurrentTimeField("time", defaultDateTimeFormat)

var defaultFields = []Field{
    defaultCurrentTimeField,
    NewLevelField(defaultLevelBracket),
    &fieldMessage{},
}

func NewLoggerWithOptions(opts ...LoggerOption) (Logger, error) {
    l := newUltraLogger()

    for _, opt := range opts {
        if err := opt(l); err != nil {
            return nil, err
        }
    }

    if len(l.destinations) == 0 {
        defaultFormatter, _ := NewFormatter(OutputFormatText, defaultFields)
        l.destinations = map[io.Writer]LogLineFormatter{os.Stdout: defaultFormatter}
    }

    return l, nil
}

// NewLogger returns a new Logger that writes to stdout with the default text output format.
func NewLogger() Logger {
    formatter, _ := NewFormatter(OutputFormatText, defaultFields)

    logger, _ := NewLoggerWithOptions(WithStdoutFormatter(formatter))

    return logger
}

//NewFileLogger returns a new Logger that writes to a file.
//
//If the filename is empty, ErrorFileNotSpecified is returned.
//If the file does not exist, ErrorFileNotFound is returned.
func NewFileLogger(filename string, outputFormat OutputFormat) (Logger, error) {
    if filename == "" {
        return nil, ErrorFileNotSpecified
    }

    var err error
    filePtr, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
    if err != nil {
        if errors.Is(err, os.ErrNotExist) {
            return nil, &ErrorFileNotFound{filename: filename}
        }
        return nil, err
    }

    formatter, err := NewFormatter(outputFormat, defaultFields)
    if err != nil {
        return nil, err
    }

    fileLogger, err := NewLoggerWithOptions(WithDestination(filePtr, formatter))
    if err != nil {
        return nil, err
    }

    return fileLogger, nil
}

// ultraLogger is standard implementation of the /ultra/log Logger interface.
type ultraLogger struct {
    minLevel          Level
    destinations      map[io.Writer]LogLineFormatter
    tag               string
    silent            bool
    fallback          bool
    panicOnPanicLevel bool
    async             bool
}

func newUltraLogger() *ultraLogger {
    return &ultraLogger{
        minLevel:          Info,
        destinations:      map[io.Writer]LogLineFormatter{},
        silent:            false,
        fallback:          true,
        panicOnPanicLevel: false,
        async:             true,
    }
}

// Log logs a message with the given level and message.
func (l *ultraLogger) Log(level Level, data any) {
    if l.silent || level < l.minLevel {
        return
    }

    args := LogLineArgs{
        Level: level,
        Tag:   l.tag,
    }

    for w, f := range l.destinations {
        if f == nil {
            continue
        }

        if l.async {
            go l.writeLogLineAsync(w, f, args, data, loglineTimeout)
            continue
        }

        l.writeLogLine(w, f, args, data)
    }
}

// Debug logs a message with the Debug level and message.
func (l *ultraLogger) Debug(data any) {
    l.Log(Debug, data)
}

// Info logs a message with the Info level and message.
func (l *ultraLogger) Info(data any) {
    l.Log(Info, data)
}

// Warn logs a message with the Warn level and message.
func (l *ultraLogger) Warn(data any) {
    l.Log(Warn, data)
}

// Error logs a message with the Error level and message.
func (l *ultraLogger) Error(data any) {
    l.Log(Error, data)
}

// Panic logs a message with the Panic level and message. If panicOnPanicLevel is true, it panics.
func (l *ultraLogger) Panic(data any) {
    l.Log(Panic, data)

    if l.panicOnPanicLevel {
        panic(data)
    }
}

func (l *ultraLogger) SetMinLevel(level Level) {
    l.minLevel = level
}

func (l *ultraLogger) SetTag(tag string) {
    l.tag = tag
}

func (l *ultraLogger) Silence(enable bool) {
    l.silent = enable
}

// handleLogWriterError handles errors that occur while writing to the output. On failure, the log will fall back to
// writing to os.Stdout.
func (l *ultraLogger) handleLogWriterError(writer io.Writer, msgLevel Level, msg any, err error) {
    if !l.fallback || writer == os.Stdout {
        panic(err)
    }

    l.destinations[writer] = nil
    l.Error(
        fmt.Sprintf("error writing to original log writer, disabling formatter for writer: %v", err),
    )
    l.Log(msgLevel, msg)
}

func (l *ultraLogger) writeLogLine(
    w io.Writer,
    f LogLineFormatter,
    args LogLineArgs,
    data any,
) {
    formatResult := f.FormatLogLine(args, data)
    if formatResult.err != nil {
        l.Error(fmt.Sprintf("failed to format log line. formatter=%v, data=%v, err=%v", f, data, formatResult.err))
        return
    }

    writeResult := write(w, formatResult.bytes)
    if writeResult != nil {
        l.handleLogWriterError(w, args.Level, data, writeResult)
    }
}

func (l *ultraLogger) writeLogLineAsync(
    w io.Writer,
    f LogLineFormatter,
    args LogLineArgs,
    data any,
    timeout time.Duration,
) {
    ctx, cancel := context.WithTimeout(context.Background(), timeout)
    defer cancel()

    fmtChan := make(chan FormatResult, 1)
    go formatLogLineAsync(ctx, fmtChan, args, f, data)

    var logBytes []byte
    select {
    case result := <-fmtChan:
        if result.err != nil {
            l.Error(fmt.Sprintf("failed to format log line. formatter=%v, data=%v, err=%v", f, data, result.err))
            return
        }

        if len(result.bytes) == 0 {
            return
        }

        logBytes = result.bytes
    case <-ctx.Done():
        return
    }

    writeChan := make(chan error, 1)
    go writeLogLineAsync(ctx, writeChan, w, logBytes)

    select {
    case err := <-writeChan:
        if err != nil {
            l.handleLogWriterError(w, args.Level, data, err)
        }
    case <-ctx.Done():
        return
    }
}

func formatLogLineAsync(
    ctx context.Context,
    resultChan chan FormatResult,
    args LogLineArgs,
    formatter LogLineFormatter,
    data any,
) {
    defer close(resultChan)

    select {
    case <-ctx.Done():
        return
    case resultChan <- formatter.FormatLogLine(args, data):
    }
}

func writeLogLineAsync(
    ctx context.Context,
    resultChan chan error,
    w io.Writer,
    b []byte,
) {
    defer close(resultChan)

    select {
    case <-ctx.Done():
        return
    case resultChan <- write(w, b):
    }
}

func write(w io.Writer, b []byte) error {
    _, err := w.Write(append(b, '\n'))
    return err
}
