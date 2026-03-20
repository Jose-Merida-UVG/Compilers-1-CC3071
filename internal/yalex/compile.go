package yalex

import (
	"fmt"
	"regexp"

	"github.com/Jose-Merida-UVG/Compilers-1-CC3071/internal/automata"
	"github.com/Jose-Merida-UVG/Compilers-1-CC3071/internal/regex"
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
			expanded := f.expandDefinitions(pat.Pattern)
			rs, err := regex.NewRegex(expanded).Preprocess()
			if err != nil {
				return nil, fmt.Errorf("rule %q pattern %q: %w", rule.Entrypoint, pat.Pattern, err)
			}
			dfas = append(dfas, automata.Compile(rs))
			actions = append(actions, pat.Action)
		}
	}

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
	// resolved[i] holds the fully-expanded regex for f.Definitions[i].
	// We process definitions in declaration order: each definition can only
	// reference identifiers that were declared earlier, so substituting
	// earlier resolved values into the current def's raw regex is enough.
	resolved := make([]string, len(f.Definitions))
	for i, def := range f.Definitions {
		val := def.Regex
		for j := 0; j < i; j++ {
			val = replaceIdent(val, f.Definitions[j].Identifier, "("+resolved[j]+")")
		}
		resolved[i] = val
	}

	result := pattern
	for i, def := range f.Definitions {
		result = replaceIdent(result, def.Identifier, "("+resolved[i]+")")
	}
	return result
}

// replaceIdent replaces whole-word occurrences of ident in s with repl.
// Uses \b word boundaries so that e.g. "digit" does not match inside "notdigit".
func replaceIdent(s, ident, repl string) string {
	re := regexp.MustCompile(`\b` + regexp.QuoteMeta(ident) + `\b`)
	return re.ReplaceAllString(s, repl)
}
