/*
parserActionTable — the full SLR(1) ACTION table as a Go map literal
parserGotoTable — the GOTO table
parserProds — head name + body length for every production (what the reduce action needs to pop)
parserIgnore — token IDs to silently skip (whitespace, comments from the .yalp IGNORE list)
Parse(l *Lexer) error — the canonical shift/reduce loop with helpful error messages (line/col + expected symbols)
tokenIDToName() — reverse map from int ID → grammar name, built from the constants emitted by GenerateCombined
main() — reads a file, runs the lexer+parser, prints OK or a parse error
*/
package codegen

import (
	"fmt"
	"sort"
	"strings"

	"github.com/Jose-Merida-UVG/Compilers-1-CC3071/internal/yapar/grammar"
	"github.com/Jose-Merida-UVG/Compilers-1-CC3071/internal/yapar/lr0"
	"github.com/Jose-Merida-UVG/Compilers-1-CC3071/internal/yapar/slr"
)

type ParserSpec struct {
	Name      string
	Grammar   *grammar.Grammar
	Automaton *lr0.Automaton
	Table     *slr.Table
}

// GenerateParser returns the full source text of a Go file that
// implements a table-driven SLR(1) parser.  It is meant to be written alongside
// the file produced by GenerateCombined:
//
// GenerateParser produces a *second* file in that package that adds only the parser tables and
// the Parse() function, keeping a clean separation from the lexer machinery.
//
// The generated Parse() function:
//   - Reads tokens from the Lexer produced by GenerateCombined.
//   - Ignores token IDs listed in spec.Grammar.IgnoreList (whitespace, comments …).
//   - Runs the canonical SLR(1) shift/reduce loop.
//   - Returns nil on successful parse or an error with line/col on failure.
func GenerateParser(spec ParserSpec) string {
	g := spec.Grammar
	table := spec.Table
	prods := spec.Automaton.Productions

	var b strings.Builder
	w := func(format string, args ...any) { fmt.Fprintf(&b, format, args...) }

	w("package main\n\n")
	w("import (\n")
	w("\t\"fmt\"\n")
	w("\t\"os\"\n")
	w("\t\"strings\"\n")
	w(")\n\n")

	// ACTION TABLE
	// Encode ActionKind as int so we can embed it without importing slr/.
	// 0=error 1=shift 2=reduce 3=accept  (matches slr.ActionKind order)
	w("// parserAction encodes one cell of the SLR(1) ACTION table.\n")
	w("// kind: 1=shift 2=reduce 3=accept; arg: target state (shift) or prod index (reduce).\n")
	w("type parserAction struct{ kind, arg int }\n\n")

	w("// parserActionTable[state][terminal] → action\n")
	w("var parserActionTable = map[int]map[string]parserAction{\n")

	stateIDs := sortedIntKeys(table.Action)
	for _, sid := range stateIDs {
		row := table.Action[sid]
		syms := sortedStringKeysAction(row)
		w("\t%d: {\n", sid)
		for _, sym := range syms {
			act := row[sym]
			kind := actionKindInt(act.Kind)
			arg := 0
			switch act.Kind {
			case slr.ActionShift:
				arg = act.State
			case slr.ActionReduce:
				arg = act.ProdIdx
			}
			w("\t\t%q: {%d, %d},\n", sym, kind, arg)
		}
		w("\t},\n")
	}
	w("}\n\n")

	// GOTO TABLE
	w("// parserGotoTable[state][nonTerminal] → next state\n")
	w("var parserGotoTable = map[int]map[string]int{\n")
	gotoStateIDs := sortedIntKeys2(table.Goto)
	for _, sid := range gotoStateIDs {
		row := table.Goto[sid]
		nts := sortedStringKeysGoto(row)
		w("\t%d: {\n", sid)
		for _, nt := range nts {
			w("\t\t%q: %d,\n", nt, row[nt])
		}
		w("\t},\n")
	}
	w("}\n\n")

	// PRODUCTION TABLE
	// head name + body symbols for each production.
	w("// parserProd describes one production: its head symbol, body symbols, and body length.\n")
	w("type parserProd struct{ head string; body string; bodyLen int }\n\n")
	w("var parserProds = []parserProd{\n")
	for i, p := range prods {
		bodySyms := make([]string, len(p.Body))
		for j, sym := range p.Body {
			bodySyms[j] = sym.Name
		}
		bodyStr := "ε"
		if len(bodySyms) > 0 {
			bodyStr = strings.Join(bodySyms, " ")
		}
		w("\t%d: {%q, %q, %d},\n", i, p.Head, bodyStr, len(p.Body))
	}
	w("}\n\n")

	// IGNORE SET
	// Token IDs to skip (whitespace, comments, …).  We compare by name because
	// the generated constants (e.g. WHITESPACE = 3) are already in the same
	// package via GenerateCombined, so we can reference them directly.
	if len(g.IgnoreList) > 0 {
		w("// parserIgnore is the set of token IDs the parser silently discards.\n")
		w("var parserIgnore = map[int]bool{\n")
		for _, name := range g.IgnoreList {
			if id, ok := g.Terminals[name]; ok {
				w("\t%d: true, // %s\n", id, name)
			}
		}
		w("}\n\n")
	} else {
		w("var parserIgnore = map[int]bool{}\n\n")
	}

	// PARSE() FUNCTION
	w(`// Parse runs the SLR(1) parse loop over the token stream produced by lexer l.
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
				"lexical error: unrecognized input '%%s' at line %%d, col %%d",
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
					"syntax error: unexpected token %%d at line %%d, col %%d",
					cur.Token,
					cur.Line,
					cur.Col,
				)
			}
		}

		row, ok := parserActionTable[state]
		if !ok {
			return fmt.Errorf(
				"line %%d col %%d: no actions in state %%d (token %%q)",
				cur.Line,
				cur.Col,
				state,
				sym,
			)
		}

		act, ok := row[sym]
		if !ok {
			return fmt.Errorf(
				"syntax error: unexpected %%s at line %%d, col %%d",
				sym,
				cur.Line,
				cur.Col,
			)
		}

		switch act.kind {

		case 1: // shift
			name := tokenName(cur, tokName)
			label := fmt.Sprintf("%%s(%%s)", name, cur.Value)
			symStk = append(symStk, label)
			fmt.Printf("Shift  %%s\n  %%s\n", label, sententialForm())
			stk = append(stk, act.arg)
			nextToken()

		case 2: // reduce
			prod := parserProds[act.arg]
			bodyStr := prodBody(act.arg)
			if prod.bodyLen > 0 {
				symStk = symStk[:len(symStk)-prod.bodyLen]
			}
			symStk = append(symStk, prod.head)
			fmt.Printf("Reduce %%s → %%s\n  %%s\n", prod.head, bodyStr, sententialForm())
			derivation = append(derivation, sententialForm())

			stk = stk[:len(stk)-prod.bodyLen]
			top := peek()

			gotoRow, ok := parserGotoTable[top]
			if !ok {
				return fmt.Errorf("state %%d: no goto row (reducing by %%q)", top, prod.head)
			}
			next, ok := gotoRow[prod.head]
			if !ok {
				return fmt.Errorf("state %%d: no goto for %%q", top, prod.head)
			}
			stk = append(stk, next)

		case 3: // accept
			fmt.Println()
			fmt.Println("── Derivation (top-down) ──")
			for i, j := 0, len(derivation)-1; i < j; i, j = i+1, j-1 {
				derivation[i], derivation[j] = derivation[j], derivation[i]
			}
			for i, form := range derivation {
				fmt.Printf("%%d. %%s\n", i+1, form)
			}

			fmt.Println()
			fmt.Println("── Symbol Table ──")
			fmt.Printf("%%-20s %%-20s %%-10s\n", "LEXEME", "TOKEN", "LINE:COL")
			for _, lex := range symbolTable {
				fmt.Printf("%%-20s %%-20s %%d:%%d\n", lex.Value, tokenName(lex, tokName), lex.Line, lex.Col)
			}
			return nil

		default:
			return fmt.Errorf(
				"line %%d col %%d: parse error (state %%d, token %%q)",
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
	return fmt.Sprintf("TOKEN%%d", lex.Token)
}

func prodBody(idx int) string {
	return parserProds[idx].body
}


`)
	// TOKENIDTONAME() HELPER
	// Builds the reverse map at runtime from the constants in the same package.
	// We emit the map literal directly so no reflection is needed.
	w("// tokenIDToName returns a map from token integer ID to its grammar name.\n")
	w("// This is the inverse of the constants emitted by GenerateCombined.\n")
	w("func tokenIDToName() map[int]string {\n")
	w("\treturn map[int]string{\n")

	type termEntry struct {
		name string
		id   int
	}
	terms := make([]termEntry, 0, len(g.Terminals))
	for name, id := range g.Terminals {
		terms = append(terms, termEntry{name, id})
	}
	sort.Slice(terms, func(i, j int) bool { return terms[i].id < terms[j].id })
	for _, t := range terms {
		w("\t\t%d: %q,\n", t.id, t.name)
	}
	w("\t}\n}\n")

	return b.String()
}

// GenerateCombinedMain returns the source for a main() that runs the full
// lexer+parser pipeline.  This replaces the placeholder emitted by
// GenerateCombined so the combined file compiles as a real parser driver.
//
// It is appended *after* GenerateCombined's output (which already defines
// Lexer, NextToken, Parse, etc.) in the same package main file.
func GenerateCombinedMain(name string) string {
	exportedName := capitalize(name)

	var b strings.Builder
	w := func(format string, args ...any) {
		fmt.Fprintf(&b, format, args...)
	}

	w("func main() {\n")
	w("\tif len(os.Args) < 2 {\n")
	w("\t\tfmt.Fprintln(os.Stderr, \"usage: %s <inputfile>\")\n", name)
	w("\t\tos.Exit(1)\n")
	w("\t}\n")

	w("\tdata, err := os.ReadFile(os.Args[1])\n")
	w("\tif err != nil {\n")
	w("\t\tfmt.Fprintln(os.Stderr, err)\n")
	w("\t\tos.Exit(1)\n")
	w("\t}\n")

	w("\tl := New%sLexer(string(data))\n", exportedName)

	w("\tif err := Parse(l); err != nil {\n")
	w("\t\tfmt.Fprintln(os.Stderr, \"parse error:\", err)\n")
	w("\t\tos.Exit(1)\n")
	w("\t}\n")

	w("\tfmt.Println(\"OK — input accepted by the grammar.\")\n")
	w("}\n")

	return b.String()
}

func actionKindInt(k slr.ActionKind) int {
	switch k {
	case slr.ActionShift:
		return 1
	case slr.ActionReduce:
		return 2
	case slr.ActionAccept:
		return 3
	}
	return 0
}

func sortedIntKeys(m map[int]map[string]slr.Action) []int {
	keys := make([]int, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	return keys
}

func sortedIntKeys2(m map[int]map[string]int) []int {
	keys := make([]int, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	return keys
}

func sortedStringKeysAction(m map[string]slr.Action) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func sortedStringKeysGoto(m map[string]int) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}