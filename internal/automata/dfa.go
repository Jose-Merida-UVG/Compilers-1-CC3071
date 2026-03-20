package automata

import (
	"fmt"
	"sort"
)

type DFA struct {
	StartState *DFAState
}

func NewDFA(startState *DFAState) *DFA {
	return &DFA{StartState: startState}
}

type DFAState struct {
	ID          int
	Transitions map[rune]*DFAState
	Accept      bool
	TokenID     int // 1-based index into the pattern list; 0 on non-accept states
}

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

func (dfa *DFA) PrintTransitionTable() {
	states := dfa.GetAllStates()
	alphabet := dfa.Alphabet()
	var symbols []rune
	for s := range alphabet {
		symbols = append(symbols, s)
	}
	sort.Slice(symbols, func(i, j int) bool { return symbols[i] < symbols[j] })

	fmt.Println("\n--- Tabla de Transiciones DFA ---")
	fmt.Printf("Estado\t")
	for _, s := range symbols {
		fmt.Printf("%c\t", s)
	}
	fmt.Println()
	for _, state := range states {
		marker := ""
		if state.Accept {
			marker = "*"
		}
		fmt.Printf("%d%s\t", state.ID, marker)
		for _, s := range symbols {
			if next, ok := state.Transitions[s]; ok {
				fmt.Printf("%d\t", next.ID)
			} else {
				fmt.Printf("-\t")
			}
		}
		fmt.Println()
	}
}

func (dfa *DFA) CountTransitions() int {
	states := dfa.GetAllStates()
	count := 0
	for _, state := range states {
		count += len(state.Transitions)
	}
	return count
}
