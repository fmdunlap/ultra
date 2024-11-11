package log

import "strings"

var Brackets = struct {
    Angle  Bracket
    Square Bracket
    Round  Bracket
    Curly  Bracket
    None   Bracket
}{
    Angle:  SimpleBracket{"<", ">"},
    Square: SimpleBracket{"[", "]"},
    Round:  SimpleBracket{"(", ")"},
    Curly:  SimpleBracket{"{", "}"},
    None:   SimpleBracket{"", ""},
}

// Bracket is a type that represents a bracket type. You can use this to wrap log fields if they accept it as an option.
type Bracket interface {
    Open() string
    Close() string
    Wrap(content string) string
}

// SimpleBracket is a simple bracket type that can be used to create custom bracket types.
type SimpleBracket struct {
    open  string
    close string
}

// Open returns the opening bracket for the bracket type.
func (sb SimpleBracket) Open() string {
    return sb.open
}

// Close returns the closing bracket for the bracket type.
func (sb SimpleBracket) Close() string {
    return sb.close
}

// Wrap wraps the content in the bracket type.
func (sb SimpleBracket) Wrap(content string) string {
    b := strings.Builder{}

    b.WriteString(sb.open)
    b.WriteString(content)
    b.WriteString(sb.close)

    return b.String()
}
