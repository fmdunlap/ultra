package log

import (
    "reflect"
    "testing"
)

func TestAllLevels(t *testing.T) {
    tests := []struct {
        name string
        want []Level
    }{
        {
            "AllLevels",
            []Level{
                Debug,
                Info,
                Warn,
                Error,
                Panic,
            },
        },
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            if got := AllLevels(); !reflect.DeepEqual(got, tt.want) {
                t.Errorf("AllLevels() = %v, want %v", got, tt.want)
            }
        })
    }
}

func TestLevel_String(t *testing.T) {
    tests := []struct {
        name string
        l    Level
        want string
    }{
        {"Debug", Debug, "DEBUG"},
        {"Info", Info, "INFO"},
        {"Warn", Warn, "WARN"},
        {"Error", Error, "ERROR"},
        {"Panic", Panic, "PANIC"},
        {"UnknownLevel", Level(42), "UNKNOWN"},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            if got := tt.l.String(); got != tt.want {
                t.Errorf("String() = %v, want %v", got, tt.want)
            }
        })
    }
}

func TestParseLevel(t *testing.T) {
    type args struct {
        levelStr string
    }
    tests := []struct {
        name    string
        args    args
        want    Level
        wantErr bool
    }{
        {"Debug", args{"debug"}, Debug, false},
        {"Info", args{"info"}, Info, false},
        {"Warn", args{"warn"}, Warn, false},
        {"Error", args{"error"}, Error, false},
        {"Panic", args{"panic"}, Panic, false},
        {"InvalidLevel", args{"invalid"}, 0, true},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := ParseLevel(tt.args.levelStr)
            if (err != nil) != tt.wantErr {
                t.Errorf("ParseLevel() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if got != tt.want {
                t.Errorf("ParseLevel() got = %v, want %v", got, tt.want)
            }
        })
    }
}
