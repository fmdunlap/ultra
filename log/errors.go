package log

import (
    "errors"
    "fmt"
)

type ErrorLoggerInitialization struct {
    err error
}

func (e *ErrorLoggerInitialization) Error() string {
    return fmt.Sprintf("error initializing logger: %v", e.err)
}

func (e *ErrorLoggerInitialization) Unwrap() error {
    return e.err
}

var ErrorFileNotSpecified = errors.New("filename not provided to NewFileLogger")

type ErrorFileNotFound struct {
    filename string
}

func (e *ErrorFileNotFound) Error() string {
    return fmt.Sprintf("file not found for FileLogger: %s", e.filename)
}

type ErrorMissingLevelColor struct {
    level Level
}

func (e *ErrorMissingLevelColor) Error() string {
    return fmt.Sprintf("missing color for level: %v", e.level)
}

type ErrorLevelParsing struct {
    level string
}

func (e *ErrorLevelParsing) Error() string {
    return fmt.Sprintf("invalid level: %s", e.level)
}

type ErrorFieldInitialization struct {
    fieldName string
    err       error
}

func (e *ErrorFieldInitialization) Error() string {
    return fmt.Sprintf("error initializing field: %v, err=%v", e.fieldName, e.err)
}

type ErrorFieldFormatterInit struct {
    field Field
    err   error
}

func (e *ErrorFieldFormatterInit) Error() string {
    return fmt.Sprintf("error formatting field: %v, err=%v", e.field, e.err)
}

func (e *ErrorFieldFormatterInit) Unwrap() error {
    return e.err
}

type ErrorInvalidOutput struct {
    outputFormat OutputFormat
}

func (e *ErrorInvalidOutput) Error() string {
    return fmt.Sprintf("invalid output format: %v", e.outputFormat)
}

type ErrorAmbiguousDestination struct{}

func (e *ErrorAmbiguousDestination) Error() string {
    return "formatters have ambiguous destinations"
}

type ErrorInvalidFieldDataType struct {
    field string
}

func (e *ErrorInvalidFieldDataType) Error() string {
    return fmt.Sprintf("invalid field data for field: %v", e.field)
}

var ErrorEmptyFieldName = errors.New("field name cannot be empty")

var ErrorNilFormatter = errors.New("formatter cannot be nil")

type ErrorMissingFieldFormatter struct {
    fieldName string
}

func (e *ErrorMissingFieldFormatter) Error() string {
    return fmt.Sprintf("missing field formatter for field: %v", e.fieldName)
}

func printSkippingFieldErr(fieldName string, err error) {
    fmt.Printf("WARNING: %s, not including field.\n", &ErrorFieldInitialization{fieldName, err})
}

type ErrorNonFatalFormatterError struct {
    fieldName string
    err       error
}

func (e *ErrorNonFatalFormatterError) Error() string {
    return fmt.Sprintf("non-fatal error formatting field: %v, err=%v", e.fieldName, e.err)
}

var ErrorTagFieldActiveButNoTag = errors.New("tag field is active but the logger has no tag set. disable the tag field, or add a tag to the logger")
