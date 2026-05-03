package automata

import (
	"fmt"
	"slices"
	"sort"
)

// Minimize reduces the DFA to its minimal equivalent using Hopcroft's
// partition-refinement algorithm. States that are indistinguishable
// are merged into a single representative state.
func (dfa *DFA) Minimize() *DFA {
	alphabet := dfa.Alphabet()

	// Initial partitioning: non-accepting states share one group & each distinct
	// TokenID gets its own group so patterns are never collapsed into each other.

	// Create token groups
	tokenGroups := make(map[int][]*DFAState)
	for _, state := range dfa.GetAllStates() {
		tokenGroups[state.TokenID] = append(tokenGroups[state.TokenID], state)
	}

	// Append non-accept states (token = 0)
	var partitions [][]*DFAState
	if g, ok := tokenGroups[0]; ok {
		partitions = append(partitions, g)
	}
	// Append token groups
	var tokenIDs []int
	for id := range tokenGroups {
		if id > 0 {
			tokenIDs = append(tokenIDs, id)
		}
	}
	// Sort for deterministic outcome
	sort.Ints(tokenIDs)
	for _, id := range tokenIDs {
		partitions = append(partitions, tokenGroups[id])
	}
	// Minimization loop
	for {
		newPartitions := minimizeRecursive(alphabet, partitions)
		if partitionsAreEqual(partitions, newPartitions) {
			break
		}
		partitions = newPartitions
	}
	return dfa.dfaFromPartitions(partitions, alphabet)
}

// dfaFromPartitions generates a DFA from the partitions & alphabet of the DFA,
// it iterates through and creates a state + its transitions based on a representative
// from each equivalence class.
func (dfa *DFA) dfaFromPartitions(partitions [][]*DFAState, alphabet map[rune]bool) *DFA {
	// Lookup table for state -> partition_idx
	stateToPartition := make(map[int]int)
	for idx, partition := range partitions {
		for _, s := range partition {
			stateToPartition[s.ID] = idx
		}
	}

	// Create state per partition
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

	// Sort alphabet in slice for deterministic transitions
	var sortedAlpha []rune
	for r := range alphabet {
		sortedAlpha = append(sortedAlpha, r)
	}
	slices.Sort(sortedAlpha)

	// For every partition / eq. new state representing the partition
	for partitionIndex, partition := range partitions {
		newState := partitionToState[partitionIndex]
		// Pick representative and add its transitions to the state
		representative := partition[0]
		for _, symbol := range sortedAlpha {
			if nextState, exists := representative.Transitions[symbol]; exists {
				nextPartition := stateToPartition[nextState.ID]
				newState.AddTransition(partitionToState[nextPartition], symbol)
			}
		}
	}
	startPartition := stateToPartition[dfa.StartState.ID]
	return NewDFA(partitionToState[startPartition])
}

// partitionsAreEqual verifies whether two partitions are the same,
// in theory we could just compare their lengths but we're comparing
// lengths of each individual subset just to make sure.
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
// stay together while diverging members form new subsets.
func minimizeRecursive(alphabet map[rune]bool, partitions [][]*DFAState) [][]*DFAState {
	newPartitions := [][]*DFAState{}
	// Map state ID → its current partition index
	partitionLookup := make(map[int]int)
	for index, partition := range partitions {
		for _, state := range partition {
			partitionLookup[state.ID] = index
		}
	}
	// For every partition
	for _, partition := range partitions {
		// Group states by their transition signature, states with the same hash
		// land in the same partition after this pass.
		hashToStates := make(map[string][]*DFAState)
		// For every state
		for _, state := range partition {
			// Generate hash & store in map
			hash := transitionHash(state, alphabet, partitionLookup)
			hashToStates[hash] = append(hashToStates[hash], state)
		}
		var hashes []string
		for h := range hashToStates {
			hashes = append(hashes, h)
		}
		// Generate new partitions based on hash
		sort.Strings(hashes)
		for _, h := range hashes {
			newPartitions = append(newPartitions, hashToStates[h])
		}
	}
	return newPartitions
}

// GetAllStates is used to get all the states in a DFA
func (dfa *DFA) GetAllStates() []*DFAState {
	visited := make(map[int]bool)
	allStates := []*DFAState{}
	getAllStatesRecursive(dfa.StartState, &allStates, visited)
	return allStates
}

// getAllStatesRecursive is the function call in GetAllStates that iterates through
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

// transitionHash encodes a state's outgoing transitions as a string so that
// states with identical behaviour produce the same hash and can be merged.
func transitionHash(state *DFAState, alphabet map[rune]bool, partitionTable map[int]int) string {
	var symbols []rune
	for symbol := range alphabet {
		symbols = append(symbols, symbol)
	}
	slices.Sort(symbols)
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
