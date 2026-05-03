package automata

import (
	"fmt"
	"slices"

	"github.com/Jose-Merida-UVG/Compilers-1-CC3071/internal/ds"
	"github.com/Jose-Merida-UVG/Compilers-1-CC3071/internal/yalex/regex"
)

// DirectTable holds the followpos data computed from the augmented regex AST.
type DirectTable struct {
	// Nextpos mapping ID -> set of IDs
	NextPos map[int]map[int]bool
	// PosToChar mapping ID -> rune
	PosToChar map[int]rune
	// Starting state for the DFA
	StartState []int
	// ID of sentinel
	AcceptID int
}

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
// post-order traversal. It also populates the DirectTable with NextPos and
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
		if node.Value.Value != regex.RuneEpsilon {
			table.PosToChar[node.Value.ID] = node.Value.Value
			// If it's the end-marker sentinel, store its ID as the accept position
			if node.Value.Value == regex.RuneEndMarker {
				table.AcceptID = node.Value.ID
			}
		}
		return

		// Handle OR
	} else if node.Value.Value == '|' {
		// Nullable = Nullable(L) || Nullable(R)
		node.Value.Nullable = node.Left.Value.Nullable || node.Right.Value.Nullable
		// FirstPos = U(L.firstPos, R.firstPos)
		node.Value.FirstPos = union(node.Left.Value.FirstPos, node.Right.Value.FirstPos)
		// LastPos = U(L.lastPos, R.lastPos)
		node.Value.LastPos = union(node.Left.Value.LastPos, node.Right.Value.LastPos)

	} else if node.Value.Value == '~' {
		// Concatenation operator

		// Nullable = Nullable(L) && Nullable(R)
		node.Value.Nullable = node.Left.Value.Nullable && node.Right.Value.Nullable

		// FirstPos depends on nullability of left child
		if node.Left.Value.Nullable {
			// FirstPos = U(L.firstPos, R.firstPos)
			node.Value.FirstPos = union(node.Left.Value.FirstPos, node.Right.Value.FirstPos)
		} else {
			// FirstPos = L.firstPos
			node.Value.FirstPos = copySlice(node.Left.Value.FirstPos)
		}

		// LastPos depends on nullability of right child
		if node.Right.Value.Nullable {
			// LastPos = U(L.lastPos, R.lastPos)
			node.Value.LastPos = union(node.Left.Value.LastPos, node.Right.Value.LastPos)
		} else {
			// LastPos = R.lastPos
			node.Value.LastPos = copySlice(node.Right.Value.LastPos)
		}

		// Update NextPos for concatenation:
		// for every position in LastPos(L), every position in FirstPos(R) is a nextpos
		for _, posInLastLeft := range node.Left.Value.LastPos {
			if _, exists := table.NextPos[posInLastLeft]; !exists {
				table.NextPos[posInLastLeft] = make(map[int]bool)
			}
			for _, posInFirstRight := range node.Right.Value.FirstPos {
				table.NextPos[posInLastLeft][posInFirstRight] = true
			}
		}

	} else if node.Value.Value == '*' {
		// Kleene star operator

		// Nullable = true by definition
		node.Value.Nullable = true
		// FirstPos = L.firstPos
		node.Value.FirstPos = copySlice(node.Left.Value.FirstPos)
		// LastPos = L.lastPos
		node.Value.LastPos = copySlice(node.Left.Value.LastPos)

		// Update NextPos for Kleene star:
		// for every position in LastPos(L), every position in FirstPos(L) is a nextpos
		// (the loop back to the start of the repeated sub-expression)
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

func copySlice(a []int) []int {
	res := make([]int, len(a))
	copy(res, a)
	return res
}

// ToDFA runs subset construction on the followpos table to produce a DFA.
// Each DFA state corresponds to a set of leaf positions; a state is accepting
// when its position set contains the end-marker position (AcceptID).
func (table *DirectTable) ToDFA() *DFA {
	// stateMap: position-set key -> DFA state.
	// Each DFA state represents a unique set of leaf positions (e.g., "{1,3,5}").
	stateMap := make(map[string]*DFAState)

	// setKey: canonical string key for a position set.
	// Sorts before encoding so {3,1} and {1,3} map to the same key.
	setKey := func(set []int) string {
		cpy := make([]int, len(set))
		copy(cpy, set)
		slices.Sort(cpy)
		key := ""
		for _, v := range cpy {
			key += fmt.Sprintf("%d,", v)
		}
		return key
	}

	// Seed the BFS with the start state (firstpos of the root).
	// startSet is the initial position set; we map it to the DFA's start state.
	queue := [][]int{}
	startSet := table.StartState
	startKey := setKey(startSet)
	startState := NewDFAState()
	// A position set is accepting if it contains the end-marker position (AcceptID).
	if slices.Contains(startSet, table.AcceptID) {
		startState.SetToken(1)
	}
	stateMap[startKey] = startState
	queue = append(queue, startSet)

	// Collect every character that appears in the regex (skip the end-marker sentinel).
	// Sorted for deterministic transition order.
	alphabet := make(map[rune]bool)
	for _, char := range table.PosToChar {
		if char != regex.RuneEndMarker {
			alphabet[char] = true
		}
	}
	var sortedAlpha []rune
	for r := range alphabet {
		sortedAlpha = append(sortedAlpha, r)
	}
	slices.Sort(sortedAlpha)

	// BFS: each iteration expands one position set into its successor sets.
	for len(queue) > 0 {
		currentSet := queue[0]
		queue = queue[1:]
		currentKey := setKey(currentSet)
		currentDFAState := stateMap[currentKey]

		for _, symbol := range sortedAlpha {
			// move(currentSet, symbol): union of followpos for every position
			// in the current set whose character matches symbol.
			newSetMap := make(map[int]bool)
			for _, pos := range currentSet {
				if table.PosToChar[pos] == symbol {
					for next := range table.NextPos[pos] {
						newSetMap[next] = true
					}
				}
			}
			if len(newSetMap) == 0 {
				// No positions match this symbol
				continue
			}
			// Flatten the set map into a slice for key computation.
			var newSet []int
			for id := range newSetMap {
				newSet = append(newSet, id)
			}
			newKey := setKey(newSet)
			// Register the new position set as a DFA state if we haven't seen it yet.
			if _, exists := stateMap[newKey]; !exists {
				newState := NewDFAState()
				if slices.Contains(newSet, table.AcceptID) {
					newState.SetToken(1)
				}
				stateMap[newKey] = newState
				queue = append(queue, newSet)
			}
			currentDFAState.AddTransition(stateMap[newKey], symbol)
		}
	}

	return NewDFA(startState)
}
