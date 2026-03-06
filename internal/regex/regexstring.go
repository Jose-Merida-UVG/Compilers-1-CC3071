package regex

import (
	"strings"
)

// Token represents a character in the regex string, indicating whether it is escaped.
type Token struct {
	Value   rune
	escaped bool
}

// IsOperator returns true if the token is an unescaped operator.
func (e Token) IsOperator() bool {
	operators := []rune{'(', ')', '*', '|', '~'}
	for _, op := range operators {
		if e.Value == op && !e.escaped {
			return true
		}
	}
	return false
}

// RegexString represents a list of characters, indicating whether they are escaped or not.
type RegexString struct {
	Chars []Token
}

// NewRegexString creates a new RegexString from the input string.
func NewRegexString(s string) *RegexString {
	rs := &RegexString{}
	runes := []rune(s) // Convert the string to a slice of runes to handle Unicode correctly
	for i := 0; i < len(runes); i++ {
		currRune := runes[i]
		// Check for escape character
		if currRune == '\\' {
			// Ensure there is a next character to escape
			if i+1 < len(runes) {
				nextRune := runes[i+1]
				rs.Chars = append(rs.Chars, Token{Value: nextRune, escaped: true})
				i++ // Skip the next rune since it's part of the escape sequence
			}
		} else {
			rs.Chars = append(rs.Chars, Token{Value: currRune, escaped: false})
		}
	}
	return rs
}

// String method for RegexString returns the regex as a string with escape sequences properly represented.
func (rs *RegexString) String() string {
	var sb strings.Builder
	for _, token := range rs.Chars {
		if token.escaped {
			sb.WriteRune('\\') // Add the escape character
		}
		sb.WriteRune(token.Value) // Write the actual character
	}
	return sb.String()
}

// AppendEndMarker adds a special end marker ('#') and a concatenation operator
// to the end of the RegexString. Since the string is in postfix notation,
// it appends '#' and then '~'.
func (rs *RegexString) AppendEndMarker() {
	// Append end marker
	rs.Chars = append(rs.Chars, Token{Value: '#', escaped: false})
	// Append concatenation operator to join the previous expression with '#'
	rs.Chars = append(rs.Chars, Token{Value: '~', escaped: false})
}

