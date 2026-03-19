package parser

import (
	"slices"

	"github.com/Jose-Merida-UVG/Compilers-1-CC3071/ds"
)

// Shunting Yard algorithm is used to turn expressions
// from infix to postfix. This implementation works
// with a RegexString to work with regular expressions,
// simplifying them to only '*', '|' and '~' (concat)
// operators

// Handles special character classes defined by [abc] -> (a|b|c)
func (regex *RegexString) HandleCharClasses() {
	result := []Token{}
	for i := 0; i < len(regex.Chars); i++ {
		currChar := regex.Chars[i]
		// If we encounter an opening bracket
		if currChar.Value == '[' && !currChar.escaped {
			result = append(result, Token{Value: '(', escaped: false})
			// Go through every character until we find a closing bracket
			for j := i + 1; j < len(regex.Chars); j++ {
				classChar := regex.Chars[j]
				if classChar.Value == ']' && !classChar.escaped {
					result = append(result, Token{Value: ')', escaped: false})
					i = j // Move the outer loop index to the closing bracket
					break
				} else {
					if j > i+1 {
						result = append(result, Token{Value: '|', escaped: false})
					}
					result = append(result, classChar)
				}
			}
		} else {
			result = append(result, currChar)
		}
	}
	regex.Chars = result
}

// Utility function to find opening parenthesis
func findOpenPar(regex []Token, closeIndex int) int {
	parenCount := 1
	for i := closeIndex - 1; i >= 0; i-- {
		currChar := regex[i]
		if currChar.Value == ')' && !currChar.escaped {
			parenCount++
		} else if currChar.Value == '(' && !currChar.escaped {
			parenCount--
			if parenCount == 0 {
				return i
			}
		}
	}
	return -1
}

// Handles special operators '?' and '+' ex.
// a? -> (a|ε)
// a+ -> a(a)*
// (a|b)+ -> (a|b)(a|b)*
// (0|1)? -> ((0|1)|ε)
func (regex *RegexString) HandleSpecialOperators() {
	result := []Token{}

	for i := 0; i < len(regex.Chars); i++ {
		currChar := regex.Chars[i]
		// If there exist a next character
		if i+1 < len(regex.Chars) {
			nextChar := regex.Chars[i+1]
			// Next char is '?' operator
			if !nextChar.escaped && nextChar.Value == '?' {
				if currChar.Value == ')' {
					// Inserts opening parenthesis at corresponding index & closes expression
					result = append(result, Token{Value: ')', escaped: false})
					openParIndex := findOpenPar(result, len(result)-1)
					result = slices.Insert(result, openParIndex, Token{Value: '(', escaped: false})
					// Append |ε) to enclose expression
					result = append(result, Token{Value: '|', escaped: false})
					result = append(result, Token{Value: 'ε', escaped: false})
					result = append(result, Token{Value: ')', escaped: false})
				} else {
					result = append(result, Token{Value: '(', escaped: false})
					result = append(result, currChar)
					result = append(result, Token{Value: '|', escaped: false})
					result = append(result, Token{Value: 'ε', escaped: false})
					result = append(result, Token{Value: ')', escaped: false})
				}
				i++ // Skip the next character as it is already processed
			} else if !nextChar.escaped && nextChar.Value == '+' {
				if currChar.Value == ')' {
					// Find the corresponding opening parenthesis
					result = append(result, Token{Value: ')', escaped: false})
					openParIndex := findOpenPar(result, len(result)-1)
					// Append the repeated expression
					resultLen := len(result) - 1
					for j := openParIndex; j <= resultLen; j++ {
						result = append(result, result[j])
					}
					result = append(result, Token{Value: '*', escaped: false})
				} else {
					// Handle the '+' operator for single characters
					result = append(result, currChar)
					result = append(result, currChar)
					result = append(result, Token{Value: '*', escaped: false})
				}
				i++ // Skip the next character as it is already processed
			} else {
				result = append(result, currChar)
			}
		} else {
			result = append(result, currChar)
		}
	}

	regex.Chars = result
}

// Handles explicit concatenation, using '~'
func (regex *RegexString) HandleExplicitConcatenation() {
	result := []Token{}
	for i, currChar := range regex.Chars {
		result = append(result, currChar)

		// Ensure i+1 is within bounds
		if i+1 < len(regex.Chars) {
			nextChar := regex.Chars[i+1]

			// currChar != '(', '|' or '~'
			isCurrCharSpecial := !currChar.escaped && (currChar.Value == '(' || currChar.Value == '|' || currChar.Value == '~')
			// nextChar != ')', '|', '*'
			isNextCharSpecial := !nextChar.escaped && (nextChar.Value == ')' || nextChar.Value == '|' || nextChar.Value == '*')

			if !isCurrCharSpecial && !isNextCharSpecial {
				result = append(result, Token{Value: '~', escaped: false})
			}
		}
	}
	regex.Chars = result
}

// Operator precedence
func precedence(op rune) int {
	switch op {
	case '(':
		return 1
	case '|':
		return 2
	case '~':
		return 3
	case '*':
		return 4
	default:
		return 10 // Default precedence for unknown characters
	}
}

// ShuntingYard converts a regex string to postfix notation
func (regex *RegexString) ShuntingYard() {
	result := []Token{}
	stack := ds.Stack[Token]{}

	for _, currChar := range regex.Chars {
		// Non-operators go to output
		if !currChar.IsOperator() {
			result = append(result, currChar)
		} else if currChar.Value == '(' {
			stack.Push(currChar)
		} else if currChar.Value == ')' {
			// Checks whether stack is empty or '(' is encountered
			for !stack.IsEmpty() {
				top, _ := stack.Peek()
				if top.Value == '(' {
					break
				}
				popped, _ := stack.Pop()
				result = append(result, popped)
			}
			stack.Pop() // Pop the '(' from the stack
		} else { // Handle operators
			for !stack.IsEmpty() {
				top, _ := stack.Peek()
				if top.Value == '(' || precedence(top.Value) < precedence(currChar.Value) {
					break
				}
				popped, _ := stack.Pop()
				result = append(result, popped)
			}
			stack.Push(currChar)
		}
	}

	// Pop all the remaining operators in the stack
	for !stack.IsEmpty() {
		popped, _ := stack.Pop()
		result = append(result, popped)
	}

	regex.Chars = result
}
