package yapar

import (
	"fmt"
	"strings"
)

// ValidateTokenCoverage cross-validates tokens between the parser spec (.yalp) and lexer spec (.yal),
// ensuring every declared parser token is actually returned by some lexer action.
func ValidateTokenCoverage(parserTokens map[string]int, lexerRules []LexerRule) error {
	lexerReturns := extractLexerTokens(lexerRules)

	// Find parser tokens with no matching lexer return statement.
	var missingInLexer []string
	for tokenName := range parserTokens {
		if !lexerReturns[tokenName] {
			missingInLexer = append(missingInLexer, tokenName)
		}
	}

	if len(missingInLexer) > 0 {
		return fmt.Errorf(
			"parser declares %d token(s) that lexer never returns: %v\n"+
				"These tokens are in %%token declarations but have no corresponding lexer rule that returns them",
			len(missingInLexer),
			missingInLexer,
		)
	}

	// Find lexer tokens not declared in the parser — likely a compilation error downstream.
	var extraInLexer []string
	for tokenName := range lexerReturns {
		if isSpecialToken(tokenName) {
			continue
		}
		if _, exists := parserTokens[tokenName]; !exists {
			extraInLexer = append(extraInLexer, tokenName)
		}
	}

	if len(extraInLexer) > 0 {
		return fmt.Errorf(
			"lexer returns %d token(s) not declared in parser: %v\n"+
				"These tokens are returned by lexer actions but not listed in %%token declarations",
			len(extraInLexer),
			extraInLexer,
		)
	}

	return nil
}

// LexerRule represents a single rule from a .yal file.
type LexerRule struct {
	Entrypoint string
	Patterns   []LexerPattern
}

// LexerPattern represents one pattern arm in a lexer rule.
type LexerPattern struct {
	Pattern string
	Action  string
}

// extractLexerTokens builds a set of token names found in "return TOKENNAME" statements
// across all lexer action blocks.
func extractLexerTokens(rules []LexerRule) map[string]bool {
	returns := make(map[string]bool)

	for _, rule := range rules {
		for _, pattern := range rule.Patterns {
			if pattern.Action == "" {
				continue
			}
			action := pattern.Action
			for {
				idx := strings.Index(action, "return")
				if idx == -1 {
					break
				}
				rest := action[idx+len("return"):]
				// skip if "return" isn't followed by whitespace (e.g. inside an identifier)
				if len(rest) == 0 || (rest[0] != ' ' && rest[0] != '\t' && rest[0] != '\n') {
					action = rest
					continue
				}
				rest = strings.TrimLeft(rest, " \t\n")
				// walk the identifier characters
				end := 0
				for end < len(rest) && isTokenChar(rest[end]) {
					end++
				}
				// only collect if it starts uppercase (token name convention)
				if end > 0 && rest[0] >= 'A' && rest[0] <= 'Z' {
					returns[rest[:end]] = true
				}
				action = rest[end:]
			}
		}
	}

	return returns
}

// isTokenChar reports whether c is a valid token name character.
func isTokenChar(c byte) bool {
	return (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '_'
}

// isSpecialToken reports whether name is an implicitly handled token (EOF, ERROR)
// that doesn't need a %token declaration.
func isSpecialToken(name string) bool {
	return name == "EOF" || name == "ERROR"
}

// ConvertYalexRulesToLexerRules is a conversion stub meant to be called from handlers.go
// where the yalex package is in scope.
func ConvertYalexRulesToLexerRules(rules any) []LexerRule {
	return nil
}
