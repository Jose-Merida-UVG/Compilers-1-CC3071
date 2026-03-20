package regex

import (
	"github.com/Jose-Merida-UVG/Compilers-1-CC3071/internal/ds"
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
		switch currChar.Kind {
		case Lit:
			data := ASTNodeData{
				Value:    currChar.Value,
				Nullable: false,
				FirstPos: []int{idCounter},
				LastPos:  []int{idCounter},
				ID:       idCounter,
			}
			idCounter++
			stack.Push(ds.Node[ASTNodeData]{Value: data, Left: nil, Right: nil})

		case Eps:
			// Epsilon leaf: nullable, no position ID
			data := ASTNodeData{
				Value:    RuneEpsilon,
				Nullable: true,
				FirstPos: []int{},
				LastPos:  []int{},
			}
			stack.Push(ds.Node[ASTNodeData]{Value: data, Left: nil, Right: nil})

		case Cat:
			right, _ := stack.Pop()
			left, _ := stack.Pop()
			stack.Push(ds.Node[ASTNodeData]{Value: ASTNodeData{Value: '~'}, Left: &left, Right: &right})

		case Alt:
			right, _ := stack.Pop()
			left, _ := stack.Pop()
			stack.Push(ds.Node[ASTNodeData]{Value: ASTNodeData{Value: '|'}, Left: &left, Right: &right})

		case Star:
			node, _ := stack.Pop()
			stack.Push(ds.Node[ASTNodeData]{Value: ASTNodeData{Value: '*'}, Left: &node, Right: nil})
		}
	}

	rootNode, _ := stack.Pop()
	return rootNode
}

