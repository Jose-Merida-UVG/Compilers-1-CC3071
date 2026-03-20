package automata

import (
	"fmt"
	"sort"
)

func (dfa *DFA) Minimize() *DFA {
	alphabet := dfa.Alphabet()

	// Initial partitioning: non-accept states (TokenID=0) in one group,
	// each distinct TokenID in its own group (so token classes are never merged).
	tokenGroups := make(map[int][]*DFAState)
	for _, state := range dfa.GetAllStates() {
		tokenGroups[state.TokenID] = append(tokenGroups[state.TokenID], state)
	}
	var partitions [][]*DFAState
	if g, ok := tokenGroups[0]; ok {
		partitions = append(partitions, g)
	}
	var tokenIDs []int
	for id := range tokenGroups {
		if id > 0 {
			tokenIDs = append(tokenIDs, id)
		}
	}
	sort.Ints(tokenIDs)
	for _, id := range tokenIDs {
		partitions = append(partitions, tokenGroups[id])
	}

	for {
		newPartitions := minimizeRecursive(alphabet, partitions)
		if partitionsAreEqual(partitions, newPartitions) {
			break
		}
		partitions = newPartitions
	}
	return dfa.dfaFromPartitions(partitions, alphabet)
}

func (dfa *DFA) dfaFromPartitions(partitions [][]*DFAState, alphabet map[rune]bool) *DFA {
	partitionToState := make(map[int]*DFAState)
	for index, partition := range partitions {
		newState := NewDFAState()
		for _, state := range partition {
			if state.Accept {
				newState.SetToken(state.TokenID)
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

func findPartition(state *DFAState, partitions [][]*DFAState) int {
	for partitionIndex, partition := range partitions {
		for _, s := range partition {
			if s.ID == state.ID {
				return partitionIndex
			}
		}
	}
	return -1
}

func partitionsAreEqual(p1, p2 [][]*DFAState) bool {
	if len(p1) != len(p2) {
		return false
	}
	for i := range p1 {
		if len(p1[i]) != len(p2[i]) {
			return false
		}
	}
	return true
}

func minimizeRecursive(alphabet map[rune]bool, partitions [][]*DFAState) [][]*DFAState {
	newPartitions := [][]*DFAState{}
	partitionLookup := make(map[int]int)
	for index, partition := range partitions {
		for _, state := range partition {
			partitionLookup[state.ID] = index
		}
	}
	for _, partition := range partitions {
		hashToStates := make(map[string][]*DFAState)
		for _, state := range partition {
			hash := transitionHash(state, alphabet, partitionLookup)
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
	getAllStatesRecursive(dfa.StartState, &allStates, visited)
	return allStates
}

func getAllStatesRecursive(state *DFAState, allStates *[]*DFAState, visited map[int]bool) {
	if visited[state.ID] {
		return
	}
	visited[state.ID] = true
	*allStates = append(*allStates, state)
	for _, nextState := range state.Transitions {
		getAllStatesRecursive(nextState, allStates, visited)
	}
}

func transitionHash(state *DFAState, alphabet map[rune]bool, partitionTable map[int]int) string {
	var symbols []rune
	for symbol := range alphabet {
		symbols = append(symbols, symbol)
	}
	sort.Slice(symbols, func(i, j int) bool { return symbols[i] < symbols[j] })
	hash := ""
	for _, symbol := range symbols {
		if nextState, exists := state.Transitions[symbol]; exists {
			hash += fmt.Sprintf("%d", partitionTable[nextState.ID])
		} else {
			hash += "X"
		}
	}
	return hash
}
