package main

import "fmt"

const (
	EOF   = 0
	ERROR = -1
	FLOAT = 1
	INT = 2
	POW = 3
	KWPOW = 4
	PLUS = 5
	MINUS = 6
	TIMES = 7
	DIV = 8
	MOD = 9
	LPAREN = 10
	RPAREN = 11
	COMMA = 12
	IDENT = 13
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

func NewArithmeticLexer(input string) *Lexer {
	return &Lexer{input: []rune(input), pos: 0, line: 1, col: 1}
}

var arithmeticTrans = map[int]map[rune]int{
	85: {
		'\t': 87,
		'\n': 87,
		'\r': 87,
		' ': 87,
		'%': 96,
		'(': 97,
		')': 98,
		'*': 94,
		'+': 92,
		',': 99,
		'-': 93,
		'/': 95,
		'0': 89,
		'1': 89,
		'2': 89,
		'3': 89,
		'4': 89,
		'5': 89,
		'6': 89,
		'7': 89,
		'8': 89,
		'9': 89,
		'A': 100,
		'B': 100,
		'C': 100,
		'D': 100,
		'E': 100,
		'F': 100,
		'G': 100,
		'H': 100,
		'I': 100,
		'J': 100,
		'K': 100,
		'L': 100,
		'M': 100,
		'N': 100,
		'O': 100,
		'P': 100,
		'Q': 100,
		'R': 100,
		'S': 100,
		'T': 100,
		'U': 100,
		'V': 100,
		'W': 100,
		'X': 100,
		'Y': 100,
		'Z': 100,
		'a': 100,
		'b': 100,
		'c': 100,
		'd': 100,
		'e': 100,
		'f': 100,
		'g': 100,
		'h': 100,
		'i': 100,
		'j': 100,
		'k': 100,
		'l': 100,
		'm': 100,
		'n': 100,
		'o': 100,
		'p': 101,
		'q': 100,
		'r': 100,
		's': 100,
		't': 100,
		'u': 100,
		'v': 100,
		'w': 100,
		'x': 100,
		'y': 100,
		'z': 100,
	},
	100: {
		'0': 100,
		'1': 100,
		'2': 100,
		'3': 100,
		'4': 100,
		'5': 100,
		'6': 100,
		'7': 100,
		'8': 100,
		'9': 100,
		'A': 100,
		'B': 100,
		'C': 100,
		'D': 100,
		'E': 100,
		'F': 100,
		'G': 100,
		'H': 100,
		'I': 100,
		'J': 100,
		'K': 100,
		'L': 100,
		'M': 100,
		'N': 100,
		'O': 100,
		'P': 100,
		'Q': 100,
		'R': 100,
		'S': 100,
		'T': 100,
		'U': 100,
		'V': 100,
		'W': 100,
		'X': 100,
		'Y': 100,
		'Z': 100,
		'a': 100,
		'b': 100,
		'c': 100,
		'd': 100,
		'e': 100,
		'f': 100,
		'g': 100,
		'h': 100,
		'i': 100,
		'j': 100,
		'k': 100,
		'l': 100,
		'm': 100,
		'n': 100,
		'o': 100,
		'p': 100,
		'q': 100,
		'r': 100,
		's': 100,
		't': 100,
		'u': 100,
		'v': 100,
		'w': 100,
		'x': 100,
		'y': 100,
		'z': 100,
	},
	87: {
		'\t': 87,
		'\n': 87,
		'\r': 87,
		' ': 87,
	},
	89: {
		'.': 86,
		'0': 89,
		'1': 89,
		'2': 89,
		'3': 89,
		'4': 89,
		'5': 89,
		'6': 89,
		'7': 89,
		'8': 89,
		'9': 89,
	},
	86: {
		'0': 88,
		'1': 88,
		'2': 88,
		'3': 88,
		'4': 88,
		'5': 88,
		'6': 88,
		'7': 88,
		'8': 88,
		'9': 88,
	},
	88: {
		'0': 88,
		'1': 88,
		'2': 88,
		'3': 88,
		'4': 88,
		'5': 88,
		'6': 88,
		'7': 88,
		'8': 88,
		'9': 88,
	},
	101: {
		'0': 100,
		'1': 100,
		'2': 100,
		'3': 100,
		'4': 100,
		'5': 100,
		'6': 100,
		'7': 100,
		'8': 100,
		'9': 100,
		'A': 100,
		'B': 100,
		'C': 100,
		'D': 100,
		'E': 100,
		'F': 100,
		'G': 100,
		'H': 100,
		'I': 100,
		'J': 100,
		'K': 100,
		'L': 100,
		'M': 100,
		'N': 100,
		'O': 100,
		'P': 100,
		'Q': 100,
		'R': 100,
		'S': 100,
		'T': 100,
		'U': 100,
		'V': 100,
		'W': 100,
		'X': 100,
		'Y': 100,
		'Z': 100,
		'a': 100,
		'b': 100,
		'c': 100,
		'd': 100,
		'e': 100,
		'f': 100,
		'g': 100,
		'h': 100,
		'i': 100,
		'j': 100,
		'k': 100,
		'l': 100,
		'm': 100,
		'n': 100,
		'o': 102,
		'p': 100,
		'q': 100,
		'r': 100,
		's': 100,
		't': 100,
		'u': 100,
		'v': 100,
		'w': 100,
		'x': 100,
		'y': 100,
		'z': 100,
	},
	102: {
		'0': 100,
		'1': 100,
		'2': 100,
		'3': 100,
		'4': 100,
		'5': 100,
		'6': 100,
		'7': 100,
		'8': 100,
		'9': 100,
		'A': 100,
		'B': 100,
		'C': 100,
		'D': 100,
		'E': 100,
		'F': 100,
		'G': 100,
		'H': 100,
		'I': 100,
		'J': 100,
		'K': 100,
		'L': 100,
		'M': 100,
		'N': 100,
		'O': 100,
		'P': 100,
		'Q': 100,
		'R': 100,
		'S': 100,
		'T': 100,
		'U': 100,
		'V': 100,
		'W': 100,
		'X': 100,
		'Y': 100,
		'Z': 100,
		'a': 100,
		'b': 100,
		'c': 100,
		'd': 100,
		'e': 100,
		'f': 100,
		'g': 100,
		'h': 100,
		'i': 100,
		'j': 100,
		'k': 100,
		'l': 100,
		'm': 100,
		'n': 100,
		'o': 100,
		'p': 100,
		'q': 100,
		'r': 100,
		's': 100,
		't': 100,
		'u': 100,
		'v': 100,
		'w': 91,
		'x': 100,
		'y': 100,
		'z': 100,
	},
	91: {
		'0': 100,
		'1': 100,
		'2': 100,
		'3': 100,
		'4': 100,
		'5': 100,
		'6': 100,
		'7': 100,
		'8': 100,
		'9': 100,
		'A': 100,
		'B': 100,
		'C': 100,
		'D': 100,
		'E': 100,
		'F': 100,
		'G': 100,
		'H': 100,
		'I': 100,
		'J': 100,
		'K': 100,
		'L': 100,
		'M': 100,
		'N': 100,
		'O': 100,
		'P': 100,
		'Q': 100,
		'R': 100,
		'S': 100,
		'T': 100,
		'U': 100,
		'V': 100,
		'W': 100,
		'X': 100,
		'Y': 100,
		'Z': 100,
		'a': 100,
		'b': 100,
		'c': 100,
		'd': 100,
		'e': 100,
		'f': 100,
		'g': 100,
		'h': 100,
		'i': 100,
		'j': 100,
		'k': 100,
		'l': 100,
		'm': 100,
		'n': 100,
		'o': 100,
		'p': 100,
		'q': 100,
		'r': 100,
		's': 100,
		't': 100,
		'u': 100,
		'v': 100,
		'w': 100,
		'x': 100,
		'y': 100,
		'z': 100,
	},
	94: {
		'*': 90,
	},
}

var arithmeticAccept = map[int]int{
	92: 6,
	100: 14,
	98: 12,
	96: 10,
	87: 1,
	89: 3,
	88: 2,
	93: 7,
	97: 11,
	101: 14,
	102: 14,
	91: 5,
	94: 8,
	90: 4,
	99: 13,
	95: 9,
}

func (l *Lexer) gettoken() int {
	for l.pos < len(l.input) {
		startPos  := l.pos
		startLine := l.line
		startCol  := l.col
		state    := 85
		lastTok  := 0
		lastPos  := l.pos
		lastLine := l.line
		lastCol  := l.col
		curLine  := l.line
		curCol   := l.col
		for l.pos < len(l.input) {
			ch := l.input[l.pos]
			row, ok := arithmeticTrans[state]
			if !ok { break }
			next, ok := row[ch]
			if !ok { break }
			l.pos++
			if ch == '\n' { curLine++; curCol = 1 } else { curCol++ }
			state = next
			if tok := arithmeticAccept[state]; tok != 0 {
				lastTok = tok; lastPos = l.pos; lastLine = curLine; lastCol = curCol
			}
		}
		if lastTok == 0 {
			errStart := startPos; errLn := l.line; errCol := l.col
			for l.pos < len(l.input) {
				ch := l.input[l.pos]
				if row, ok := arithmeticTrans[85]; ok { if _, ok := row[ch]; ok { break } }
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
			fmt.Print("FLOAT  "); return FLOAT
		case 3:
			fmt.Print("INT    "); return INT
		case 4:
			fmt.Print("POW    "); return POW
		case 5:
			fmt.Print("KWPOW  "); return KWPOW
		case 6:
			fmt.Print("PLUS   "); return PLUS
		case 7:
			fmt.Print("MINUS  "); return MINUS
		case 8:
			fmt.Print("TIMES  "); return TIMES
		case 9:
			fmt.Print("DIV    "); return DIV
		case 10:
			fmt.Print("MOD    "); return MOD
		case 11:
			fmt.Print("LPAREN "); return LPAREN
		case 12:
			fmt.Print("RPAREN "); return RPAREN
		case 13:
			fmt.Print("COMMA  "); return COMMA
		case 14:
			fmt.Print("IDENT  "); return IDENT
		}
	}
	return EOF
}

func (l *Lexer) NextToken() Lexeme {
	tok := l.gettoken()
	return Lexeme{Token: tok, Value: l.Lxm, Line: l.Ln, Col: l.Col}
}
