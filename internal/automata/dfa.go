package automata

// DFA is the top-level automaton. Only the start state is stored; all other
// states are reachable by following Transitions recursively.
type DFA struct {
	StartState *DFAState
}

func NewDFA(startState *DFAState) *DFA {
	return &DFA{StartState: startState}
}

// DFAState is a node in the automaton.
// TokenID is set on accepting states to record which pattern matched (1-based).
// Non-accepting states carry TokenID = 0.
type DFAState struct {
	ID          int
	Transitions map[rune]*DFAState
	Accept      bool
	TokenID     int // 1-based index into the pattern list; 0 on non-accept states
}

// DFAStateCounter is a global monotonic counter used to assign unique IDs.
// All DFAState instances share this counter regardless of which DFA they belong to.
var DFAStateCounter int

func NewDFAState() *DFAState {
	DFAStateCounter++
	return &DFAState{
		ID:          DFAStateCounter,
		Transitions: make(map[rune]*DFAState),
		Accept:      false,
	}
}

func (from *DFAState) AddTransition(to *DFAState, char rune) {
	from.Transitions[char] = to
}

func (state *DFAState) SetAccepting() {
	state.Accept = true
}

// SetToken marks the state as accepting for the given 1-based token index.
func (state *DFAState) SetToken(id int) {
	state.Accept = true
	state.TokenID = id
}

func (dfa *DFA) Alphabet() map[rune]bool {
	visited := make(map[int]bool)
	alphabet := make(map[rune]bool)
	dfaAlphabetRecursive(dfa.StartState, alphabet, visited)
	return alphabet
}

func dfaAlphabetRecursive(state *DFAState, alphabet map[rune]bool, visited map[int]bool) {
	if visited[state.ID] {
		return
	}
	visited[state.ID] = true
	for symbol, nextState := range state.Transitions {
		alphabet[symbol] = true
		dfaAlphabetRecursive(nextState, alphabet, visited)
	}
}
