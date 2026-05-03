package yalex

import (
	"fmt"

	"github.com/Jose-Merida-UVG/Compilers-1-CC3071/internal/yalex/automata"
	"github.com/Jose-Merida-UVG/Compilers-1-CC3071/internal/yalex/regex"
)

// CompiledLexer is the result of compiling a YalFile, a single minimized DFA
// that recognizes all token patterns, and the ordered action strings to execute
// on a match. Actions[i] corresponds to TokenID i+1 in the DFA (1-based).
type CompiledLexer struct {
	DFA     *automata.DFA
	Actions []string
}

// Compile expands let-definitions into each rule pattern, builds a minimized
// DFA per pattern, merges them all into one combined minimized DFA, and returns
// the DFA alongside its action table. This is the main orchestration method that
// codegen consumes.
func (f *YalFile) Compile() (*CompiledLexer, error) {
	var dfas []*automata.DFA
	var actions []string

	for _, rule := range f.Rules {
		for _, pat := range rule.Patterns {
			// Expand definitions
			expanded := f.expandDefinitions(pat.Pattern)
			// Create DFA
			rs, err := regex.NewRegex(expanded).Preprocess()
			if err != nil {
				return nil, fmt.Errorf("rule %q pattern %q: %w", rule.Entrypoint, pat.Pattern, err)
			}
			// Store DFA + actions
			dfas = append(dfas, automata.Compile(rs))
			actions = append(actions, pat.Action)
		}
	}

	// Returned merged + minimized DFA + actions based on token ID.
	return &CompiledLexer{
		DFA:     automata.Merge(dfas),
		Actions: actions,
	}, nil
}

// expandDefinitions replaces every let-identifier in a pattern with its
// corresponding regex, wrapped in parentheses to preserve precedence.
// It first builds a fully-resolved value for each definition (so that
// references between defs are transitively expanded), then applies
// word-boundary-aware replacement to the rule pattern.
func (f *YalFile) expandDefinitions(pattern string) string {
	// Each definition can only reference identifiers that were declared earlier,
	// so substituting earlier resolved values into the current def's raw regex is enough.
	resolved := make([]string, len(f.Definitions))
	for i, def := range f.Definitions {
		val := def.Regex
		for j := range i {
			val = replaceIdent(val, f.Definitions[j].Identifier, "("+resolved[j]+")")
		}
		resolved[i] = val
	}

	// Replace the definitions by their respective patterns
	result := pattern
	for i, def := range f.Definitions {
		result = replaceIdent(result, def.Identifier, "("+resolved[i]+")")
	}
	return result
}

// replaceIdent replaces whole-word occurrences of ident in s with repl.
// Skips occurrences inside double-quoted strings and only matches when the
// identifier is not adjacent to a word character on either side. Eg.
// identifiera != identifier but identifier* or identifier_ does.
func replaceIdent(s, ident, repl string) string {
	var out []byte
	i := 0
	for i < len(s) {
		// Skip quoted strings/chars verbatim so identifiers inside are not replaced.
		if s[i] == '"' || s[i] == '\'' {
			q := s[i]
			out = append(out, s[i])
			i++
			for i < len(s) {
				out = append(out, s[i])
				if s[i] == '\\' && i+1 < len(s) {
					i++
					out = append(out, s[i])
				} else if s[i] == q {
					i++
					break
				}
				i++
			}
			continue
		}
		// Check if ident starts here as a whole word.
		end := i + len(ident)
		if end <= len(s) && s[i:end] == ident {
			isWord := func(c byte) bool {
				return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9')
			}
			if (i == 0 || !isWord(s[i-1])) && (end == len(s) || !isWord(s[end])) {
				out = append(out, repl...)
				i = end
				continue
			}
		}
		out = append(out, s[i])
		i++
	}
	return string(out)
}
