package ds

import "errors"

// Stack implementation, LIFO
type Stack[T any] []T

func (s *Stack[T]) Push(v T) {
	*s = append(*s, v)
}

func (s *Stack[T]) Pop() (T, error) {
	size := len(*s)
	if size == 0 {
		var zeroValue T
		return zeroValue, errors.New("Stack Vacio")
	}
	top := (*s)[size-1]
	*s = (*s)[:size-1]
	return top, nil
}

func (s *Stack[T]) Peek() (T, bool) {
	size := len(*s)
	if size == 0 {
		var zeroValue T
		return zeroValue, false
	}
	top := (*s)[size-1]
	return top, true
}

func (s *Stack[T]) IsEmpty() bool {
	return len(*s) == 0
}
