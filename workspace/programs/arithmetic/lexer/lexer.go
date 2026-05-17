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
	359: {
		'\t': 361,
		'\n': 361,
		'\r': 361,
		' ': 361,
		'%': 370,
		'(': 371,
		')': 372,
		'*': 368,
		'+': 366,
		',': 373,
		'-': 367,
		'/': 369,
		'0': 363,
		'1': 363,
		'2': 363,
		'3': 363,
		'4': 363,
		'5': 363,
		'6': 363,
		'7': 363,
		'8': 363,
		'9': 363,
		'A': 374,
		'B': 374,
		'C': 374,
		'D': 374,
		'E': 374,
		'F': 374,
		'G': 374,
		'H': 374,
		'I': 374,
		'J': 374,
		'K': 374,
		'L': 374,
		'M': 374,
		'N': 374,
		'O': 374,
		'P': 374,
		'Q': 374,
		'R': 374,
		'S': 374,
		'T': 374,
		'U': 374,
		'V': 374,
		'W': 374,
		'X': 374,
		'Y': 374,
		'Z': 374,
		'a': 374,
		'b': 374,
		'c': 374,
		'd': 374,
		'e': 374,
		'f': 374,
		'g': 374,
		'h': 374,
		'i': 374,
		'j': 374,
		'k': 374,
		'l': 374,
		'm': 374,
		'n': 374,
		'o': 374,
		'p': 375,
		'q': 374,
		'r': 374,
		's': 374,
		't': 374,
		'u': 374,
		'v': 374,
		'w': 374,
		'x': 374,
		'y': 374,
		'z': 374,
	},
	374: {
		'0': 374,
		'1': 374,
		'2': 374,
		'3': 374,
		'4': 374,
		'5': 374,
		'6': 374,
		'7': 374,
		'8': 374,
		'9': 374,
		'A': 374,
		'B': 374,
		'C': 374,
		'D': 374,
		'E': 374,
		'F': 374,
		'G': 374,
		'H': 374,
		'I': 374,
		'J': 374,
		'K': 374,
		'L': 374,
		'M': 374,
		'N': 374,
		'O': 374,
		'P': 374,
		'Q': 374,
		'R': 374,
		'S': 374,
		'T': 374,
		'U': 374,
		'V': 374,
		'W': 374,
		'X': 374,
		'Y': 374,
		'Z': 374,
		'a': 374,
		'b': 374,
		'c': 374,
		'd': 374,
		'e': 374,
		'f': 374,
		'g': 374,
		'h': 374,
		'i': 374,
		'j': 374,
		'k': 374,
		'l': 374,
		'm': 374,
		'n': 374,
		'o': 374,
		'p': 374,
		'q': 374,
		'r': 374,
		's': 374,
		't': 374,
		'u': 374,
		'v': 374,
		'w': 374,
		'x': 374,
		'y': 374,
		'z': 374,
	},
	361: {
		'\t': 361,
		'\n': 361,
		'\r': 361,
		' ': 361,
	},
	363: {
		'.': 360,
		'0': 363,
		'1': 363,
		'2': 363,
		'3': 363,
		'4': 363,
		'5': 363,
		'6': 363,
		'7': 363,
		'8': 363,
		'9': 363,
	},
	360: {
		'0': 362,
		'1': 362,
		'2': 362,
		'3': 362,
		'4': 362,
		'5': 362,
		'6': 362,
		'7': 362,
		'8': 362,
		'9': 362,
	},
	362: {
		'0': 362,
		'1': 362,
		'2': 362,
		'3': 362,
		'4': 362,
		'5': 362,
		'6': 362,
		'7': 362,
		'8': 362,
		'9': 362,
	},
	368: {
		'*': 364,
	},
	375: {
		'0': 374,
		'1': 374,
		'2': 374,
		'3': 374,
		'4': 374,
		'5': 374,
		'6': 374,
		'7': 374,
		'8': 374,
		'9': 374,
		'A': 374,
		'B': 374,
		'C': 374,
		'D': 374,
		'E': 374,
		'F': 374,
		'G': 374,
		'H': 374,
		'I': 374,
		'J': 374,
		'K': 374,
		'L': 374,
		'M': 374,
		'N': 374,
		'O': 374,
		'P': 374,
		'Q': 374,
		'R': 374,
		'S': 374,
		'T': 374,
		'U': 374,
		'V': 374,
		'W': 374,
		'X': 374,
		'Y': 374,
		'Z': 374,
		'a': 374,
		'b': 374,
		'c': 374,
		'd': 374,
		'e': 374,
		'f': 374,
		'g': 374,
		'h': 374,
		'i': 374,
		'j': 374,
		'k': 374,
		'l': 374,
		'm': 374,
		'n': 374,
		'o': 376,
		'p': 374,
		'q': 374,
		'r': 374,
		's': 374,
		't': 374,
		'u': 374,
		'v': 374,
		'w': 374,
		'x': 374,
		'y': 374,
		'z': 374,
	},
	376: {
		'0': 374,
		'1': 374,
		'2': 374,
		'3': 374,
		'4': 374,
		'5': 374,
		'6': 374,
		'7': 374,
		'8': 374,
		'9': 374,
		'A': 374,
		'B': 374,
		'C': 374,
		'D': 374,
		'E': 374,
		'F': 374,
		'G': 374,
		'H': 374,
		'I': 374,
		'J': 374,
		'K': 374,
		'L': 374,
		'M': 374,
		'N': 374,
		'O': 374,
		'P': 374,
		'Q': 374,
		'R': 374,
		'S': 374,
		'T': 374,
		'U': 374,
		'V': 374,
		'W': 374,
		'X': 374,
		'Y': 374,
		'Z': 374,
		'a': 374,
		'b': 374,
		'c': 374,
		'd': 374,
		'e': 374,
		'f': 374,
		'g': 374,
		'h': 374,
		'i': 374,
		'j': 374,
		'k': 374,
		'l': 374,
		'm': 374,
		'n': 374,
		'o': 374,
		'p': 374,
		'q': 374,
		'r': 374,
		's': 374,
		't': 374,
		'u': 374,
		'v': 374,
		'w': 365,
		'x': 374,
		'y': 374,
		'z': 374,
	},
	365: {
		'0': 374,
		'1': 374,
		'2': 374,
		'3': 374,
		'4': 374,
		'5': 374,
		'6': 374,
		'7': 374,
		'8': 374,
		'9': 374,
		'A': 374,
		'B': 374,
		'C': 374,
		'D': 374,
		'E': 374,
		'F': 374,
		'G': 374,
		'H': 374,
		'I': 374,
		'J': 374,
		'K': 374,
		'L': 374,
		'M': 374,
		'N': 374,
		'O': 374,
		'P': 374,
		'Q': 374,
		'R': 374,
		'S': 374,
		'T': 374,
		'U': 374,
		'V': 374,
		'W': 374,
		'X': 374,
		'Y': 374,
		'Z': 374,
		'a': 374,
		'b': 374,
		'c': 374,
		'd': 374,
		'e': 374,
		'f': 374,
		'g': 374,
		'h': 374,
		'i': 374,
		'j': 374,
		'k': 374,
		'l': 374,
		'm': 374,
		'n': 374,
		'o': 374,
		'p': 374,
		'q': 374,
		'r': 374,
		's': 374,
		't': 374,
		'u': 374,
		'v': 374,
		'w': 374,
		'x': 374,
		'y': 374,
		'z': 374,
	},
}

// arithmeticAccept maps accepting-state IDs to their 1-based TokenID.
var arithmeticAccept = map[int]int{
	374: 14,
	361: 1,
	372: 12,
	363: 3,
	362: 2,
	368: 8,
	364: 4,
	371: 11,
	373: 13,
	369: 9,
	370: 10,
	366: 6,
	367: 7,
	375: 14,
	376: 14,
	365: 5,
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

		state    := 359  // start state of the DFA
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
				if row, ok := arithmeticTrans[359]; ok {
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
