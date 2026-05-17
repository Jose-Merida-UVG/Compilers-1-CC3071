package main

import (
	"fmt"
	"os"
)

// Sentinel return values for Yylex.
const (
	EOF   = 0  // end of input
	ERROR = -1 // unrecognised character
)

// Token ID constants — one per pattern, in the order they appear in the spec.
const (
	TOKEN_1 = 1
	TOKEN_2 = 2 // fmt.Print("FLOAT  "); return FLOAT
	TOKEN_3 = 3 // fmt.Print("INT    "); return INT
	TOKEN_4 = 4 // fmt.Print("POW    "); return POW
	TOKEN_5 = 5 // fmt.Print("KWPOW  "); return KWPOW
	TOKEN_6 = 6 // fmt.Print("PLUS   "); return PLUS
	TOKEN_7 = 7 // fmt.Print("MINUS  "); return MINUS
	TOKEN_8 = 8 // fmt.Print("TIMES  "); return TIMES
	TOKEN_9 = 9 // fmt.Print("DIV    "); return DIV
	TOKEN_10 = 10 // fmt.Print("MOD    "); return MOD
	TOKEN_11 = 11 // fmt.Print("LPAREN "); return LPAREN
	TOKEN_12 = 12 // fmt.Print("RPAREN "); return RPAREN
	TOKEN_13 = 13 // fmt.Print("COMMA  "); return COMMA
	TOKEN_14 = 14 // fmt.Print("IDENT  "); return IDENT
)

// Lexeme is the unit the parser consumes: a token ID plus the matched text and position.
type Lexeme struct {
	Token int
	Value string
	Line  int
	Col   int
}

// Lexer holds the scanning state between calls to Scan.
type Lexer struct {
	input []rune
	pos   int
	line  int
	col   int
	Lxm   string // matched lexeme
	Ln    int    // 1-based line where the current token starts
	Col   int    // 1-based column where the current token starts
}

// NewArithmeticLexer creates a Lexer ready to scan input.
func NewArithmeticLexer(input string) *Lexer {
	return &Lexer{input: []rune(input), pos: 0, line: 1, col: 1}
}

// arithmeticTrans is the DFA transition table: stateID → rune → nextStateID.
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
	94: {
		'*': 90,
	},
	87: {
		'\t': 87,
		'\n': 87,
		'\r': 87,
		' ': 87,
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
}

// arithmeticAccept maps accepting-state IDs to their 1-based TokenID.
var arithmeticAccept = map[int]int{
	100: 14,
	89: 3,
	88: 2,
	94: 8,
	90: 4,
	92: 6,
	87: 1,
	98: 12,
	99: 13,
	95: 9,
	96: 10,
	97: 11,
	101: 14,
	102: 14,
	91: 5,
	93: 7,
}

// gettoken advances to the next token and returns its ID.
// Lxm, Ln, and Col are set before the action runs.
// Actions that return emit the token to the caller.
// Actions that don't return let the scan loop continue (skip).
// Consecutive unrecognised characters are grouped into one ERROR token.
// Returns EOF (0) at end of input, ERROR (-1) for unrecognised characters.
func (l *Lexer) gettoken() int {
	for l.pos < len(l.input) {
		// Snapshot the position at the start of each token attempt.
		// If the inner loop overshoots, we backtrack here via lastPos/lastLine/lastCol.
		startPos  := l.pos
		startLine := l.line
		startCol  := l.col

		state    := 85  // start state of the DFA
		lastTok  := 0   // TokenID of the last accepting state seen (0 = none yet)
		lastPos  := l.pos
		lastLine := l.line
		lastCol  := l.col
		curLine  := l.line
		curCol   := l.col

		// Inner loop: maximal munch — keep advancing as long as the DFA has a
		// transition. Every time we land on an accepting state, snapshot it.
		// When no transition exists we stop and backtrack to the last snapshot.
		for l.pos < len(l.input) {
			ch := l.input[l.pos]
			row, ok := arithmeticTrans[state]
			if !ok {
				break
			}
			next, ok := row[ch]
			if !ok {
				break
			}
			l.pos++
			if ch == '\n' {
				curLine++
				curCol = 1
			} else {
				curCol++
			}
			state = next
			if tok := arithmeticAccept[state]; tok != 0 {
				// Accepting state — snapshot position and token ID.
				// A later, longer match may overwrite this snapshot.
				lastTok  = tok
				lastPos  = l.pos
				lastLine = curLine
				lastCol  = curCol
			}
		}

		if lastTok == 0 {
			// No accepting state was ever reached from startPos — unrecognised input.
			// Consume characters one-by-one until the DFA's start state can fire again,
			// grouping them all into a single ERROR token rather than one per character.
			errStart := startPos
			errLn    := l.line
			errCol   := l.col
			for l.pos < len(l.input) {
				ch := l.input[l.pos]
				// Stop as soon as the start state has a transition for this char —
				// the next call to Scan will pick up from here cleanly.
				if row, ok := arithmeticTrans[85]; ok {
					if _, ok := row[ch]; ok {
						break
					}
				}
				if ch == '\n' {
					l.line++
					l.col = 1
				} else {
					l.col++
				}
				l.pos++
			}
			l.Lxm = string(l.input[errStart:l.pos])
			l.Ln  = errLn
			l.Col = errCol
			return ERROR
		}

		// Backtrack to the last accepting snapshot and commit that match.
		l.pos = lastPos
		l.line = lastLine
		l.col  = lastCol
		l.Lxm  = string(l.input[startPos:lastPos])
		l.Ln   = startLine
		l.Col  = startCol

		// Dispatch to the verbatim action for the winning token.
		// Actions that execute a return statement exit Scan() with that token ID.
		// Actions without a return fall through and the outer loop retries (skip).
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

// NextToken advances to the next token and returns a Lexeme.
func (l *Lexer) NextToken() Lexeme {
	tok := l.gettoken()
	return Lexeme{Token: tok, Value: l.Lxm, Line: l.Ln, Col: l.Col}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: arithmetic <inputfile>")
		os.Exit(1)
	}
	data, err := os.ReadFile(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	l := NewArithmeticLexer(string(data))
	for {
		tok := l.gettoken()
		if tok == EOF {
			break
		}
		if tok == ERROR {
			fmt.Printf("ERROR  %-20q  ln=%d col=%d-%d\n", l.Lxm, l.Ln, l.Col, l.Col+len([]rune(l.Lxm))-1)
			continue
		}
		fmt.Printf("%d  %-20q  ln=%d col=%d\n", tok, l.Lxm, l.Ln, l.Col)
	}
}
