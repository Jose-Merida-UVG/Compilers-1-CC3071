package automata

import (
	"errors"

	"github.com/TonitoMC/TDLCGo/internal/ds"
	"github.com/TonitoMC/TDLCGo/internal/regex"
)

// Thompson's Construction for NFA based on regular expression.
// Takes a RegexString of Tokens and returns an NFA, expression
// MUST already be simplified to only include '*', '|' or concat.
func ThompsonConstruction(regex *regex.RegexString) (*NFA, error) {
	stack := ds.Stack[*NFA]{}
	for _, char := range regex.Chars {
		if !char.IsOperator() {
			// If char isn't an operator a new FSM gets added with
			// the character as it's sole transition from accept to end
			startState := NewState()
			acceptState := NewState()
			if char.Value == 'ε' {
				startState.AddEpsilonTransition(acceptState)
			} else {
				startState.AddTransition(acceptState, char.Value)
			}
			fsm := NewFSM(startState, acceptState)
			stack.Push(fsm)
		} else if char.Value == '~' {
			fsm2, err := stack.Pop()
			if err != nil {
				return nil, errors.New("error en construccion de NFA, regex no valido")
			}
			fsm1, err := stack.Pop()
			if err != nil {
				return nil, errors.New("error en construccion de NFA, regex no valido")
			}
			// Add an epsilon transition from fsm1's accept state to fsm2's start state
			MergeStates(fsm1.AcceptState, fsm2.StartState)
			fsm := NewFSM(fsm1.StartState, fsm2.AcceptState)
			stack.Push(fsm)
		} else if char.Value == '|' {
			fsm2, err := stack.Pop()
			if err != nil {
				return nil, errors.New("error en construccion de NFA, regex no valido")
			}
			fsm1, err := stack.Pop()
			if err != nil {
				return nil, errors.New("error en construccion de NFA, regex no valido")
			}
			// Build a new initial state that transitions via epsilon to either start
			// state and a new final state that is transitioned to via epsilon from
			// either accept state
			startState := NewState()
			acceptState := NewState()
			startState.AddEpsilonTransition(fsm1.StartState)
			startState.AddEpsilonTransition(fsm2.StartState)
			fsm1.AcceptState.AddEpsilonTransition(acceptState)
			fsm2.AcceptState.AddEpsilonTransition(acceptState)
			fsm := NewFSM(startState, acceptState)
			stack.Push(fsm)
		} else if char.Value == '*' {
			fsm1, err := stack.Pop()
			if err != nil {
				return nil, errors.New("error en construccion de NFA, regex no valido")
			}
			// Build a new initial state that transitions via epsilon to the
			// start state and a new final state that is transitioned to via epsilon
			// from accept state. Also add epsilon transition from original end
			// state to original start state
			startState := NewState()
			acceptState := NewState()
			fsm1.AcceptState.AddEpsilonTransition(fsm1.StartState)
			startState.AddEpsilonTransition(fsm1.StartState)
			startState.AddEpsilonTransition(acceptState)
			fsm1.AcceptState.AddEpsilonTransition(acceptState)
			fsm := NewFSM(startState, acceptState)
			stack.Push(fsm)
		}
	}
	toReturn, _ := stack.Pop()
	return toReturn, nil
}
