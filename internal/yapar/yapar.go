package yapar

import (
	"fmt"
	"strings"
	"unicode"
)

type TokenDef struct {
	Name string `json:"name"`
	ID   int    `json:"id"`
}

type Production struct {
	Name  string     `json:"name"`
	Rules [][]string `json:"rules"`
}

type YalpFile struct {
	Tokens       []TokenDef
	TokenMap     map[string]int
	IgnoreList   []string
	Productions  []Production
	NonTerminals map[string]bool // set of all production head names
	StartSymbol  string          // head of the first production
}

func ParseYalpContent(content string) (*YalpFile, error) {
	content = removeComments(content)

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

func removeComments(s string) string {
	var b strings.Builder
	i := 0
	for i < len(s) {
		if i+1 < len(s) && s[i] == '/' && s[i+1] == '*' {
			i += 2
			for i < len(s) {
				if i+1 < len(s) && s[i] == '*' && s[i+1] == '/' {
					i += 2
					break
				}
				i++
			}
		} else {
			b.WriteByte(s[i])
			i++
		}
	}
	return b.String()
}

func parseTokenSection(yf *YalpFile, section string) error {
	id := 1
	seen := make(map[string]bool)

	for _, line := range strings.Split(section, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		if strings.HasPrefix(line, "%token") {
			rest := strings.TrimSpace(strings.TrimPrefix(line, "%token"))
			for _, name := range strings.Fields(rest) {
				if err := validateTokenName(name); err != nil {
					return err
				}
				if seen[name] {
					return fmt.Errorf("duplicate token declaration %q", name)
				}
				seen[name] = true
				yf.Tokens = append(yf.Tokens, TokenDef{Name: name, ID: id})
				yf.TokenMap[name] = id
				id++
			}
		} else if strings.HasPrefix(line, "IGNORE") {
			name := strings.TrimSpace(strings.TrimPrefix(line, "IGNORE"))
			if name == "" {
				return fmt.Errorf("IGNORE requires a token name")
			}
			yf.IgnoreList = append(yf.IgnoreList, name)
		}
	}
	return nil
}

func parseProductionSection(yf *YalpFile, section string) error {
	lines := strings.Split(section, "\n")
	var current *Production

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// New production: a single word ending with ':'
		if strings.HasSuffix(line, ":") && !strings.ContainsAny(strings.TrimSuffix(line, ":"), " \t") {
			if current != nil {
				return fmt.Errorf("production %q not terminated with ;", current.Name)
			}
			name := strings.TrimSuffix(line, ":")
			current = &Production{Name: name}
			yf.NonTerminals[name] = true
			if yf.StartSymbol == "" {
				yf.StartSymbol = name
			}
			continue
		}

		if line == ";" {
			if current == nil {
				return fmt.Errorf("unexpected ;")
			}
			yf.Productions = append(yf.Productions, *current)
			current = nil
			continue
		}

		if current == nil {
			continue
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

func validate(yf *YalpFile) error {
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
				_, isToken := yf.TokenMap[sym]
				if !isToken && !nonTerminals[sym] {
					return fmt.Errorf("production %q references unknown symbol %q", p.Name, sym)
				}
			}
		}
	}

	return nil
}

func validateTokenName(name string) error {
	if !isValidGoIdent(name) {
		return fmt.Errorf("invalid token name %q: must be a valid Go identifier", name)
	}
	if strings.ToUpper(name) != name {
		return fmt.Errorf("token name %q must be uppercase", name)
	}
	return nil
}

func isValidGoIdent(s string) bool {
	if s == "" {
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
