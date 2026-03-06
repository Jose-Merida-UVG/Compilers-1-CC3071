package automata

// State class to represent different states
// of a Finite State Machine
type State struct {
	ID                 int
	Transitions        map[rune]*State // Maps character transitions
	EpsilonTransitions []*State
}

// Keeping State ID's unique
var stateCounter int

// Constructor for state
func NewState() *State {
	stateCounter++
	return &State{
		ID:                 stateCounter,
		Transitions:        make(map[rune]*State),
		EpsilonTransitions: []*State{},
	}
}

// Calculate EpsilonClosure for a State
func (state *State) EpsilonClosure() *StateSet {
	closure := NewStateSet()
	EpsilonClosureRecursive(state, closure)
	return closure
}

// Recursive Epsilon Closure
func EpsilonClosureRecursive(state *State, closure *StateSet) {
	if closure.Contains(state) {
		return
	}
	closure.Add(state)
	for _, nextState := range state.EpsilonTransitions {
		EpsilonClosureRecursive(nextState, closure)
	}
}

func (from *State) AddTransition(to *State, char rune) {
	from.Transitions[char] = to
}

func (from *State) AddEpsilonTransition(to *State) {
	from.EpsilonTransitions = append(from.EpsilonTransitions, to)
}

// NFA class to represent a finite state machine
// via its start state and accept state
type NFA struct {
	StartState  *State
	AcceptState *State
}

// Constructor for FSM struct
func NewFSM(start *State, accept *State) *NFA {
	fsm := NFA{StartState: start,
		AcceptState: accept}
	return &fsm
}

// Simulates a FSM, returns whether a string 'input'
// belongs to the regular languange represented by
// the FSM
func (fsm *NFA) Sim(input string) bool {
	currentStates := fsm.StartState.EpsilonClosure()
	for _, char := range input {
		currentStates = currentStates.SingleTransition(char)
	}
	return currentStates.Contains(fsm.AcceptState)
}

// StateSet struct to represent a set of states, handles utility
// functions such as calculating the epsilon closure and the
// transitions following a symbol (symbolε*) to compute all
// reachable states
type StateSet struct {
	States map[int]*State
}

// Constructor for Stateset
func NewStateSet() *StateSet {
	return &StateSet{
		States: make(map[int]*State),
	}
}

// MergeStates merges two states into one
func MergeStates(state1 *State, state2 *State) {
	// Merge transitions
	for symbol, nextState := range state2.Transitions {
		state1.Transitions[symbol] = nextState
	}
	// Merge epsilon transitions
	state1.EpsilonTransitions = append(state1.EpsilonTransitions, state2.EpsilonTransitions...)
}

// Add a state to a StateSet
func (set *StateSet) Add(state *State) {
	set.States[state.ID] = state
}

func (firstSet *StateSet) Concat(secondSet *StateSet) {
	for id, state := range secondSet.States {
		firstSet.States[id] = state
	}
}

// Returns whether a specific state is contained within
// a Stateset
func (stateSet *StateSet) Contains(state *State) bool {
	_, exists := stateSet.States[state.ID]
	return exists
}

// Returns the states reachable via a single character transition
// for any state in a StateSet
func (stateSet *StateSet) SingleTransition(char rune) *StateSet {
	nextStates := NewStateSet()
	for _, state := range stateSet.States {
		if nextState, ok := state.Transitions[char]; ok {
			nextStates.Concat(nextState.EpsilonClosure())
		}
	}
	return nextStates
}

func (firstStateSet *StateSet) CompareTo(secondStateSet *StateSet) bool {
	if len(firstStateSet.States) != len(secondStateSet.States) {
		return false
	}
	for id := range firstStateSet.States {
		if _, exists := secondStateSet.States[id]; !exists {
			return false
		}
	}
	return true
}

func (stateSet *StateSet) IsEmpty() bool {
	return len(stateSet.States) == 0
}
