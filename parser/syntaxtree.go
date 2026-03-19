package parser

import (
	"github.com/Jose-Merida-UVG/Compilers-1-CC3071/ds"
)

// ASTNodeData holds the data needed for direct DFA conversion algorithms.
type ASTNodeData struct {
	Value    rune
	ID       int // Position identifier for leaf nodes
	Nullable bool
	FirstPos []int
	LastPos  []int
}

// BuildDirectAST constructs the syntax tree for direct conversion,
// values for leaf nodes are gathered here & operation nodes just
// follow the same rules as a normal AST + extra empty fields
func BuildDirectAST(regex *RegexString) ds.Node[ASTNodeData] {
	// This stack keeps track of our nodes
	stack := ds.Stack[ds.Node[ASTNodeData]]{}
	// This counter helps us enumerate our leaf nodes
	idCounter := 1

	for _, currChar := range regex.Chars {
		// Non operators pushed into stack
		if !currChar.IsOperator() {
			// Initialize empty data
			data := ASTNodeData{
				Value:    currChar.Value,
				Nullable: currChar.Value == 'ε',
				FirstPos: []int{},
				LastPos:  []int{},
			}

			// Initialize ID + FirstPos & LastPos for non-epsilons
			if currChar.Value != 'ε' {
				data.ID = idCounter
				data.FirstPos = []int{idCounter}
				data.LastPos = []int{idCounter}
				idCounter++
			}

			// Push node  into stack
			stack.Push(ds.Node[ASTNodeData]{Value: data, Left: nil, Right: nil})

			// Handle concatenation
		} else if currChar.Value == '~' {
			right, _ := stack.Pop()
			left, _ := stack.Pop()
			data := ASTNodeData{Value: '~'}
			stack.Push(ds.Node[ASTNodeData]{Value: data, Left: &left, Right: &right})

			// Handle OR
		} else if currChar.Value == '|' {
			right, _ := stack.Pop()
			left, _ := stack.Pop()
			data := ASTNodeData{Value: '|'}
			stack.Push(ds.Node[ASTNodeData]{Value: data, Left: &left, Right: &right})

			// Handle Kleene Star
		} else if currChar.Value == '*' {
			node, _ := stack.Pop()
			data := ASTNodeData{Value: '*'}
			stack.Push(ds.Node[ASTNodeData]{Value: data, Left: &node, Right: nil})
		}
	}

	// Return root node
	rootNode, _ := stack.Pop()
	return rootNode
}

// Builds AST for RegexString
func BuildAST(regex *RegexString) ds.Node[rune] {
	// Initialize stack of nodes with value rune
	stack := ds.Stack[ds.Node[rune]]{}
	for _, currChar := range regex.Chars {
		// Non operators pushed into stack
		if !currChar.IsOperator() {
			stack.Push(ds.Node[rune]{Value: currChar.Value, Left: nil, Right: nil})
			// Handle concatenation
		} else if currChar.Value == '~' {
			right, _ := stack.Pop()
			left, _ := stack.Pop()
			stack.Push(ds.Node[rune]{Value: '~', Left: &left, Right: &right})
			// Handle OR
		} else if currChar.Value == '|' {
			right, _ := stack.Pop()
			left, _ := stack.Pop()
			stack.Push(ds.Node[rune]{Value: '|', Left: &left, Right: &right})
		} else if currChar.Value == '*' {
			node, _ := stack.Pop()
			stack.Push(ds.Node[rune]{Value: '*', Left: &node, Right: nil})
		}
	}
	rootNode, _ := stack.Pop()
	return rootNode
}
