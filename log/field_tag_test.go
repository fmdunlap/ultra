package log

import (
    "testing"
)

func TestTagField_FieldPrinter(t *testing.T) {
    tests := []struct {
        name     string
        tagField Field
        args     LogLineArgs
        want     string
        wantErr  bool
    }{
        {
            name:     "Default",
            tagField: NewDefaultTagField(),
            args: LogLineArgs{
                Level:        Info,
                Tag:          "test",
                OutputFormat: OutputFormatText,
            },
            want: "[test]",
        },
        {
            name:     "With Bracket",
            tagField: NewTagField(Brackets.Round, nil),
            args: LogLineArgs{
                Level:        Info,
                Tag:          "test",
                OutputFormat: OutputFormatText,
            },
            want: "(test)",
        },
        {
            name:     "With Padding",
            tagField: NewTagField(Brackets.Square, &TagPadSettings{PrefixPadSize: 1, SuffixPadSize: 2}),
            args: LogLineArgs{
                Level:        Info,
                Tag:          "test",
                OutputFormat: OutputFormatText,
            },
            want: " [test]  ",
        },
        {
            name:     "With Prefix Pad",
            tagField: NewTagField(Brackets.Square, &TagPadSettings{PrefixPadSize: 5}),
            args: LogLineArgs{
                Level:        Info,
                Tag:          "test",
                OutputFormat: OutputFormatText,
            },
            want: "     [test]",
        },
        {
            name:     "With Suffix Pad",
            tagField: NewTagField(Brackets.Square, &TagPadSettings{SuffixPadSize: 5}),
            args: LogLineArgs{
                Level:        Info,
                Tag:          "test",
                OutputFormat: OutputFormatText,
            },
            want: "[test]     ",
        },
        {
            name:     "With Pad Char",
            tagField: NewTagField(Brackets.Square, &TagPadSettings{PadChar: "!", PrefixPadSize: 1, SuffixPadSize: 2}),
            args: LogLineArgs{
                Level:        Info,
                Tag:          "test",
                OutputFormat: OutputFormatText,
            },
            want: "![test]!!",
        },
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            formatter, err := tt.tagField.NewFieldFormatter()
            if (err != nil) != tt.wantErr {
                t.Errorf("NewFieldFormatter() error = %v, wantErr %v", err, tt.wantErr)
                return
            }

            res, err := formatter(tt.args, nil)
            if (err != nil) != tt.wantErr {
                t.Errorf("NewFieldFormatter() error = %v, wantErr %v", err, tt.wantErr)
                return
            }

            if res.Data != tt.want {
                t.Errorf("NewFieldFormatter() formatted result = %v, want %v", res.Data, tt.want)
            }
        })
    }
}
