package automata

import (
	"fmt"
	"sort"
)

// Minimize reduces the DFA to its minimal equivalent using Hopcroft's
// partition-refinement algorithm. States that are indistinguishable —
// same acceptance class and identical transition targets for every symbol —
// are merged into a single representative state.
func (dfa *DFA) Minimize() *DFA {
	alphabet := dfa.Alphabet()

	// Initial partitioning: non-accepting states share one group; each distinct
	// TokenID gets its own group so patterns are never collapsed into each other.
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
	// O(1) lookup: state ID → partition index
	stateToPartition := make(map[int]int)
	for idx, partition := range partitions {
		for _, s := range partition {
			stateToPartition[s.ID] = idx
		}
	}

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

	// Sort alphabet for deterministic transitions
	var sortedAlpha []rune
	for r := range alphabet {
		sortedAlpha = append(sortedAlpha, r)
	}
	sort.Slice(sortedAlpha, func(i, j int) bool { return sortedAlpha[i] < sortedAlpha[j] })

	for partitionIndex, partition := range partitions {
		newState := partitionToState[partitionIndex]
		for _, state := range partition {
			for _, symbol := range sortedAlpha {
				if nextState, exists := state.Transitions[symbol]; exists {
					nextPartition := stateToPartition[nextState.ID]
					newState.AddTransition(partitionToState[nextPartition], symbol)
				}
			}
		}
	}
	startPartition := stateToPartition[dfa.StartState.ID]
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

// minimizeRecursive performs one refinement pass over the current partitions.
// States within a group are re-split if they disagree on which partition a
// symbol leads to. Groups whose members all produce the same transition hash
// stay together; diverging members form new sub-groups.
func minimizeRecursive(alphabet map[rune]bool, partitions [][]*DFAState) [][]*DFAState {
	newPartitions := [][]*DFAState{}
	// Map state ID → its current partition index for O(1) lookups in transitionHash.
	partitionLookup := make(map[int]int)
	for index, partition := range partitions {
		for _, state := range partition {
			partitionLookup[state.ID] = index
		}
	}
	for _, partition := range partitions {
		// Group states by their transition signature — states with the same hash
		// land in the same partition after this pass.
		hashToStates := make(map[string][]*DFAState)
		for _, state := range partition {
			hash := transitionHash(state, alphabet, partitionLookup)
			hashToStates[hash] = append(hashToStates[hash], state)
		}
		var hashes []string
		for h := range hashToStates {
			hashes = append(hashes, h)
		}
		sort.Strings(hashes)
		for _, h := range hashes {
			newPartitions = append(newPartitions, hashToStates[h])
		}
	}
	return newPartitions
}

// transitionHash encodes a state's outgoing transitions as a string so that
// states with identical behaviour produce the same hash and can be merged.

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
