package automata

import (
	"fmt"
	"sort"
)

// Minimization for DFA
func (dfa *DFA) Minimize() *DFA {
	// fmt.Println("Starting DFA minimization...")

	alphabet := dfa.Alphabet()
	partitions := [][]*DFAState{}
	// Initial partition creation, division between acceptance and non-acceptance states
	acceptStates := []*DFAState{}
	nonAcceptStates := []*DFAState{}

	for _, state := range dfa.GetAllStates() {
		if state.Accept {
			acceptStates = append(acceptStates, state)
		} else {
			nonAcceptStates = append(nonAcceptStates, state)
		}
	}

	partitions = append(partitions, nonAcceptStates)
	partitions = append(partitions, acceptStates)

	iteration := 0
	for {
		newPartitions := MinimizeRecursive(alphabet, partitions)
		if partitionsAreEqual(partitions, newPartitions) {
			break
		}
		partitions = newPartitions
		iteration++
	}

	return dfa.DFAFromPartitions(partitions, alphabet)
}

func (dfa *DFA) DFAFromPartitions(partitions [][]*DFAState, alphabet map[rune]bool) *DFA {
	partitionToState := make(map[int]*DFAState)
	for index, partition := range partitions {
		newState := NewDFAState()

		for _, state := range partition {
			if state.Accept {
				newState.SetAccepting()
				break
			}
		}
		partitionToState[index] = newState
	}

	for partitionIndex, partition := range partitions {
		newState := partitionToState[partitionIndex]
		for _, state := range partition {
			for symbol := range alphabet {
				if nextState, exists := state.Transitions[symbol]; exists {
					nextPartition := findPartition(nextState, partitions)
					newState.AddTransition(partitionToState[nextPartition], symbol)
				}
			}
		}
	}

	startPartition := findPartition(dfa.StartState, partitions)
	return NewDFA(partitionToState[startPartition])
}

// Method to find partition based on state
func findPartition(state *DFAState, partitions [][]*DFAState) int {
	for partitionIndex, partition := range partitions {
		for _, s := range partition {
			if s.ID == state.ID {
				return partitionIndex
			}
		}
	}
	return -1 // Should not happen if the state is valid
}

// Method to determine if partitions are equal
func partitionsAreEqual(p1, p2 [][]*DFAState) bool {
	if len(p1) != len(p2) {
		return false
	}

	// Since the contents of the partitions aren't interchanged and simply split
	// we can compare two partitions simply by the length of each of each
	// partition
	for i := range p1 {
		if len(p1[i]) != len(p2[i]) {
			return false
		}
	}
	return true
}

func MinimizeRecursive(alphabet map[rune]bool, partitions [][]*DFAState) [][]*DFAState {
	newPartitions := [][]*DFAState{}
	partitionLookup := make(map[int]int)

	for index, partition := range partitions {
		for _, state := range partition {
			partitionLookup[state.ID] = index
		}
	}

	// Refine each partition based on transition behavior
	for _, partition := range partitions {
		hashToStates := make(map[string][]*DFAState)

		for _, state := range partition {
			hash := state.TransitionHash(alphabet, partitionLookup)
			hashToStates[hash] = append(hashToStates[hash], state)
		}

		for _, states := range hashToStates {
			newPartitions = append(newPartitions, states)
		}
	}

	return newPartitions
}

func (dfa *DFA) GetAllStates() []*DFAState {
	visited := make(map[int]bool)
	allStates := []*DFAState{}
	GetAllStatesRecursive(dfa.StartState, &allStates, visited)
	return allStates
}

func GetAllStatesRecursive(state *DFAState, allStates *[]*DFAState, visited map[int]bool) {
	if visited[state.ID] {
		return
	}
	visited[state.ID] = true
	*allStates = append(*allStates, state)
	for _, nextState := range state.Transitions {
		GetAllStatesRecursive(nextState, allStates, visited)
	}
}

func (state *DFAState) TransitionHash(alphabet map[rune]bool, partitionTable map[int]int) string {
	var symbols []rune

	// Collect all symbols from the alphabet map
	for symbol := range alphabet {
		symbols = append(symbols, symbol)
	}

	// Sort the symbols to ensure consistent order
	sort.Slice(symbols, func(i, j int) bool {
		return symbols[i] < symbols[j]
	})

	hash := ""
	for _, symbol := range symbols {
		if nextState, exists := state.Transitions[symbol]; exists {
			partitionIndex := partitionTable[nextState.ID]
			hash += fmt.Sprintf("%d", partitionIndex)
		} else {
			hash += "X" // No transition for this symbol, use a placeholder
		}
	}

	return hash
}
