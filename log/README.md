# ultra/log

[![Logo](/repo/logo.png)](https://github.com/fmdunlap/ultra/tree/main/log)

[![Go Reference](https://pkg.go.dev/badge/github.com/fmdunlap/go-ultralogger.svg)](https://pkg.go.dev/github.com/fmdunlap/go-ultralogger)

## Overview
UltraLogger (`ultra/log`) is a high-performance, structured logging package for Go with support for concurrent writers,
multiple output formats, and extensive field customization. It is designed to be easy to use while providing advanced
features such as colorized output, customizable log levels, and detailed error handling.

`ultra/log` only uses the stdlib and is written in pure Go.

## Features

- 🚀 Concurrent log writing with configurable timeouts
- 📝 Multiple output formats (JSON, Text)
- 🎨 Terminal color support
- 🔄 Multiple writer destinations
- 🏷️ Rich field system for structured logging
- ⚡ Async logging by default
- 🎯 Log level filtering
- 🔧 Highly configurable options

## Installation

```sh
go get github.com/fmdunlap/ultra/log
```

## Quick Start

```go
import "github.com/fmdunlap/ultra/log"

// Create a basic logger
logger := log.NewLogger()

// Log some messages
logger.Info("Hello, World!") // Output: 2023-01-01 12:00:00 <INFO> Hello, World!
logger.Debug("Debug message") // Output: 2023-01-01 12:00:00 <DEBUG> Debug message
logger.Error("Something went wrong") // Output: 2023-01-01 12:00:00 <ERROR> Something went wrong

// Log with custom fields
logger, err := log.NewLoggerWithOptions(
    log.WithFields(os.Stdout, []log.Field{
        log.NewLevelField(log.Brackets.Angle),
        log.NewMessageField(),
        log.NewBoolField("isAdmin"),
    }),
)

logger.Info(true, "Hello, World!") // Output: <INFO> isAdmin=true Hello, World!
logger.Debug(false, "Debug message") // Output: <DEBUG> isAdmin=false Debug message
```

## More Usage Examples

### Custom Formatting

```go
package main

import "github.com/fmdunlap/ultra/log"

type User struct {
    Name string
    Admin bool
    SomeLargeStruct BigStruct
}

type UserLogEntry struct {
    Name string
    Admin bool
}

func main() {
    formatter, _ := log.NewFormatter(log.OutputFormatText, []log.Field{
        log.NewLevelField(log.Brackets.Angle),
        log.NewMessageField(),
        log.NewObjectField[User]("user", func(args log.LogLineArgs, data User) any {
            if args.OutputFormat == log.OutputFormatText {
                return fmt.Sprintf("'Name: %s, Admin: %v'", data.Name, data.Admin)
            }
            return UserLogEntry{
                Name: data.Name,
                Admin: data.Admin,
            }
        }),
    })

    logger, _ := log.NewLoggerWithOptions(
        log.WithDestination(os.Stdout, formatter),
    )

    user := User{
        Name: "John Doe",
        Admin: true,
        SomeLargeStruct: BigStruct{
            Field1: "Some value",
            Field2: "Another value",
            Field3: "A third value",
        },
    }

    logger.Info("Some message", user) // Output: <INFO> Some message user='Name: John Doe, Admin: true'    
}
```

### JSON Logging

```go
package main

import "github.com/fmdunlap/ultra/log"

formatter, _ := log.NewFormatter(log.OutputFormatJSON, []log.Field{
    log.NewLevelField(log.Brackets.Angle),
    log.NewMessageField(),
    log.NewObjectField[User]("user", func(args log.LogLineArgs, data User) any {
        if args.OutputFormat == log.OutputFormatText {
            return fmt.Sprintf("'Name: %s, Admin: %v'", data.Name, data.Admin)
        }
        return UserLogEntry{
            Name: data.Name,
            Admin: data.Admin,
        }
    }),
})

logger, _ := log.NewLoggerWithOptions(
    log.WithDestination(os.Stdout, formatter),
)

user := User{
    Name: "John Doe",
    Admin: true,
    SomeLargeStruct: BigStruct{
        Field1: "Some value",
        Field2: "Another value",
        Field3: "A third value",
    },
}

logger.Info("Some message", user) // Output: {"level":"INFO","message":"Some message","user":{"Name":"John Doe","Admin":true}}
```

## Key Features

### Static Structured Logging

Ultralogger, unlike some other logging libraries, is designed to be statically structured. That means that you set up
the structure of your log lines beforehand, and then you can use the logger to log your data in a consistent way.

This makes it easier to read and understand your logs, and also makes it easier to use the logger in a multi-threaded
environment. It also makes ultra/log *really fast.*

## TODO

- [ ] Provide a dynamic structured logging interface that allows for more flexibility in logging data.*
- [ ] Add more output formats + improve extensibility of the file formatter interface (CLF, XML, etc.)
- [ ] Improve docs, tests, and examples.
- [ ] Add more examples
- [x] Review processor implementation & improve readability
- [ ] General optimizations; we've got 2 allocs per log line, and disabled levels are taking a little longer than I'd
      like them to.
- [ ] Provide a common benchmark suite for internal benchmarks & comparison w/ other logging libraries.
- [x] Provide a mechanism for allowing the user to flush-on-panic. (E.g. by defining a defer in their main w/ 
      logger.Flush())
- [ ] Make the logger async timeout configurable.

*Work in progress. Plan is for the client to provide a field that defines how a data type should be logged, and then
re-use that field to log all copies of that data type passed to the logger.