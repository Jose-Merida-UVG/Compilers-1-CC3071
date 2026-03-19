package graph

import (
	"fmt"

	"github.com/Jose-Merida-UVG/Compilers-1-CC3071/automata"
)

// Builds a .DOT file that represents a DFA
func BuildDFA(dfa *automata.DFA, dir string, filename string) {
	content := CreateHeader()
	visited := make(map[int]bool)
	content = AddDFAState(dfa.StartState, visited, content)
	WriteFile(content, dir, filename)
}

// Recursive function that draws state and it's transitions
func AddDFAState(state *automata.DFAState, visited map[int]bool, content string) string {
	if state == nil {
		panic("Encountered a nil DFA state")
	}

	if visited[state.ID] {
		return content
	}

	stateLabel := fmt.Sprintf("%d", len(visited))
	stateID := fmt.Sprintf("%d", state.ID)
	visited[state.ID] = true
	if state.Accept {
		content = AddAcceptNode(content, stateID, stateLabel)
	} else {
		content = AddNode(content, stateID, stateLabel)
	}

	for char, nextState := range state.Transitions {
		if nextState == nil {
			panic(fmt.Sprintf("State %d has a nil transition on symbol %c", state.ID, char))
		}
		nextStateID := fmt.Sprintf("%d", nextState.ID)
		content = AddDFAState(nextState, visited, content)
		transitionLabel := fmt.Sprintf("%s", string(char))
		content = AddEdge(content, stateID, nextStateID, transitionLabel)
	}
	return content
}
