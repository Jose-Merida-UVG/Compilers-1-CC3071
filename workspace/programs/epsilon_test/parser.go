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
		"A": {1, 2},
	},
	1: {
		"$": {3, 0},
	},
	2: {
		"C": {1, 4},
	},
	3: {
		"F": {2, 7},
		"G": {1, 7},
		"H": {2, 7},
	},
	4: {
		"B": {1, 9},
		"F": {2, 4},
		"G": {2, 4},
		"H": {2, 4},
	},
	5: {
		"H": {1, 10},
	},
	6: {
		"F": {1, 12},
		"H": {2, 9},
	},
	7: {
		"F": {2, 6},
		"H": {2, 6},
	},
	8: {
		"F": {2, 2},
		"G": {2, 2},
		"H": {2, 2},
	},
	9: {
		"B": {1, 9},
		"F": {2, 4},
		"G": {2, 4},
		"H": {2, 4},
	},
	10: {
		"$": {2, 1},
	},
	11: {
		"H": {2, 5},
	},
	12: {
		"H": {2, 8},
	},
	13: {
		"F": {2, 3},
		"G": {2, 3},
		"H": {2, 3},
	},
}

// parserGotoTable[state][nonTerminal] → next state
var parserGotoTable = map[int]map[string]int{
	0: {
		"s": 1,
	},
	2: {
		"bprod": 3,
	},
	3: {
		"dprod": 5,
		"eprod": 6,
	},
	4: {
		"cprod": 8,
	},
	6: {
		"fprod": 11,
	},
	9: {
		"cprod": 13,
	},
}

// parserProd describes one production: its head symbol, body symbols, and body length.
type parserProd struct{ head string; body string; bodyLen int }

var parserProds = []parserProd{
	0: {"s'", "s", 1},
	1: {"s", "A bprod dprod H", 4},
	2: {"bprod", "C cprod", 2},
	3: {"cprod", "B cprod", 2},
	4: {"cprod", "ε", 0},
	5: {"dprod", "eprod fprod", 2},
	6: {"eprod", "G", 1},
	7: {"eprod", "ε", 0},
	8: {"fprod", "F", 1},
	9: {"fprod", "ε", 0},
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
		1: "A",
		2: "B",
		3: "C",
		4: "F",
		5: "G",
		6: "H",
	}
}


func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: epsilon_test <inputfile>")
		os.Exit(1)
	}
	data, err := os.ReadFile(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	l := NewEpsilon_testLexer(string(data))
	if err := Parse(l); err != nil {
		fmt.Fprintln(os.Stderr, "parse error:", err)
		os.Exit(1)
	}
	fmt.Println("OK — input accepted by the grammar.")
}
