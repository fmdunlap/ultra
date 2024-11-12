package log

import "errors"

type fieldProcessingResult struct {
    fieldName     string
    fieldData     any
    fieldSettings FieldSettings
    err           error
}

func processFieldsWithData(
    resultChan chan fieldProcessingResult,
    args LogLineArgs,
    fields []Field,
    fieldFormatters map[string]FieldFormatter,
    data []any,
) {
    defer close(resultChan)

    processor := &fieldProcessor{
        args:        args,
        fields:      fields,
        formatters:  fieldFormatters,
        data:        data,
        matchedData: make([]bool, len(data)),
        resultChan:  resultChan,
    }

    processor.processAllFields()
}

type fieldProcessor struct {
    args        LogLineArgs
    fields      []Field
    formatters  map[string]FieldFormatter
    data        []any
    matchedData []bool
    resultChan  chan fieldProcessingResult
}

// TODO: Currently O(nlogn) for n fields. Worse if the user sends a ton of unmatchable data (more data than fields). Can
//  probably be optimized to O(n) by preprocessing matches on the data and then iterating over the fields in order. Need
//  to add better matching logic to determine which fields match which data.

func (p *fieldProcessor) processAllFields() {
    for _, field := range p.fields {
        if err := p.processField(field); err != nil {
            p.sendError(field.Name(), err)
            return
        }
    }
}

func (p *fieldProcessor) processField(field Field) error {
    formatter, err := p.getFormatter(field)
    if err != nil {
        return err
    }

    if field.Settings().AlwaysMatch {
        return p.processAlwaysMatchField(field, formatter)
    }

    return p.processDataMatchingField(field, formatter)
}

func (p *fieldProcessor) getFormatter(field Field) (FieldFormatter, error) {
    formatter, exists := p.formatters[field.Name()]
    if !exists {
        return nil, &ErrorFieldFormatterInit{field: field}
    }
    return formatter, nil
}

func (p *fieldProcessor) processAlwaysMatchField(field Field, formatter FieldFormatter) error {
    result, err := formatter(p.args, struct{}{})
    if err != nil {
        if p.handleProcessorError(field, err) {
            return nil
        }
        return err
    }

    if result != nil {
        p.sendResult(field, result)
    }
    return nil
}

func (p *fieldProcessor) processDataMatchingField(field Field, formatter FieldFormatter) error {
    for i, datum := range p.data {
        if p.matchedData[i] {
            continue
        }

        result, err := formatter(p.args, datum)
        if err != nil {
            if p.handleProcessorError(field, err) {
                continue
            }
            return err
        }

        // TODO: Add a mechanism for a field to disclaim a match even if the data type is a match. E.g. a field that
        //  matches on a string with a specific prefix. Currently it'll match to the first string field. Not always the
        //  desired behavior.

        if result != nil {
            p.matchedData[i] = true
            p.sendResult(field, result)
        }
    }
    return nil
}

func (p *fieldProcessor) handleProcessorError(field Field, err error) bool {
    nonFatalError := &ErrorNonFatalFormatterError{}
    InvalidFieldDataTypeError := &ErrorInvalidFieldDataType{}

    switch {
    case errors.As(err, &nonFatalError):
        p.sendResult(field, err.Error())
        return true
    case errors.As(err, &InvalidFieldDataTypeError):
        return true
    default:
        return false
    }
}

func (p *fieldProcessor) sendResult(field Field, data any) {
    p.resultChan <- fieldProcessingResult{
        fieldName:     field.Name(),
        fieldSettings: field.Settings(),
        fieldData:     data,
    }
}

func (p *fieldProcessor) sendError(fieldName string, err error) {
    p.resultChan <- fieldProcessingResult{
        fieldName: fieldName,
        err:       err,
    }
}

//func processFieldsWithData(lineProcChannel chan fieldProcessingResult, args LogLineArgs, fields []Field, fieldFormatters map[string]FieldFormatter, data []any) {
//    defer close(lineProcChannel)
//
//    dataAlreadyMatched := make([]bool, len(data))
//
//    for _, field := range fields {
//        fName := field.Name()
//        fSettings := field.Settings()
//        formatter, ok := fieldFormatters[fName]
//        if !ok {
//            lineProcChannel <- fieldProcessingResult{fieldName: fName, err: &ErrorFieldFormatterInit{field: field}}
//            return
//        }
//
//        // TODO: Currently O(nlogn) for n fields. Can probably be optimized to O(n) by preprocessing matches on the data
//        //  and then iterating over the fields in order. Need to add better matching logic to determine which fields
//        //  match which data.
//
//        // TODO: This also feels a bit hacky. There's a fair deal of repeated logic here if when we hit an AlwaysMatch.
//        //  Maybe it can be DRY'd up a bit, but not sure since it's just error handling. Could always add the field name
//        //  back to the result, but that feels a bit backwards since the field name is already present on the field
//        //  itself.
//        if fSettings.AlwaysMatch {
//            formattedData, err := formatter(args, struct{}{})
//
//            if err != nil {
//                var nonFatalError *ErrorNonFatalFormatterError
//                if errors.As(err, &nonFatalError) {
//                    lineProcChannel <- fieldProcessingResult{fieldName: fName, fieldData: err.Error()}
//                    continue
//                }
//                var invalidDataTypeError *ErrorInvalidFieldDataType
//                if errors.As(err, &invalidDataTypeError) {
//                    continue
//                }
//                lineProcChannel <- fieldProcessingResult{
//                    fieldName: fName,
//                    err:       err,
//                }
//                return
//            }
//
//            if formattedData == nil {
//                continue
//            }
//
//            lineProcChannel <- fieldProcessingResult{
//                fieldName:     fName,
//                fieldSettings: fSettings,
//                fieldData:     formattedData,
//            }
//            continue
//        }
//
//        for i, datum := range data {
//            if dataAlreadyMatched[i] {
//                continue
//            }
//
//            formattedData, err := formatter(args, datum)
//
//            if err != nil {
//                // If a non-fatal error occurrs, we'll return the error as the field data -- It's kinda a weird solution
//                // but I mean hey. This is a logger. We should log our errors using our own output.
//                var nonFatalError *ErrorNonFatalFormatterError
//                if errors.As(err, &nonFatalError) {
//                    lineProcChannel <- fieldProcessingResult{fieldName: fName, fieldData: err.Error()}
//                    continue
//                }
//                var invalidDataTypeError *ErrorInvalidFieldDataType
//                if errors.As(err, &invalidDataTypeError) {
//                    continue
//                }
//
//                lineProcChannel <- fieldProcessingResult{fieldName: fName, err: err}
//                return
//            }
//
//            // Assuming that the fields matched (even if the formattedData is nil)
//            dataAlreadyMatched[i] = true
//            lineProcChannel <- fieldProcessingResult{
//                fieldName:     fName,
//                fieldSettings: fSettings,
//                fieldData:     formattedData,
//            }
//        }
//    }
//}
