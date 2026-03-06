package regex

import (
	"github.com/TonitoMC/TDLCGo/internal/ds"
)

// Checks if the string has valid parenthesis
func Balanced(rs RegexString) bool {
	// Define the open and closed parenthesis, brackets and braces
	charMap := map[rune]rune{')': '(', ']': '[', '}': '{'}
	stack := ds.Stack[rune]{}
	// Iterate over the string
	for _, char := range rs.Chars {
		matchingChar, isOperator := charMap[char.Value]
		// Handle closing parenthesis
		if isOperator && !char.escaped {
			// Pop the stack
			poppedChar, valid := stack.Pop()
			if valid != nil {
				// Return false if not matching
				return false
				// Return false if empty stack
			}
			if poppedChar != matchingChar {
				return false
			}
			// Add open parenthesis to stack
		} else if (char.Value == '(' || char.Value == '[' || char.Value == '{') && !char.escaped {
			stack.Push(char.Value)
		}
	}
	return stack.IsEmpty()
}
