package regex

// Regex represents a single YALex pattern string to be compiled.
type Regex struct {
	Pattern string
}

func NewRegex(pattern string) *Regex {
	return &Regex{Pattern: pattern}
}

// Preprocess normalizes and transforms the YALex pattern into a postfix token
// sequence ready for AST/DFA construction. The pipeline is:
//
//  1. NormalizePattern  — resolve all YALex syntax (quoted chars, string
//     literals, char classes, negation, set-difference, wildcard, eof) into
//     a flat stream of canonical tokens.
//  2. HandleSpecialOperators — desugar + and ? into * and |.
//  3. HandleExplicitConcatenation — insert Cat (~) tokens.
//  4. ShuntingYard — convert infix to postfix (RPN).
//  5. AppendEndMarker — append the # end-marker leaf.
func (r *Regex) Preprocess() (*RegexString, error) {
	rs, err := NormalizePattern(r.Pattern)
	if err != nil {
		return nil, err
	}
	rs.HandleSpecialOperators()
	rs.HandleExplicitConcatenation()
	rs.ShuntingYard()
	rs.AppendEndMarker()
	return rs, nil
}
