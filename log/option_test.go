package log

import (
    "bytes"
    "fmt"
    "io"
    "os"
)

func ExampleWithMinLevel() {
    // Note: were setting WithAsync(false) here just to ensure that the output is synchronous in the example.
    // In a real application, you *could* do this, but it will make your logging block the main thread until the log
    // has been written to the output.
    logger, _ := NewLoggerWithOptions(WithFields(os.Stdout, []Field{NewDefaultLevelField(), NewMessageField()}), WithMinLevel(Warn), WithAsync(false))

    logger.Info("This is an info message.")
    logger.Debug("This is a debug message.")
    logger.Warn("This is a warning message.")
    logger.Error("This is an error message.")
    logger.Panic("This is a panic message.")
    // Output:
    // <WARN> This is a warning message.
    // <ERROR> This is an error message.
    // <PANIC> This is a panic message.
}

func ExampleWithFields() {
    // Note: were setting WithAsync(false) here just to ensure that the output is synchronous in the example.
    // In a real application, you *could* do this, but it will make your logging block the main thread until the log
    // has been written to the output.
    logger, _ := NewLoggerWithOptions(
        WithFields(os.Stdout, []Field{NewDefaultLevelField(), NewMessageField()}),
        WithAsync(false),
        WithAsync(false),
    )

    logger.Debug("This is a debug message.") // Below min level, won't be output.
    logger.Info("This is an info message.")
    logger.Warn("This is a warning message.")
    logger.Error("This is an error message.")
    logger.Panic("This is a panic message.")
    // Output:
    // <INFO> This is an info message.
    // <WARN> This is a warning message.
    // <ERROR> This is an error message.
    // <PANIC> This is a panic message.
}

// ExampleWithDestination shows how to use WithDestination to log to multiple writers with different formatters
func ExampleWithDestination() {
    bufOne := &bytes.Buffer{}
    bufTwo := &bytes.Buffer{}

    formatterOne, err := NewFormatter(OutputFormatText, []Field{NewDefaultLevelField(), NewMessageField()})
    if err != nil {
        panic(err)
    }

    formatterTwo, err := NewFormatter(OutputFormatJSON, []Field{NewDefaultTagField(), NewMessageField()})
    if err != nil {
        panic(err)
    }

    // Note: were setting WithAsync(false) here just to ensure that the output is synchronous in the example.
    // In a real application, you *could* do this, but it will make your logging block the main thread until the log
    // has been written to the output.
    logger, _ := NewLoggerWithOptions(
        WithDestination(bufOne, formatterOne),
        WithDestination(bufTwo, formatterTwo),
        WithTag("TAG"),
        WithAsync(false),
    )

    logger.Info("This is an info message.")

    fmt.Print(bufOne.String())
    fmt.Print(bufTwo.String())
    // Output:
    // <INFO> This is an info message.
    // {"message":"This is an info message.","tag":"TAG"}
}

// ExampleWithDestination_sharedFormatter shows how to use WithDestination to log to multiple writers using a single
// formatter. Note that changes to teh formatter will be reflected in all destinations.
func ExampleWithDestination_sharedFormatter() {
    bufOne := bytes.NewBufferString("")
    bufTwo := bytes.NewBufferString("")

    formatter, err := NewFormatter(OutputFormatText, []Field{
        NewDefaultTagField(),
        NewDefaultLevelField(),
        NewMessageField(),
    })
    if err != nil {
        panic(err)
    }

    // Note: were setting WithAsync(false) here just to ensure that the output is synchronous in the example.
    // In a real application, you *could* do this, but it will make your logging block the main thread until the log
    // has been written to the output.
    logger, _ := NewLoggerWithOptions(
        WithDestination(bufOne, formatter),
        WithDestination(bufTwo, formatter),
        WithTag("TAG"),
        WithAsync(false),
    )

    logger.Info("This is an info message.")

    fmt.Print(bufOne.String())
    fmt.Print(bufTwo.String())
    // Output:
    // [TAG] <INFO> This is an info message.
    // [TAG] <INFO> This is an info message.
}

// ExampleWithDestinations shows how to use WithDestinations to log to multiple writers with different formatters
func ExampleWithDestinations() {
    bufOne := &bytes.Buffer{}
    bufTwo := &bytes.Buffer{}

    formatterOne, err := NewFormatter(OutputFormatText, []Field{NewDefaultLevelField(), NewMessageField()})
    if err != nil {
        panic(err)
    }

    formatterTwo, err := NewFormatter(OutputFormatJSON, []Field{NewDefaultTagField(), NewMessageField()})
    if err != nil {
        panic(err)
    }

    destinations := map[io.Writer]LogLineFormatter{
        bufOne: formatterOne,
        bufTwo: formatterTwo,
    }

    // Note: were setting WithAsync(false) here just to ensure that the output is synchronous in the example.
    // In a real application, you *could* do this, but it will make your logging block the main thread until the log
    // has been written to the output.
    logger, _ := NewLoggerWithOptions(
        WithDestinations(destinations),
        WithTag("TAG"),
        WithAsync(false),
    )

    logger.Info("This is an info message.")

    fmt.Print(bufOne.String())

    fmt.Print(bufTwo.String())
    // Output:
    // <INFO> This is an info message.
    // {"message":"This is an info message.","tag":"TAG"}
}

func ExampleWithSilent() {
    buf := &bytes.Buffer{}
    // Note: were setting WithAsync(false) here just to ensure that the output is synchronous in the example.
    // In a real application, you *could* do this, but it will make your logging block the main thread until the log
    // has been written to the output.
    logger, _ := NewLoggerWithOptions(
        WithFields(buf, []Field{NewDefaultLevelField(), NewMessageField()}),
        WithSilent(true),
        WithAsync(false),
    )

    logger.Info("This is an info message.")

    fmt.Print(buf.String())
    // Output:
    //
}

// ExampleWithDefaultColorizationEnabled shows how to use WithDefaultColorizationEnabled to enable colorization for
// the default formatter.
func ExampleWithDefaultColorizationEnabled() {
    buf := &bytes.Buffer{}
    // Note: were setting WithAsync(false) here just to ensure that the output is synchronous in the example.
    // In a real application, you *could* do this, but it will make your logging block the main thread until the log
    // has been written to the output.
    logger, _ := NewLoggerWithOptions(
        WithFields(buf, []Field{NewDefaultLevelField(), NewMessageField()}),
        WithDefaultColorizationEnabled(buf),
        WithAsync(false),
    )

    logger.Info("This is an info message.")

    fmt.Print(buf.Bytes())
    // Output:
    // [27 91 51 55 109 60 73 78 70 79 62 32 84 104 105 115 32 105 115 32 97 110 32 105 110 102 111 32 109 101 115 115 97 103 101 46 27 91 48 109 10]
}

func ExampleWithCustomColorization() {
    testColors := map[Level]Color{
        Debug: ColorAnsiRGB(235, 216, 52),
        Info:  ColorAnsiRGB(12, 240, 228),
        Warn:  ColorAnsiRGB(237, 123, 0),
        Error: ColorAnsiRGB(237, 0, 0),
        Panic: ColorAnsiRGB(237, 0, 0),
    }

    formatter, _ := NewFormatter(OutputFormatText, []Field{NewDefaultLevelField(), NewMessageField()}, WithColorization(testColors))

    buf := &bytes.Buffer{}

    // Note: were setting WithAsync(false) here just to ensure that the output is synchronous in the example.
    // In a real application, you *could* do this, but it will make your logging block the main thread until the log
    // has been written to the output.
    logger, err := NewLoggerWithOptions(WithDestination(buf, formatter), WithMinLevel(Debug), WithAsync(false))
    if err != nil {
        panic(err)
    }

    logger.Debug("Debug")
    logger.Info("Info")
    logger.Warn("Warn")
    logger.Error("Error")
    // NOTE: Colorization breaks Golang's default output formatting, so you'll need to use the following to see the
    // colors:
    //
    // go test -v -run TestFormatter_customColors
    //
    // <DEBUG> Debug // Yellowish
    // <INFO> Info // Cyanish
    // <WARN> Warn // Orangish
    // <ERROR> Error // Reddish

    fmt.Print(buf.Bytes())
    // Output:
    // [27 91 51 56 59 50 59 50 51 53 59 50 49 54 59 53 50 109 60 68 69 66 85 71 62 32 68 101 98 117 103 27 91 48 109 10 27 91 51 56 59 50 59 49 50 59 50 52 48 59 50 50 56 109 60 73 78 70 79 62 32 73 110 102 111 27 91 48 109 10 27 91 51 56 59 50 59 50 51 55 59 49 50 51 59 48 109 60 87 65 82 78 62 32 87 97 114 110 27 91 48 109 10 27 91 51 56 59 50 59 50 51 55 59 48 59 48 109 60 69 82 82 79 82 62 32 69 114 114 111 114 27 91 48 109 10]
}

// ExampleWithTag shows how to use WithTag to set the tag for the logger.
func ExampleWithTag() {
    buf := &bytes.Buffer{}
    logger, _ := NewLoggerWithOptions(
        WithFields(buf, []Field{NewDefaultTagField(), NewDefaultLevelField(), NewMessageField()}),
        WithTag("TAG"),
        WithAsync(false),
    )

    logger.Info("This is an info message.")

    fmt.Print(buf.String())
    // Output:
    // [TAG] <INFO> This is an info message.
}
