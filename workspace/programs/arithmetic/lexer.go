package main

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
	1053: {
		'\t': 1055,
		'\n': 1055,
		'\r': 1055,
		' ': 1055,
		'%': 1064,
		'(': 1065,
		')': 1066,
		'*': 1062,
		'+': 1060,
		',': 1067,
		'-': 1061,
		'/': 1063,
		'0': 1057,
		'1': 1057,
		'2': 1057,
		'3': 1057,
		'4': 1057,
		'5': 1057,
		'6': 1057,
		'7': 1057,
		'8': 1057,
		'9': 1057,
		'A': 1068,
		'B': 1068,
		'C': 1068,
		'D': 1068,
		'E': 1068,
		'F': 1068,
		'G': 1068,
		'H': 1068,
		'I': 1068,
		'J': 1068,
		'K': 1068,
		'L': 1068,
		'M': 1068,
		'N': 1068,
		'O': 1068,
		'P': 1068,
		'Q': 1068,
		'R': 1068,
		'S': 1068,
		'T': 1068,
		'U': 1068,
		'V': 1068,
		'W': 1068,
		'X': 1068,
		'Y': 1068,
		'Z': 1068,
		'a': 1068,
		'b': 1068,
		'c': 1068,
		'd': 1068,
		'e': 1068,
		'f': 1068,
		'g': 1068,
		'h': 1068,
		'i': 1068,
		'j': 1068,
		'k': 1068,
		'l': 1068,
		'm': 1068,
		'n': 1068,
		'o': 1068,
		'p': 1069,
		'q': 1068,
		'r': 1068,
		's': 1068,
		't': 1068,
		'u': 1068,
		'v': 1068,
		'w': 1068,
		'x': 1068,
		'y': 1068,
		'z': 1068,
	},
	1068: {
		'0': 1068,
		'1': 1068,
		'2': 1068,
		'3': 1068,
		'4': 1068,
		'5': 1068,
		'6': 1068,
		'7': 1068,
		'8': 1068,
		'9': 1068,
		'A': 1068,
		'B': 1068,
		'C': 1068,
		'D': 1068,
		'E': 1068,
		'F': 1068,
		'G': 1068,
		'H': 1068,
		'I': 1068,
		'J': 1068,
		'K': 1068,
		'L': 1068,
		'M': 1068,
		'N': 1068,
		'O': 1068,
		'P': 1068,
		'Q': 1068,
		'R': 1068,
		'S': 1068,
		'T': 1068,
		'U': 1068,
		'V': 1068,
		'W': 1068,
		'X': 1068,
		'Y': 1068,
		'Z': 1068,
		'a': 1068,
		'b': 1068,
		'c': 1068,
		'd': 1068,
		'e': 1068,
		'f': 1068,
		'g': 1068,
		'h': 1068,
		'i': 1068,
		'j': 1068,
		'k': 1068,
		'l': 1068,
		'm': 1068,
		'n': 1068,
		'o': 1068,
		'p': 1068,
		'q': 1068,
		'r': 1068,
		's': 1068,
		't': 1068,
		'u': 1068,
		'v': 1068,
		'w': 1068,
		'x': 1068,
		'y': 1068,
		'z': 1068,
	},
	1057: {
		'.': 1054,
		'0': 1057,
		'1': 1057,
		'2': 1057,
		'3': 1057,
		'4': 1057,
		'5': 1057,
		'6': 1057,
		'7': 1057,
		'8': 1057,
		'9': 1057,
	},
	1054: {
		'0': 1056,
		'1': 1056,
		'2': 1056,
		'3': 1056,
		'4': 1056,
		'5': 1056,
		'6': 1056,
		'7': 1056,
		'8': 1056,
		'9': 1056,
	},
	1056: {
		'0': 1056,
		'1': 1056,
		'2': 1056,
		'3': 1056,
		'4': 1056,
		'5': 1056,
		'6': 1056,
		'7': 1056,
		'8': 1056,
		'9': 1056,
	},
	1055: {
		'\t': 1055,
		'\n': 1055,
		'\r': 1055,
		' ': 1055,
	},
	1062: {
		'*': 1058,
	},
	1069: {
		'0': 1068,
		'1': 1068,
		'2': 1068,
		'3': 1068,
		'4': 1068,
		'5': 1068,
		'6': 1068,
		'7': 1068,
		'8': 1068,
		'9': 1068,
		'A': 1068,
		'B': 1068,
		'C': 1068,
		'D': 1068,
		'E': 1068,
		'F': 1068,
		'G': 1068,
		'H': 1068,
		'I': 1068,
		'J': 1068,
		'K': 1068,
		'L': 1068,
		'M': 1068,
		'N': 1068,
		'O': 1068,
		'P': 1068,
		'Q': 1068,
		'R': 1068,
		'S': 1068,
		'T': 1068,
		'U': 1068,
		'V': 1068,
		'W': 1068,
		'X': 1068,
		'Y': 1068,
		'Z': 1068,
		'a': 1068,
		'b': 1068,
		'c': 1068,
		'd': 1068,
		'e': 1068,
		'f': 1068,
		'g': 1068,
		'h': 1068,
		'i': 1068,
		'j': 1068,
		'k': 1068,
		'l': 1068,
		'm': 1068,
		'n': 1068,
		'o': 1070,
		'p': 1068,
		'q': 1068,
		'r': 1068,
		's': 1068,
		't': 1068,
		'u': 1068,
		'v': 1068,
		'w': 1068,
		'x': 1068,
		'y': 1068,
		'z': 1068,
	},
	1070: {
		'0': 1068,
		'1': 1068,
		'2': 1068,
		'3': 1068,
		'4': 1068,
		'5': 1068,
		'6': 1068,
		'7': 1068,
		'8': 1068,
		'9': 1068,
		'A': 1068,
		'B': 1068,
		'C': 1068,
		'D': 1068,
		'E': 1068,
		'F': 1068,
		'G': 1068,
		'H': 1068,
		'I': 1068,
		'J': 1068,
		'K': 1068,
		'L': 1068,
		'M': 1068,
		'N': 1068,
		'O': 1068,
		'P': 1068,
		'Q': 1068,
		'R': 1068,
		'S': 1068,
		'T': 1068,
		'U': 1068,
		'V': 1068,
		'W': 1068,
		'X': 1068,
		'Y': 1068,
		'Z': 1068,
		'a': 1068,
		'b': 1068,
		'c': 1068,
		'd': 1068,
		'e': 1068,
		'f': 1068,
		'g': 1068,
		'h': 1068,
		'i': 1068,
		'j': 1068,
		'k': 1068,
		'l': 1068,
		'm': 1068,
		'n': 1068,
		'o': 1068,
		'p': 1068,
		'q': 1068,
		'r': 1068,
		's': 1068,
		't': 1068,
		'u': 1068,
		'v': 1068,
		'w': 1059,
		'x': 1068,
		'y': 1068,
		'z': 1068,
	},
	1059: {
		'0': 1068,
		'1': 1068,
		'2': 1068,
		'3': 1068,
		'4': 1068,
		'5': 1068,
		'6': 1068,
		'7': 1068,
		'8': 1068,
		'9': 1068,
		'A': 1068,
		'B': 1068,
		'C': 1068,
		'D': 1068,
		'E': 1068,
		'F': 1068,
		'G': 1068,
		'H': 1068,
		'I': 1068,
		'J': 1068,
		'K': 1068,
		'L': 1068,
		'M': 1068,
		'N': 1068,
		'O': 1068,
		'P': 1068,
		'Q': 1068,
		'R': 1068,
		'S': 1068,
		'T': 1068,
		'U': 1068,
		'V': 1068,
		'W': 1068,
		'X': 1068,
		'Y': 1068,
		'Z': 1068,
		'a': 1068,
		'b': 1068,
		'c': 1068,
		'd': 1068,
		'e': 1068,
		'f': 1068,
		'g': 1068,
		'h': 1068,
		'i': 1068,
		'j': 1068,
		'k': 1068,
		'l': 1068,
		'm': 1068,
		'n': 1068,
		'o': 1068,
		'p': 1068,
		'q': 1068,
		'r': 1068,
		's': 1068,
		't': 1068,
		'u': 1068,
		'v': 1068,
		'w': 1068,
		'x': 1068,
		'y': 1068,
		'z': 1068,
	},
}

var arithmeticAccept = map[int]int{
	1068: 14,
	1057: 3,
	1056: 2,
	1067: 13,
	1055: 1,
	1062: 8,
	1058: 4,
	1069: 14,
	1070: 14,
	1059: 5,
	1066: 12,
	1060: 6,
	1064: 10,
	1061: 7,
	1065: 11,
	1063: 9,
}

func (l *Lexer) gettoken() int {
	for l.pos < len(l.input) {
		startPos  := l.pos
		startLine := l.line
		startCol  := l.col
		state    := 1053
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
				if row, ok := arithmeticTrans[1053]; ok { if _, ok := row[ch]; ok { break } }
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
			return FLOAT
		case 3:
			return INT
		case 4:
			return POW
		case 5:
			return KWPOW
		case 6:
			return PLUS
		case 7:
			return MINUS
		case 8:
			return TIMES
		case 9:
			return DIV
		case 10:
			return MOD
		case 11:
			return LPAREN
		case 12:
			return RPAREN
		case 13:
			return COMMA
		case 14:
			return IDENT
		}
	}
	return EOF
}

func (l *Lexer) NextToken() Lexeme {
	tok := l.gettoken()
	return Lexeme{Token: tok, Value: l.Lxm, Line: l.Ln, Col: l.Col}
}
