package regex

// This package manages the flow of getting a pattern (regex) and
// turning it into an AST. This AST is then consumed downstream for
// the automata generation via Direct Method. The process is a bit
// convoluted, since this project was built top of a different implementation
// for different specs. However, general file structure is the following:
// - regexstring.go: The struct definition of regex, a list of tokens.
// - normalize.go: Construction of a regexstring from a valid string pattern
// - shuntingyard.go: ShuntingYard + other preprocessing to go from infix -> postfix -> AST
// This file specifically just handles the whole process

// Regex represents a single YALex pattern string to be compiled.
type Regex struct {
	Pattern string
}

func NewRegex(pattern string) *Regex {
	return &Regex{Pattern: pattern}
}

// Preprocess normalizes and transforms the YALex pattern into a postfix token
// sequence ready for AST/DFA construction.
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
