# UltraLogger

[![Logo](/git/logo.png)](https://github.com/fmdunlap/ultralogger)

[![Go Reference](https://pkg.go.dev/badge/github.com/fmdunlap/go-ultralogger.svg)](https://pkg.go.dev/github.com/fmdunlap/go-ultralogger)

## Overview
UltraLogger is a versatile, flexible, and efficient logging library for Go that supports multiple log levels, custom
formatting, and output destinations. It is designed to be easy to use while providing advanced features such as
colorized output, customizable log levels, and detailed error handling.

UltraLogger only uses the stdlib and is written in pure Go.

## Features
- **Multiple Log Levels**: DEBUG, INFO, WARN, ERROR, PANIC, and FATAL levels.
- **Concurrent-Safe Logging**: Supports logging from multiple goroutines, and provides asynchronous logging when writing
  to file or network writers.
- **Flexible Formatting**: Multiple build-in field types, extensible interface for custom field types.
- **Custom Formatting**: Flexible formatting options including custom date/time formats, tag padding, and bracket types.
- **Colorization**: Configurable colorization can change the color of terminal output based on severity level for better 
  visibility.
- **Output Redirection**: Supports writing logs to various io.Writer destinations such as files or standard output using
  a single logger.
- **Silent Mode**: Allows disabling all logging when set.
- **Error Handling**: Implements robust error handling and fallback mechanisms for logging errors.

## Installation
To install UltraLogger, use `go get`:
```sh
go get github.com/fmdunlap/ultralogger
```

## Usage
Here's a basic example of how to use UltraLogger:

```go
package main

import (
    "os"
    "github.com/fmdunlap/ultralogger"
)

func main() {
    logger := ultralogger.NewUltraLogger(os.Stdout)
    
    logger.Info("This is an info message.")  // Output: 2006-01-02 15:04:05 <INFO> This is an info message.
    logger.Debug("This is a debug message.") // Output: 2006-01-02 15:04:05 <DEBUG> This is a debug message.
    logger.Warn("This is a warning message.") // Output: 2006-01-02 15:04:05 <WARN> This is a warning message.
    
    logger.Infof("This is an info message with %s!", "formatting") // Output: 2006-01-02 15:04:05 <INFO> This is an info message with formatting!
}
```

## Configuration

Ultralogger provides various configuration options to customize logging behaviour.

### Minimum Log Level

```go
logger := ultralogger.NewUltraLogger(os.Stdout).MinLogLevel(ultralogger.LogLevelDebug)
```

### Tags

```go
logger.SetTag("MyTag")
logger.Info("Message")   // -> 2006-01-02 15:04:05 [MyTag] <INFO> Message

logger.SetTag("Another")
logger.Infof("Ultralogger is %s!", "super cool") // -> 2006-01-02 15:04:05 [Another] <INFO> Ultralogger is super cool!
```

### Formatting

#### Date and Time

```go
logger.SetDateFormat("01/02/2006")
logger.ShowTime(false)

logger.Info("Message") // -> 01/02/2006 <INFO> Message

_, err := logger.ShowTime(true).SetTimeFormat("15|04|05")
if err != nil {
    // Handle error if the time format is invalid
}

logger.Info("Message") // -> 01/02/2006 15|04|05 <INFO> Message

logger.SetDateTimeSeparator("@")
logger.Info("Message") // -> 01/02/2006@15|04|05 <INFO> Message
```

#### Tag Padding
```go
// Many style config funcs can be chained together. This is just an example.
logger.SetTag("MyTag").EnabledTagPadding(true).SetTagPadSize(10)

logger.Info("Message") // -> 2006-01-02 15:04:05 [MyTag]   <INFO> Message
logger.Error("Error!") // -> 2006-01-02 15:04:05 [MyTag]   <ERROR> Warning!
```

#### Tag Bracket Type
```go
logger.SetTag("MyTag").SetTagBracketType(ultralogger.BracketTypeSquare)
logger.Info("Message") // -> 2006-01-02 15:04:05 [MyTag] <INFO> Message
```

#### Log Bracket Type
```go
logger.SetLogBracketType(ultralogger.BracketTypeRound)
logger.Info("Message") // -> 2006-01-02 15:04:05 (INFO) Message
```

#### Bracket Types

```go
ultralogger.BracketTypeNone // "tag"
ultralogger.BracketTypeSquare // "[tag]"
ultralogger.BracketTypeRound // "(tag)"
ultralogger.BracketTypeCurly // "{tag}"
ultralogger.BracketTypeAngle // "<tag>"
```

### Terminal Colorization

```go
logger := ultralogger.NewStdoutLogger()
logger.SetMinLevel(ultralogger.DebugLevel)
logger.SetColorize(true)

logger.Debug("Debug")
logger.Info("Message")
logger.Warn("Warning")
logger.Error("Error!")
logger.Panic("Panic!")
```

![Colorized Output](/git/colorized.png)

You can also set custom colors for each log level!

```go
logger := ultralogger.NewStdoutLogger()
logger.SetMinLevel(ultralogger.DebugLevel)
logger.SetColorize(true)

// Set custom colors for each log level
logger.SetLevelColor(ultralogger.DebugLevel, ultralogger.ColorCyan)
logger.Debug("Cyan is peaceful. Like the messages ad the debug level!")

// Set custom colors for all log levels
logger.SetLevelColors(map[ultralogger.Level]ultralogger.Color{
    ultralogger.WarnLevel: ultralogger.ColorGreen,
    ultralogger.ErrorLevel: ultralogger.ColorBlue,
    ultralogger.PanicLevel: ultralogger.ColorRed,
})

logger.Warn("Green might not be the best warning color...")
logger.Error("I'll be blue if I see this error!")
logger.Panic("RED IS THE BEST PANIC COLOR!")
```

![Custom Colors](/git/custom_colors.png)

### Output Destinations

```go
// Files
logger := ultralogger.NewFileLogger("somefile.log")

// Stdout
logger := ultralogger.NewStdoutLogger()

// ByteBuffer
buf := new(bytes.Buffer)
logger := ultralogger.NewUltraLogger(buf)

logger.Info("Debug to byte buffer")

fmt.Printf("Buffer: %s\n", buf.String()) // Buffer: 2006-01-02 15:04:05 <INFO> Debug to byte buffer
```

### Silent Mode

```go
logger := ultralogger.NewStdoutLogger()
logger.SetMinLevel(ultralogger.WarnLevel)

logger.Info("Message") // -> 2006-01-02 15:04:05 <INFO> Message
logger.Warn("Message") // -> Nothing will be printed to the output
```

### PanicOnPanicLevel Mode

```go
logger := ultralogger.NewStdoutLogger()
logger.SetMinLevel(ultralogger.WarnLevel)
logger.SetPanicOnPanicLevel(true)

logger.Panic("Panic!") // -> 2006-01-02 15:04:05 <PANIC> Panic! (and then panics)

logger.SetPanicOnPanicLevel(false)
logger.Panic("Panic!") // -> 2006-01-02 15:04:05 <PANIC> Panic! (and then does not panic)
```

## TODO

- [ ] Improve Color handling. Colors are currently using ASNI color cods, which isn' the most modern way to do this.
- [ ] Optimize logging speed.
- [ ] Make logging fields configurable (ability to add fields, change content, change field color, change field order,etc.)
- [ ] Add the ability to pass option funcs to the constructors to allow for more flexibility.
- [ ] Custom bracket types!
- [ ] OnLogLevel handlers to allow for more advanced logging.

## License

MIT. See LICENSE for more details.

But also, it's MIT. Go nuts.

## Contributing

Feel free to open an issue or PR if you have any suggestions or improvements!