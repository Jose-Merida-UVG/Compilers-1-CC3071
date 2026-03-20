package yalex

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// YalFile is the fully parsed representation of a .yal spec file.
type YalFile struct {
	Header      string
	Definitions []LetDefinition
	Rules       []Rule
	Trailer     string
}

// LetDefinition holds a single 'let ident = regexp' binding.
type LetDefinition struct {
	Identifier string
	Regex      string
}

// Rule is one 'rule entrypoint = …' block with its pattern arms.
type Rule struct {
	Entrypoint string
	Args       []string
	Patterns   []RulePattern
}

// RulePattern is a single 'regexp { action }' arm inside a rule.
type RulePattern struct {
	Pattern string
	Action  string // empty when no action block is present
}

// ParseYalFile reads a .yal file from disk and returns the parsed structure.
func ParseYalFile(filename string) (*YalFile, error) {
	// Open file
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Read & parse lines
	scanner := bufio.NewScanner(file)
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err = scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}
	return parseLines(lines)
}

// ParseYalContent parses a .yal spec from an in-memory string,
// this function is only used for internal testing
func ParseYalContent(content string) (*YalFile, error) {
	return parseLines(strings.Split(content, "\n"))
}

// parseLines is the full pre-processing pipeline that maps
// the .yal file to a struct.
func parseLines(lines []string) (*YalFile, error) {
	// Validate only printable ASCII + whitespace
	if err := validateASCII(strings.Join(lines, "\n")); err != nil {
		return nil, err
	}
	// Remove comments delimited by (* ... *)
	stripped, err := removeComments(lines)
	if err != nil {
		return nil, err
	}
	// Locate header / let-defs / rules / trailer boundaries
	bounds, err := identifySections(stripped)
	if err != nil {
		return nil, err
	}
	// Load sections into the YalFile struct
	f, err := parseSections(stripped, bounds)
	if err != nil {
		return nil, err
	}
	// Validate structure of the file
	return f, validateStructure(f)
}

// parseSections extracts header, let-defs, rules, and trailer from the
// already-located section boundaries.
func parseSections(lines []string, bounds *sectionBoundaries) (*YalFile, error) {
	f := &YalFile{}

	// Extract header if well defined bounds
	if bounds.headerStart >= 0 && bounds.headerEnd > bounds.headerStart {
		f.Header = strings.Join(lines[bounds.headerStart:bounds.headerEnd], "\n")
	}

	// Each non-empty 'let ident = regexp' line in the let section.
	if bounds.letStart >= 0 && bounds.letEnd > bounds.letStart {
		for i := bounds.letStart; i < bounds.letEnd; i++ {
			line := strings.TrimSpace(lines[i])
			if line == "" || !strings.HasPrefix(line, "let ") {
				continue
			}
			parts := strings.SplitN(line[4:], "=", 2)
			if len(parts) == 2 {
				f.Definitions = append(f.Definitions, LetDefinition{
					Identifier: strings.TrimSpace(parts[0]),
					Regex:      strings.TrimSpace(parts[1]),
				})
			}
		}
	}

	// Parse rules through rule range found by scanner
	for _, rr := range bounds.ruleRanges {
		rule, err := parseRule(lines[rr.start:rr.end])
		if err != nil {
			return nil, fmt.Errorf("rule at line %d: %w", rr.start, err)
		}
		f.Rules = append(f.Rules, *rule)
	}

	// Extrat trailer if well defined bounds
	if bounds.trailerStart >= 0 && bounds.trailerEnd > bounds.trailerStart {
		f.Trailer = strings.Join(lines[bounds.trailerStart:bounds.trailerEnd], "\n")
	}

	return f, nil
}

// parseRule parses a single rule block. The first line must be `rule name =`;
// subsequent lines are pattern arms, optionally prefixed with `|`.
func parseRule(lines []string) (*Rule, error) {
	if len(lines) == 0 {
		return nil, fmt.Errorf("empty rule section")
	}

	rule := &Rule{}
	firstLine := strings.TrimSpace(lines[0])
	if !strings.HasPrefix(firstLine, "rule ") {
		return nil, fmt.Errorf("expected 'rule' keyword, got %q", firstLine)
	}

	// Extract entrypoint name from `rule <name> =`.
	ruleLine := firstLine[5:]
	if eqIdx := strings.Index(ruleLine, "="); eqIdx > 0 {
		rule.Entrypoint = strings.TrimSpace(ruleLine[:eqIdx])
	}

	// Each remaining non-empty line is one pattern arm.
	for i := 1; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])
		if line == "" {
			continue
		}
		if strings.HasPrefix(line, "|") {
			line = strings.TrimSpace(line[1:])
		}

		pat := RulePattern{}
		aStart, aEnd := findActionBounds(line)
		if aStart >= 0 && aEnd > aStart {
			pat.Pattern = strings.TrimSpace(line[:aStart])
			pat.Action = strings.TrimSpace(line[aStart+1 : aEnd])
		} else {
			pat.Pattern = strings.TrimSpace(line)
		}

		if pat.Pattern != "" {
			rule.Patterns = append(rule.Patterns, pat)
		}
	}

	return rule, nil
}

// findActionBounds returns the indices of the outermost `{` and `}` in a
// pattern line, skipping braces inside single- or double-quoted literals.
// Returns (-1, -1) when no action block is present.
func findActionBounds(line string) (int, int) {
	inSingleQuote := false
	inDoubleQuote := false
	start := -1
	end := -1

	for i := 0; i < len(line); i++ {
		ch := line[i]
		if (inSingleQuote || inDoubleQuote) && ch == '\\' && i+1 < len(line) {
			i++
			continue
		}
		if ch == '\'' && !inDoubleQuote {
			inSingleQuote = !inSingleQuote
			continue
		}
		if ch == '"' && !inSingleQuote {
			inDoubleQuote = !inDoubleQuote
			continue
		}
		if !inSingleQuote && !inDoubleQuote {
			if ch == '{' && start == -1 {
				start = i
			}
			if ch == '}' {
				end = i
			}
		}
	}
	return start, end
}

// validateStructure checks high-level invariants after parsing:
//   - at least one rule must be present
//   - every let-def identifier must be a valid ASCII identifier
//   - every rule must have an entrypoint and at least one pattern
func validateStructure(f *YalFile) error {
	// Have at least one rule
	if len(f.Rules) == 0 {
		return fmt.Errorf("no rule found: a .yal file must contain at least one rule block")
	}
	// Every identifier must have a regex
	for _, def := range f.Definitions {
		if def.Regex == "" {
			return fmt.Errorf("let-def %q has an empty regex", def.Identifier)
		}
	}
	// Rule must have entrypoint name and at least one pattern
	for _, rule := range f.Rules {
		if rule.Entrypoint == "" {
			return fmt.Errorf("rule is missing an entrypoint name")
		}
		if len(rule.Patterns) == 0 {
			return fmt.Errorf("rule %q has no patterns", rule.Entrypoint)
		}
	}
	return nil
}
