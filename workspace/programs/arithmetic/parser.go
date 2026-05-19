package main

import (
	"fmt"
	"os"
	"strings"
)

// parserAction encodes one cell of the SLR(1) ACTION table.
// kind: 1=shift 2=reduce 3=accept; arg: target state (shift) or prod index (reduce).
type parserAction struct{ kind, arg int }

// parserActionTable[state][terminal] → action
var parserActionTable = map[int]map[string]parserAction{
	0: {
		"FLOAT": {1, 8},
		"IDENT": {1, 5},
		"INT": {1, 4},
		"KWPOW": {1, 1},
		"LPAREN": {1, 9},
	},
	1: {
		"LPAREN": {1, 10},
	},
	2: {
		"$": {2, 7},
		"COMMA": {2, 7},
		"DIV": {2, 7},
		"MINUS": {2, 7},
		"MOD": {2, 7},
		"PLUS": {2, 7},
		"RPAREN": {2, 7},
		"TIMES": {2, 7},
	},
	3: {
		"$": {2, 9},
		"COMMA": {2, 9},
		"DIV": {2, 9},
		"MINUS": {2, 9},
		"MOD": {2, 9},
		"PLUS": {2, 9},
		"POW": {1, 11},
		"RPAREN": {2, 9},
		"TIMES": {2, 9},
	},
	4: {
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
	5: {
		"$": {2, 12},
		"COMMA": {2, 12},
		"DIV": {2, 12},
		"LPAREN": {1, 12},
		"MINUS": {2, 12},
		"MOD": {2, 12},
		"PLUS": {2, 12},
		"POW": {2, 12},
		"RPAREN": {2, 12},
		"TIMES": {2, 12},
	},
	6: {
		"$": {3, 0},
		"MINUS": {1, 13},
		"PLUS": {1, 14},
	},
	7: {
		"$": {2, 3},
		"COMMA": {2, 3},
		"DIV": {1, 16},
		"MINUS": {2, 3},
		"MOD": {1, 17},
		"PLUS": {2, 3},
		"RPAREN": {2, 3},
		"TIMES": {1, 15},
	},
	8: {
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
	9: {
		"FLOAT": {1, 8},
		"IDENT": {1, 5},
		"INT": {1, 4},
		"KWPOW": {1, 1},
		"LPAREN": {1, 9},
	},
	10: {
		"FLOAT": {1, 8},
		"IDENT": {1, 5},
		"INT": {1, 4},
		"KWPOW": {1, 1},
		"LPAREN": {1, 9},
	},
	11: {
		"FLOAT": {1, 8},
		"IDENT": {1, 5},
		"INT": {1, 4},
		"KWPOW": {1, 1},
		"LPAREN": {1, 9},
	},
	12: {
		"FLOAT": {1, 8},
		"IDENT": {1, 5},
		"INT": {1, 4},
		"KWPOW": {1, 1},
		"LPAREN": {1, 9},
		"RPAREN": {2, 17},
	},
	13: {
		"FLOAT": {1, 8},
		"IDENT": {1, 5},
		"INT": {1, 4},
		"KWPOW": {1, 1},
		"LPAREN": {1, 9},
	},
	14: {
		"FLOAT": {1, 8},
		"IDENT": {1, 5},
		"INT": {1, 4},
		"KWPOW": {1, 1},
		"LPAREN": {1, 9},
	},
	15: {
		"FLOAT": {1, 8},
		"IDENT": {1, 5},
		"INT": {1, 4},
		"KWPOW": {1, 1},
		"LPAREN": {1, 9},
	},
	16: {
		"FLOAT": {1, 8},
		"IDENT": {1, 5},
		"INT": {1, 4},
		"KWPOW": {1, 1},
		"LPAREN": {1, 9},
	},
	17: {
		"FLOAT": {1, 8},
		"IDENT": {1, 5},
		"INT": {1, 4},
		"KWPOW": {1, 1},
		"LPAREN": {1, 9},
	},
	18: {
		"MINUS": {1, 13},
		"PLUS": {1, 14},
		"RPAREN": {1, 29},
	},
	19: {
		"COMMA": {1, 30},
		"MINUS": {1, 13},
		"PLUS": {1, 14},
	},
	20: {
		"$": {2, 8},
		"COMMA": {2, 8},
		"DIV": {2, 8},
		"MINUS": {2, 8},
		"MOD": {2, 8},
		"PLUS": {2, 8},
		"RPAREN": {2, 8},
		"TIMES": {2, 8},
	},
	21: {
		"COMMA": {1, 31},
		"RPAREN": {2, 16},
	},
	22: {
		"COMMA": {2, 18},
		"MINUS": {1, 13},
		"PLUS": {1, 14},
		"RPAREN": {2, 18},
	},
	23: {
		"RPAREN": {1, 32},
	},
	24: {
		"$": {2, 2},
		"COMMA": {2, 2},
		"DIV": {1, 16},
		"MINUS": {2, 2},
		"MOD": {1, 17},
		"PLUS": {2, 2},
		"RPAREN": {2, 2},
		"TIMES": {1, 15},
	},
	25: {
		"$": {2, 1},
		"COMMA": {2, 1},
		"DIV": {1, 16},
		"MINUS": {2, 1},
		"MOD": {1, 17},
		"PLUS": {2, 1},
		"RPAREN": {2, 1},
		"TIMES": {1, 15},
	},
	26: {
		"$": {2, 4},
		"COMMA": {2, 4},
		"DIV": {2, 4},
		"MINUS": {2, 4},
		"MOD": {2, 4},
		"PLUS": {2, 4},
		"RPAREN": {2, 4},
		"TIMES": {2, 4},
	},
	27: {
		"$": {2, 5},
		"COMMA": {2, 5},
		"DIV": {2, 5},
		"MINUS": {2, 5},
		"MOD": {2, 5},
		"PLUS": {2, 5},
		"RPAREN": {2, 5},
		"TIMES": {2, 5},
	},
	28: {
		"$": {2, 6},
		"COMMA": {2, 6},
		"DIV": {2, 6},
		"MINUS": {2, 6},
		"MOD": {2, 6},
		"PLUS": {2, 6},
		"RPAREN": {2, 6},
		"TIMES": {2, 6},
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
		"FLOAT": {1, 8},
		"IDENT": {1, 5},
		"INT": {1, 4},
		"KWPOW": {1, 1},
		"LPAREN": {1, 9},
	},
	31: {
		"FLOAT": {1, 8},
		"IDENT": {1, 5},
		"INT": {1, 4},
		"KWPOW": {1, 1},
		"LPAREN": {1, 9},
	},
	32: {
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
	33: {
		"MINUS": {1, 13},
		"PLUS": {1, 14},
		"RPAREN": {1, 35},
	},
	34: {
		"COMMA": {2, 19},
		"MINUS": {1, 13},
		"PLUS": {1, 14},
		"RPAREN": {2, 19},
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
		"atom": 3,
		"expr": 6,
		"power": 2,
		"term": 7,
	},
	9: {
		"atom": 3,
		"expr": 18,
		"power": 2,
		"term": 7,
	},
	10: {
		"atom": 3,
		"expr": 19,
		"power": 2,
		"term": 7,
	},
	11: {
		"atom": 3,
		"power": 20,
	},
	12: {
		"arglist": 21,
		"args": 23,
		"atom": 3,
		"expr": 22,
		"power": 2,
		"term": 7,
	},
	13: {
		"atom": 3,
		"power": 2,
		"term": 24,
	},
	14: {
		"atom": 3,
		"power": 2,
		"term": 25,
	},
	15: {
		"atom": 3,
		"power": 26,
	},
	16: {
		"atom": 3,
		"power": 27,
	},
	17: {
		"atom": 3,
		"power": 28,
	},
	30: {
		"atom": 3,
		"expr": 33,
		"power": 2,
		"term": 7,
	},
	31: {
		"atom": 3,
		"expr": 34,
		"power": 2,
		"term": 7,
	},
}

// parserProd describes one production: its head symbol, body symbols, and body length.
type parserProd struct{ head string; body string; bodyLen int }

var parserProds = []parserProd{
	0: {"expr'", "expr", 1},
	1: {"expr", "expr PLUS term", 3},
	2: {"expr", "expr MINUS term", 3},
	3: {"expr", "term", 1},
	4: {"term", "term TIMES power", 3},
	5: {"term", "term DIV power", 3},
	6: {"term", "term MOD power", 3},
	7: {"term", "power", 1},
	8: {"power", "atom POW power", 3},
	9: {"power", "atom", 1},
	10: {"atom", "INT", 1},
	11: {"atom", "FLOAT", 1},
	12: {"atom", "IDENT", 1},
	13: {"atom", "LPAREN expr RPAREN", 3},
	14: {"atom", "KWPOW LPAREN expr COMMA expr RPAREN", 6},
	15: {"atom", "IDENT LPAREN args RPAREN", 4},
	16: {"args", "arglist", 1},
	17: {"args", "ε", 0},
	18: {"arglist", "expr", 1},
	19: {"arglist", "arglist COMMA expr", 3},
}

var parserIgnore = map[int]bool{}

// Parse runs the SLR(1) parse loop over the token stream produced by lexer l.
// It returns nil on a successful parse, or a descriptive error on failure.
// Tokens whose IDs appear in parserIgnore are silently skipped.
func Parse(l *Lexer) error {
	// State stack — start in state 0.
	stk := []int{0}
	peek := func() int { return stk[len(stk)-1] }

	// Symbol stack — tracks the sentential form for derivation display.
	var symStk []string

	// Fetch the first non-ignored token.
	var cur Lexeme

	var symbolTable []Lexeme
	var sententialForms []string

	nextToken := func() {
		for {
			cur = l.NextToken()

			if cur.Token != EOF && cur.Token != ERROR {
				symbolTable = append(symbolTable, cur)
			}

			if !parserIgnore[cur.Token] {
				break
			}
		}
	}
	nextToken()

	// Map token ID → terminal name for table look-ups.
	tokName := tokenIDToName()

	fmt.Println("── Parse Actions ──")

	for {
		state := peek()

		// If we hit an ERROR token from the lexer, immediately fail.
		if cur.Token == ERROR {
			return fmt.Errorf(
				"lexical error: unrecognized input '%s' at line %d, col %d",
				cur.Value,
				cur.Line,
				cur.Col,
			)
		}

		sym := "$"
		if cur.Token != EOF {
			if name, ok := tokName[cur.Token]; ok {
				sym = name
			} else {
				return fmt.Errorf(
					"syntax error: unexpected token %d at line %d, col %d",
					cur.Token,
					cur.Line,
					cur.Col,
				)
			}
		}

		row, ok := parserActionTable[state]
		if !ok {
			return fmt.Errorf(
				"line %d col %d: no actions in state %d (token %q)",
				cur.Line,
				cur.Col,
				state,
				sym,
			)
		}

		act, ok := row[sym]
		if !ok {
			return fmt.Errorf(
				"syntax error: unexpected %s at line %d, col %d",
				sym,
				cur.Line,
				cur.Col,
			)
		}

		switch act.kind {

		case 1: // shift
			symName := tokenName(cur, tokName)
			fmt.Printf("Shift:  %s\n", symName)
			symStk = append(symStk, symName)
			stk = append(stk, act.arg)
			nextToken()

		case 2: // reduce
			prod := parserProds[act.arg]
			bodyStr := prodBody(act.arg)
			fmt.Printf("Reduce: %s → %s\n", prod.head, bodyStr)

			// Pop body symbols and push head.
			if prod.bodyLen > 0 {
				symStk = symStk[:len(symStk)-prod.bodyLen]
			}
			symStk = append(symStk, prod.head)

			// Record the current sentential form.
			sententialForms = append(sententialForms, joinSymbols(symStk))

			// Pop |body| states off the stack.
			stk = stk[:len(stk)-prod.bodyLen]

			// Look up Goto[top][head] to find the new state.
			top := peek()

			gotoRow, ok := parserGotoTable[top]
			if !ok {
				return fmt.Errorf(
					"state %d: no goto row (reducing by %q)",
					top,
					prod.head,
				)
			}

			next, ok := gotoRow[prod.head]
			if !ok {
				return fmt.Errorf(
					"state %d: no goto for %q",
					top,
					prod.head,
				)
			}

			stk = append(stk, next)

		case 3: // accept
			fmt.Println("\n\n── Sentential Forms (Rightmost Derivation in Reverse) ──")
			for i, form := range sententialForms {
				fmt.Printf("%d. %s\n", i+1, form)
			}

			fmt.Println("\n── Symbol Table ──")
			fmt.Printf("%-20s %-20s %-10s\n", "LEXEME", "TOKEN", "LINE:COL")

			for _, lex := range symbolTable {
				tokenStr := tokenName(lex, tokName)
				fmt.Printf(
					"%-20s %-20s %d:%d\n",
					lex.Value,
					tokenStr,
					lex.Line,
					lex.Col,
				)
			}

			return nil

		default:
			return fmt.Errorf(
				"line %d col %d: parse error (state %d, token %q)",
				cur.Line,
				cur.Col,
				state,
				sym,
			)
		}
	}
}

func tokenName(lex Lexeme, tokName map[int]string) string {
	if lex.Token == EOF {
		return "$"
	}
	if name, ok := tokName[lex.Token]; ok {
		return name
	}
	return fmt.Sprintf("TOKEN%d", lex.Token)
}

func prodBody(idx int) string {
	return parserProds[idx].body
}

func joinSymbols(syms []string) string {
	if len(syms) == 0 {
		return "ε"
	}
	return strings.Join(syms, " ")
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
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	l := NewArithmeticLexer(string(data))
	if err := Parse(l); err != nil {
		fmt.Fprintln(os.Stderr, "parse error:", err)
		os.Exit(1)
	}
	fmt.Println("OK — input accepted by the grammar.")
}
