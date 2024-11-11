package log

import (
    "strings"
)

// Level is a type representing the level of a log message.
//
// It can be one of the following:
//   - Debug
//   - Info
//   - Warn
//   - Error
//   - Panic
//
// Levels determine the priority of a log message, and can be hidden if a logger's minimum level is set to a higher
// level than the message's level.
//
// For example, if a logger's minimum level is set to Warn, then a message with a level of Info will not be
// written to the output.
type Level int

const (
    Debug Level = iota
    Info
    Warn
    Error
    Panic
)

// AllLevels returns a slice of all available levels.
func AllLevels() []Level {
    return []Level{
        Debug,
        Info,
        Warn,
        Error,
        Panic,
    }
}

func (l Level) String() string {
    switch l {
    case Debug:
        return "DEBUG"
    case Info:
        return "INFO"
    case Warn:
        return "WARN"
    case Error:
        return "ERROR"
    case Panic:
        return "PANIC"
    default:
        return "UNKNOWN"
    }
}

// ParseLevel parses a string into a Level. Returns an error if the string is not a valid Level.
func ParseLevel(levelStr string) (Level, error) {
    switch strings.ToLower(levelStr) {
    case "debug":
        return Debug, nil
    case "info":
        return Info, nil
    case "warn":
        return Warn, nil
    case "error":
        return Error, nil
    case "panic":
        return Panic, nil
    default:
        return 0, &ErrorLevelParsing{level: levelStr}
    }
}
