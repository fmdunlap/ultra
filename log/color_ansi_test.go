package log

import (
    "bytes"
    "fmt"
    "testing"
)

func ExampleColorAnsi_Colorize() {
    // Colorize a string of text with the Red color.
    red := Colors.Red
    colorizedBytes := red.Colorize([]byte("This is red!"))

    // \033[31mThis is red!\033[0m
    fmt.Println(colorizedBytes)

    // Output:
    // [27 91 51 49 109 84 104 105 115 32 105 115 32 114 101 100 33 27 91 48 109]
}

func ExampleColorAnsi_Colorize_multiple() {
    // Colorize a string of text with the Red color, the Bold setting, and a yellow background.
    redBoldYellowBackground := Colors.Red.Bold().SetBackground(BackgroundColors.Yellow)
    colorizedBytes := redBoldYellowBackground.Colorize([]byte("This is BOLD red with a yellow background!"))

    // \033[1;48;2;255;255;0;31mThis is BOLD red with a yellow background!\033[0m
    fmt.Println(colorizedBytes)

    // Output:
    // [27 91 49 59 52 51 59 51 49 109 84 104 105 115 32 105 115 32 66 79 76 68 32 114 101 100 32 119 105 116 104 32 97 32 121 101 108 108 111 119 32 98 97 99 107 103 114 111 117 110 100 33 27 91 48 109]
}

func ExampleColorAnsi_Colorize_rgb() {
    // Colorize a string of text with an RGB color and an RGB background.
    redForegroundGreenBackground := ColorAnsiRGB(255, 0, 0).SetBackground(BackgroundRGB(0, 255, 0))
    colorized := redForegroundGreenBackground.Colorize([]byte("This is red text on a green background!"))

    // \033[48;2;0;255;0;38;2;255;0;0mThis is red text on a green background!\033[0m
    fmt.Println(colorized)

    // Output:
    // [27 91 52 56 59 50 59 48 59 50 53 53 59 48 59 51 56 59 50 59 50 53 53 59 48 59 48 109 84 104 105 115 32 105 115 32 114 101 100 32 116 101 120 116 32 111 110 32 97 32 103 114 101 101 110 32 98 97 99 107 103 114 111 117 110 100 33 27 91 48 109]

}

func TestAnsiColor_Colorize(t *testing.T) {
    tests := []struct {
        name string
        msg  []byte
        c    ColorAnsi
        want []byte
    }{
        {
            name: "ColorRed",
            msg:  []byte("test"),
            c:    Colors.Red,
            want: []byte("\033[31mtest\033[0m"),
        },
        {
            name: "Bold",
            msg:  []byte("test"),
            c: ColorAnsi{
                Code:     []byte("31"),
                Settings: []AnsiSetting{ColorSettings.Bold},
            },
            want: []byte("\033[1;31mtest\033[0m"),
        },
        {
            name: "Dim",
            msg:  []byte("test"),
            c: ColorAnsi{
                Code:     []byte("31"),
                Settings: []AnsiSetting{ColorSettings.Dim},
            },
            want: []byte("\033[2;31mtest\033[0m"),
        },
        {
            name: "Italic",
            msg:  []byte("test"),
            c: ColorAnsi{
                Code:     []byte("31"),
                Settings: []AnsiSetting{ColorSettings.Italic},
            },
            want: []byte("\033[3;31mtest\033[0m"),
        },
        {
            name: "Underline",
            msg:  []byte("test"),
            c: ColorAnsi{
                Code:     []byte("31"),
                Settings: []AnsiSetting{ColorSettings.Underline},
            },
            want: []byte("\033[4;31mtest\033[0m"),
        },
        {
            name: "Blink",
            msg:  []byte("test"),
            c: ColorAnsi{
                Code:     []byte("31"),
                Settings: []AnsiSetting{ColorSettings.Blink},
            },
            want: []byte("\033[5;31mtest\033[0m"),
        },
        {
            name: "Strikethrough",
            msg:  []byte("test"),
            c: ColorAnsi{
                Code:     []byte("31"),
                Settings: []AnsiSetting{ColorSettings.Strikethrough},
            },
            want: []byte("\033[9;31mtest\033[0m"),
        },
        {
            name: "Multiple Settings",
            msg:  []byte("test"),
            c: ColorAnsi{
                Code:     []byte("31"),
                Settings: []AnsiSetting{ColorSettings.Bold, ColorSettings.Italic, ColorSettings.Underline, ColorSettings.Blink, ColorSettings.Strikethrough},
            },
            want: []byte("\033[1;3;4;5;9;31mtest\033[0m"),
        },
        {
            name: "ColorAnsiRGB",
            msg:  []byte("test"),
            c:    ColorAnsiRGB(138, 206, 0),
            want: []byte("\033[38;2;138;206;0mtest\033[0m"),
        },
        {
            name: "BackgroundRed",
            msg:  []byte("test"),
            c: ColorAnsi{
                Code:       []byte("30"),
                Settings:   []AnsiSetting{},
                Background: BackgroundColors.Red,
            },
            want: []byte("\033[41;30mtest\033[0m"),
        },
        {
            name: "BackgroundRGB",
            msg:  []byte("test"),
            c: ColorAnsi{
                Code:       []byte("30"),
                Settings:   []AnsiSetting{},
                Background: BackgroundRGB(138, 206, 0),
            },
            want: []byte("\033[48;2;138;206;0;30mtest\033[0m"),
        },
        {
            name: "BackgroundRGB + Bold",
            msg:  []byte("test"),
            c: ColorAnsi{
                Code:       []byte("30"),
                Settings:   []AnsiSetting{ColorSettings.Bold},
                Background: BackgroundRGB(138, 206, 0),
            },
            want: []byte("\033[1;48;2;138;206;0;30mtest\033[0m"),
        },
        {
            name: "BackgroundRed + Multiple Settings",
            msg:  []byte("test"),
            c: ColorAnsi{
                Code:       []byte("30"),
                Settings:   []AnsiSetting{ColorSettings.Bold, ColorSettings.Italic, ColorSettings.Underline, ColorSettings.Blink, ColorSettings.Strikethrough},
                Background: BackgroundColors.Red,
            },
            want: []byte("\033[1;3;4;5;9;41;30mtest\033[0m"),
        },
        {
            name: "BackgroundRGB + Multiple Settings",
            msg:  []byte("test"),
            c: ColorAnsi{
                Code:       []byte("30"),
                Settings:   []AnsiSetting{ColorSettings.Bold, ColorSettings.Italic, ColorSettings.Underline, ColorSettings.Blink, ColorSettings.Strikethrough},
                Background: BackgroundRGB(138, 206, 0),
            },
            want: []byte("\033[1;3;4;5;9;48;2;138;206;0;30mtest\033[0m"),
        },
        {
            name: "ColorRGB + BackgroundRGB",
            msg:  []byte("test"),
            c:    ColorAnsiRGB(138, 206, 0).SetBackground(BackgroundRGB(255, 0, 0)),
            want: []byte("\033[48;2;255;0;0;38;2;138;206;0mtest\033[0m"),
        },
        {
            name: "ColorRGB + BackgroundRGB + Bold",
            msg:  []byte("test"),
            c:    ColorAnsiRGB(138, 206, 0).SetBackground(BackgroundRGB(255, 0, 0)).Bold(),
            want: []byte("\033[1;48;2;255;0;0;38;2;138;206;0mtest\033[0m"),
        },
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got := tt.c.Colorize(tt.msg)
            if !bytes.Equal(got, tt.want) {
                fmt.Println("Got:  ", got)
                fmt.Println("Want: ", tt.want)
                t.Errorf("Colorize() = %v, want %v", got, tt.want)
            }
        })
    }
}

func TestAnsiColor_totalBufferLength(t *testing.T) {
    tests := []struct {
        name  string
        c     ColorAnsi
        input []byte
        want  int
    }{
        {
            name: "No Settings",
            c: ColorAnsi{
                Code:       []byte("31"),
                Settings:   []AnsiSetting{},
                Background: nil,
                // output:     "\033[31mtest\033[0m",
            },
            input: []byte("test"),
            want:  13,
        },
        {
            name: "Bold",
            c: ColorAnsi{
                Code:       []byte("31"),
                Settings:   []AnsiSetting{ColorSettings.Bold},
                Background: nil,
                // output:     "\033[1;31mtest\033[0m",
            },
            input: []byte("test"),
            want:  15,
        },
        {
            name: "Multiple Settings",
            c: ColorAnsi{
                Code:       []byte("31"),
                Settings:   []AnsiSetting{ColorSettings.Bold, ColorSettings.Italic, ColorSettings.Underline, ColorSettings.Blink, ColorSettings.Strikethrough},
                Background: nil,
                // output:     "\033[1;3;4;5;9;31mtest\033[0m",
            },
            input: []byte("test"),
            want:  23,
        },
        {
            name: "ColorAnsiRGB",
            c:    ColorAnsiRGB(138, 206, 0),
            // output: "\033[38;2;138;206;0mtest\033[0m",
            input: []byte("test"),
            want:  25,
        },
        {
            name: "BackgroundRed",
            c: ColorAnsi{
                Code:       []byte("30"),
                Settings:   []AnsiSetting{},
                Background: BackgroundColors.Red,
                // output:     "\033[41;30mtest\033[0m",
            },
            input: []byte("test"),
            want:  16,
        },
        {
            name: "BackgroundRGB",
            c: ColorAnsi{
                Code:       []byte("30"),
                Settings:   []AnsiSetting{},
                Background: BackgroundRGB(138, 206, 0),
                // output:     "\033[48;2;138;206;0;30mtest\033[0m",
            },
            input: []byte("test"),
            want:  28,
        },
        {
            name:  "ColorRGB + BackgroundRGB",
            c:     ColorAnsiRGB(138, 206, 0).SetBackground(BackgroundRGB(255, 0, 0)),
            input: []byte("test"),
            want:  38,
        },
        {
            name:  "ColorRGB + BackgroundRGB + Bold",
            c:     ColorAnsiRGB(138, 206, 0).SetBackground(BackgroundRGB(255, 0, 0)).Bold(),
            input: []byte("test"),
            want:  40,
        },
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            if got := tt.c.totalBufferLength(tt.input); got != tt.want {
                t.Errorf("totalBufferLength() = %v, want %v", got, tt.want)
            }
        })
    }
}
