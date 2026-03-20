package yalex

import (
	"fmt"
	"strings"
)

// validateASCII rejects any character that is not printable ASCII
// or one of the standard whitespace control codes (\t, \n, \r). This
// is an additional constraint we've implemented to make sure things like
// character class negations don't go crazy.
func validateASCII(content string) error {
	lineNum := 1
	for _, ch := range content {
		if ch == '\n' {
			lineNum++
			continue
		}
		if ch == '\t' || ch == '\r' {
			continue
		}
		if ch < 0x20 || ch > 0x7E {
			return fmt.Errorf("non-printable or non-ASCII character U+%04X at line %d", ch, lineNum)
		}
	}
	return nil
}

// removeComments strips (* ... *) block comments from the source. Handles
// nested / multi-line comments and ignores '(*' character combinations present in string
// literals. This is another constraint we have, where our 'destination language'
// code blocks CANNOT have this combination of characters, otherwise the input file is
// deemed invalid.
func removeComments(lines []string) ([]string, error) {
	fullText := strings.Join(lines, "\n")
	var result strings.Builder
	i := 0
	depth := 0
	lineNum := 1
	inSingleQuote := false
	inDoubleQuote := false

	for i < len(fullText) {
		ch := fullText[i]
		if ch == '\n' {
			lineNum++
		}

		// For any escaped character inside a literal we can consume both
		if (inSingleQuote || inDoubleQuote) && ch == '\\' && i+1 < len(fullText) {
			if depth == 0 {
				result.WriteByte(ch)
				result.WriteByte(fullText[i+1])
			}
			i += 2
			continue
		}

		// Quote tracking (only matters outside comments).
		if depth == 0 {
			if ch == '\'' && !inDoubleQuote {
				inSingleQuote = !inSingleQuote
				result.WriteByte(ch)
				i++
				continue
			}
			if ch == '"' && !inSingleQuote {
				inDoubleQuote = !inDoubleQuote
				result.WriteByte(ch)
				i++
				continue
			}
		}

		// Comment open / close, only outside single or double quotes
		if !inSingleQuote && !inDoubleQuote && i+1 < len(fullText) {
			if fullText[i:i+2] == "(*" {
				depth++
				i += 2
				continue
			}
			if fullText[i:i+2] == "*)" {
				if depth == 0 {
					return nil, fmt.Errorf("unexpected comment close '*)'  at line %d", lineNum)
				}
				depth--
				i += 2
				continue
			}
		}

		if depth == 0 {
			result.WriteByte(ch)
		}
		i++
	}

	if depth > 0 {
		return nil, fmt.Errorf("unclosed comment '(*' (depth %d at end of file)", depth)
	}
	return strings.Split(result.String(), "\n"), nil
}

// sectionBoundaries holds the line-index ranges for each logical section of a
// .yal file after comment removal. -1 means the section is absent.
type sectionBoundaries struct {
	headerStart  int
	headerEnd    int
	letStart     int
	letEnd       int
	ruleRanges   []ruleRange
	trailerStart int
	trailerEnd   int
}

type ruleRange struct {
	start int
	end   int
}

// identifySections scans the post-comment lines and records where each section
// lives so the parser can extract content without re-scanning.
func identifySections(lines []string) (*sectionBoundaries, error) {
	// Initialize all sections as empty
	bounds := &sectionBoundaries{
		headerStart:  -1,
		headerEnd:    -1,
		letStart:     -1,
		letEnd:       -1,
		trailerStart: -1,
		trailerEnd:   -1,
	}
	i := 0

	// Header block, starts with '{' on its own line
	// Skip empty lines
	for i < len(lines) && strings.TrimSpace(lines[i]) == "" {
		i++
	}

	// Match the '{' (if present) to its closing '}'
	if i < len(lines) && strings.TrimSpace(lines[i]) == "{" {
		bounds.headerStart = i + 1
		i++
		depth := 1
		for i < len(lines) && depth > 0 {
			depth += countBraces(lines[i])
			if depth == 0 {
				bounds.headerEnd = i
				i++
				break
			}
			i++
		}
		if depth > 0 {
			return nil, fmt.Errorf("unclosed header block '{'")
		}
	}

	// Let-definition section, all lines between header and first `rule` keyword.
	// Skip empty lines
	for i < len(lines) && strings.TrimSpace(lines[i]) == "" {
		i++
	}
	// Every line is a let-def until we see the keyword 'rule'
	firstRule := -1
	for j := i; j < len(lines); j++ {
		if strings.HasPrefix(strings.TrimSpace(lines[j]), "rule ") {
			firstRule = j
			break
		}
	}
	if firstRule > i {
		bounds.letStart = i
		bounds.letEnd = firstRule
		i = firstRule
	} else if firstRule == i {
		// Rules start immediately, no let-defs.
		bounds.letStart = -1
		bounds.letEnd = -1
	}

	// Rule blocks scan from first `rule` to the trailer `{` (or EOF).
	if firstRule >= 0 {
		var starts []int
		for j := firstRule; j < len(lines); j++ {
			line := strings.TrimSpace(lines[j])
			if line == "{" {
				// Bare '{' marks the start of the trailer.
				i = j
				break
			}
			if strings.HasPrefix(line, "rule ") {
				starts = append(starts, j)
			}
		}

		for idx, start := range starts {
			rr := ruleRange{start: start}
			if idx+1 < len(starts) {
				rr.end = starts[idx+1]
			} else {
				// Last rule ends at the trailer '{' or EOF.
				rr.end = len(lines)
				i = len(lines)
				for j := start + 1; j < len(lines); j++ {
					if strings.TrimSpace(lines[j]) == "{" {
						rr.end = j
						i = j
						break
					}
				}
			}
			bounds.ruleRanges = append(bounds.ruleRanges, rr)
		}
	}

	// Trailer block '{' on its own like, brace count is also checked
	if i < len(lines) && strings.TrimSpace(lines[i]) == "{" {
		bounds.trailerStart = i + 1
		i++
		depth := 1
		for i < len(lines) && depth > 0 {
			depth += countBraces(lines[i])
			if depth == 0 {
				bounds.trailerEnd = i
				break
			}
			i++
		}
		if depth > 0 {
			return nil, fmt.Errorf("unclosed trailer block '{'")
		}
	}

	return bounds, nil
}

// countBraces returns the net brace delta (+1 per `{`, -1 per `}`) for one
// line, ignoring braces inside single or double quoted literals. Function
// mainly just cleans up logic for some of the other functions as they work
// on line basis.
func countBraces(line string) int {
	count := 0
	inSingleQuote := false
	inDoubleQuote := false
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
			if ch == '{' {
				count++
			} else if ch == '}' {
				count--
			}
		}
	}
	return count
}
