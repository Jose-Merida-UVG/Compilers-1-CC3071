// Package yapar handles parsing and validation of .yalp grammar spec files.
// A .yalp file has two sections separated by %%:
//   - token section: %token declarations and IGNORE directives
//   - productions section: CFG rules in BNF form
//
// This file is responsible for parsing only: ParseYalpContent reads a raw .yalp
// string and returns a validated YalpFile. End-to-end orchestration (building
// the grammar, computing FIRST/FOLLOW, constructing the LR(0) automaton, and
// deriving the SLR(1) table) lives in compile.go via YalpFile.Compile() and
// is handled via different packages.
package yapar

import (
	"fmt"
	"strings"
	"unicode"
)

// TokenDef (terminals) are defined by ID (lexer output) + their names, non-terminals
// don't need this translation layer since they're strictly grammar-internal.
type TokenDef struct {
	Name string `json:"name"`
	ID   int    `json:"id"`
}

// Production is defined by its non-terminal and its derivations
type Production struct {
	Name  string     `json:"name"`  // LHS of production
	Rules [][]string `json:"rules"` // RHS of production
}

// YalpFile is the fully parsed + pre-processed representation of a .yalp file / spec
type YalpFile struct {
	Tokens       []TokenDef
	TokenMap     map[string]int // Maps string -> int for token name (string) -> ID lookups
	IgnoreList   []string       // Tokens that are ignored by the parser
	Productions  []Production
	NonTerminals map[string]bool // Set of all production head names
	StartSymbol  string          // Head of the first production
}

// ParseYalpContent parses a .yalp spec from an in-memory string. The frontend sends file
// contents directly, so we don't interact with the filesystem. This function is also the
// entrypoint, so it's basically a high-level overview calling auxiliary functions defined
// in this file.
func ParseYalpContent(content string) (*YalpFile, error) {
	var err error
	content, err = removeComments(content)
	if err != nil {
		return nil, err
	}

	parts := strings.SplitN(content, "%%", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("missing %%%% separator between token and productions sections")
	}

	yf := &YalpFile{TokenMap: make(map[string]int), NonTerminals: make(map[string]bool)}

	if err := parseTokenSection(yf, parts[0]); err != nil {
		return nil, err
	}
	if err := parseProductionSection(yf, parts[1]); err != nil {
		return nil, err
	}
	if err := validate(yf); err != nil {
		return nil, err
	}

	return yf, nil
}

// removeComments strips /* ... */ block comments from the source string.
// Supports nested comments by tracking depth. Returns error if comment is unclosed.
func removeComments(s string) (string, error) {
	var b strings.Builder
	i := 0
	depth := 0

	// Go through the in-memory string
	for i < len(s) {
		// Check for '/*' - opening a comment
		if i+1 < len(s) && s[i] == '/' && s[i+1] == '*' {
			depth++
			i += 2
			continue
		}

		// Check for '*/' - closing a comment
		if i+1 < len(s) && s[i] == '*' && s[i+1] == '/' {
			if depth == 0 {
				return "", fmt.Errorf("unmatched */ at position %d", i)
			}
			depth--
			i += 2
			continue
		}

		// Only write to output if we're not inside a comment
		if depth == 0 {
			b.WriteByte(s[i])
		}
		i++
	}

	if depth > 0 {
		return "", fmt.Errorf("unclosed comment (missing %d closing */)", depth)
	}

	return b.String(), nil
}

// parseTokenSection reads the token section (before %%) and populates yf.Tokens,
// yf.TokenMap, and yf.IgnoreList. Each %token line may declare multiple names.
func parseTokenSection(yf *YalpFile, section string) error {
	// Initialize ID to 1, so we can assign increasing unique ID's to each token
	// Note: These are used to 'translate' Lexer output, we get a Lexemme struct
	// with Token ID + extra information.
	id := 1
	seen := make(map[string]bool)

	for _, line := range strings.Split(section, "\n") {
		// Ignore leading / trailing whitespace
		line = strings.TrimSpace(line)
		// Skip empty lines
		if line == "" {
			continue
		}

		// Check for %token entries
		if strings.HasPrefix(line, "%token") {
			rest := strings.TrimSpace(strings.TrimPrefix(line, "%token"))
			for _, name := range strings.Fields(rest) {
				// Validate name
				if err := validateTokenName(name); err != nil {
					return err
				}
				// Validate duplicates
				if seen[name] {
					return fmt.Errorf("duplicate token declaration %q", name)
				}
				// Add to token set + populate ID
				seen[name] = true
				yf.Tokens = append(yf.Tokens, TokenDef{Name: name, ID: id})
				yf.TokenMap[name] = id
				id++
			}
			// Check for IGNORE entries - must have whitespace after IGNORE keyword
		} else if strings.HasPrefix(line, "IGNORE ") || line == "IGNORE" {
			rest := strings.TrimSpace(strings.TrimPrefix(line, "IGNORE"))
			if rest == "" {
				return fmt.Errorf("IGNORE requires at least one token name")
			}
			// Support multiple tokens on one IGNORE line, like %token does
			for _, name := range strings.Fields(rest) {
				yf.IgnoreList = append(yf.IgnoreList, name)
			}
		} else {
			return fmt.Errorf("unexpected token section entry %q: expected %%token or IGNORE", line)
		}
	}
	return nil
}

// parseProductionSection reads the grammar rules section (after %%) and populates
// yf.Productions, yf.NonTerminals, and yf.StartSymbol.
// Productions are written as:
// name:
// body
// | alt;
// An empty body line means an ε production.
func parseProductionSection(yf *YalpFile, section string) error {
	// Iterate through lines, keeping track of current production
	lines := strings.Split(section, "\n")
	var current *Production

	for _, line := range lines {
		line = strings.TrimSpace(line)
		// Skip empty lines
		if line == "" {
			continue
		}

		// Check for line ending in ':' (new production)
		if strings.HasSuffix(line, ":") {
			name := strings.TrimSuffix(line, ":")
			// Production name must not contain whitespace
			if len(strings.Fields(name)) != 1 {
				return fmt.Errorf("production name %q cannot contain whitespace", name)
			}
			// If we're already 'inside' of a production then the section was not ended properly
			if current != nil {
				return fmt.Errorf("production %q not terminated with ;", current.Name)
			}
			// Non-terminal names must be lowercase per spec
			if strings.ToLower(name) != name {
				return fmt.Errorf("non-terminal name %q must be lowercase", name)
			}
			current = &Production{Name: name}
			yf.NonTerminals[name] = true
			if yf.StartSymbol == "" {
				yf.StartSymbol = name
			}
			continue
		}
		// Check for ';' line (ends production)
		if line == ";" {
			// If not inside production, it's an input file error
			if current == nil {
				return fmt.Errorf("unexpected ; outside of production")
			}
			// 'Close' the production section & append to struct
			yf.Productions = append(yf.Productions, *current)
			current = nil
			continue
		}

		// If we're not in a production, this line is unexpected
		if current == nil {
			return fmt.Errorf("unexpected line %q: expected production declaration (name:)", line)
		}

		// Strip leading pipe for alternate rules
		if strings.HasPrefix(line, "|") {
			line = strings.TrimSpace(line[1:])
		}

		// Always append — empty symbols slice means ε production.
		current.Rules = append(current.Rules, strings.Fields(line))
	}

	// Per spec: last production may omit the trailing ;
	if current != nil {
		yf.NonTerminals[current.Name] = true
		if yf.StartSymbol == "" {
			yf.StartSymbol = current.Name
		}
		yf.Productions = append(yf.Productions, *current)
	}

	return nil
}

// validate checks cross-section consistency: IGNORE tokens must be declared,
// and every symbol referenced in a production body must be a known token or non-terminal.
func validate(yf *YalpFile) error {
	// Grammar must have at least one production
	if len(yf.Productions) == 0 {
		return fmt.Errorf("grammar must have at least one production")
	}

	// Each production must have at least one rule
	for _, p := range yf.Productions {
		if len(p.Rules) == 0 {
			return fmt.Errorf("production %q has no rules (expected at least one, or an empty line for ε)", p.Name)
		}
	}

	// IGNORE tokens must be declared
	for _, name := range yf.IgnoreList {
		if _, ok := yf.TokenMap[name]; !ok {
			return fmt.Errorf("IGNORE references undeclared token %q", name)
		}
	}

	// Build set of known non-terminal names
	nonTerminals := make(map[string]bool)
	for _, p := range yf.Productions {
		nonTerminals[p.Name] = true
	}

	// Every symbol in a production body must be a declared token or a known non-terminal
	for _, p := range yf.Productions {
		for _, rule := range p.Rules {
			for _, sym := range rule { // empty rule = ε, skip validation
				if strings.ContainsRune(sym, 'ε') {
					return fmt.Errorf("production %q: use an empty body line to express ε, not the literal character", p.Name)
				}
				_, isToken := yf.TokenMap[sym]
				// Tokens must be uppercase, non-terminals must be lowercase.
				// Note: Symbols with only underscores/digits (e.g., "_", "123") are
				// case-neutral: ToUpper("_") == ToLower("_") == "_", so they pass both
				// checks. This is correct behavior - they're valid in both contexts.
				if isToken && strings.ToUpper(sym) != sym {
					return fmt.Errorf("production %q: token %q must be uppercase", p.Name, sym)
				}
				if !isToken && strings.ToLower(sym) != sym {
					return fmt.Errorf("production %q: non-terminal %q must be lowercase", p.Name, sym)
				}
				if !isToken && !nonTerminals[sym] {
					return fmt.Errorf("production %q references unknown symbol %q", p.Name, sym)
				}
			}
		}
	}

	// Rule 3:
	// Tokens declared under IGNORE cannot appear in productions.
	ignoreSet := make(map[string]bool)
	for _, tok := range yf.IgnoreList {
		ignoreSet[tok] = true
	}

	for _, p := range yf.Productions {
		for _, rule := range p.Rules {
			for _, sym := range rule {
				if ignoreSet[sym] {
					return fmt.Errorf(
						"token %q is declared IGNORE but used in production %q",
						sym,
						p.Name,
					)
				}
			}
		}
	}

	return nil
}

// validateTokenName checks that a token name follows .yalp spec rules:
// must be a non-empty uppercase identifier (letters, digits, underscores;
// first character must be a letter or underscore).
func validateTokenName(name string) error {
	if !isValidSymbolName(name) {
		return fmt.Errorf("invalid token name %q: must start with a letter or underscore and contain only letters, digits, or underscores", name)
	}
	if strings.ToUpper(name) != name {
		return fmt.Errorf("token name %q must be uppercase", name)
	}
	return nil
}

// isValidSymbolName reports whether s is a valid .yalp symbol name:
// non-empty, starts with a letter or underscore, followed by letters, digits, or underscores.
// Rejects the literal ε character — use an empty body line for epsilon productions instead.
func isValidSymbolName(s string) bool {
	if s == "" {
		return false
	}
	// Explicitly forbid the epsilon character to give a clear error message
	if strings.ContainsRune(s, 'ε') {
		return false
	}
	for i, r := range s {
		if i == 0 {
			if !unicode.IsLetter(r) && r != '_' {
				return false
			}
		} else {
			if !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '_' {
				return false
			}
		}
	}
	return true
}
