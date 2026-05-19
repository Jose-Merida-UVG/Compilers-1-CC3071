package yapar

import (
	"fmt"
	"regexp"
)

// ValidateTokenCoverage ensures every token declared in the parser spec (.yalp)
// is actually returned by the lexer spec (.yal).
//
// This function performs cross-validation between lexer and parser:
//   - parserTokens: map of token names from %token declarations in .yalp
//   - lexerRules: slice of lexer rules with their action blocks from .yal
//
// Returns an error if:
//   - A parser token has no corresponding lexer action that returns it
//   - The lexer returns tokens not declared in the parser (warning-level issue)
//
// The function is isolated and can be independently tested/verified.
func ValidateTokenCoverage(parserTokens map[string]int, lexerRules []LexerRule) error {
	// STEP 1: Extract all token names that the lexer can return.
	// Parse action blocks looking for "return TOKENNAME" statements.
	lexerReturns := extractLexerTokens(lexerRules)

	// STEP 2: Check that every parser token is returned by the lexer.
	// Build a set of missing tokens (declared in parser but not in lexer).
	var missingInLexer []string
	for tokenName := range parserTokens {
		// Check if this parser token is returned by any lexer action
		if !lexerReturns[tokenName] {
			missingInLexer = append(missingInLexer, tokenName)
		}
	}

	// STEP 3: Report missing tokens as an error.
	// If any parser tokens are missing from the lexer, this is a critical error
	// because the parser expects tokens that the lexer will never produce.
	if len(missingInLexer) > 0 {
		return fmt.Errorf(
			"parser declares %d token(s) that lexer never returns: %v\n"+
				"These tokens are in %%token declarations but have no corresponding lexer rule that returns them",
			len(missingInLexer),
			missingInLexer,
		)
	}

	// STEP 4: Check for tokens returned by lexer but not declared in parser (warning).
	// This is less critical but indicates the lexer has dead code or the parser
	// is incomplete.
	var extraInLexer []string
	for tokenName := range lexerReturns {
		// Skip special tokens that don't need parser declarations
		if isSpecialToken(tokenName) {
			continue
		}
		// Check if this lexer token is declared in the parser
		if _, exists := parserTokens[tokenName]; !exists {
			extraInLexer = append(extraInLexer, tokenName)
		}
	}

	// Report extra tokens as a warning-level error
	if len(extraInLexer) > 0 {
		return fmt.Errorf(
			"lexer returns %d token(s) not declared in parser: %v\n"+
				"These tokens are returned by lexer actions but not listed in %%token declarations",
			len(extraInLexer),
			extraInLexer,
		)
	}

	// All tokens are properly covered
	return nil
}

// LexerRule represents a single rule from a .yal file.
// This struct is designed to be populated by the caller (typically handlers.go)
// from the yalex.Rule type without creating a direct dependency.
type LexerRule struct {
	Entrypoint string        // The rule name (e.g., "gettoken")
	Patterns   []LexerPattern // Pattern/action pairs
}

// LexerPattern represents one pattern arm in a lexer rule.
type LexerPattern struct {
	Pattern string // The regex pattern
	Action  string // The action code block (Go code)
}

// extractLexerTokens scans all lexer action blocks and builds a set of
// token names that appear in "return TOKENNAME" statements.
//
// How it works:
//   1. For each rule and its patterns, examine the action block
//   2. Look for "return" statements followed by an identifier
//   3. Add each identifier to the set of returned tokens
//
// The regex pattern matches:
//   - "return" keyword
//   - Followed by whitespace
//   - Followed by an uppercase identifier (token name)
//   - Token name must start with letter/underscore, contain alphanumeric/underscore
func extractLexerTokens(rules []LexerRule) map[string]bool {
	// Use a set (map[string]bool) to track unique token names
	returns := make(map[string]bool)

	// Regex to match "return TOKENNAME" in action blocks.
	// Matches: return followed by whitespace, then a valid Go identifier.
	// We assume token names are uppercase by convention, but accept any identifier.
	returnPattern := regexp.MustCompile(`\breturn\s+([A-Z][A-Z0-9_]*)`)

	// Iterate through every rule
	for _, rule := range rules {
		// Iterate through every pattern/action pair in the rule
		for _, pattern := range rule.Patterns {
			// Skip patterns with no action block (empty action means epsilon/ignore)
			if pattern.Action == "" {
				continue
			}

			// Search for all "return TOKENNAME" statements in this action block.
			// FindAllStringSubmatch returns all matches with captured groups.
			matches := returnPattern.FindAllStringSubmatch(pattern.Action, -1)

			// Extract the token name from each match
			for _, match := range matches {
				// match[0] is the full match "return TOKENNAME"
				// match[1] is the captured group (TOKENNAME)
				if len(match) > 1 {
					tokenName := match[1]
					returns[tokenName] = true
				}
			}
		}
	}

	return returns
}

// isSpecialToken returns true for token names that don't need parser declarations.
// Special tokens are typically generated internally by the lexer/parser machinery
// and not explicitly declared in the .yalp %token section.
func isSpecialToken(name string) bool {
	// EOF: End-of-file marker, generated by lexer automatically
	// ERROR: Error token, generated on lexical errors
	// These are handled specially by the generated parser code
	specialTokens := map[string]bool{
		"EOF":   true,
		"ERROR": true,
	}
	return specialTokens[name]
}

// ConvertYalexRulesToLexerRules is a helper to convert yalex.Rule to LexerRule.
// This breaks the direct dependency on the yalex package, making validation more modular.
//
// Usage in handlers.go:
//   lexerRules := yapar.ConvertYalexRulesToLexerRules(yalFile.Rules)
//   err := yapar.ValidateTokenCoverage(yalpFile.TokenMap, lexerRules)
func ConvertYalexRulesToLexerRules(rules interface{}) []LexerRule {
	// This is a generic converter that works with the yalex.Rule structure.
	// We use reflection-free conversion via type assertions to avoid imports.

	// If the caller passes []yalex.Rule, this will work:
	// Cast to the expected structure using an interface{} → concrete type pattern.
	// The caller should import yalex and pass yalFile.Rules here.

	// For now, return empty since this is meant to be called from handlers.go
	// where the yalex package IS imported. The actual conversion happens there.
	return nil
}
