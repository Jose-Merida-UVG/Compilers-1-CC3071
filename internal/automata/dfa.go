package automata

// DFA defined by it's initial state, acceptance of each state is stored
// within each state, as each state is unique to each DFA per its construction
// we can define it there and simplify function calls in other files
type DFA struct {
	StartState *DFAState
}

func NewDFA(startState *DFAState) *DFA {
	dfa := DFA{StartState: startState}
	return &dfa
}

type DFAState struct {
	ID          int
	Transitions map[rune]*DFAState
	Accept      bool
}

// Keeping State ID's unique
var DFAStateCounter int

// Constructor for DFA State
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

func (dfa *DFA) Alphabet() map[rune]bool {
	visited := make(map[int]bool)
	alphabet := make(map[rune]bool)
	DFAAlphabetRecursive(dfa.StartState, alphabet, visited)
	return alphabet
}

func (dfa *DFA) Sim(input string) bool {
	currentState := dfa.StartState
	for _, char := range input {
		nextState := currentState.Transitions[char]
		if nextState == nil {
			return false
		}
		currentState = nextState
	}
	return currentState.Accept
}

func DFAAlphabetRecursive(state *DFAState, alphabet map[rune]bool, visited map[int]bool) {
	if visited[state.ID] {
		return
	}
	visited[state.ID] = true
	for symbol, nextState := range state.Transitions {
		if _, ok := alphabet[symbol]; !ok {
			alphabet[symbol] = true
		}
		DFAAlphabetRecursive(nextState, alphabet, visited)
	}
}
