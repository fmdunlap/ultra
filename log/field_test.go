package log

import (
    "bytes"
    "fmt"
    "os"
    "testing"
    "time"
)

func ExampleField() {
    formatter, _ := NewFormatter(OutputFormatText, []Field{
        NewLevelField(Brackets.Angle),
        NewMessageField(),
    })

    // Note: were setting WithAsync(false) here just to ensure that the output is synchronous in the example.
    // In a real application, you *could* do this, but it will make your logging block the main thread until the log
    // has been written to the output.
    logger, _ := NewLoggerWithOptions(WithDestination(os.Stdout, formatter), WithAsync(false))

    logger.Info("This is an info message.")
    // Output: <INFO> This is an info message.
}

func ExampleField_jSON() {
    formatter, _ := NewFormatter(OutputFormatJSON, []Field{
        NewLevelField(Brackets.Angle),
        NewMessageField(),
    })

    // Note: were setting WithAsync(false) here just to ensure that the output is synchronous in the example.
    // In a real application, you *could* do this, but it will make your logging block the main thread until the log
    // has been written to the output.
    logger, _ := NewLoggerWithOptions(WithDestination(os.Stdout, formatter), WithAsync(false))

    logger.Info("This is an info message.")
    // Output: {"level":"INFO","message":"This is an info message."}
}

func ExampleNewArrayField() {
    type Person struct {
        Name string
        Age  int
    }

    stringArrayField, _ := NewArrayField[Person](
        "people",
        func(args LogLineArgs, data Person) any {
            if args.OutputFormat == OutputFormatText {
                return fmt.Sprintf("%s:%d", data.Name, data.Age)
            }
            return data
        },
    )

    formatter, _ := NewFormatter(OutputFormatText, []Field{
        NewLevelField(Brackets.Angle),
        stringArrayField,
    })

    // Note: were setting WithAsync(false) here just to ensure that the output is synchronous in the example.
    // In a real application, you *could* do this, but it will make your logging block the main thread until the log
    // has been written to the output.
    logger, _ := NewLoggerWithOptions(WithDestination(os.Stdout, formatter), WithAsync(false))

    logger.Info([]Person{Person{"John", 25}, Person{"Jane", 30}})
    // Output: <INFO> [John:25, Jane:30]
}

func ExampleNewArrayField_jSON() {
    type Person struct {
        Name string
        Age  int
    }

    stringArrayField, _ := NewArrayField[Person](
        "people",
        func(args LogLineArgs, data Person) any {
            if args.OutputFormat == OutputFormatText {
                return fmt.Sprintf("%s:%d", data.Name, data.Age)
            }
            return data
        },
    )

    formatter, _ := NewFormatter(OutputFormatJSON, []Field{
        NewLevelField(Brackets.Angle),
        stringArrayField,
    })

    // Note: were setting WithAsync(false) here just to ensure that the output is synchronous in the example.
    // In a real application, you *could* do this, but it will make your logging block the main thread until the log
    // has been written to the output.
    logger, _ := NewLoggerWithOptions(WithDestination(os.Stdout, formatter), WithAsync(false))

    logger.Info([]Person{Person{"John", 25}, Person{"Jane", 30}})
    // Output: {"level":"INFO","people":[{"Name":"John","Age":25},{"Name":"Jane","Age":30}]}
}

// ExampleNewObjectField demonstrates how to create a custom field that formats a struct into a different struct before
// logging it. This is particularly useful if you need to manipulate a struct before collecting it into logs, or if you
// want to log only a subset of fields on a struct.
func ExampleNewObjectField() {
    type Person struct {
        Name string
        Age  int
    }

    // This is the struct that will be logged. It contains a description of the person.
    // Note: when using a marshalled output formatter like JSON, output will follow standard marshalling rules. That
    // means that field tags are supported, and the logged struct must export its fields.
    type PersonLogData struct {
        Description string `json:"description"`
    }

    // This is the field that actually formats the Person struct into a PersonLogData struct.
    personDescriptionField, _ := NewObjectField[Person](
        "person", // The name of the field in the resulting log line.
        func(args LogLineArgs, data Person) any {
            description := fmt.Sprintf("%s is %d years old", data.Name, data.Age)
            if args.OutputFormat == OutputFormatText {
                return description
            }
            return PersonLogData{
                Description: description,
            }
        },
    )

    formatter, _ := NewFormatter(OutputFormatJSON, []Field{NewLevelField(Brackets.None), personDescriptionField})

    // Note: We're setting WithAsync(false) here just to ensure that the output is synchronous in the example.
    // In a real application, you *could* do this, but it will make your logging block the main thread until the log
    // has been written to the output.
    logger, _ := NewLoggerWithOptions(WithDestination(os.Stdout, formatter), WithAsync(false))

    john := Person{
        Name: "John",
        Age:  25,
    }

    logger.Info(john)
    // Output: {"level":"INFO","person":{"description":"John is 25 years old"}}
}

type mockClock struct{}

func (c mockClock) Now() time.Time {
    return time.Date(2024, time.November, 7, 19, 30, 0, 0, time.UTC)
}

func TestLevelField(t *testing.T) {
    tests := []struct {
        name       string
        levelField Field
        args       LogLineArgs
        want       string
        wantErr    bool
    }{
        {
            name:       "Default",
            levelField: NewLevelField(Brackets.Angle),
            args: LogLineArgs{
                Level:        Info,
                OutputFormat: OutputFormatText,
            },
            want: "<INFO>",
        },
        {
            name:       "Round Bracket",
            levelField: NewLevelField(Brackets.Round),
            args: LogLineArgs{
                Level:        Info,
                OutputFormat: OutputFormatText,
            },
            want: "(INFO)",
        },
        {
            name:       "Debug",
            levelField: NewLevelField(Brackets.Angle),
            args: LogLineArgs{
                Level:        Debug,
                OutputFormat: OutputFormatText,
            },
            want: "<DEBUG>",
        },
        {
            name:       "Warn",
            levelField: NewLevelField(Brackets.Angle),
            args: LogLineArgs{
                Level:        Warn,
                OutputFormat: OutputFormatText,
            },
            want: "<WARN>",
        },
        {
            name:       "Error",
            levelField: NewLevelField(Brackets.Angle),
            args: LogLineArgs{
                Level:        Error,
                OutputFormat: OutputFormatText,
            },
            want: "<ERROR>",
        },
        {
            name:       "Panic",
            levelField: NewLevelField(Brackets.Angle),
            args: LogLineArgs{
                Level:        Panic,
                OutputFormat: OutputFormatText,
            },
            want: "<PANIC>",
        },
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            formatter, err := tt.levelField.NewFieldFormatter()
            if (err != nil) != tt.wantErr {
                t.Errorf("NewFieldFormatter() error = %v, wantErr %v", err, tt.wantErr)
                return
            }

            result, err := formatter(tt.args, nil)
            if (err != nil) != tt.wantErr {
                t.Errorf("NewFieldFormatter() error = %v, wantErr %v", err, tt.wantErr)
                return
            }

            if result.Data != tt.want {
                t.Errorf("NewFieldFormatter() formatter = %v, want %v", result.Data, tt.want)
            }
        })
    }
}

func TestDateTimeField(t *testing.T) {
    tests := []struct {
        name          string
        dateTimeField *currentTimeField
        args          LogLineArgs
        want          string
        wantErr       bool
    }{
        {
            name: "Default",
            dateTimeField: &currentTimeField{
                fmtString: "2006-01-02 15:04:05",
                clock:     mockClock{},
            },
            args: LogLineArgs{
                Level:        Info,
                OutputFormat: OutputFormatText,
            },
            want: "2024-11-07 19:30:00",
        },
        {
            name: "Only Time",
            dateTimeField: &currentTimeField{
                fmtString: "15:04:05",
                clock:     mockClock{},
            },
            args: LogLineArgs{
                Level:        Info,
                OutputFormat: OutputFormatText,
            },
            want: "19:30:00",
        },
        {
            name: "Only Date",
            dateTimeField: &currentTimeField{
                fmtString: "2006-01-02",
                clock:     mockClock{},
            },
            args: LogLineArgs{
                Level:        Info,
                OutputFormat: OutputFormatText,
            },
            want: "2024-11-07",
        },
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            formatter, err := tt.dateTimeField.NewFieldFormatter()
            if (err != nil) != tt.wantErr {
                t.Errorf("NewFieldFormatter() error = %v, wantErr %v", err, tt.wantErr)
                return
            }

            result, err := formatter(tt.args, nil)
            if (err != nil) != tt.wantErr {
                t.Errorf("NewFieldFormatter() error = %v, wantErr %v", err, tt.wantErr)
                return
            }

            if result.Data != tt.want {
                t.Errorf("formatter() got = %v, want %v", result.Data, tt.want)
            }
        })
    }
}

type ComplexMapKey struct {
    Key string
    B   bool
    I   int
}

type ComplexMapValue struct {
    Val string
    B   bool
    I   int
}

func Test_QuickTest(t *testing.T) {
    type tStruct struct {
        Val string
        B   bool
    }

    stringArrayField, _ := NewArrayField[tStruct](
        "stringArray",
        func(args LogLineArgs, data tStruct) any {
            if args.OutputFormat == OutputFormatText {
                return fmt.Sprintf("%s&@&%v", data.Val, data.B)
            }
            return data
        },
    )

    stringField, _ := NewStringField("string")

    boolField, _ := NewBoolField("bool")

    currentTimeField, _ := NewCurrentTimeField("CurrentTime", "2006-01-02 15:04:05")

    responseField, _ := NewResponseField("response", ResponseFieldSettings{
        LogStatus: true,
        LogPath:   true,
    })

    mapField, _ := NewMapField[string, string]("map", func(args LogLineArgs, data string) any {
        return data
    }, func(args LogLineArgs, data string) any {
        return data
    })

    complexMapField, _ := NewMapField[ComplexMapKey, ComplexMapValue]("complexMap",
        func(args LogLineArgs, data ComplexMapKey) any {
            if args.OutputFormat != OutputFormatText {
                return fmt.Sprintf("%v:%v", data.Key, data.I)
            }
            return data
        },
        func(args LogLineArgs, data ComplexMapValue) any {
            return data
        },
    )

    testColors := map[Level]Color{
        Debug: ColorAnsiRGB(235, 216, 52),
        Info:  ColorAnsiRGB(12, 240, 228),
        Warn:  ColorAnsiRGB(237, 123, 0),
        Error: ColorAnsiRGB(237, 0, 0),
        Panic: ColorAnsiRGB(237, 0, 0),
    }

    testFormatter, _ := NewFormatter(OutputFormatText, []Field{stringArrayField, stringField, boolField, currentTimeField, mapField, responseField, complexMapField}, WithColorization(testColors))

    buf := &bytes.Buffer{}

    logger, err := NewLoggerWithOptions(WithDestination(buf, testFormatter), WithMinLevel(Debug))
    if err != nil {
        panic(err)
    }

    t.Run("Test", func(t *testing.T) {
        complexMap := map[ComplexMapKey]ComplexMapValue{
            {Key: "testAlpha", B: true, I: 10}: {Val: "ValAlpha", B: true, I: 1},
            {Key: "testBeta", B: false, I: 20}: {Val: "ValBeta", B: false, I: 2},
        }

        logger.Debug(complexMap)
        logger.Info(complexMap)
        logger.Warn(complexMap)
        logger.Error(complexMap)

        time.Sleep(time.Millisecond * 100)

        fmt.Println(buf.String())
    })
}
