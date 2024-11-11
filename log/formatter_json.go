package log

import (
    "encoding/json"
)

// JSONFormatter is a formatter that formats log lines as JSON.
type JSONFormatter struct {
    Fields                 []Field
    destinationInitialized bool
}

// TODO: Provide a way to specify behavior on nil data. I.e. if the field should be omitted, or if we should include
//  a zero value, or something else. This is a bit tricky, because we don't know the type of the data, and we don't
//  know the type of the field.

// FormatLogLine formats the log line using the provided data and returns a FormatResult which contains the formatted
// log line and any errors that may have occurred.
func (f *JSONFormatter) FormatLogLine(args LogLineArgs, data any) FormatResult {
    jsonMap := make(map[string]any)

    args.OutputFormat = OutputFormatJSON

    for _, field := range f.Fields {
        fieldResult, err := computeFieldResult(field, args, data)
        if err != nil {
            return FormatResult{nil, err}
        }

        // Throw away fields that are nil or have nil data.
        if fieldResult == nil || fieldResult.Data == nil {
            continue
        }

        jsonMap[fieldResult.Name] = fieldResult.Data
    }

    jBytes, err := json.Marshal(jsonMap)
    return FormatResult{jBytes, err}
}
