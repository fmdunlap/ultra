package log

import (
    "fmt"
)

// TextFormatter is a formatter that formats log lines as text.
type TextFormatter struct {
    Fields         []Field
    FieldSeparator string
}

// TODO: Provide a way to specify the separator between fields.
// TODO: Provide a way to specify behavior on nil data.

// FormatLogLine formats the log line using the provided data and returns a FormatResult which contains the formatted
// log line and any errors that may have occurred.
func (f *TextFormatter) FormatLogLine(args LogLineArgs, data any) FormatResult {
    line := make([]byte, 0)
    args.OutputFormat = OutputFormatText

    for i, field := range f.Fields {
        fieldResult, err := computeFieldResult(field, args, data)
        if err != nil {
            return FormatResult{nil, &ErrorFieldFormatterInit{field: field, err: err}}
        }

        if fieldResult == nil {
            continue
        }

        resultBytes := fieldResult.Data
        if fieldResult.Data == nil {
            resultBytes = "<nil>"
        }

        if i < len(f.Fields)-1 {
            line = fmt.Append(line, resultBytes, " ")
        } else {
            line = fmt.Append(line, resultBytes)
        }
    }

    return FormatResult{line, nil}
}
