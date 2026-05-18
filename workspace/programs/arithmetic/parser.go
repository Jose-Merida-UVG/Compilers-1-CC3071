package main

import (
	"fmt"
	"os"
)

// parserAction encodes one cell of the SLR(1) ACTION table.
// kind: 1=shift 2=reduce 3=accept; arg: target state (shift) or prod index (reduce).
type parserAction struct{ kind, arg int }

// parserActionTable[state][terminal] → action
var parserActionTable = map[int]map[string]parserAction{
	0: {
		"FLOAT": {1, 6},
		"IDENT": {1, 7},
		"INT": {1, 5},
		"KWPOW": {1, 8},
		"LPAREN": {1, 9},
	},
	1: {
		"$": {3, 0},
		"MINUS": {1, 11},
		"PLUS": {1, 10},
	},
	2: {
		"$": {2, 3},
		"COMMA": {2, 3},
		"DIV": {1, 12},
		"MINUS": {2, 3},
		"MOD": {1, 13},
		"PLUS": {2, 3},
		"RPAREN": {2, 3},
		"TIMES": {1, 14},
	},
	3: {
		"$": {2, 7},
		"COMMA": {2, 7},
		"DIV": {2, 7},
		"MINUS": {2, 7},
		"MOD": {2, 7},
		"PLUS": {2, 7},
		"RPAREN": {2, 7},
		"TIMES": {2, 7},
	},
	4: {
		"$": {2, 9},
		"COMMA": {2, 9},
		"DIV": {2, 9},
		"MINUS": {2, 9},
		"MOD": {2, 9},
		"PLUS": {2, 9},
		"POW": {1, 15},
		"RPAREN": {2, 9},
		"TIMES": {2, 9},
	},
	5: {
		"$": {2, 10},
		"COMMA": {2, 10},
		"DIV": {2, 10},
		"MINUS": {2, 10},
		"MOD": {2, 10},
		"PLUS": {2, 10},
		"POW": {2, 10},
		"RPAREN": {2, 10},
		"TIMES": {2, 10},
	},
	6: {
		"$": {2, 11},
		"COMMA": {2, 11},
		"DIV": {2, 11},
		"MINUS": {2, 11},
		"MOD": {2, 11},
		"PLUS": {2, 11},
		"POW": {2, 11},
		"RPAREN": {2, 11},
		"TIMES": {2, 11},
	},
	7: {
		"$": {2, 12},
		"COMMA": {2, 12},
		"DIV": {2, 12},
		"LPAREN": {1, 16},
		"MINUS": {2, 12},
		"MOD": {2, 12},
		"PLUS": {2, 12},
		"POW": {2, 12},
		"RPAREN": {2, 12},
		"TIMES": {2, 12},
	},
	8: {
		"LPAREN": {1, 17},
	},
	9: {
		"FLOAT": {1, 6},
		"IDENT": {1, 7},
		"INT": {1, 5},
		"KWPOW": {1, 8},
		"LPAREN": {1, 9},
	},
	10: {
		"FLOAT": {1, 6},
		"IDENT": {1, 7},
		"INT": {1, 5},
		"KWPOW": {1, 8},
		"LPAREN": {1, 9},
	},
	11: {
		"FLOAT": {1, 6},
		"IDENT": {1, 7},
		"INT": {1, 5},
		"KWPOW": {1, 8},
		"LPAREN": {1, 9},
	},
	12: {
		"FLOAT": {1, 6},
		"IDENT": {1, 7},
		"INT": {1, 5},
		"KWPOW": {1, 8},
		"LPAREN": {1, 9},
	},
	13: {
		"FLOAT": {1, 6},
		"IDENT": {1, 7},
		"INT": {1, 5},
		"KWPOW": {1, 8},
		"LPAREN": {1, 9},
	},
	14: {
		"FLOAT": {1, 6},
		"IDENT": {1, 7},
		"INT": {1, 5},
		"KWPOW": {1, 8},
		"LPAREN": {1, 9},
	},
	15: {
		"FLOAT": {1, 6},
		"IDENT": {1, 7},
		"INT": {1, 5},
		"KWPOW": {1, 8},
		"LPAREN": {1, 9},
	},
	16: {
		"FLOAT": {1, 6},
		"IDENT": {1, 7},
		"INT": {1, 5},
		"KWPOW": {1, 8},
		"LPAREN": {1, 9},
		"RPAREN": {2, 17},
	},
	17: {
		"FLOAT": {1, 6},
		"IDENT": {1, 7},
		"INT": {1, 5},
		"KWPOW": {1, 8},
		"LPAREN": {1, 9},
	},
	18: {
		"MINUS": {1, 11},
		"PLUS": {1, 10},
		"RPAREN": {1, 29},
	},
	19: {
		"$": {2, 1},
		"COMMA": {2, 1},
		"DIV": {1, 12},
		"MINUS": {2, 1},
		"MOD": {1, 13},
		"PLUS": {2, 1},
		"RPAREN": {2, 1},
		"TIMES": {1, 14},
	},
	20: {
		"$": {2, 2},
		"COMMA": {2, 2},
		"DIV": {1, 12},
		"MINUS": {2, 2},
		"MOD": {1, 13},
		"PLUS": {2, 2},
		"RPAREN": {2, 2},
		"TIMES": {1, 14},
	},
	21: {
		"$": {2, 5},
		"COMMA": {2, 5},
		"DIV": {2, 5},
		"MINUS": {2, 5},
		"MOD": {2, 5},
		"PLUS": {2, 5},
		"RPAREN": {2, 5},
		"TIMES": {2, 5},
	},
	22: {
		"$": {2, 6},
		"COMMA": {2, 6},
		"DIV": {2, 6},
		"MINUS": {2, 6},
		"MOD": {2, 6},
		"PLUS": {2, 6},
		"RPAREN": {2, 6},
		"TIMES": {2, 6},
	},
	23: {
		"$": {2, 4},
		"COMMA": {2, 4},
		"DIV": {2, 4},
		"MINUS": {2, 4},
		"MOD": {2, 4},
		"PLUS": {2, 4},
		"RPAREN": {2, 4},
		"TIMES": {2, 4},
	},
	24: {
		"$": {2, 8},
		"COMMA": {2, 8},
		"DIV": {2, 8},
		"MINUS": {2, 8},
		"MOD": {2, 8},
		"PLUS": {2, 8},
		"RPAREN": {2, 8},
		"TIMES": {2, 8},
	},
	25: {
		"COMMA": {1, 30},
		"RPAREN": {2, 16},
	},
	26: {
		"RPAREN": {1, 31},
	},
	27: {
		"COMMA": {2, 18},
		"MINUS": {1, 11},
		"PLUS": {1, 10},
		"RPAREN": {2, 18},
	},
	28: {
		"COMMA": {1, 32},
		"MINUS": {1, 11},
		"PLUS": {1, 10},
	},
	29: {
		"$": {2, 13},
		"COMMA": {2, 13},
		"DIV": {2, 13},
		"MINUS": {2, 13},
		"MOD": {2, 13},
		"PLUS": {2, 13},
		"POW": {2, 13},
		"RPAREN": {2, 13},
		"TIMES": {2, 13},
	},
	30: {
		"FLOAT": {1, 6},
		"IDENT": {1, 7},
		"INT": {1, 5},
		"KWPOW": {1, 8},
		"LPAREN": {1, 9},
	},
	31: {
		"$": {2, 15},
		"COMMA": {2, 15},
		"DIV": {2, 15},
		"MINUS": {2, 15},
		"MOD": {2, 15},
		"PLUS": {2, 15},
		"POW": {2, 15},
		"RPAREN": {2, 15},
		"TIMES": {2, 15},
	},
	32: {
		"FLOAT": {1, 6},
		"IDENT": {1, 7},
		"INT": {1, 5},
		"KWPOW": {1, 8},
		"LPAREN": {1, 9},
	},
	33: {
		"COMMA": {2, 19},
		"MINUS": {1, 11},
		"PLUS": {1, 10},
		"RPAREN": {2, 19},
	},
	34: {
		"MINUS": {1, 11},
		"PLUS": {1, 10},
		"RPAREN": {1, 35},
	},
	35: {
		"$": {2, 14},
		"COMMA": {2, 14},
		"DIV": {2, 14},
		"MINUS": {2, 14},
		"MOD": {2, 14},
		"PLUS": {2, 14},
		"POW": {2, 14},
		"RPAREN": {2, 14},
		"TIMES": {2, 14},
	},
}

// parserGotoTable[state][nonTerminal] → next state
var parserGotoTable = map[int]map[string]int{
	0: {
		"atom": 4,
		"expr": 1,
		"power": 3,
		"term": 2,
	},
	9: {
		"atom": 4,
		"expr": 18,
		"power": 3,
		"term": 2,
	},
	10: {
		"atom": 4,
		"power": 3,
		"term": 19,
	},
	11: {
		"atom": 4,
		"power": 3,
		"term": 20,
	},
	12: {
		"atom": 4,
		"power": 21,
	},
	13: {
		"atom": 4,
		"power": 22,
	},
	14: {
		"atom": 4,
		"power": 23,
	},
	15: {
		"atom": 4,
		"power": 24,
	},
	16: {
		"arglist": 25,
		"args": 26,
		"atom": 4,
		"expr": 27,
		"power": 3,
		"term": 2,
	},
	17: {
		"atom": 4,
		"expr": 28,
		"power": 3,
		"term": 2,
	},
	30: {
		"atom": 4,
		"expr": 33,
		"power": 3,
		"term": 2,
	},
	32: {
		"atom": 4,
		"expr": 34,
		"power": 3,
		"term": 2,
	},
}

// parserProd describes one production: its head symbol and body length.
type parserProd struct{ head string; bodyLen int }

var parserProds = []parserProd{
	0: {"expr'", 1},
	1: {"expr", 3},
	2: {"expr", 3},
	3: {"expr", 1},
	4: {"term", 3},
	5: {"term", 3},
	6: {"term", 3},
	7: {"term", 1},
	8: {"power", 3},
	9: {"power", 1},
	10: {"atom", 1},
	11: {"atom", 1},
	12: {"atom", 1},
	13: {"atom", 3},
	14: {"atom", 6},
	15: {"atom", 4},
	16: {"args", 1},
	17: {"args", 0},
	18: {"arglist", 1},
	19: {"arglist", 3},
}

var parserIgnore = map[int]bool{}

// Parse runs the SLR(1) parse loop over the token stream produced by lexer l.
// It returns nil on a successful parse, or a descriptive error on failure.
// Tokens whose IDs appear in parserIgnore are silently skipped.
func Parse(l *Lexer) error {
	// State stack — start in state 0.
	stk := []int{0}
	peek := func() int { return stk[len(stk)-1] }

	// Fetch the first non-ignored token.
	var cur Lexeme

	symbolTable := map[string][]Lexeme{}

	nextToken := func() {
		for {
			cur = l.NextToken()
			if cur.Token != EOF && cur.Token != ERROR {
				symbolTable[cur.Value] = append(symbolTable[cur.Value], cur)
			}
			if !parserIgnore[cur.Token] {
				break
			}
		}
	}
	nextToken()

	// Map token ID → terminal name for table look-ups.
	// The reverse map is built from the constant values embedded in the file.
	tokName := tokenIDToName()

	for {
		state := peek()
		sym := "$"
		if cur.Token != EOF {
			if name, ok := tokName[cur.Token]; ok {
				sym = name
			} else {
				return fmt.Errorf(
					"syntax error: unexpected %!s(MISSING) at line %!d(MISSING), col %!d(MISSING)",
					sym,
					cur.Line,
					cur.Col,
				)
			}
		}

		row, ok := parserActionTable[state]
		if !ok {
			return fmt.Errorf("line %!d(MISSING) col %!d(MISSING): no actions in state %!d(MISSING) (token %!q(MISSING))",
				cur.Line, cur.Col, state, sym)
		}
		act, ok := row[sym]
		if !ok {
			// Collect expected symbols for a helpful message.
			expected := make([]string, 0, len(row))
			for s := range row { expected = append(expected, s) }
			return fmt.Errorf("line %!d(MISSING) col %!d(MISSING): unexpected %!q(MISSING), expected one of %!v(MISSING)",
				cur.Line, cur.Col, sym, expected)
		}

		switch act.kind {
		case 1: // shift
			stk = append(stk, act.arg)
			nextToken()

		case 2: // reduce
			prod := parserProds[act.arg]
			// Pop |body| states off the stack.
			stk = stk[:len(stk)-prod.bodyLen]
			// Look up Goto[top][head] to find the new state.
			top := peek()
			gotoRow, ok := parserGotoTable[top]
			if !ok {
				return fmt.Errorf("state %!d(MISSING): no goto row (reducing by %!q(MISSING))", top, prod.head)
			}
			next, ok := gotoRow[prod.head]
			if !ok {
				return fmt.Errorf("state %!d(MISSING): no goto for %!q(MISSING)", top, prod.head)
			}
			stk = append(stk, next)

		case 3: // accept
			fmt.Println("\n── Tabla de símbolos ──")
			fmt.Printf("%!s(MISSING) %!s(MISSING) %!s(MISSING)\n", "LEXEMA", "TOKEN", "LÍNEA:COL")

			for lexeme, occs := range symbolTable {
				for _, occ := range occs {
					fmt.Printf(
						"%!s(MISSING) %!d(MISSING) %!d(MISSING):%!d(MISSING)\n",
						lexeme,
						occ.Token,
						occ.Line,
						occ.Col,
					)
				}
			}
			return nil
		default:
			return fmt.Errorf("line %!d(MISSING) col %!d(MISSING): parse error (state %!d(MISSING), token %!q(MISSING))",
				cur.Line, cur.Col, state, sym)
		}
	}
}

// tokenIDToName returns a map from token integer ID to its grammar name.
// This is the inverse of the constants emitted by GenerateCombined.
func tokenIDToName() map[int]string {
	return map[int]string{
		1: "FLOAT",
		2: "INT",
		3: "POW",
		4: "KWPOW",
		5: "PLUS",
		6: "MINUS",
		7: "TIMES",
		8: "DIV",
		9: "MOD",
		10: "LPAREN",
		11: "RPAREN",
		12: "COMMA",
		13: "IDENT",
	}
}


func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: arithmetic <inputfile>")
		os.Exit(1)
	}
	data, err := os.ReadFile(os.Args[1])
	if err != nil { fmt.Fprintln(os.Stderr, err); os.Exit(1) }
	l := NewArithmeticLexer(string(data))
	if err := Parse(l); err != nil {
		fmt.Fprintln(os.Stderr, "parse error:", err)
		os.Exit(1)
	}
	fmt.Println("OK — input accepted by the grammar.")
}
