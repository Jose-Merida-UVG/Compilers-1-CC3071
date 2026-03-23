package codegen

import (
	"fmt"
	"sort"
	"strings"

	"github.com/Jose-Merida-UVG/Compilers-1-CC3071/internal/automata"
	"github.com/Jose-Merida-UVG/Compilers-1-CC3071/internal/yalex"
)

// GenerateLexer returns the full source text of a Go file (package lexers)
// that implements a Lex-style streaming scanner driven by the given minimized DFA.
//
//   name    — base identifier; the constructor will be called New<Name>Lexer.
//   yf      — already-parsed YalFile, used to extract header / trailer sections.
//   dfa     — minimized, merged DFA from automata.Merge + Minimize.
//   actions — parallel action strings; actions[i] maps to TokenID i+1.
//             Actions are embedded verbatim as Go code. Actions that return
//             exit Yylex (token emitted); actions that don't return let the
//             scan loop continue (skip/whitespace).
func GenerateLexer(name string, yf *yalex.YalFile, dfa *automata.DFA, actions []string) string {
	exportedName := capitalize(name)

	header := yf.Header
	trailer := yf.Trailer

	states := dfa.GetAllStates()
	startID := dfa.StartState.ID

	var b strings.Builder
	w := func(format string, args ...any) { fmt.Fprintf(&b, format, args...) }

	w("package main\n\n")
	w("import (\n\t\"fmt\"\n\t\"os\"\n)\n\n")

	w("// Sentinel return values for Yylex.\n")
	w("const (\n")
	w("\tEOF   = 0  // end of input\n")
	w("\tERROR = -1 // unrecognised character\n")
	w(")\n\n")

	// TOKEN_N constants mirror the 1-based TokenID values stored in the DFA's
	// accept table. They are emitted for reference — the user's named constants
	// from the header (FLOAT, INT, etc.) are what actions actually return.
	w("// Token ID constants — one per pattern, in the order they appear in the spec.\n")
	w("const (\n")
	for i, action := range actions {
		a := strings.TrimSpace(action)
		if a != "" {
			w("\tTOKEN_%d = %d // %s\n", i+1, i+1, a)
		} else {
			w("\tTOKEN_%d = %d\n", i+1, i+1)
		}
	}
	w(")\n\n")

	if header != "" {
		w("// --- header ---\n%s\n\n", header)
	}

	w(`// Lexer holds the scanning state between calls to Scan.
type Lexer struct {
	input []rune
	pos   int
	line  int
	col   int
	Lxm   string // matched lexeme
	Ln    int    // 1-based line where the current token starts
	Col   int    // 1-based column where the current token starts
}

`)

	w("// New%sLexer creates a Lexer ready to scan input.\n", exportedName)
	w("func New%sLexer(input string) *Lexer {\n", exportedName)
	w("\treturn &Lexer{input: []rune(input), pos: 0, line: 1, col: 1}\n")
	w("}\n\n")

	w("// %sTrans is the DFA transition table: stateID → rune → nextStateID.\n", name)
	w("var %sTrans = map[int]map[rune]int{\n", name)
	for _, s := range states {
		if len(s.Transitions) == 0 {
			continue
		}
		w("\t%d: {\n", s.ID)
		for _, r := range sortedRunes(s.Transitions) {
			w("\t\t%s: %d,\n", runeLit(r), s.Transitions[r].ID)
		}
		w("\t},\n")
	}
	w("}\n\n")

	w("// %sAccept maps accepting-state IDs to their 1-based TokenID.\n", name)
	w("var %sAccept = map[int]int{\n", name)
	for _, s := range states {
		if s.Accept {
			w("\t%d: %d,\n", s.ID, s.TokenID)
		}
	}
	w("}\n\n")

	// NOTE: all format verbs inside this raw string are for the *generated* file's
	// fmt.Printf calls, not for this w() call. %% in the template becomes a literal
	// % in the output; %d / %s here are consumed by w() to splice in startID / name.
	w(`// Scan advances to the next token and returns its ID.
// Lxm, Ln, and Col are set before the action runs.
// Actions that return emit the token to the caller.
// Actions that don't return let the scan loop continue (skip).
// Consecutive unrecognised characters are grouped into one ERROR token.
// Returns EOF (0) at end of input, ERROR (-1) for unrecognised characters.
func (l *Lexer) Scan() int {
	for l.pos < len(l.input) {
		// Snapshot the position at the start of each token attempt.
		// If the inner loop overshoots, we backtrack here via lastPos/lastLine/lastCol.
		startPos  := l.pos
		startLine := l.line
		startCol  := l.col

		state    := %d  // start state of the DFA
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
			row, ok := %sTrans[state]
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
			if tok := %sAccept[state]; tok != 0 {
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
				if row, ok := %sTrans[%d]; ok {
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
`, startID, name, name, name, startID)

	for i, action := range actions {
		w("\t\tcase %d:\n", i+1)
		a := strings.TrimSpace(action)
		if a == "" {
			w("\t\t\t// no action\n")
		} else {
			for _, ln := range strings.Split(a, "\n") {
				w("\t\t\t%s\n", ln)
			}
		}
	}

	w("\t\t}\n\t}\n\treturn EOF\n}\n\n")

	if trailer != "" {
		w("\n// --- trailer ---\n%s\n", trailer)
	}

	w(`func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: %s <inputfile>")
		os.Exit(1)
	}
	data, err := os.ReadFile(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	l := New%sLexer(string(data))
	for {
		tok := l.Scan()
		if tok == EOF {
			break
		}
		if tok == ERROR {
			fmt.Printf("ERROR  %%-20q  ln=%%d col=%%d-%%d\n", l.Lxm, l.Ln, l.Col, l.Col+len([]rune(l.Lxm))-1)
			continue
		}
		fmt.Printf("%%d  %%-20q  ln=%%d col=%%d\n", tok, l.Lxm, l.Ln, l.Col)
	}
}
`, name, exportedName)

	return b.String()
}

// capitalize upper-cases the first rune of s.
func capitalize(s string) string {
	if s == "" {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}


// sortedRunes returns the rune keys of a DFAState transition map in ascending
// order for deterministic output.
func sortedRunes(m map[rune]*automata.DFAState) []rune {
	runes := make([]rune, 0, len(m))
	for r := range m {
		runes = append(runes, r)
	}
	sort.Slice(runes, func(i, j int) bool { return runes[i] < runes[j] })
	return runes
}

// runeLit returns a valid Go rune literal for r.
func runeLit(r rune) string {
	switch r {
	case '\n':
		return `'\n'`
	case '\t':
		return `'\t'`
	case '\r':
		return `'\r'`
	case '\\':
		return `'\\'`
	case '\'':
		return `'\''`
	}
	if r >= 32 && r <= 126 {
		return fmt.Sprintf("'%c'", r)
	}
	return fmt.Sprintf("'\\u%04X'", r)
}