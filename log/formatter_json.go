package log

import (
    "encoding/json"
)

// jsonFormatter is a formatter that formats log lines as JSON.
type jsonFormatter struct {
    Fields          []Field // Keep these in an array to preserve the order of the fields.
    FieldFormatters map[string]FieldFormatter
}

// TODO: Provide a way to specify behavior on nil data. I.e. if the field should be omitted, or if we should include
//  a zero value, or something else. This is a bit tricky, because we don't know the type of the data, and we don't
//  know the type of the field.

// FormatLogLine formats the log line using the provided data and returns a FormatResult which contains the formatted
// log line and any errors that may have occurred.
func (f *jsonFormatter) FormatLogLine(args LogLineArgs, data []any) FormatResult {
    args.OutputFormat = OutputFormatJSON

    jsonMap := make(map[string]any)
    fieldResultChan := make(chan fieldProcessingResult)

    // Guaranteed to close on error result and once all fields have been processed.
    go processFieldsWithData(fieldResultChan, args, f.Fields, f.FieldFormatters, data)

    for {
        result, ok := <-fieldResultChan
        if !ok {
            break
        }

        if result.err != nil {
            return FormatResult{nil, result.err}
        }

        jsonMap[result.fieldName] = result.fieldData
    }

    jBytes, err := json.Marshal(jsonMap)
    return FormatResult{jBytes, err}
}
