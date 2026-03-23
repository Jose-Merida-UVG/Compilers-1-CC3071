package codegen

import (
	"fmt"
	"sort"
	"strings"

	"github.com/Jose-Merida-UVG/Compilers-1-CC3071/internal/automata"
	"github.com/Jose-Merida-UVG/Compilers-1-CC3071/internal/yalex"
)

// GenerateLexer returns the full source text of a Go file (package lexers)
// that implements a maximal-munch scanner driven by the given minimized DFA.
//
//   name    — base identifier; the exported function will be called <Name>Lex.
//   yf      — already-parsed YalFile, used to extract header / trailer sections.
//   dfa     — minimized, merged DFA from automata.Merge + Minimize.
//   actions — parallel action strings; actions[i] maps to TokenID i+1.
func GenerateLexer(name string, yf *yalex.YalFile, dfa *automata.DFA, actions []string) string {
	exportedName := capitalize(name)

	header := yf.Header
	trailer := yf.Trailer

	// Do any actions use int(lxm)? Only import strconv when needed.
	needsStrconv := false
	for _, a := range actions {
		if strings.Contains(a, "int(lxm)") {
			needsStrconv = true
			break
		}
	}

	states := dfa.GetAllStates()
	startID := dfa.StartState.ID

	var b strings.Builder
	w := func(format string, args ...any) { fmt.Fprintf(&b, format, args...) }

	// ── file preamble ──────────────────────────────────────────────────────
	w("package lexers\n\n")

	if needsStrconv {
		w("import \"strconv\"\n\n")
	}

	// ── user header (copied verbatim from the .yal file) ──────────────────
	if header != "" {
		w("// --- header from .yal ---\n%s\n\n", header)
	}

	// ── Token struct ───────────────────────────────────────────────────────
	w(`// Token is one lexical unit produced by the generated scanner.
type Token struct {
	Kind   int    // 1-based TokenID; -1 = unrecognised character; -2 = EOF/raise
	Lexeme string // matched text  (yytext / lxm)
	Value  any    // semantic value set by the action (yyval)
	Line   int    // 1-based line where the token starts
	Col    int    // 1-based column where the token starts
}

`)

	// ── DFA transition table ───────────────────────────────────────────────
	// map[stateID]map[rune]nextStateID — generated as a Go var so the compiler
	// bakes it into the binary; no runtime DFA object needed.
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

	// ── accept table ──────────────────────────────────────────────────────
	// map[stateID]tokenID (0 = non-accepting, but 0 is never stored here).
	w("// %sAccept maps accepting-state IDs to their 1-based TokenID.\n", name)
	w("var %sAccept = map[int]int{\n", name)
	for _, s := range states {
		if s.Accept {
			w("\t%d: %d,\n", s.ID, s.TokenID)
		}
	}
	w("}\n\n")

	// ── strconv helper (only when int(lxm) actions are present) ──────────
	if needsStrconv {
		w(`// yalexInt parses a decimal integer string (used by int(lxm) actions).
func yalexInt(s string) int { n, _ := strconv.Atoi(s); return n }

`)
	}

	// ── main Lex function ─────────────────────────────────────────────────
	w(`// %sLex tokenises input with greedy longest-match (maximal munch).
// An error Token (Kind=-1) is emitted for every unrecognised character so
// scanning always continues rather than panicking.
func %sLex(input string) []Token {
	runes := []rune(input)
	var tokens []Token
	pos := 0
	line := 1
	col := 1

	for pos < len(runes) {
		startPos := pos
		startLine := line
		startCol := col

		state := %d  // DFA start state
		lastTok := 0 // TokenID at the most recent accepting state (0 = none)
		lastPos := pos
		lastLine := line
		lastCol := col
		curLine := line
		curCol := col

		// Speculative forward scan — keep going as long as the DFA has a
		// transition, recording the furthest accepting position seen.
		for pos < len(runes) {
			ch := runes[pos]
			row, ok := %sTrans[state]
			if !ok {
				break
			}
			next, ok := row[ch]
			if !ok {
				break
			}
			pos++
			if ch == '\n' {
				curLine++
				curCol = 1
			} else {
				curCol++
			}
			state = next
			if tok := %sAccept[state]; tok != 0 {
				lastTok = tok
				lastPos = pos
				lastLine = curLine
				lastCol = curCol
			}
		}

		if lastTok == 0 {
			// No rule matched the character at startPos — emit an error token
			// for that single rune and advance past it.
			tokens = append(tokens, Token{
				Kind:   -1,
				Lexeme: string(runes[startPos : startPos+1]),
				Line:   startLine,
				Col:    startCol,
			})
			if runes[startPos] == '\n' {
				line++
				col = 1
			} else {
				col++
			}
			pos = startPos + 1
			continue
		}

		// Commit to the longest match: restore position and line/col to the
		// snapshot taken at the last accepting state.
		pos = lastPos
		line = lastLine
		col = lastCol
		lxm := string(runes[startPos:lastPos])
		_ = lxm // may go unused when all matching actions do "return lexbuf"

		switch lastTok {
`, exportedName, exportedName, startID, name, name)

	// ── one case per action ───────────────────────────────────────────────
	for i, action := range actions {
		tokenID := i + 1
		w("\t\tcase %d:\n", tokenID)
		for _, ln := range strings.Split(genAction(action, tokenID), "\n") {
			w("\t\t\t%s\n", ln)
		}
	}

	w(`		}
	}
	return tokens
}
`)

	// ── user trailer (copied verbatim from the .yal file) ─────────────────
	if trailer != "" {
		w("\n// --- trailer from .yal ---\n%s\n", trailer)
	}

	return b.String()
}

// ─── internal helpers ─────────────────────────────────────────────────────────

// capitalize upper-cases the first rune of s.
func capitalize(s string) string {
	if s == "" {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

// genAction translates one raw YALex action string into Go statement(s) that
// run inside the scanner switch. The following names are in scope at the call
// site: lxm (string), startLine, startCol (int), tokens ([]Token).
func genAction(action string, tokenID int) string {
	a := strings.TrimSpace(action)

	switch {

	// ── return lexbuf ─────────────────────────────────────────────────────
	// The lexeme is consumed but no token is emitted (whitespace, comments…).
	case a == "return lexbuf":
		return "// skip — consume lexeme without emitting a token"

	// ── raise(...) ────────────────────────────────────────────────────────
	// Signals EOF or a fatal condition; we emit a sentinel token and stop.
	case strings.HasPrefix(a, "raise("):
		return fmt.Sprintf(
			"tokens = append(tokens, Token{Kind: -2, Lexeme: lxm, Line: startLine, Col: startCol})\n"+
				"return tokens // %s", a)

	// ── return int(lxm) ──────────────────────────────────────────────────
	// Numeric literal: parse the matched text as a decimal integer.
	case a == "return int(lxm)":
		return fmt.Sprintf(
			`tokens = append(tokens, Token{Kind: %d, Lexeme: lxm, Value: yalexInt(lxm), Line: startLine, Col: startCol})`,
			tokenID)

	// ── return <expr> ─────────────────────────────────────────────────────
	// General case: expr is embedded verbatim as the Value field.
	// If expr is an identifier (e.g. PLUS, EOL) it must be defined in the
	// header or an imported package so the generated file compiles.
	case strings.HasPrefix(a, "return "):
		expr := strings.TrimSpace(a[len("return "):])
		return fmt.Sprintf(
			`tokens = append(tokens, Token{Kind: %d, Lexeme: lxm, Value: %s, Line: startLine, Col: startCol})`,
			tokenID, expr)

	// ── fallback ──────────────────────────────────────────────────────────
	// Unknown action form: emit a basic token and leave the raw action as a
	// comment so the developer can see what needs adapting.
	default:
		return fmt.Sprintf(
			"// action: %s\n"+
				"tokens = append(tokens, Token{Kind: %d, Lexeme: lxm, Line: startLine, Col: startCol})",
			a, tokenID)
	}
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