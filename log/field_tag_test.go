package log

import (
    "testing"
)

func TestTagField_FieldPrinter(t *testing.T) {
    tests := []struct {
        name             string
        tagFieldSettings *TagFieldSettings
        args             LogLineArgs
        want             string
        wantErr          bool
    }{
        {
            name:             "Default",
            tagFieldSettings: nil,
            args: LogLineArgs{
                Level:        Info,
                Tag:          "test",
                OutputFormat: OutputFormatText,
            },
            want: "[test]",
        },
        {
            name: "With Bracket",
            tagFieldSettings: &TagFieldSettings{
                Bracket: Brackets.Round,
            },
            args: LogLineArgs{
                Level:        Info,
                Tag:          "test",
                OutputFormat: OutputFormatText,
            },
            want: "(test)",
        },
        {
            name: "With Padding",
            tagFieldSettings: &TagFieldSettings{
                Bracket: Brackets.Square,
                PadSettings: &TagPadSettings{
                    PrefixPadSize: 1,
                    SuffixPadSize: 2,
                },
            },
            args: LogLineArgs{
                Level:        Info,
                Tag:          "test",
                OutputFormat: OutputFormatText,
            },
            want: " [test]  ",
        },
        {
            name: "With Prefix Pad",
            tagFieldSettings: &TagFieldSettings{
                Bracket: Brackets.Square,
                PadSettings: &TagPadSettings{
                    PrefixPadSize: 5,
                },
            },
            args: LogLineArgs{
                Level:        Info,
                Tag:          "test",
                OutputFormat: OutputFormatText,
            },
            want: "     [test]",
        },
        {
            name: "With Suffix Pad",
            tagFieldSettings: &TagFieldSettings{
                Bracket: Brackets.Square,
                PadSettings: &TagPadSettings{
                    SuffixPadSize: 5,
                },
            },
            args: LogLineArgs{
                Level:        Info,
                Tag:          "test",
                OutputFormat: OutputFormatText,
            },
            want: "[test]     ",
        },
        {
            name: "With Pad Char",
            tagFieldSettings: &TagFieldSettings{
                Bracket: Brackets.Square,
                PadSettings: &TagPadSettings{
                    PadChar:       "!",
                    PrefixPadSize: 1,
                    SuffixPadSize: 2,
                },
            },
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
            tagField, err := NewTagField(tt.tagFieldSettings)
            if err != nil {
                t.Errorf("NewTagField() error = %v", err)
                return
            }

            formatter, err := tagField.NewFieldFormatter()
            if (err != nil) != tt.wantErr {
                t.Errorf("NewFieldFormatter() error = %v, wantErr %v", err, tt.wantErr)
                return
            }

            res, err := formatter(tt.args, struct{}{})
            if (err != nil) != tt.wantErr {
                t.Errorf("NewFieldFormatter() error = %v, wantErr %v", err, tt.wantErr)
                return
            }

            if res != tt.want {
                t.Errorf("NewFieldFormatter() formatted result = %v, want %v", res, tt.want)
            }
        })
    }
}
