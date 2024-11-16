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
	// TODO: the formatter could panic... we should handle that nicely by logging an error about it, and exposing a
	//  setting to allow the user to squelch formatter panics. Hmmmm... Generally, I think the error handling of the
	//  logger should be configurable by the user. Do you want to allow a panic to propagate? Do you want to squelch and
	//  log? Do you want to disable a destination on panic? Two options: either add predefined set of 'on-panic'
	//  behaviors, create a panic handler interface that allows the user to define their own behavior. Leaning towards
	//  the former b/c we don't need every possible behavior to be configurable. The latter is more flexible, but
	//  requires more work, and adds complexity.
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

		// TODO: See above comment about processor panic handling.
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
