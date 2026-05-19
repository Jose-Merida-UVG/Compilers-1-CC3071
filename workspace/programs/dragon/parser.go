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
		"ID": {1, 1},
		"LPAREN": {1, 5},
	},
	1: {
		"$": {2, 6},
		"PLUS": {2, 6},
		"RPAREN": {2, 6},
		"TIMES": {2, 6},
	},
	2: {
		"$": {3, 0},
		"PLUS": {1, 6},
	},
	3: {
		"$": {2, 2},
		"PLUS": {2, 2},
		"RPAREN": {2, 2},
		"TIMES": {1, 7},
	},
	4: {
		"$": {2, 4},
		"PLUS": {2, 4},
		"RPAREN": {2, 4},
		"TIMES": {2, 4},
	},
	5: {
		"ID": {1, 1},
		"LPAREN": {1, 5},
	},
	6: {
		"ID": {1, 1},
		"LPAREN": {1, 5},
	},
	7: {
		"ID": {1, 1},
		"LPAREN": {1, 5},
	},
	8: {
		"PLUS": {1, 6},
		"RPAREN": {1, 11},
	},
	9: {
		"$": {2, 1},
		"PLUS": {2, 1},
		"RPAREN": {2, 1},
		"TIMES": {1, 7},
	},
	10: {
		"$": {2, 3},
		"PLUS": {2, 3},
		"RPAREN": {2, 3},
		"TIMES": {2, 3},
	},
	11: {
		"$": {2, 5},
		"PLUS": {2, 5},
		"RPAREN": {2, 5},
		"TIMES": {2, 5},
	},
}

// parserGotoTable[state][nonTerminal] → next state
var parserGotoTable = map[int]map[string]int{
	0: {
		"e": 2,
		"f": 4,
		"t": 3,
	},
	5: {
		"e": 8,
		"f": 4,
		"t": 3,
	},
	6: {
		"f": 4,
		"t": 9,
	},
	7: {
		"f": 10,
	},
}

// parserProd describes one production: its head symbol, body symbols, and body length.
type parserProd struct{ head string; body string; bodyLen int }

var parserProds = []parserProd{
	0: {"e'", "e", 1},
	1: {"e", "e PLUS t", 3},
	2: {"e", "t", 1},
	3: {"t", "t TIMES f", 3},
	4: {"t", "f", 1},
	5: {"f", "LPAREN e RPAREN", 3},
	6: {"f", "ID", 1},
}

var parserIgnore = map[int]bool{}

// Parse runs the SLR(1) parse loop over the token stream produced by lexer l.
// It returns nil on a successful parse, or a descriptive error on failure.
// Tokens whose IDs appear in parserIgnore are silently skipped.
func Parse(l *Lexer) error {
	stk := []int{0}
	peek := func() int { return stk[len(stk)-1] }

	// symStk stores display labels: terminals as TOKEN(value), non-terminals as their name.
	var symStk []string

	var cur Lexeme
	var symbolTable []Lexeme
	var derivation []string

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

	tokName := tokenIDToName()

	sententialForm := func() string {
		if len(symStk) == 0 {
			return "ε"
		}
		return strings.Join(symStk, " ")
	}

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
			name := tokenName(cur, tokName)
			label := fmt.Sprintf("%s(%s)", name, cur.Value)
			symStk = append(symStk, label)
			fmt.Printf("Shift  %s\n  %s\n", label, sententialForm())
			stk = append(stk, act.arg)
			nextToken()

		case 2: // reduce
			prod := parserProds[act.arg]
			if prod.bodyLen > 0 {
				symStk = symStk[:len(symStk)-prod.bodyLen]
			}
			symStk = append(symStk, prod.head)
			fmt.Printf("Reduce %s → %s\n  %s\n", prod.head, prod.body, sententialForm())
			derivation = append(derivation, sententialForm())

			stk = stk[:len(stk)-prod.bodyLen]
			top := peek()

			gotoRow, ok := parserGotoTable[top]
			if !ok {
				return fmt.Errorf("state %d: no goto row (reducing by %q)", top, prod.head)
			}
			next, ok := gotoRow[prod.head]
			if !ok {
				return fmt.Errorf("state %d: no goto for %q", top, prod.head)
			}
			stk = append(stk, next)

		case 3: // accept
			fmt.Println()
			fmt.Println("── Derivation (top-down) ──")
			for i, j := 0, len(derivation)-1; i < j; i, j = i+1, j-1 {
				derivation[i], derivation[j] = derivation[j], derivation[i]
			}
			for i, form := range derivation {
				fmt.Printf("%d. %s\n", i+1, form)
			}

			fmt.Println()
			fmt.Println("── Symbol Table ──")
			fmt.Printf("%-20s %-20s %-10s\n", "LEXEME", "TOKEN", "LINE:COL")
			for _, lex := range symbolTable {
				fmt.Printf("%-20s %-20s %d:%d\n", lex.Value, tokenName(lex, tokName), lex.Line, lex.Col)
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



// tokenIDToName returns a map from token integer ID to its grammar name.
// This is the inverse of the constants emitted by GenerateCombined.
func tokenIDToName() map[int]string {
	return map[int]string{
		1: "ID",
		2: "PLUS",
		3: "TIMES",
		4: "LPAREN",
		5: "RPAREN",
	}
}


func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: dragon <inputfile>")
		os.Exit(1)
	}
	data, err := os.ReadFile(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	l := NewDragonLexer(string(data))
	if err := Parse(l); err != nil {
		fmt.Fprintln(os.Stderr, "parse error:", err)
		os.Exit(1)
	}
	fmt.Println("OK — input accepted by the grammar.")
}
