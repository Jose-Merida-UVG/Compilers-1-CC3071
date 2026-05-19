package main

const (
	EOF   = 0
	ERROR = -1
	ID = 1
	PLUS = 2
	TIMES = 3
	LPAREN = 4
	RPAREN = 5
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

func NewDragonLexer(input string) *Lexer {
	return &Lexer{input: []rune(input), pos: 0, line: 1, col: 1}
}

var dragonTrans = map[int]map[rune]int{
	1968: {
		'\t': 1969,
		'\n': 1969,
		'\r': 1969,
		' ': 1969,
		'(': 1973,
		')': 1974,
		'*': 1972,
		'+': 1971,
		'A': 1970,
		'B': 1970,
		'C': 1970,
		'D': 1970,
		'E': 1970,
		'F': 1970,
		'G': 1970,
		'H': 1970,
		'I': 1970,
		'J': 1970,
		'K': 1970,
		'L': 1970,
		'M': 1970,
		'N': 1970,
		'O': 1970,
		'P': 1970,
		'Q': 1970,
		'R': 1970,
		'S': 1970,
		'T': 1970,
		'U': 1970,
		'V': 1970,
		'W': 1970,
		'X': 1970,
		'Y': 1970,
		'Z': 1970,
		'a': 1970,
		'b': 1970,
		'c': 1970,
		'd': 1970,
		'e': 1970,
		'f': 1970,
		'g': 1970,
		'h': 1970,
		'i': 1970,
		'j': 1970,
		'k': 1970,
		'l': 1970,
		'm': 1970,
		'n': 1970,
		'o': 1970,
		'p': 1970,
		'q': 1970,
		'r': 1970,
		's': 1970,
		't': 1970,
		'u': 1970,
		'v': 1970,
		'w': 1970,
		'x': 1970,
		'y': 1970,
		'z': 1970,
	},
	1969: {
		'\t': 1969,
		'\n': 1969,
		'\r': 1969,
		' ': 1969,
	},
	1970: {
		'0': 1970,
		'1': 1970,
		'2': 1970,
		'3': 1970,
		'4': 1970,
		'5': 1970,
		'6': 1970,
		'7': 1970,
		'8': 1970,
		'9': 1970,
		'A': 1970,
		'B': 1970,
		'C': 1970,
		'D': 1970,
		'E': 1970,
		'F': 1970,
		'G': 1970,
		'H': 1970,
		'I': 1970,
		'J': 1970,
		'K': 1970,
		'L': 1970,
		'M': 1970,
		'N': 1970,
		'O': 1970,
		'P': 1970,
		'Q': 1970,
		'R': 1970,
		'S': 1970,
		'T': 1970,
		'U': 1970,
		'V': 1970,
		'W': 1970,
		'X': 1970,
		'Y': 1970,
		'Z': 1970,
		'a': 1970,
		'b': 1970,
		'c': 1970,
		'd': 1970,
		'e': 1970,
		'f': 1970,
		'g': 1970,
		'h': 1970,
		'i': 1970,
		'j': 1970,
		'k': 1970,
		'l': 1970,
		'm': 1970,
		'n': 1970,
		'o': 1970,
		'p': 1970,
		'q': 1970,
		'r': 1970,
		's': 1970,
		't': 1970,
		'u': 1970,
		'v': 1970,
		'w': 1970,
		'x': 1970,
		'y': 1970,
		'z': 1970,
	},
}

var dragonAccept = map[int]int{
	1969: 1,
	1970: 2,
	1972: 4,
	1974: 6,
	1971: 3,
	1973: 5,
}

func (l *Lexer) gettoken() int {
	for l.pos < len(l.input) {
		startPos  := l.pos
		startLine := l.line
		startCol  := l.col
		state    := 1968
		lastTok  := 0
		lastPos  := l.pos
		lastLine := l.line
		lastCol  := l.col
		curLine  := l.line
		curCol   := l.col
		for l.pos < len(l.input) {
			ch := l.input[l.pos]
			row, ok := dragonTrans[state]
			if !ok { break }
			next, ok := row[ch]
			if !ok { break }
			l.pos++
			if ch == '\n' { curLine++; curCol = 1 } else { curCol++ }
			state = next
			if tok := dragonAccept[state]; tok != 0 {
				lastTok = tok; lastPos = l.pos; lastLine = curLine; lastCol = curCol
			}
		}
		if lastTok == 0 {
			errStart := startPos; errLn := l.line; errCol := l.col
			for l.pos < len(l.input) {
				ch := l.input[l.pos]
				if row, ok := dragonTrans[1968]; ok { if _, ok := row[ch]; ok { break } }
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
			return ID
		case 3:
			return PLUS
		case 4:
			return TIMES
		case 5:
			return LPAREN
		case 6:
			return RPAREN
		}
	}
	return EOF
}

func (l *Lexer) NextToken() Lexeme {
	tok := l.gettoken()
	return Lexeme{Token: tok, Value: l.Lxm, Line: l.Ln, Col: l.Col}
}
