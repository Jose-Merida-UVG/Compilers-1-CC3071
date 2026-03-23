package regex

import (
	"fmt"
	"sort"
)

// printableASCII holds all printable ASCII runes (0x20–0x7E).
// Used to expand the wildcard _ and to compute negated/difference char classes.
var printableASCII []rune

func init() {
	for r := rune(0x20); r <= 0x7E; r++ {
		printableASCII = append(printableASCII, r)
	}
}

// NormalizePattern converts a raw YALex pattern string into a RegexString whose
// token slice uses only the defined kinds (Lit, Eps, Star, Plus, Quest, Alt,
// LParen, RParen). All YALex-specific syntax is fully resolved here so that the
// rest of the pipeline operates on a clean regex.
func NormalizePattern(pattern string) (*RegexString, error) {
	p := &patParser{runes: []rune(pattern)}
	tokens, err := p.scanAll()
	if err != nil {
		return nil, fmt.Errorf("normalize %q: %w", pattern, err)
	}
	return &RegexString{Chars: tokens}, nil
}

// Parser tracks all chars (runes) + current position
type patParser struct {
	runes []rune
	pos   int
}

// eof when already gone through all chars
func (p *patParser) eof() bool { return p.pos >= len(p.runes) }

// peek gets next rune w/o consuming it
func (p *patParser) peek() rune {
	if p.eof() {
		return 0
	}
	return p.runes[p.pos]
}

// next consumes next rune
func (p *patParser) next() rune { ch := p.runes[p.pos]; p.pos++; return ch }

// skipWS skips whitespace
func (p *patParser) skipWS() {
	for !p.eof() && (p.peek() == ' ' || p.peek() == '\t') {
		p.pos++
	}
}

// scanAll is the main loop, iterate through runes -> append output into Token
func (p *patParser) scanAll() ([]Token, error) {
	var out []Token
	for {
		p.skipWS()
		if p.eof() {
			break
		}
		// This function branches off into different ones
		toks, err := p.scanOne()
		if err != nil {
			return nil, err
		}
		out = append(out, toks...)
	}
	return out, nil
}

// scanOne is the main function called on every iteration, we peek
// to see next char then branch off depending on what this character is.
// Here's where we call the auxiliary functions like reading literals or
// character classes.
func (p *patParser) scanOne() ([]Token, error) {
	switch p.peek() {
	// Read literal 'c'
	case '\'':
		r, err := p.readCharLit()
		if err != nil {
			return nil, err
		}
		return []Token{litToken(r)}, nil

	// Read string literal "string-literal"
	case '"':
		return p.readStringLit()

	// Read [set] or [^set], can also be followed by # for [setA] # [setB]
	case '[':
		return p.readCharClass()

	// Wildcard - any printable ASCII (our alphabet)
	case '_':
		p.next()
		return runeSetToTokens(printableASCII), nil

	// Reading operators, since we've already exhausted the 'special cases' where
	// we could be inside of string or rune literals / character classes we can
	// go ahead and just create these tokens on a switch case
	case '(':
		p.next()
		return []Token{{Kind: LParen, Value: '('}}, nil
	case ')':
		p.next()
		return []Token{{Kind: RParen, Value: ')'}}, nil
	case '*':
		p.next()
		return []Token{{Kind: Star, Value: '*'}}, nil
	case '+':
		p.next()
		return []Token{{Kind: Plus, Value: '+'}}, nil
	case '?':
		p.next()
		return []Token{{Kind: Quest, Value: '?'}}, nil
	case '|':
		p.next()
		return []Token{{Kind: Alt, Value: '|'}}, nil

	// Anything that isn't one of the cases above MUST be eof or
	// an error.
	default:
		ch := p.peek()
		if isIdentStart(ch) {
			id := p.readIdent()
			switch id {
			case "eof":
				return []Token{litToken(RuneEOF)}, nil
			//
			default:
				return nil, fmt.Errorf("unknown identifier %q — let-defs must be expanded before normalization", id)
			}
		}
		return nil, fmt.Errorf("unexpected character %q at position %d", ch, p.pos)
	}
}

// readCharLit parses 'c' or '\t' (positioned at opening ') and returns the rune.
func (p *patParser) readCharLit() (rune, error) {
	// Consume the opening '
	p.next()
	// Read the body
	r, err := p.readCharBody()
	if err != nil {
		return 0, err
	}
	// Check for closing '
	if p.eof() || p.peek() != '\'' {
		return 0, fmt.Errorf("missing closing ' in char literal")
	}
	// Consume closing '
	p.next()
	return r, nil
}

// readCharBody reads one rune (or backslash-escape) from inside a quoted context.
func (p *patParser) readCharBody() (rune, error) {
	if p.eof() {
		return 0, fmt.Errorf("unexpected end of input inside literal")
	}
	// Return the char if not an escape sequence
	ch := p.next()
	if ch != '\\' {
		return ch, nil
	}
	return p.readEscape()
}

// readEscape resolves the character after a backslash, only the six sequences
// defined in the spec are valid.
func (p *patParser) readEscape() (rune, error) {
	if p.eof() {
		return 0, fmt.Errorf("backslash at end of pattern")
	}
	// The 6 defined sequences
	ch := p.next()
	switch ch {
	case 'n':
		return '\n', nil
	case 't':
		return '\t', nil
	case 'r':
		return '\r', nil
	case '\\':
		return '\\', nil
	case '\'':
		return '\'', nil
	case '"':
		return '"', nil
	default:
		return 0, fmt.Errorf("invalid escape sequence \\%c", ch)
	}
}

// readStringLit parses "str" and emits one Lit token per resolved char.
func (p *patParser) readStringLit() ([]Token, error) {
	// Consume opening "
	p.next()
	var toks []Token
	// Iterate through string literal
	for {
		if p.eof() {
			return nil, fmt.Errorf("unclosed string literal")
		}
		// Break once closing " is found
		if p.peek() == '"' {
			p.next()
			break
		}
		// Read char / escape sequence
		r, err := p.readCharBody()
		if err != nil {
			return nil, err
		}
		// Append read char / escape sequence
		toks = append(toks, litToken(r))
	}
	if len(toks) == 0 {
		return nil, fmt.Errorf("empty string literal \"\"")
	}
	return toks, nil
}

// readCharClass parses [set] or [^set], then handles the optional # [set]
// set-difference operator and emits (a|b|c|...) for the resolved rune set.
func (p *patParser) readCharClass() ([]Token, error) {
	// Consume opening [
	p.next()
	// Check negation
	negated := !p.eof() && p.peek() == '^'
	if negated {
		p.next()
	}
	// Get raw contents
	rawSet, err := p.readCharSetContent()
	if err != nil {
		return nil, err
	}
	if p.eof() || p.peek() != ']' {
		return nil, fmt.Errorf("missing ] in char class")
	}
	// Consume closing ]
	p.next()

	// If negated, its alphabet - set
	set := rawSet
	if negated {
		set = runeSetDiff(printableASCII, toRuneSet(rawSet))
	}

	// Set-difference '#': [set1] # [set2] → set1 − set2.
	// Per the spec this has the highest precedence, and both operands must be
	// char classes. We handle it immediately after parsing the left class.
	p.skipWS()
	if !p.eof() && p.peek() == '#' {
		// Consume #
		p.next()
		p.skipWS()
		if p.eof() || p.peek() != '[' {
			return nil, fmt.Errorf("'#' must be followed by a char class [...]")
		}
		rhsToks, err := p.readCharClass() // recursively handles negation/# on the right too
		if err != nil {
			return nil, err
		}
		// Collect the rune set from the already-expanded sequence.
		rhsSet := map[rune]bool{}
		for _, t := range rhsToks {
			if t.Kind == Lit {
				rhsSet[t.Value] = true
			}
		}
		set = runeSetDiff(set, rhsSet)
	}

	if len(set) == 0 {
		return nil, fmt.Errorf("character class resolves to an empty set")
	}
	return runeSetToTokens(set), nil
}

// readCharSetContent parses the interior of [...] and returns the collected rune set.
// Handles:
//   - 'c'         — single char literal (with escape resolution)
//   - 'c1'-'c2'  — inclusive range
//   - "abc"       — string (each char added individually)
func (p *patParser) readCharSetContent() ([]rune, error) {
	seen := map[rune]bool{}
	var set []rune
	add := func(r rune) {
		if !seen[r] {
			seen[r] = true
			set = append(set, r)
		}
	}

	for {
		p.skipWS()
		if p.eof() || p.peek() == ']' {
			break
		}
		switch p.peek() {
		case '\'':
			r1, err := p.readCharLit()
			if err != nil {
				return nil, err
			}
			p.skipWS()
			// Check for 'c1'-'c2' range: '-' must be followed by another char literal.
			if !p.eof() && p.peek() == '-' {
				saved := p.pos
				p.next() // tentatively consume '-'
				p.skipWS()
				if !p.eof() && p.peek() == '\'' {
					r2, err := p.readCharLit()
					if err != nil {
						return nil, err
					}
					if r1 > r2 {
						return nil, fmt.Errorf("invalid range '%c'-'%c': start > end", r1, r2)
					}
					for r := r1; r <= r2; r++ {
						add(r)
					}
				} else {
					// '-' was not part of a range — backtrack and add r1 and '-' literally
					p.pos = saved
					add(r1)
				}
			} else {
				add(r1)
			}

		case '"':
			p.next() // consume opening "
			for {
				if p.eof() {
					return nil, fmt.Errorf("unclosed string inside char class")
				}
				if p.peek() == '"' {
					p.next()
					break
				}
				r, err := p.readCharBody()
				if err != nil {
					return nil, err
				}
				add(r)
			}

		default:
			return nil, fmt.Errorf("unexpected %q inside char class", p.peek())
		}
	}
	return set, nil
}

func isIdentStart(ch rune) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z')
}

func (p *patParser) readIdent() string {
	start := p.pos
	for !p.eof() && (isIdentStart(p.peek()) || (p.peek() >= '0' && p.peek() <= '9') || p.peek() == '_') {
		p.pos++
	}
	return string(p.runes[start:p.pos])
}

// toRuneSet converts a []rune slice to a map for O(1) membership tests.
func toRuneSet(rs []rune) map[rune]bool {
	m := make(map[rune]bool, len(rs))
	for _, r := range rs {
		m[r] = true
	}
	return m
}

// runeSetDiff returns the runes in `a` that are absent from `exclude`,
// preserving the original order of `a`.
func runeSetDiff(a []rune, exclude map[rune]bool) []rune {
	var out []rune
	for _, r := range a {
		if !exclude[r] {
			out = append(out, r)
		}
	}
	return out
}

// runeSetToTokens expands a rune set to LP Lit(r0) Alt Lit(r1) … RP,
// sorted by rune value for deterministic output.
func runeSetToTokens(rs []rune) []Token {
	sorted := make([]rune, len(rs))
	copy(sorted, rs)
	sort.Slice(sorted, func(i, j int) bool { return sorted[i] < sorted[j] })

	out := make([]Token, 0, 2+len(sorted)*2)
	out = append(out, Token{Kind: LParen, Value: '('})
	for i, r := range sorted {
		if i > 0 {
			out = append(out, Token{Kind: Alt, Value: '|'})
		}
		out = append(out, litToken(r))
	}
	out = append(out, Token{Kind: RParen, Value: ')'})
	return out
}
