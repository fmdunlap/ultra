package log

import (
    "fmt"
    "strings"
)

// textFormatter is a formatter that formats log lines as text.
type textFormatter struct {
    Fields          []Field                   // Keep these in an array to preserve the order of the fields.
    FieldFormatters map[string]FieldFormatter // Map of the field name to its formatter
    FieldSeparator  string
}

// TODO: Provide a way to specify the separator between fields.
// TODO: Provide a way to specify behavior on nil data.

// FormatLogLine formats the log line using the provided data and returns a FormatResult which contains the formatted
// log line and any errors that may have occurred.
func (f *textFormatter) FormatLogLine(args LogLineArgs, data []any) FormatResult {
    args.OutputFormat = OutputFormatText

    line := make([]byte, 0)
    procResChan := make(chan fieldProcessingResult)

    go processFieldsWithData(procResChan, args, f.Fields, f.FieldFormatters, data)
    for {
        result, ok := <-procResChan
        if !ok {
            break
        }

        if result.err != nil {
            return FormatResult{nil, result.err}
        }

        line = f.addDataToLogLine(line, result.fieldData, result.fieldName, result.fieldSettings)
    }

    if len(line) > 0 {
        line = line[:len(line)-1]
    }

    return FormatResult{line, nil}
}

func (f *textFormatter) addDataToLogLine(line []byte, resultBytes any, fName string, fSettings FieldSettings) []byte {
    b := strings.Builder{}

    if !fSettings.HideKey {
        b.WriteString(fName)
        b.WriteString("=")
    }

    b.WriteString(fmt.Sprintf("%v", resultBytes))

    b.WriteString(" ")

    return fmt.Append(line, b.String())
}
