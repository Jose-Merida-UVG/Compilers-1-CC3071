package parser

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// YALexFile represents the parsed structure of a .yal file
type YALexFile struct {
	Header      string
	Definitions []LetDefinition
	Rules       []Rule
	Trailer     string
}

// LetDefinition represents: let ident = regexp
type LetDefinition struct {
	Identifier string
	Regex      string
}

// Rule represents: rule entrypoint [args] = pattern { action } | ...
type Rule struct {
	Entrypoint string
	Args       []string
	Patterns   []RulePattern
}

// RulePattern represents a single pattern with its action
type RulePattern struct {
	Pattern string
	Action  string
}

// ParseYALexFile reads and parses a .yal file
func ParseYALexFile(filename string) (*YALexFile, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err = scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	// Step 1: Remove all (* ... *) comments from the file
	lines = removeComments(lines)

	// Step 2: Identify where each section starts and ends
	boundaries, err := identifySections(lines)
	if err != nil {
		return nil, err
	}

	// Step 3: Parse content from each identified section
	yalFile, err := parseSections(lines, boundaries)
	if err != nil {
		return nil, err
	}

	return yalFile, nil
}

// SectionBoundaries holds the line ranges for each section.
// Each section is defined by [start, end) where end is exclusive.
type SectionBoundaries struct {
	HeaderStart  int
	HeaderEnd    int
	LetStart     int
	LetEnd       int
	RuleRanges   []RuleRange // Support multiple rules
	TrailerStart int
	TrailerEnd   int
}

// RuleRange defines the line range for a single rule
type RuleRange struct {
	Start int
	End   int
}

// removeComments strips all (* ... *) comments from the file.
// Strategy: Track whether we're in a comment, single-quote string, or double-quote string.
// Only remove (* *) when not inside string literals to preserve patterns like "(*" or '*)'.
// removeComments strips (* ... *) comments from the entire file.
// Supports nested comments like (* outer (* inner *) outer *).
// Tracks string context to avoid (* inside strings being treated as comments.
func removeComments(lines []string) []string {
	// Join all lines to handle multi-line comments
	fullText := strings.Join(lines, "\n")

	result := ""
	i := 0
	depth := 0 // Track nested comments
	inSingleQuote := false
	inDoubleQuote := false

	for i < len(fullText) {
		ch := fullText[i]

		// Skip escaped characters inside strings (simplified: just skip next char)
		if (inSingleQuote || inDoubleQuote) && ch == '\\' && i+1 < len(fullText) {
			if depth == 0 {
				result += string(ch)
				result += string(fullText[i+1])
			}
			i += 2
			continue
		}

		// Track string context (only when not in comment)
		if depth == 0 {
			if ch == '\'' && !inDoubleQuote {
				inSingleQuote = !inSingleQuote
				result += string(ch)
				i++
				continue
			}
			if ch == '"' && !inSingleQuote {
				inDoubleQuote = !inDoubleQuote
				result += string(ch)
				i++
				continue
			}
		}

		// Only process comments when not inside a string
		if !inSingleQuote && !inDoubleQuote {
			// Check for comment start
			if i+1 < len(fullText) && fullText[i:i+2] == "(*" {
				depth++
				i += 2
				continue
			}

			// Check for comment end
			if i+1 < len(fullText) && fullText[i:i+2] == "*)" {
				depth--
				i += 2
				continue
			}
		}

		// Only add characters if we're not in a comment
		if depth == 0 {
			result += string(ch)
		}
		i++
	}

	// Split back into lines
	return strings.Split(result, "\n")
}

// identifySections finds the boundaries of header, let definitions, rules, and trailer.
// Strategy:
// 1. Header: First { ... } block
// 2. Let definitions: All "let" lines before first "rule"
// 3. Rules: All "rule" declarations until trailer or EOF (supports multiple rules)
// 4. Trailer: Last { ... } block
func identifySections(lines []string) (*SectionBoundaries, error) {
	bounds := &SectionBoundaries{
		HeaderStart:  -1,
		HeaderEnd:    -1,
		LetStart:     -1,
		LetEnd:       -1,
		TrailerStart: -1,
		TrailerEnd:   -1,
	}

	i := 0

	// Skip initial blank lines
	for i < len(lines) && strings.TrimSpace(lines[i]) == "" {
		i++
	}

	// Look for header section (first { ... })
	if i < len(lines) && strings.TrimSpace(lines[i]) == "{" {
		bounds.HeaderStart = i + 1
		i++
		braceCount := 1
		for i < len(lines) && braceCount > 0 {
			braceCount += countBraces(lines[i])
			if braceCount == 0 {
				bounds.HeaderEnd = i
				i++
				break
			}
			i++
		}
	}

	// Skip blank lines after header
	for i < len(lines) && strings.TrimSpace(lines[i]) == "" {
		i++
	}

	// Find first 'rule' keyword - everything before it is let definitions
	firstRuleLine := -1
	for j := i; j < len(lines); j++ {
		if strings.HasPrefix(strings.TrimSpace(lines[j]), "rule ") {
			firstRuleLine = j
			break
		}
	}

	// Let definitions: from current position to first rule
	if firstRuleLine > i {
		bounds.LetStart = i
		bounds.LetEnd = firstRuleLine
		i = firstRuleLine
	} else if firstRuleLine == i {
		bounds.LetStart = -1
		bounds.LetEnd = -1
	}

	// Find ALL rule declarations and their boundaries
	// Each rule ends when we hit: next "rule", trailer "{", or EOF
	if firstRuleLine >= 0 {
		// Scan through and find all lines that start with "rule "
		ruleStarts := []int{}
		for j := firstRuleLine; j < len(lines); j++ {
			line := strings.TrimSpace(lines[j])
			// Stop if we hit trailer section
			if line == "{" {
				i = j
				break
			}
			if strings.HasPrefix(line, "rule ") {
				ruleStarts = append(ruleStarts, j)
			}
		}

		// For each rule start, determine its end
		for idx, start := range ruleStarts {
			ruleRange := RuleRange{Start: start}

			// Rule ends at next rule start, or trailer, or EOF
			if idx+1 < len(ruleStarts) {
				ruleRange.End = ruleStarts[idx+1]
			} else {
				// This is the last rule - find where it ends
				foundEnd := false
				for j := start + 1; j < len(lines); j++ {
					if strings.TrimSpace(lines[j]) == "{" {
						ruleRange.End = j
						i = j
						foundEnd = true
						break
					}
				}
				if !foundEnd {
					ruleRange.End = len(lines)
					i = len(lines)
				}
			}

			bounds.RuleRanges = append(bounds.RuleRanges, ruleRange)
		}
	}

	// Look for trailer section (last { ... })
	if i < len(lines) && strings.TrimSpace(lines[i]) == "{" {
		bounds.TrailerStart = i + 1
		i++
		braceCount := 1
		for i < len(lines) && braceCount > 0 {
			braceCount += countBraces(lines[i])
			if braceCount == 0 {
				bounds.TrailerEnd = i
				break
			}
			i++
		}
	}

	return bounds, nil
}

// countBraces counts net brace changes in a line, respecting string literals.
// Returns +1 for each {, -1 for each }, ignoring braces inside strings.
func countBraces(line string) int {
	count := 0
	inSingleQuote := false
	inDoubleQuote := false

	for i := 0; i < len(line); i++ {
		ch := line[i]

		// Handle escape sequences inside strings
		if (inSingleQuote || inDoubleQuote) && ch == '\\' && i+1 < len(line) {
			i++ // Skip escaped character
			continue
		}

		// Track string context
		if ch == '\'' && !inDoubleQuote {
			inSingleQuote = !inSingleQuote
			continue
		}
		if ch == '"' && !inSingleQuote {
			inDoubleQuote = !inDoubleQuote
			continue
		}

		// Count braces only outside of strings
		if !inSingleQuote && !inDoubleQuote {
			if ch == '{' {
				count++
			} else if ch == '}' {
				count--
			}
		}
	}

	return count
}

// parseSections extracts and parses content from each identified section
func parseSections(lines []string, bounds *SectionBoundaries) (*YALexFile, error) {
	yalFile := &YALexFile{}

	// Parse header
	if bounds.HeaderStart >= 0 && bounds.HeaderEnd > bounds.HeaderStart {
		headerLines := lines[bounds.HeaderStart:bounds.HeaderEnd]
		yalFile.Header = strings.Join(headerLines, "\n")
	}

	// Parse let definitions
	if bounds.LetStart >= 0 && bounds.LetEnd > bounds.LetStart {
		for i := bounds.LetStart; i < bounds.LetEnd; i++ {
			line := strings.TrimSpace(lines[i])

			if line == "" {
				continue
			}

			if strings.HasPrefix(line, "let ") {
				parts := strings.SplitN(line[4:], "=", 2)
				if len(parts) == 2 {
					def := LetDefinition{
						Identifier: strings.TrimSpace(parts[0]),
						Regex:      strings.TrimSpace(parts[1]),
					}
					yalFile.Definitions = append(yalFile.Definitions, def)
				}
			}
		}
	}

	// Parse all rules
	for _, ruleRange := range bounds.RuleRanges {
		rule, err := parseRule(lines[ruleRange.Start:ruleRange.End])
		if err != nil {
			return nil, fmt.Errorf("error parsing rule at line %d: %w", ruleRange.Start, err)
		}
		yalFile.Rules = append(yalFile.Rules, *rule)
	}

	// Parse trailer
	if bounds.TrailerStart >= 0 && bounds.TrailerEnd > bounds.TrailerStart {
		trailerLines := lines[bounds.TrailerStart:bounds.TrailerEnd]
		yalFile.Trailer = strings.Join(trailerLines, "\n")
	}

	return yalFile, nil
}

// parseRule parses a single rule section.
// Format: rule entrypoint [args] =
//
//	  pattern { action }
//	| pattern { action }
//	...
func parseRule(lines []string) (*Rule, error) {
	if len(lines) == 0 {
		return nil, fmt.Errorf("empty rule section")
	}

	rule := &Rule{}

	// Parse first line: rule entrypoint [args] =
	firstLine := strings.TrimSpace(lines[0])
	if !strings.HasPrefix(firstLine, "rule ") {
		return nil, fmt.Errorf("rule must start with 'rule' keyword")
	}

	ruleLine := firstLine[5:]
	eqIndex := strings.Index(ruleLine, "=")
	if eqIndex > 0 {
		entrySig := strings.TrimSpace(ruleLine[:eqIndex])
		rule.Entrypoint = entrySig
	}

	// Parse pattern-action pairs from remaining lines
	for i := 1; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])

		if line == "" {
			continue
		}

		// Remove leading | if present
		if strings.HasPrefix(line, "|") {
			line = strings.TrimSpace(line[1:])
		}

		pattern := RulePattern{}

		// Find action block { ... } while respecting string literals
		actionStart, actionEnd := findActionBounds(line)

		if actionStart >= 0 && actionEnd > actionStart {
			pattern.Pattern = strings.TrimSpace(line[:actionStart])
			pattern.Action = strings.TrimSpace(line[actionStart+1 : actionEnd])
		} else {
			pattern.Pattern = strings.TrimSpace(line)
		}

		if pattern.Pattern != "" {
			rule.Patterns = append(rule.Patterns, pattern)
		}
	}

	return rule, nil
}

// findActionBounds locates the { and } delimiters of an action block,
// skipping over braces that appear inside string literals.
func findActionBounds(line string) (int, int) {
	inSingleQuote := false
	inDoubleQuote := false
	actionStart := -1
	actionEnd := -1

	for i := 0; i < len(line); i++ {
		ch := line[i]

		// Handle escape sequences inside strings
		if (inSingleQuote || inDoubleQuote) && ch == '\\' && i+1 < len(line) {
			i++ // Skip escaped character
			continue
		}

		// Track string context
		if ch == '\'' && !inDoubleQuote {
			inSingleQuote = !inSingleQuote
			continue
		}
		if ch == '"' && !inSingleQuote {
			inDoubleQuote = !inDoubleQuote
			continue
		}

		// Look for { } only outside of strings
		if !inSingleQuote && !inDoubleQuote {
			if ch == '{' && actionStart == -1 {
				actionStart = i
			}
			if ch == '}' {
				actionEnd = i
			}
		}
	}

	return actionStart, actionEnd
}

// PrintStructure prints the parsed structure for debugging
func (y *YALexFile) PrintStructure() {
	fmt.Println("=== YALex File Structure ===")

	if y.Header != "" {
		fmt.Println("\n--- Header ---")
		fmt.Println(y.Header)
	}

	if len(y.Definitions) > 0 {
		fmt.Println("\n--- Let Definitions ---")
		for i, def := range y.Definitions {
			fmt.Printf("%d. let %s = %s\n", i+1, def.Identifier, def.Regex)
		}
	}

	if len(y.Rules) > 0 {
		fmt.Println("\n--- Rules ---")
		for i, rule := range y.Rules {
			fmt.Printf("\nRule %d: %s\n", i+1, rule.Entrypoint)
			if len(rule.Args) > 0 {
				fmt.Printf("  Args: %v\n", rule.Args)
			}
			fmt.Printf("  Patterns (%d):\n", len(rule.Patterns))
			for j, pat := range rule.Patterns {
				fmt.Printf("    %d. Pattern: %s\n", j+1, pat.Pattern)
				if pat.Action != "" {
					fmt.Printf("       Action: %s\n", pat.Action)
				}
			}
		}
	}

	if y.Trailer != "" {
		fmt.Println("\n--- Trailer ---")
		fmt.Println(y.Trailer)
	}

	fmt.Println("\n=== End of Structure ===")
}
