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
	919: {
		'\t': 921,
		'\n': 921,
		'\r': 921,
		' ': 921,
		'%': 930,
		'(': 931,
		')': 932,
		'*': 928,
		'+': 926,
		',': 933,
		'-': 927,
		'/': 929,
		'0': 923,
		'1': 923,
		'2': 923,
		'3': 923,
		'4': 923,
		'5': 923,
		'6': 923,
		'7': 923,
		'8': 923,
		'9': 923,
		'A': 934,
		'B': 934,
		'C': 934,
		'D': 934,
		'E': 934,
		'F': 934,
		'G': 934,
		'H': 934,
		'I': 934,
		'J': 934,
		'K': 934,
		'L': 934,
		'M': 934,
		'N': 934,
		'O': 934,
		'P': 934,
		'Q': 934,
		'R': 934,
		'S': 934,
		'T': 934,
		'U': 934,
		'V': 934,
		'W': 934,
		'X': 934,
		'Y': 934,
		'Z': 934,
		'a': 934,
		'b': 934,
		'c': 934,
		'd': 934,
		'e': 934,
		'f': 934,
		'g': 934,
		'h': 934,
		'i': 934,
		'j': 934,
		'k': 934,
		'l': 934,
		'm': 934,
		'n': 934,
		'o': 934,
		'p': 935,
		'q': 934,
		'r': 934,
		's': 934,
		't': 934,
		'u': 934,
		'v': 934,
		'w': 934,
		'x': 934,
		'y': 934,
		'z': 934,
	},
	923: {
		'.': 920,
		'0': 923,
		'1': 923,
		'2': 923,
		'3': 923,
		'4': 923,
		'5': 923,
		'6': 923,
		'7': 923,
		'8': 923,
		'9': 923,
	},
	920: {
		'0': 922,
		'1': 922,
		'2': 922,
		'3': 922,
		'4': 922,
		'5': 922,
		'6': 922,
		'7': 922,
		'8': 922,
		'9': 922,
	},
	922: {
		'0': 922,
		'1': 922,
		'2': 922,
		'3': 922,
		'4': 922,
		'5': 922,
		'6': 922,
		'7': 922,
		'8': 922,
		'9': 922,
	},
	934: {
		'0': 934,
		'1': 934,
		'2': 934,
		'3': 934,
		'4': 934,
		'5': 934,
		'6': 934,
		'7': 934,
		'8': 934,
		'9': 934,
		'A': 934,
		'B': 934,
		'C': 934,
		'D': 934,
		'E': 934,
		'F': 934,
		'G': 934,
		'H': 934,
		'I': 934,
		'J': 934,
		'K': 934,
		'L': 934,
		'M': 934,
		'N': 934,
		'O': 934,
		'P': 934,
		'Q': 934,
		'R': 934,
		'S': 934,
		'T': 934,
		'U': 934,
		'V': 934,
		'W': 934,
		'X': 934,
		'Y': 934,
		'Z': 934,
		'a': 934,
		'b': 934,
		'c': 934,
		'd': 934,
		'e': 934,
		'f': 934,
		'g': 934,
		'h': 934,
		'i': 934,
		'j': 934,
		'k': 934,
		'l': 934,
		'm': 934,
		'n': 934,
		'o': 934,
		'p': 934,
		'q': 934,
		'r': 934,
		's': 934,
		't': 934,
		'u': 934,
		'v': 934,
		'w': 934,
		'x': 934,
		'y': 934,
		'z': 934,
	},
	921: {
		'\t': 921,
		'\n': 921,
		'\r': 921,
		' ': 921,
	},
	935: {
		'0': 934,
		'1': 934,
		'2': 934,
		'3': 934,
		'4': 934,
		'5': 934,
		'6': 934,
		'7': 934,
		'8': 934,
		'9': 934,
		'A': 934,
		'B': 934,
		'C': 934,
		'D': 934,
		'E': 934,
		'F': 934,
		'G': 934,
		'H': 934,
		'I': 934,
		'J': 934,
		'K': 934,
		'L': 934,
		'M': 934,
		'N': 934,
		'O': 934,
		'P': 934,
		'Q': 934,
		'R': 934,
		'S': 934,
		'T': 934,
		'U': 934,
		'V': 934,
		'W': 934,
		'X': 934,
		'Y': 934,
		'Z': 934,
		'a': 934,
		'b': 934,
		'c': 934,
		'd': 934,
		'e': 934,
		'f': 934,
		'g': 934,
		'h': 934,
		'i': 934,
		'j': 934,
		'k': 934,
		'l': 934,
		'm': 934,
		'n': 934,
		'o': 936,
		'p': 934,
		'q': 934,
		'r': 934,
		's': 934,
		't': 934,
		'u': 934,
		'v': 934,
		'w': 934,
		'x': 934,
		'y': 934,
		'z': 934,
	},
	936: {
		'0': 934,
		'1': 934,
		'2': 934,
		'3': 934,
		'4': 934,
		'5': 934,
		'6': 934,
		'7': 934,
		'8': 934,
		'9': 934,
		'A': 934,
		'B': 934,
		'C': 934,
		'D': 934,
		'E': 934,
		'F': 934,
		'G': 934,
		'H': 934,
		'I': 934,
		'J': 934,
		'K': 934,
		'L': 934,
		'M': 934,
		'N': 934,
		'O': 934,
		'P': 934,
		'Q': 934,
		'R': 934,
		'S': 934,
		'T': 934,
		'U': 934,
		'V': 934,
		'W': 934,
		'X': 934,
		'Y': 934,
		'Z': 934,
		'a': 934,
		'b': 934,
		'c': 934,
		'd': 934,
		'e': 934,
		'f': 934,
		'g': 934,
		'h': 934,
		'i': 934,
		'j': 934,
		'k': 934,
		'l': 934,
		'm': 934,
		'n': 934,
		'o': 934,
		'p': 934,
		'q': 934,
		'r': 934,
		's': 934,
		't': 934,
		'u': 934,
		'v': 934,
		'w': 925,
		'x': 934,
		'y': 934,
		'z': 934,
	},
	925: {
		'0': 934,
		'1': 934,
		'2': 934,
		'3': 934,
		'4': 934,
		'5': 934,
		'6': 934,
		'7': 934,
		'8': 934,
		'9': 934,
		'A': 934,
		'B': 934,
		'C': 934,
		'D': 934,
		'E': 934,
		'F': 934,
		'G': 934,
		'H': 934,
		'I': 934,
		'J': 934,
		'K': 934,
		'L': 934,
		'M': 934,
		'N': 934,
		'O': 934,
		'P': 934,
		'Q': 934,
		'R': 934,
		'S': 934,
		'T': 934,
		'U': 934,
		'V': 934,
		'W': 934,
		'X': 934,
		'Y': 934,
		'Z': 934,
		'a': 934,
		'b': 934,
		'c': 934,
		'd': 934,
		'e': 934,
		'f': 934,
		'g': 934,
		'h': 934,
		'i': 934,
		'j': 934,
		'k': 934,
		'l': 934,
		'm': 934,
		'n': 934,
		'o': 934,
		'p': 934,
		'q': 934,
		'r': 934,
		's': 934,
		't': 934,
		'u': 934,
		'v': 934,
		'w': 934,
		'x': 934,
		'y': 934,
		'z': 934,
	},
	928: {
		'*': 924,
	},
}

var arithmeticAccept = map[int]int{
	923: 3,
	922: 2,
	934: 14,
	926: 6,
	931: 11,
	921: 1,
	933: 13,
	929: 9,
	935: 14,
	936: 14,
	925: 5,
	932: 12,
	930: 10,
	927: 7,
	928: 8,
	924: 4,
}

func (l *Lexer) gettoken() int {
	for l.pos < len(l.input) {
		startPos  := l.pos
		startLine := l.line
		startCol  := l.col
		state    := 919
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
				if row, ok := arithmeticTrans[919]; ok { if _, ok := row[ch]; ok { break } }
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
