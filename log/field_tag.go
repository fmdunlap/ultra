package log

import (
    "strings"
)

type tagField struct {
    bracket     Bracket
    padSettings *TagPadSettings
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

var defaultTagPadSettings = &TagPadSettings{
    PadChar:       " ",
    PrefixPadSize: 0,
    SuffixPadSize: 0,
}

// NewDefaultTagField returns a new tag field with the default settings.
//
// The default settings are square brackets, and no padding.
func NewDefaultTagField() Field {
    return NewTagField(Brackets.Square, defaultTagPadSettings)
}

// NewTagField returns a new tag field with the provided settings.
func NewTagField(bracket Bracket, padSettings *TagPadSettings) Field {
    if padSettings == nil {
        padSettings = defaultTagPadSettings
    }

    if padSettings.PadChar == "" {
        padSettings.PadChar = " "
    }

    tf := &tagField{
        bracket:     bracket,
        padSettings: padSettings,
    }

    return tf
}

func (f *tagField) NewFieldFormatter() (FieldFormatter, error) {
    return f.format, nil
}

func (f *tagField) format(args LogLineArgs, _ any) (FieldResult, error) {
    result := FieldResult{
        Name: "tag",
    }

    switch args.OutputFormat {
    case OutputFormatText:
        result.Data = f.tagString(args.Tag)
    case OutputFormatJSON:
        result.Data = args.Tag
    }

    return result, nil
}

func (f *tagField) tagString(tag string) string {
    if tag == "" {
        return ""
    }

    b := strings.Builder{}

    b.WriteString(strings.Repeat(f.padSettings.PadChar, f.padSettings.PrefixPadSize))

    b.WriteString(f.bracket.Open())
    b.WriteString(tag)
    b.WriteString(f.bracket.Close())

    b.WriteString(strings.Repeat(f.padSettings.PadChar, f.padSettings.SuffixPadSize))

    return b.String()
}
