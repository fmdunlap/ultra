package log

import (
    "bytes"
    "errors"
    "fmt"
    "os"
    "testing"
)

func ExampleNewFormatter() {
    formatter, _ := NewFormatter(OutputFormatText, []Field{
        NewLevelField(Brackets.Angle),
        NewMessageField(),
    })

    logger, _ := NewLoggerWithOptions(WithDestination(os.Stdout, formatter), WithAsync(false))

    logger.Info("This is an info message.")
    // Output: <INFO> This is an info message.
}

func ExampleNewFormatter_jSON() {
    formatter, _ := NewFormatter(OutputFormatJSON, []Field{
        NewLevelField(Brackets.Angle),
        NewMessageField(),
    })

    logger, _ := NewLoggerWithOptions(WithDestination(os.Stdout, formatter), WithAsync(false))

    logger.Info("This is an info message.")
    // Output: {"level":"INFO","message":"This is an info message."}
}

func ExampleWithDefaultColorization() {
    formatter, _ := NewFormatter(OutputFormatText, []Field{
        NewLevelField(Brackets.Angle),
        NewMessageField(),
    }, WithDefaultColorization())

    buf := &bytes.Buffer{}
    logger, _ := NewLoggerWithOptions(WithDestination(buf, formatter), WithAsync(false))

    logger.Warn("This is an info message.")

    // NOTE: Colorization breaks Golang's default output formatting, so you'll need to run this example in a terminal
    // that supports ANSI colors.

    fmt.Println(buf.Bytes())
    // Output: [27 91 51 51 109 60 87 65 82 78 62 32 84 104 105 115 32 105 115 32 97 110 32 105 110 102 111 32 109 101 115 115 97 103 101 46 27 91 48 109 10]
}

type invalidField struct{}

func (f invalidField) NewFieldFormatter() (FieldFormatter, error) {
    return nil, errors.New("invalid field")
}

func Test_ultraFormatter_Format(t *testing.T) {
    type args struct {
        level Level
        msg   string
    }
    tests := []struct {
        name        string
        fields      []Field
        enableColor bool
        levelColors map[Level]Color
        args        args
        want        []byte
        wantErr     bool
    }{
        {
            name: "Default",
            args: args{
                level: Info,
                msg:   "test",
            },
            want: []byte("[tag] <INFO> test"),
            fields: []Field{
                NewDefaultTagField(),
                NewLevelField(Brackets.Angle),
                &fieldMessage{},
            },
        },
        {
            name: "No Fields",
            args: args{
                level: Info,
                msg:   "test",
            },
            want: []byte(""),
        },
        {
            name: "Invalid prefix field throws error",
            args: args{
                level: Info,
                msg:   "test",
            },
            fields: []Field{
                invalidField{},
            },
            wantErr: true,
        },
        {
            name: "Colorize",
            args: args{
                level: Info,
                msg:   "test",
            },
            fields: []Field{
                &fieldMessage{},
            },
            enableColor: true,
            levelColors: map[Level]Color{
                Debug: Colors.White,
                Info:  Colors.Green,
                Warn:  Colors.Yellow,
                Error: Colors.Red,
                Panic: Colors.Magenta,
            },
            want: Colors.Green.Colorize([]byte("test")),
        },
        {
            name: "Colorize fields",
            args: args{
                level: Error,
                msg:   "test",
            },
            enableColor: true,
            levelColors: map[Level]Color{
                Debug: Colors.White,
                Info:  Colors.Green,
                Warn:  Colors.Yellow,
                Error: Colors.Red,
                Panic: Colors.Magenta,
            },
            fields: []Field{
                NewDefaultTagField(),
                NewLevelField(Brackets.Angle),
                &fieldMessage{},
                NewDefaultTagField(),
                NewLevelField(Brackets.Angle),
            },
            want: Colors.Red.Colorize([]byte("[tag] <ERROR> test [tag] <ERROR>")),
        },
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            var f LogLineFormatter
            f = &TextFormatter{
                Fields:         tt.fields,
                FieldSeparator: " ",
            }

            if tt.enableColor {
                f = NewColorizedFormatter(f, tt.levelColors)
            }

            lineArgs := LogLineArgs{
                Level: tt.args.level,
                Tag:   "tag",
            }

            if got := f.FormatLogLine(lineArgs, tt.args.msg); !bytes.Equal(got.bytes, tt.want) {
                fmt.Println("Got:  ", string(got.bytes))
                fmt.Println("Got:  ", got.bytes)
                fmt.Println("Want: ", tt.want)
                t.Errorf("Format() = %v, want %v", got, tt.want)
            }
        })
    }
}
