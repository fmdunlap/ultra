package log

import (
    "bytes"
    "errors"
    "fmt"
    "io"
    "net/http"
    "net/url"
    "os"
    "testing"
    "time"
)

func ExampleField() {

    formatter, _ := NewFormatter(OutputFormatText, []Field{
        NewLevelField(nil),
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
        NewDefaultLevelField(),
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
        func(args LogLineArgs, data Person) (any, error) {
            if args.OutputFormat == OutputFormatText {
                return fmt.Sprintf("%s:%d", data.Name, data.Age), nil
            }
            return data, nil
        },
    )

    formatter, _ := NewFormatter(OutputFormatText, []Field{
        NewDefaultLevelField(),
        NewMessageField(),
        stringArrayField,
    })

    // Note: were setting WithAsync(false) here just to ensure that the output is synchronous in the example.
    // In a real application, you *could* do this, but it will make your logging block the main thread until the log
    // has been written to the output.
    logger, _ := NewLoggerWithOptions(WithDestination(os.Stdout, formatter), WithAsync(false))

    logger.Info("People", []Person{{"John", 25}, {"Jane", 30}})
    // Output: <INFO> People people=[John:25, Jane:30]
}

func ExampleNewArrayField_jSON() {
    type Person struct {
        Name string
        Age  int
    }

    stringArrayField, _ := NewArrayField[Person](
        "people",
        func(args LogLineArgs, data Person) (any, error) {
            if args.OutputFormat == OutputFormatText {
                return fmt.Sprintf("%s:%d", data.Name, data.Age), nil
            }
            return data, nil
        },
    )

    formatter, _ := NewFormatter(OutputFormatJSON, []Field{
        NewDefaultLevelField(),
        stringArrayField,
        NewMessageField(),
    })

    // Note: were setting WithAsync(false) here just to ensure that the output is synchronous in the example.
    // In a real application, you *could* do this, but it will make your logging block the main thread until the log
    // has been written to the output.
    logger, _ := NewLoggerWithOptions(WithDestination(os.Stdout, formatter), WithAsync(false))

    people := []Person{{"John", 25}, {"Jane", 30}}

    logger.Info("Did a thing", people)
    // Output: {"level":"INFO","message":"Did a thing","people":[{"Name":"John","Age":25},{"Name":"Jane","Age":30}]}
}

// ExampleNewObjectField demonstrates how to create a custom field that formats a struct into a different struct before
// logging it. This is particularly useful if you need to manipulate a struct before collecting it into logs, or if you
// want to log only a subset of fields on a struct.
func ExampleNewObjectField() {
    type BigStruct struct {
        Field1 string
        Field2 string
        Field3 string
    }

    type User struct {
        ID          string
        Name        string
        Age         int
        IsAdmin     bool
        LargeStruct BigStruct
    }

    // This is the struct that will be logged. It contains a description of the person.
    // Note: when using a marshalled output formatter like JSON, output will follow standard marshalling rules. That
    // means that field tags are supported, and the logged struct must export its fields.
    type UserLogData struct {
        ID      string `json:"ID"`
        Name    string `json:"Name"`
        Age     int    `json:"Age"`
        IsAdmin bool   `json:"IsAdmin"`
    }

    userDescription := func(u User) string {
        description := fmt.Sprintf("ID: %s, Name: %s, Age: %d", u.ID, u.Name, u.Age)
        if u.IsAdmin {
            description += "[ADMIN]"
        }
        return description
    }

    // This field formats the User struct into a UserLogData struct.
    personDescriptionField, _ := NewObjectField[User](
        "user", // The name of the field in the resulting log line.
        func(args LogLineArgs, data User) (any, error) {
            if args.OutputFormat == OutputFormatText {
                return fmt.Sprintf("'%s'", userDescription(data)), nil
            }

            return UserLogData{
                ID:      data.ID,
                Name:    data.Name,
                Age:     data.Age,
                IsAdmin: data.IsAdmin,
            }, nil
        },
    )

    jsonFormatter, _ := NewFormatter(OutputFormatJSON, []Field{NewDefaultLevelField(), personDescriptionField, NewMessageField()})
    textFormatter, _ := NewFormatter(OutputFormatText, []Field{NewDefaultLevelField(), personDescriptionField, NewMessageField()})

    // Note: We're setting WithAsync(false) here just to ensure that the output is synchronous in the example.
    // In a real application, you *could* do this, but it will make your logging block the main thread until the log
    // has been written to the output.
    textBuffer := &bytes.Buffer{}
    jsonBuffer := &bytes.Buffer{}
    logger, _ := NewLoggerWithOptions(
        WithDestination(jsonBuffer, jsonFormatter),
        WithDestination(textBuffer, textFormatter),
        WithAsync(false),
    )

    john := User{
        Name: "John",
        Age:  25,
        LargeStruct: BigStruct{
            Field1: "Some value",
            Field2: "Another value",
            Field3: "A third value",
        },
    }

    logger.Info("message about john", john)

    fmt.Print(jsonBuffer.String())
    fmt.Print(textBuffer.String())
    // Output:
    // {"level":"INFO","message":"message about john","user":{"ID":"","Name":"John","Age":25,"IsAdmin":false}}
    // <INFO> user='ID: , Name: John, Age: 25' message about john
}

func TestObjectField(t *testing.T) {
    type newObjectFieldArgs struct {
        name      string
        formatter ObjectFieldFormatter[string]
        options   []FieldOption
    }

    type formatterArgs struct {
        args LogLineArgs
        data any
    }

    tests := []struct {
        name               string
        newObjectFieldArgs newObjectFieldArgs
        formatterArgs      formatterArgs
        want               string
        wantFieldInitErr   bool
        wantFormatErr      bool
        wantHideKey        bool
        wantAlwaysMatch    bool
    }{
        {
            name: "Default",
            newObjectFieldArgs: newObjectFieldArgs{
                name: "test",
                formatter: func(args LogLineArgs, data string) (any, error) {
                    return data, nil
                },
            },
            formatterArgs: formatterArgs{
                args: LogLineArgs{
                    Level: Info,
                },
                data: "test",
            },
            want: "test",
        },
        {
            name: "Formatter intercepts data",
            newObjectFieldArgs: newObjectFieldArgs{
                name: "test",
                formatter: func(args LogLineArgs, data string) (any, error) {
                    return data + "!", nil
                },
            },
            formatterArgs: formatterArgs{
                args: LogLineArgs{
                    Level: Info,
                },
                data: "test",
            },
            want: "test!",
        },
        {
            name: "HideKey Is Set",
            newObjectFieldArgs: newObjectFieldArgs{
                name: "test",
                formatter: func(args LogLineArgs, data string) (any, error) {
                    return data, nil
                },
                options: []FieldOption{
                    WithHideKey(true),
                },
            },
            wantHideKey: true,
            formatterArgs: formatterArgs{
                args: LogLineArgs{
                    Level: Info,
                },
                data: "test",
            },
            want: "test",
        },
        {
            name: "AlwaysMatch Is Set",
            newObjectFieldArgs: newObjectFieldArgs{
                name: "test",
                formatter: func(args LogLineArgs, data string) (any, error) {
                    return data, nil
                },
                options: []FieldOption{
                    WithAlwaysMatch(true),
                },
            },
            formatterArgs: formatterArgs{
                args: LogLineArgs{
                    Level: Info,
                },
                data: "test",
            },
            wantAlwaysMatch: true,
            want:            "test",
        },
        {
            name: "Error On Field Init - Empty Name",
            newObjectFieldArgs: newObjectFieldArgs{
                name: "",
                formatter: func(args LogLineArgs, data string) (any, error) {
                    return data, nil
                },
            },
            formatterArgs: formatterArgs{
                args: LogLineArgs{
                    Level: Info,
                },
                data: "test",
            },
            wantFieldInitErr: true,
        },
        {
            name: "Error On Field Init - Nil Formatter",
            newObjectFieldArgs: newObjectFieldArgs{
                name:      "test",
                formatter: nil,
            },
            formatterArgs: formatterArgs{
                args: LogLineArgs{
                    Level: Info,
                },
                data: "test",
            },
            wantFieldInitErr: true,
        },
        {
            name: "Error On Field Init - Failed Option",
            newObjectFieldArgs: newObjectFieldArgs{
                name: "test",
                formatter: func(args LogLineArgs, data string) (any, error) {
                    return data, nil
                },
                options: []FieldOption{
                    func(f *FieldSettings) error {
                        return errors.New("test")
                    },
                },
            },
            formatterArgs: formatterArgs{
                args: LogLineArgs{
                    Level: Info,
                },
                data: "test",
            },
            wantFieldInitErr: true,
        },
        {
            name: "Error On Format - Unmatched Data Type",
            newObjectFieldArgs: newObjectFieldArgs{
                name: "test",
                formatter: func(args LogLineArgs, data string) (any, error) {
                    return data, nil
                },
            },
            formatterArgs: formatterArgs{
                args: LogLineArgs{
                    Level: Info,
                },
                data: struct{}{},
            },
            wantFormatErr: true,
        },
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            objectField, err := NewObjectField[string](
                tt.newObjectFieldArgs.name,
                tt.newObjectFieldArgs.formatter,
                tt.newObjectFieldArgs.options...,
            )

            if err != nil {
                if tt.wantFieldInitErr {
                    return
                }
                t.Errorf("NewObjectField() error = %v, wantErr %v", err, tt.wantFieldInitErr)
                return
            }

            if objectField.Name() != tt.newObjectFieldArgs.name {
                t.Errorf("NewObjectField() name = %v, want %v", objectField.Name(), tt.newObjectFieldArgs.name)
            }

            if objectField.Settings().HideKey != tt.wantHideKey {
                t.Errorf("NewObjectField() HideKey = %v, want %v", objectField.Settings().HideKey, tt.wantHideKey)
            }

            if objectField.Settings().AlwaysMatch != tt.wantAlwaysMatch {
                t.Errorf("NewObjectField() AlwaysMatch = %v, want %v", objectField.Settings().AlwaysMatch, tt.wantAlwaysMatch)
            }

            formatter, _ := objectField.NewFieldFormatter()

            result, err := formatter(tt.formatterArgs.args, tt.formatterArgs.data)
            if err != nil {
                if tt.wantFormatErr {
                    return
                }
                t.Errorf("NewFieldFormatter() error = %v, wantErr %v", err, tt.wantFormatErr)
                return
            }

            if result != tt.want {
                t.Errorf("NewFieldFormatter() formatter = %v, want %v", result, tt.want)
            }
        })
    }
}

func TestLevelField(t *testing.T) {
    tests := []struct {
        name               string
        levelFieldSettings *LevelFieldSettings
        args               LogLineArgs
        want               string
        wantErr            bool
    }{
        {
            name: "Default",
            args: LogLineArgs{
                Level:        Info,
                OutputFormat: OutputFormatText,
            },
            want: "<INFO>",
        },
        {
            name: "Round Bracket",
            levelFieldSettings: &LevelFieldSettings{
                Bracket: Brackets.Round,
            },
            args: LogLineArgs{
                Level:        Info,
                OutputFormat: OutputFormatText,
            },
            want: "(INFO)",
        },
        {
            name: "Default - Debug",
            args: LogLineArgs{
                Level:        Debug,
                OutputFormat: OutputFormatText,
            },
            want: "<DEBUG>",
        },
        {
            name: "Default - Warn",
            args: LogLineArgs{
                Level:        Warn,
                OutputFormat: OutputFormatText,
            },
            want: "<WARN>",
        },
        {
            name: "Default - Error",
            args: LogLineArgs{
                Level:        Error,
                OutputFormat: OutputFormatText,
            },
            want: "<ERROR>",
        },
        {
            name: "Panic",
            args: LogLineArgs{
                Level:        Panic,
                OutputFormat: OutputFormatText,
            },
            want: "<PANIC>",
        },
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            levelField := NewLevelField(tt.levelFieldSettings)
            formatter, err := levelField.NewFieldFormatter()
            if (err != nil) != tt.wantErr {
                t.Errorf("NewFieldFormatter() error = %v, wantErr %v", err, tt.wantErr)
                return
            }

            result, err := formatter(tt.args, struct{}{})
            if (err != nil) != tt.wantErr {
                t.Errorf("NewFieldFormatter() error = %v, wantErr %v", err, tt.wantErr)
                return
            }

            if result != tt.want {
                t.Errorf("NewFieldFormatter() formatter = %v, want %v", result, tt.want)
            }
        })
    }
}

func TestDateTimeField(t *testing.T) {
    tests := []struct {
        name                     string
        currentTimeFieldSettings *CurrentTimeFieldSettings
        args                     LogLineArgs
        want                     string
        wantErr                  bool
    }{
        {
            name: "Default",
            currentTimeFieldSettings: &CurrentTimeFieldSettings{
                Name:   "currentTime",
                Format: "2006-01-02 15:04:05",
            },
            args: LogLineArgs{
                Level:        Info,
                OutputFormat: OutputFormatText,
            },
            want: "2024-11-07 19:30:00",
        },
        {
            name: "Only Time",
            currentTimeFieldSettings: &CurrentTimeFieldSettings{
                Name:   "currentTime",
                Format: "15:04:05",
            },
            args: LogLineArgs{
                Level:        Info,
                OutputFormat: OutputFormatText,
            },
            want: "19:30:00",
        },
        {
            name: "Only Date",
            currentTimeFieldSettings: &CurrentTimeFieldSettings{
                Name:   "currentTime",
                Format: "2006-01-02",
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
            fakeNow := time.Date(2024, time.November, 7, 19, 30, 0, 0, time.UTC)
            tt.currentTimeFieldSettings.fakeNow = &fakeNow
            currentTimeField := NewCurrentTimeField(tt.currentTimeFieldSettings)

            formatter, err := currentTimeField.NewFieldFormatter()
            if (err != nil) != tt.wantErr {
                t.Errorf("NewFieldFormatter() error = %v, wantErr %v", err, tt.wantErr)
                return
            }

            result, err := formatter(tt.args, struct{}{})
            if (err != nil) != tt.wantErr {
                t.Errorf("NewFieldFormatter() error = %v, wantErr %v", err, tt.wantErr)
                return
            }

            if result != tt.want {
                t.Errorf("formatter() got = %v, want %v", result, tt.want)
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

    tStructArrayField, _ := NewArrayField[tStruct](
        "tStructArray",
        func(args LogLineArgs, data tStruct) (any, error) {
            if args.OutputFormat == OutputFormatText {
                return fmt.Sprintf("%s&@&%v", data.Val, data.B), nil
            }
            return data, nil
        },
    )

    stringArrayField, _ := NewArrayField[string](
        "stringArray",
        func(args LogLineArgs, data string) (any, error) {
            return data, nil
        },
    )

    stringField, _ := NewStringField("string")

    boolField, _ := NewBoolField("bool")

    currentTimeField := NewCurrentTimeField(&CurrentTimeFieldSettings{
        Name:   "CurrentTime",
        Format: "2006-01-02 15:04:05",
    })

    responseField, _ := NewResponseField(&ResponseFieldSettings{
        LogStatus: true,
        LogPath:   true,
    })

    mapField, _ := NewMapField[string, string]("map", func(args LogLineArgs, data string) (any, error) {
        return data, nil
    }, func(args LogLineArgs, data string) (any, error) {
        return data, nil
    })

    complexMapField, _ := NewMapField[ComplexMapKey, ComplexMapValue]("complexMap",
        func(args LogLineArgs, data ComplexMapKey) (any, error) {
            if args.OutputFormat != OutputFormatText {
                return fmt.Sprintf("%v:%v", data.Key, data.I), nil
            }
            return data, nil
        },
        func(args LogLineArgs, data ComplexMapValue) (any, error) {
            return data, nil
        },
    )

    testColors := map[Level]Color{
        Debug: ColorAnsiRGB(235, 216, 52),
        Info:  ColorAnsiRGB(12, 240, 228),
        Warn:  ColorAnsiRGB(237, 123, 0),
        Error: ColorAnsiRGB(237, 0, 0),
        Panic: ColorAnsiRGB(237, 0, 0),
    }

    testFormatter, _ := NewFormatter(OutputFormatText, []Field{tStructArrayField, stringArrayField, stringField, boolField, currentTimeField, mapField, responseField, complexMapField}, WithColorization(testColors))

    buf := &bytes.Buffer{}

    logger, err := NewLoggerWithOptions(WithDestination(io.Discard, testFormatter), WithMinLevel(Debug), WithAsync(false))
    if err != nil {
        panic(err)
    }

    t.Run("Test", func(t *testing.T) {
        complexMap := map[ComplexMapKey]ComplexMapValue{
            {Key: "testAlpha", B: true, I: 10}: {Val: "ValAlpha", B: true, I: 1},
            {Key: "testBeta", B: false, I: 20}: {Val: "ValBeta", B: false, I: 2},
        }

        stringMap := map[string]string{
            "test":  "test",
            "hello": "world",
        }

        stringArray := []string{
            "test",
            "hello",
            "world",
        }

        tStructArray := []tStruct{
            {Val: "golang", B: true},
            {Val: "is", B: false},
            {Val: "awesome", B: true},
        }

        data := []any{
            complexMap,
            tStructArray,
            "hello",
            true,
            stringMap,
            stringArray,
            &http.Response{
                Status:     "OK",
                StatusCode: 200,
                Request: &http.Request{
                    URL: &url.URL{
                        Path: "/test",
                    },
                },
            },
        }

        now := time.Now()
        logger.Debug(data...)
        fmt.Println(time.Since(now))

        now = time.Now()
        logger.Info(data...)
        fmt.Println(time.Since(now))

        now = time.Now()
        logger.Warn(data...)
        fmt.Println(time.Since(now))

        now = time.Now()
        logger.Error(data...)
        fmt.Println(time.Since(now))

        fmt.Println(buf.String())
    })
}
