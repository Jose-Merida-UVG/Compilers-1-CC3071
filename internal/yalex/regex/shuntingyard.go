package regex

import (
	"slices"

	"github.com/Jose-Merida-UVG/Compilers-1-CC3071/internal/ds"
)

// findOpenPar scans leftward from closeIndex to find the matching `(`.
// Just a utility func :)
func findOpenPar(regex []Token, closeIndex int) int {
	parenCount := 1
	for i := closeIndex - 1; i >= 0; i-- {
		currChar := regex[i]
		if currChar.Kind == RParen {
			parenCount++
		} else if currChar.Kind == LParen {
			parenCount--
			if parenCount == 0 {
				return i
			}
		}
	}
	return -1
}

// HandleSpecialOperators substitutes '?' and '+' with their primitive equivalents.
//   - x? → (x|ε)
//   - x+ → xx*
func (regex *RegexString) HandleSpecialOperators() {
	result := []Token{}
	for i := 0; i < len(regex.Chars); i++ {
		currChar := regex.Chars[i]
		if i+1 < len(regex.Chars) {
			nextChar := regex.Chars[i+1]
			if nextChar.Kind == Quest {
				if currChar.Kind == RParen {
					result = append(result, Token{Kind: RParen, Value: ')'})
					openParIndex := findOpenPar(result, len(result)-1)
					if openParIndex < 0 {
						result = append(result, currChar)
						i++
						continue
					}
					result = slices.Insert(result, openParIndex, Token{Kind: LParen, Value: '('})
					result = append(result, Token{Kind: Alt, Value: '|'})
					result = append(result, epsToken())
					result = append(result, Token{Kind: RParen, Value: ')'})
				} else {
					result = append(result, Token{Kind: LParen, Value: '('})
					result = append(result, currChar)
					result = append(result, Token{Kind: Alt, Value: '|'})
					result = append(result, epsToken())
					result = append(result, Token{Kind: RParen, Value: ')'})
				}
				i++
			} else if nextChar.Kind == Plus {
				if currChar.Kind == RParen {
					result = append(result, Token{Kind: RParen, Value: ')'})
					openParIndex := findOpenPar(result, len(result)-1)
					if openParIndex < 0 {
						result = append(result, currChar)
						i++
						continue
					}
					resultLen := len(result) - 1
					for j := openParIndex; j <= resultLen; j++ {
						result = append(result, result[j])
					}
					result = append(result, Token{Kind: Star, Value: '*'})
				} else {
					result = append(result, currChar)
					result = append(result, currChar)
					result = append(result, Token{Kind: Star, Value: '*'})
				}
				i++
			} else {
				result = append(result, currChar)
			}
		} else {
			result = append(result, currChar)
		}
	}
	regex.Chars = result
}

// HandleExplicitConcatenation inserts the `~` Cat operator between tokens
// that should be concatenated but have no explicit operator between them.
func (regex *RegexString) HandleExplicitConcatenation() {
	result := []Token{}
	for i, currChar := range regex.Chars {
		result = append(result, currChar)
		if i+1 < len(regex.Chars) {
			nextChar := regex.Chars[i+1]
			// Don't insert Cat if the current token "opens" something or
			// the next token "closes" or is itself a binary operator.
			isCurrOpen := currChar.Kind == LParen || currChar.Kind == Alt || currChar.Kind == Cat
			isNextClose := nextChar.Kind == RParen || nextChar.Kind == Alt || nextChar.Kind == Star || nextChar.Kind == Plus || nextChar.Kind == Quest
			if !isCurrOpen && !isNextClose {
				result = append(result, Token{Kind: Cat, Value: '~'})
			}
		}
	}
	regex.Chars = result
}

// precedence returns the Shunting-Yard precedence for operator tokens.
func precedence(t Token) int {
	switch t.Kind {
	case LParen:
		return 1
	case Alt:
		return 2
	case Cat:
		return 3
	case Star:
		return 4
	default:
		return 10
	}
}

// ShuntingYard converts the infix token sequence to postfix.
func (regex *RegexString) ShuntingYard() {
	result := []Token{}
	stack := ds.Stack[Token]{}
	for _, currChar := range regex.Chars {
		if !currChar.IsOperator() {
			result = append(result, currChar)
		} else if currChar.Kind == LParen {
			stack.Push(currChar)
		} else if currChar.Kind == RParen {
			for !stack.IsEmpty() {
				top, _ := stack.Peek()
				if top.Kind == LParen {
					break
				}
				popped, _ := stack.Pop()
				result = append(result, popped)
			}
			stack.Pop() // discard the LParen
		} else {
			for !stack.IsEmpty() {
				top, _ := stack.Peek()
				if top.Kind == LParen || precedence(top) < precedence(currChar) {
					break
				}
				popped, _ := stack.Pop()
				result = append(result, popped)
			}
			stack.Push(currChar)
		}
	}
	for !stack.IsEmpty() {
		popped, _ := stack.Pop()
		result = append(result, popped)
	}
	regex.Chars = result
}
