package regex

import "strings"

// Sentinel runes used internally. All live in Unicode's Private Use Area
// (U+E000–U+F8FF) so they can never appear in validated YALex input (which
// is restricted to printable ASCII 0x20–0x7E plus \t \n \r).
const (
	RuneEndMarker rune = '\uE000' // end-of-regex marker appended by AppendEndMarker
	RuneEpsilon   rune = '\uE001' // empty-string leaf inserted by HandleSpecialOperators
	RuneEOF       rune = '\uE002' // end-of-buffer sentinel from the `eof` keyword
)

// TokenKind classifies the semantic role of a Token after preprocessing.
// This replaces the old "escaped bool" approach: the kind carries all the
// operator/operand information so downstream passes never need to inspect
// the raw rune to decide how to treat a token.
type TokenKind int

const (
	Lit    TokenKind = iota // literal operand: char, end-marker (#)
	Eps                     // ε empty string (inserted during desugaring, never from user input)
	Star                    // * Kleene closure
	Plus                    // + one-or-more (desugared to xx* by HandleSpecialOperators)
	Quest                   // ? optional     (desugared to (x|ε) by HandleSpecialOperators)
	Alt                     // | alternation
	Cat                     // ~ explicit concatenation (inserted by preprocessing)
	LParen                  // ( grouping open
	RParen                  // ) grouping close
)

// Token is the minimal unit in a RegexString. Escaping is resolved at
// construction time: Kind carries the operator/operand role, Value carries
// the resolved rune (only meaningful for Lit tokens, though operators store
// their rune too for debug convenience).
type Token struct {
	Kind  TokenKind
	Value rune
}

// IsOperator reports whether this token acts as an operator in the regex grammar.
// Lit and Eps are both operands (leaves in the AST); everything else is an operator.
func (t Token) IsOperator() bool { return t.Kind != Lit && t.Kind != Eps }

// epsToken is a convenience constructor for the epsilon (empty string) token.
func epsToken() Token { return Token{Kind: Eps, Value: RuneEpsilon} }

// litToken is a convenience constructor for a literal token.
func litToken(ch rune) Token { return Token{Kind: Lit, Value: ch} }

type RegexString struct {
	Chars []Token
}

// String returns a human-readable representation of the token sequence,
// useful for printing pipeline stages during debugging.
func (rs *RegexString) String() string {
	var sb strings.Builder
	for _, t := range rs.Chars {
		sb.WriteRune(t.Value)
	}
	return sb.String()
}

// AppendEndMarker appends the `#` end-marker literal and a Cat operator,
// forming `<regex> ~ #` which anchors the rightmost leaf for direct DFA.
func (rs *RegexString) AppendEndMarker() {
	rs.Chars = append(rs.Chars, litToken(RuneEndMarker))
	rs.Chars = append(rs.Chars, Token{Kind: Cat, Value: '~'})
}
