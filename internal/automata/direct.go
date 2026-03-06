package automata

import (
	"github.com/TonitoMC/TDLCGo/internal/ds"
	"github.com/TonitoMC/TDLCGo/internal/regex"
)

// This file implements the 'meat & potatoes' of the direct method
// for convering a regular expression into a DFA. It takes a tree
// made in regex/syntaxtree.go (would've been nice to have it here
// in hindsight) and gives us a table with information we need
// to construct the DFA

// DirectTable holds all the information extracted from the syntax tree
// required to build a DFA directly.
type DirectTable struct {
	// NextPos mapping: ID -> Set of IDs. Using map[int]bool for set representation.
	NextPos map[int]map[int]bool
	// PosToChar mapping: ID -> rune (to know what character an ID represents)
	PosToChar map[int]rune
	// The starting state for the DFA (FirstPos of the root node)
	StartState []int
	// The ID of the '#' character. Any DFA state containing this ID is an accept state.
	AcceptID int
}

// BuildDirectTable calculates Nullable, FirstPos, LastPos & NextPos
// for the 'direct' conversion method. Since we're doing a post-order
// traversal of the tree, we can go ahead & calculate every value
// in one run because for every node it's children will have their
// attributes updated already.

// This function also populates the DirectTable with NextPos and PosToChar mappings,
// and sets the initial StartState and AcceptID.
func BuildDirectTable(root *ds.Node[regex.ASTNodeData]) *DirectTable {
	if root == nil {
		return nil
	}

	table := &DirectTable{
		NextPos:   make(map[int]map[int]bool),
		PosToChar: make(map[int]rune),
		AcceptID:  0,
	}

	computeAttributes(root, table)

	table.StartState = root.Value.FirstPos

	return table
}

// computeAttributes calculates Nullable, FirstPos, and LastPos for each node
// in the AST built for direct DFA conversion using a recursive approach for
// post-order trasversal. It also populates the DirectTable with NextPos and
// PosToChar mappings.
func computeAttributes(node *ds.Node[regex.ASTNodeData], table *DirectTable) {
	if node == nil {
		return
	}

	// Recursive call for left child
	computeAttributes(node.Left, table)

	// Recursive call for right child
	computeAttributes(node.Right, table)

	// Base case: we're at a leaf node
	if node.Left == nil && node.Right == nil {
		// Nullable, FirstPos, and LastPos are already correctly initialized
		// during AST construction in BuildDirectAST.
		if node.Value.Value != 'ε' {
			table.PosToChar[node.Value.ID] = node.Value.Value
			// If it's the '#' character, store its ID
			if node.Value.Value == '#' {
				table.AcceptID = node.Value.ID
			}
		}
		return
		// Handle OR
	} else if node.Value.Value == '|' {

		// Nullable = Nullable(L | R)
		node.Value.Nullable = node.Left.Value.Nullable || node.Right.Value.Nullable
		// FirstPos = U(L.firstPos | R.firstPos)
		node.Value.FirstPos = union(node.Left.Value.FirstPos, node.Right.Value.FirstPos)
		// LastPos = U(L.lastPos | R.lastPos)
		node.Value.LastPos = union(node.Left.Value.LastPos, node.Right.Value.LastPos)
	} else if node.Value.Value == '~' {
		// Concatenation Operator

		// Nullable = Nullable (L & R)
		node.Value.Nullable = node.Left.Value.Nullable && node.Right.Value.Nullable

		// Firstpos, depends of nullability
		if node.Left.Value.Nullable {
			// FirstPos = U(L.firstPos, R.firstPos)
			node.Value.FirstPos = union(node.Left.Value.FirstPos, node.Right.Value.FirstPos)
		} else {
			// FirstPos = L.firstPos
			node.Value.FirstPos = copySlice(node.Left.Value.FirstPos)
		}

		// Lastpos, depends on nullability
		if node.Right.Value.Nullable {
			// LastPos = U(L.firstPos, R.firstPos)
			node.Value.LastPos = union(node.Left.Value.LastPos, node.Right.Value.LastPos)
		} else {
			node.Value.LastPos = copySlice(node.Right.Value.LastPos)
		}

		// Update NextPos for concatenation
		// For every position in LastPos(Left child), every position in FirstPos(Right child) is a nextpos.
		for _, posInLastLeft := range node.Left.Value.LastPos {
			if _, exists := table.NextPos[posInLastLeft]; !exists {
				table.NextPos[posInLastLeft] = make(map[int]bool)
			}
			for _, posInFirstRight := range node.Right.Value.FirstPos {
				table.NextPos[posInLastLeft][posInFirstRight] = true
			}
		}

	} else if node.Value.Value == '*' {
		// Kleene Star Operator

		// Nullable = true by default
		node.Value.Nullable = true
		// FirstPos = L.firstPos
		node.Value.FirstPos = copySlice(node.Left.Value.FirstPos)
		// LastPos = L.firstPos
		node.Value.LastPos = copySlice(node.Left.Value.LastPos)

		// Update NextPos for Kleene Star
		// For every position in LastPos(Left child), every position in FirstPos(Left child) is a nextpos.
		for _, posInLastLeft := range node.Left.Value.LastPos {
			if _, exists := table.NextPos[posInLastLeft]; !exists {
				table.NextPos[posInLastLeft] = make(map[int]bool)
			}
			for _, posInFirstLeft := range node.Left.Value.FirstPos {
				table.NextPos[posInLastLeft][posInFirstLeft] = true
			}
		}
	}
}

// union is a helper function to merge two integer slices without duplicates.
func union(a, b []int) []int {
	m := make(map[int]bool)
	for _, item := range a {
		m[item] = true
	}
	for _, item := range b {
		m[item] = true
	}
	var res []int
	for k := range m {
		res = append(res, k)
	}
	return res
}

// copySlice is a utility to avoid mutating original slices by reference.
func copySlice(a []int) []int {
	res := make([]int, len(a))
	copy(res, a)
	return res
}
