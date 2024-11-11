package log

// Color is a type that represents a color. It can be used to colorize arbitrary []byte content.
type Color interface {
    Colorize(str []byte) []byte
}
