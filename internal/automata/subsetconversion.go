package automata

// Returns a map with every symbol contained within the transitions
// of a NFA
func (NFA *NFA) Alphabet() map[rune]bool {
	visitedStates := NewStateSet()
	alphabet := make(map[rune]bool)
	AlphabetRecursive(NFA.StartState, visitedStates, alphabet)
	return alphabet
}

// Recursive function to calculate Alphabet, adds State to
// visited stateSet, adds symbols for which there exists transition
// to alphabet
func AlphabetRecursive(currentState *State, visitedStates *StateSet, alphabet map[rune]bool) {
	if visitedStates.Contains(currentState) {
		return
	}
	visitedStates.Add(currentState)
	for symbol, nextState := range currentState.Transitions {
		if _, ok := alphabet[symbol]; !ok {
			alphabet[symbol] = true
		}
		AlphabetRecursive(nextState, visitedStates, alphabet)
	}
	for _, nextState := range currentState.EpsilonTransitions {
		AlphabetRecursive(nextState, visitedStates, alphabet)
	}
}

func (nfa *NFA) ToDFA() *DFA {
	alphabet := nfa.Alphabet()
	existingStates := make(map[*StateSet]*DFAState)
	initialState := nfa.StartState.EpsilonClosure()
	DFA := NewDFA(ToDFARecursive(initialState, alphabet, existingStates))
	for stateSet, dfaState := range existingStates {
		if stateSet.Contains(nfa.AcceptState) {
			dfaState.SetAccepting()
		}
	}
	return DFA
}

// For every new StateSet Calculated
// Check if it matches any existing state
// If it doesn't
// Create new state and add transition
func ToDFARecursive(currentStateSet *StateSet, alphabet map[rune]bool, existingStates map[*StateSet]*DFAState) *DFAState {
	// Verify if state already exists, if it does it returns the instance
	// of the already created state
	for existingStateSet := range existingStates {
		if existingStateSet.CompareTo(currentStateSet) {
			return existingStates[existingStateSet]
		}
	}
	// The state we have is not in the StateSet, we create a new state corresponding to our
	// StateSet and add it to the existingStates map
	existingStates[currentStateSet] = NewDFAState()
	// Now we work through the transitions for it
	for symbol := range alphabet {
		nextStateSet := currentStateSet.SingleTransition(symbol)
		if !nextStateSet.IsEmpty() {
			nextState := ToDFARecursive(nextStateSet, alphabet, existingStates)
			existingStates[currentStateSet].AddTransition(nextState, symbol)
		}
	}
	return existingStates[currentStateSet]
}
