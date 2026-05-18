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
	187: {
		'\t': 189,
		'\n': 189,
		'\r': 189,
		' ': 189,
		'%': 198,
		'(': 199,
		')': 200,
		'*': 196,
		'+': 194,
		',': 201,
		'-': 195,
		'/': 197,
		'0': 191,
		'1': 191,
		'2': 191,
		'3': 191,
		'4': 191,
		'5': 191,
		'6': 191,
		'7': 191,
		'8': 191,
		'9': 191,
		'A': 202,
		'B': 202,
		'C': 202,
		'D': 202,
		'E': 202,
		'F': 202,
		'G': 202,
		'H': 202,
		'I': 202,
		'J': 202,
		'K': 202,
		'L': 202,
		'M': 202,
		'N': 202,
		'O': 202,
		'P': 202,
		'Q': 202,
		'R': 202,
		'S': 202,
		'T': 202,
		'U': 202,
		'V': 202,
		'W': 202,
		'X': 202,
		'Y': 202,
		'Z': 202,
		'a': 202,
		'b': 202,
		'c': 202,
		'd': 202,
		'e': 202,
		'f': 202,
		'g': 202,
		'h': 202,
		'i': 202,
		'j': 202,
		'k': 202,
		'l': 202,
		'm': 202,
		'n': 202,
		'o': 202,
		'p': 203,
		'q': 202,
		'r': 202,
		's': 202,
		't': 202,
		'u': 202,
		'v': 202,
		'w': 202,
		'x': 202,
		'y': 202,
		'z': 202,
	},
	189: {
		'\t': 189,
		'\n': 189,
		'\r': 189,
		' ': 189,
	},
	202: {
		'0': 202,
		'1': 202,
		'2': 202,
		'3': 202,
		'4': 202,
		'5': 202,
		'6': 202,
		'7': 202,
		'8': 202,
		'9': 202,
		'A': 202,
		'B': 202,
		'C': 202,
		'D': 202,
		'E': 202,
		'F': 202,
		'G': 202,
		'H': 202,
		'I': 202,
		'J': 202,
		'K': 202,
		'L': 202,
		'M': 202,
		'N': 202,
		'O': 202,
		'P': 202,
		'Q': 202,
		'R': 202,
		'S': 202,
		'T': 202,
		'U': 202,
		'V': 202,
		'W': 202,
		'X': 202,
		'Y': 202,
		'Z': 202,
		'a': 202,
		'b': 202,
		'c': 202,
		'd': 202,
		'e': 202,
		'f': 202,
		'g': 202,
		'h': 202,
		'i': 202,
		'j': 202,
		'k': 202,
		'l': 202,
		'm': 202,
		'n': 202,
		'o': 202,
		'p': 202,
		'q': 202,
		'r': 202,
		's': 202,
		't': 202,
		'u': 202,
		'v': 202,
		'w': 202,
		'x': 202,
		'y': 202,
		'z': 202,
	},
	191: {
		'.': 188,
		'0': 191,
		'1': 191,
		'2': 191,
		'3': 191,
		'4': 191,
		'5': 191,
		'6': 191,
		'7': 191,
		'8': 191,
		'9': 191,
	},
	188: {
		'0': 190,
		'1': 190,
		'2': 190,
		'3': 190,
		'4': 190,
		'5': 190,
		'6': 190,
		'7': 190,
		'8': 190,
		'9': 190,
	},
	190: {
		'0': 190,
		'1': 190,
		'2': 190,
		'3': 190,
		'4': 190,
		'5': 190,
		'6': 190,
		'7': 190,
		'8': 190,
		'9': 190,
	},
	196: {
		'*': 192,
	},
	203: {
		'0': 202,
		'1': 202,
		'2': 202,
		'3': 202,
		'4': 202,
		'5': 202,
		'6': 202,
		'7': 202,
		'8': 202,
		'9': 202,
		'A': 202,
		'B': 202,
		'C': 202,
		'D': 202,
		'E': 202,
		'F': 202,
		'G': 202,
		'H': 202,
		'I': 202,
		'J': 202,
		'K': 202,
		'L': 202,
		'M': 202,
		'N': 202,
		'O': 202,
		'P': 202,
		'Q': 202,
		'R': 202,
		'S': 202,
		'T': 202,
		'U': 202,
		'V': 202,
		'W': 202,
		'X': 202,
		'Y': 202,
		'Z': 202,
		'a': 202,
		'b': 202,
		'c': 202,
		'd': 202,
		'e': 202,
		'f': 202,
		'g': 202,
		'h': 202,
		'i': 202,
		'j': 202,
		'k': 202,
		'l': 202,
		'm': 202,
		'n': 202,
		'o': 204,
		'p': 202,
		'q': 202,
		'r': 202,
		's': 202,
		't': 202,
		'u': 202,
		'v': 202,
		'w': 202,
		'x': 202,
		'y': 202,
		'z': 202,
	},
	204: {
		'0': 202,
		'1': 202,
		'2': 202,
		'3': 202,
		'4': 202,
		'5': 202,
		'6': 202,
		'7': 202,
		'8': 202,
		'9': 202,
		'A': 202,
		'B': 202,
		'C': 202,
		'D': 202,
		'E': 202,
		'F': 202,
		'G': 202,
		'H': 202,
		'I': 202,
		'J': 202,
		'K': 202,
		'L': 202,
		'M': 202,
		'N': 202,
		'O': 202,
		'P': 202,
		'Q': 202,
		'R': 202,
		'S': 202,
		'T': 202,
		'U': 202,
		'V': 202,
		'W': 202,
		'X': 202,
		'Y': 202,
		'Z': 202,
		'a': 202,
		'b': 202,
		'c': 202,
		'd': 202,
		'e': 202,
		'f': 202,
		'g': 202,
		'h': 202,
		'i': 202,
		'j': 202,
		'k': 202,
		'l': 202,
		'm': 202,
		'n': 202,
		'o': 202,
		'p': 202,
		'q': 202,
		'r': 202,
		's': 202,
		't': 202,
		'u': 202,
		'v': 202,
		'w': 193,
		'x': 202,
		'y': 202,
		'z': 202,
	},
	193: {
		'0': 202,
		'1': 202,
		'2': 202,
		'3': 202,
		'4': 202,
		'5': 202,
		'6': 202,
		'7': 202,
		'8': 202,
		'9': 202,
		'A': 202,
		'B': 202,
		'C': 202,
		'D': 202,
		'E': 202,
		'F': 202,
		'G': 202,
		'H': 202,
		'I': 202,
		'J': 202,
		'K': 202,
		'L': 202,
		'M': 202,
		'N': 202,
		'O': 202,
		'P': 202,
		'Q': 202,
		'R': 202,
		'S': 202,
		'T': 202,
		'U': 202,
		'V': 202,
		'W': 202,
		'X': 202,
		'Y': 202,
		'Z': 202,
		'a': 202,
		'b': 202,
		'c': 202,
		'd': 202,
		'e': 202,
		'f': 202,
		'g': 202,
		'h': 202,
		'i': 202,
		'j': 202,
		'k': 202,
		'l': 202,
		'm': 202,
		'n': 202,
		'o': 202,
		'p': 202,
		'q': 202,
		'r': 202,
		's': 202,
		't': 202,
		'u': 202,
		'v': 202,
		'w': 202,
		'x': 202,
		'y': 202,
		'z': 202,
	},
}

var arithmeticAccept = map[int]int{
	189: 1,
	202: 14,
	201: 13,
	191: 3,
	190: 2,
	200: 12,
	194: 6,
	196: 8,
	192: 4,
	203: 14,
	204: 14,
	193: 5,
	195: 7,
	199: 11,
	198: 10,
	197: 9,
}

func (l *Lexer) gettoken() int {
	for l.pos < len(l.input) {
		startPos  := l.pos
		startLine := l.line
		startCol  := l.col
		state    := 187
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
				if row, ok := arithmeticTrans[187]; ok { if _, ok := row[ch]; ok { break } }
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
