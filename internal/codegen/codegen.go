package codegen

import "github.com/Jose-Merida-UVG/Compilers-1-CC3071/internal/automata"

func GenerateLexer(name, sourcePath string, dfa *automata.DFA, actions []string) string {
	// TODO: implement code generation.
	//
	// Use dfa.GetAllStates() to iterate states.
	// Each state has:
	//   state.ID          int
	//   state.Accept      bool
	//   state.TokenID     int          (1-based; 0 on non-accept states)
	//   state.Transitions map[rune]*DFAState
	//
	// actions[tokenID-1] is the raw action string for that token.
	//
	// Return the full text of the generated .go file as a string.
	panic("codegen.GenerateLexer: not implemented — see the doc comment above")
}
