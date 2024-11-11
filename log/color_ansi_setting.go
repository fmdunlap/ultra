package log

// ColorSettings are the default Settings supported by Ultralogger. These Settings have a mixed support environment,
// and are only supported by some terminals. They can be used in a ColorizedFormatter to colorize log lines by level.
//
// See https://en.wikipedia.org/wiki/ANSI_escape_code#3-bit_and_4-bit for more info on ANSI Settings.
var ColorSettings = struct {
    Bold          AnsiSetting
    Dim           AnsiSetting
    Italic        AnsiSetting
    Underline     AnsiSetting
    Blink         AnsiSetting
    Strikethrough AnsiSetting
}{
    Bold:          AnsiSetting("1"),
    Dim:           AnsiSetting("2"),
    Italic:        AnsiSetting("3"),
    Underline:     AnsiSetting("4"),
    Blink:         AnsiSetting("5"),
    Strikethrough: AnsiSetting("9"),
}

// AnsiSetting is a type that represents an ANSI setting. It can be applied to arbitrary []byte content through the
// ColorAnsi.SetSetting() method.
type AnsiSetting = []byte
