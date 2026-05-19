package main

const (
	EOF   = 0
	ERROR = -1
	A = 1
	B = 2
	C = 3
	F = 4
	G = 5
	H = 6
)

// Lexeme is the unit the parser consumes: a token ID plus the matched text and position.
type Lexeme struct {
	Token int
	Value string
	Line  int
	Col   int
}

type Lexer struct {
	input []rune
	pos   int
	line  int
	col   int
	Lxm   string
	Ln    int
	Col   int
}

func NewEpsilon_testLexer(input string) *Lexer {
	return &Lexer{input: []rune(input), pos: 0, line: 1, col: 1}
}

var epsilon_testTrans = map[int]map[rune]int{
	2011: {
		'\t': 2012,
		'\n': 2012,
		'\r': 2012,
		' ': 2012,
		'a': 2013,
		'b': 2014,
		'c': 2015,
		'f': 2016,
		'g': 2017,
		'h': 2018,
	},
	2012: {
		'\t': 2012,
		'\n': 2012,
		'\r': 2012,
		' ': 2012,
	},
}

var epsilon_testAccept = map[int]int{
	2012: 1,
	2015: 4,
	2016: 5,
	2013: 2,
	2014: 3,
	2017: 6,
	2018: 7,
}

func (l *Lexer) gettoken() int {
	for l.pos < len(l.input) {
		startPos  := l.pos
		startLine := l.line
		startCol  := l.col
		state    := 2011
		lastTok  := 0
		lastPos  := l.pos
		lastLine := l.line
		lastCol  := l.col
		curLine  := l.line
		curCol   := l.col
		for l.pos < len(l.input) {
			ch := l.input[l.pos]
			row, ok := epsilon_testTrans[state]
			if !ok { break }
			next, ok := row[ch]
			if !ok { break }
			l.pos++
			if ch == '\n' { curLine++; curCol = 1 } else { curCol++ }
			state = next
			if tok := epsilon_testAccept[state]; tok != 0 {
				lastTok = tok; lastPos = l.pos; lastLine = curLine; lastCol = curCol
			}
		}
		if lastTok == 0 {
			errStart := startPos; errLn := l.line; errCol := l.col
			for l.pos < len(l.input) {
				ch := l.input[l.pos]
				if row, ok := epsilon_testTrans[2011]; ok { if _, ok := row[ch]; ok { break } }
				if ch == '\n' { l.line++; l.col = 1 } else { l.col++ }
				l.pos++
			}
			l.Lxm = string(l.input[errStart:l.pos]); l.Ln = errLn; l.Col = errCol
			return ERROR
		}
		l.pos = lastPos; l.line = lastLine; l.col = lastCol
		l.Lxm = string(l.input[startPos:lastPos]); l.Ln = startLine; l.Col = startCol
		switch lastTok {
		case 1:
			// no action
		case 2:
			return A
		case 3:
			return B
		case 4:
			return C
		case 5:
			return F
		case 6:
			return G
		case 7:
			return H
		}
	}
	return EOF
}

func (l *Lexer) NextToken() Lexeme {
	tok := l.gettoken()
	return Lexeme{Token: tok, Value: l.Lxm, Line: l.Ln, Col: l.Col}
}
