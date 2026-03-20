package automata

import (
	"fmt"
	"sort"

	"github.com/Jose-Merida-UVG/Compilers-1-CC3071/internal/ds"
	"github.com/Jose-Merida-UVG/Compilers-1-CC3071/internal/regex"
)

type DirectTable struct {
	NextPos    map[int]map[int]bool
	PosToChar  map[int]rune
	StartState []int
	AcceptID   int
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

func computeAttributes(node *ds.Node[regex.ASTNodeData], table *DirectTable) {
	if node == nil {
		return
	}
	computeAttributes(node.Left, table)
	computeAttributes(node.Right, table)

	if node.Left == nil && node.Right == nil {
		if node.Value.Value != regex.RuneEpsilon {
			table.PosToChar[node.Value.ID] = node.Value.Value
			if node.Value.Value == regex.RuneEndMarker {
				table.AcceptID = node.Value.ID
			}
		}
		return
	} else if node.Value.Value == '|' {
		node.Value.Nullable = node.Left.Value.Nullable || node.Right.Value.Nullable
		node.Value.FirstPos = union(node.Left.Value.FirstPos, node.Right.Value.FirstPos)
		node.Value.LastPos = union(node.Left.Value.LastPos, node.Right.Value.LastPos)
	} else if node.Value.Value == '~' {
		node.Value.Nullable = node.Left.Value.Nullable && node.Right.Value.Nullable
		if node.Left.Value.Nullable {
			node.Value.FirstPos = union(node.Left.Value.FirstPos, node.Right.Value.FirstPos)
		} else {
			node.Value.FirstPos = copySlice(node.Left.Value.FirstPos)
		}
		if node.Right.Value.Nullable {
			node.Value.LastPos = union(node.Left.Value.LastPos, node.Right.Value.LastPos)
		} else {
			node.Value.LastPos = copySlice(node.Right.Value.LastPos)
		}
		for _, posInLastLeft := range node.Left.Value.LastPos {
			if _, exists := table.NextPos[posInLastLeft]; !exists {
				table.NextPos[posInLastLeft] = make(map[int]bool)
			}
			for _, posInFirstRight := range node.Right.Value.FirstPos {
				table.NextPos[posInLastLeft][posInFirstRight] = true
			}
		}
	} else if node.Value.Value == '*' {
		node.Value.Nullable = true
		node.Value.FirstPos = copySlice(node.Left.Value.FirstPos)
		node.Value.LastPos = copySlice(node.Left.Value.LastPos)
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

func (table *DirectTable) ToDFA() *DFA {
	stateMap := make(map[string]*DFAState)

	setKey := func(set []int) string {
		cpy := make([]int, len(set))
		copy(cpy, set)
		sort.Ints(cpy)
		key := ""
		for _, v := range cpy {
			key += fmt.Sprintf("%d,", v)
		}
		return key
	}

	contains := func(slice []int, value int) bool {
		for _, v := range slice {
			if v == value {
				return true
			}
		}
		return false
	}

	queue := [][]int{}
	startSet := table.StartState
	startKey := setKey(startSet)
	startState := NewDFAState()
	if contains(startSet, table.AcceptID) {
		startState.SetToken(1)
	}
	stateMap[startKey] = startState
	queue = append(queue, startSet)

	for len(queue) > 0 {
		currentSet := queue[0]
		queue = queue[1:]
		currentKey := setKey(currentSet)
		currentDFAState := stateMap[currentKey]

		alphabet := make(map[rune]bool)
		for _, char := range table.PosToChar {
			if char != regex.RuneEndMarker {
				alphabet[char] = true
			}
		}

		for symbol := range alphabet {
			newSetMap := make(map[int]bool)
			for _, pos := range currentSet {
				if table.PosToChar[pos] == symbol {
					for next := range table.NextPos[pos] {
						newSetMap[next] = true
					}
				}
			}
			if len(newSetMap) == 0 {
				continue
			}
			var newSet []int
			for id := range newSetMap {
				newSet = append(newSet, id)
			}
			newKey := setKey(newSet)
			if _, exists := stateMap[newKey]; !exists {
				newState := NewDFAState()
				if contains(newSet, table.AcceptID) {
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
