package ds

import "errors"

// Stack implementation, LIFO
type Stack[T any] []T

// Push a new element to the top of the stack
func (s *Stack[T]) Push(v T) {
	*s = append(*s, v)
}

// Pop the top element from the stack, returns element's value
func (s *Stack[T]) Pop() (T, error) {
	// Check if the stack is empty
	size := len(*s)
	// Return false if stack is empty
	if size == 0 {
		var zeroValue T
		return zeroValue, errors.New("Stack Vacio")
	}
	// Pop the top element if stack is not empty
	top := (*s)[size-1]
	*s = (*s)[:size-1] // Remove the top element
	return top, nil
}

// Peek the top element from the stack, returns element's value
func (s *Stack[T]) Peek() (T, bool) {
	// Check if the stack is empty
	size := len(*s)
	// Return false if stack is empty
	if size == 0 {
		var zeroValue T
		return zeroValue, false
	}
	// Peek the top element if stack is not empty
	top := (*s)[size-1]
	return top, true
}

// Checks if the stack is empty
func (s *Stack[T]) IsEmpty() bool {
	return len(*s) == 0
}
