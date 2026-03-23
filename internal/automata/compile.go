// Package automata converts preprocessed regex token sequences into minimized
// DFAs and combines them into a single scanner automaton.
//
// The pipeline for a single pattern is:
//   - compile.go: orchestrates the full flow — regex → AST → table → DFA → minimize
//   - direct.go:  builds the DFA directly from the syntax tree (direct method),
//                 computing nullable/firstpos/lastpos/followpos on the AST nodes
//   - minimization.go: Hopcroft partition-refinement minimization + state rebuild
//   - dfa.go:     DFA/DFAState structs, transitions, alphabet, simulation helpers
//
// For multi-pattern lexers, Merge() in compile.go runs parallel simulation over
// all per-pattern DFAs to produce one combined minimized DFA with first-rule-wins
// priority baked in.
package automata

import (
	"strconv"
	"strings"

	"github.com/Jose-Merida-UVG/Compilers-1-CC3071/internal/regex"
)

// Compile builds a minimized DFA from a preprocessed RegexString.
func Compile(rs *regex.RegexString) *DFA {
	root := regex.BuildDirectAST(rs)
	table := BuildDirectTable(&root)
	return table.ToDFA().Minimize()
}

// Merge combines multiple per-pattern DFAs into one minimized DFA using
// parallel simulation (also called the product construction).
//
// The core idea: instead of running each DFA independently and picking the
// winner afterward, we simulate ALL of them simultaneously. Each state in the
// merged DFA is a multiState — a tuple of state IDs, one per input DFA.
// Component i holds the current state of DFA i, or -1 if that DFA has died
// (no valid transition for the symbol seen so far).
//
// At any accepting multiState, the component with the lowest index (i.e. the
// pattern listed earliest in the .yal file) wins — this is first-rule-wins.
// The merged DFA is minimized afterward to collapse any equivalent states.
func Merge(dfas []*DFA) *DFA {
	if len(dfas) == 0 {
		return nil
	}

	n := len(dfas)

	// Build ID → *DFAState lookup for each input DFA so we can follow
	// transitions by ID during the BFS without chasing pointers across DFAs.
	stateByID := make([]map[int]*DFAState, n)
	for i, dfa := range dfas {
		stateByID[i] = make(map[int]*DFAState)
		for _, s := range dfa.GetAllStates() {
			stateByID[i][s.ID] = s
		}
	}

	// Union of all alphabets — the merged DFA must handle every symbol any
	// individual DFA could ever see.
	alphabet := make(map[rune]bool)
	for _, dfa := range dfas {
		for r := range dfa.Alphabet() {
			alphabet[r] = true
		}
	}

	// A multiState is a slice of n state IDs, one per input DFA.
	// -1 means that DFA has no transition for the current path (dead).
	type multiState = []int

	// keyOf serialises a multiState to a string so it can be used as a map key.
	// e.g. DFA 0 in state 3, DFA 1 in state 7, DFA 2 dead → "3,7,-1,"
	keyOf := func(ms multiState) string {
		var b strings.Builder
		for _, id := range ms {
			b.WriteString(strconv.Itoa(id))
			b.WriteByte(',')
		}
		return b.String()
	}

	// tokenOf scans the multiState left-to-right and returns the 1-based token
	// ID of the first accepting component. Left-to-right order is pattern order,
	// so the earliest pattern in the .yal file wins when multiple DFAs accept
	// the same string — first-rule-wins priority.
	// Returns 0 if no component is currently in an accepting state.
	tokenOf := func(ms multiState) int {
		for i, id := range ms {
			if id == -1 {
				continue
			}
			if s, ok := stateByID[i][id]; ok && s.Accept {
				return i + 1
			}
		}
		return 0
	}

	// Seed the BFS: the start multiState has every DFA at its own start state.
	startMS := make(multiState, n)
	for i, dfa := range dfas {
		startMS[i] = dfa.StartState.ID
	}

	// resultMap tracks which multiStates we've already created a DFA state for,
	// keyed by their serialised form so we don't create duplicates.
	resultMap := make(map[string]*DFAState)
	queue := []multiState{startMS}

	startState := NewDFAState()
	if tok := tokenOf(startMS); tok > 0 {
		startState.SetToken(tok)
	}
	resultMap[keyOf(startMS)] = startState

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]
		currentState := resultMap[keyOf(current)]

		for sym := range alphabet {
			// Advance every component DFA by sym.
			// If a component is already dead (-1) or has no transition for sym,
			// it stays dead (-1) in the next multiState.
			next := make(multiState, n)
			anyAlive := false
			for i, id := range current {
				if id == -1 {
					next[i] = -1
					continue
				}
				s, ok := stateByID[i][id]
				if !ok {
					next[i] = -1
					continue
				}
				if ns, ok2 := s.Transitions[sym]; ok2 {
					next[i] = ns.ID
					anyAlive = true
				} else {
					// This DFA dies on sym — mark it dead for the rest of this path
					next[i] = -1
				}
			}
			// If every component is dead, this symbol leads nowhere — skip it.
			if !anyAlive {
				continue
			}
			nextKey := keyOf(next)
			if _, exists := resultMap[nextKey]; !exists {
				ns := NewDFAState()
				if tok := tokenOf(next); tok > 0 {
					ns.SetToken(tok)
				}
				resultMap[nextKey] = ns
				queue = append(queue, next)
			}
			currentState.AddTransition(resultMap[nextKey], sym)
		}
	}

	// Minimize the merged DFA — parallel simulation can leave many equivalent
	// states that Hopcroft will collapse, keeping the final automaton small.
	return NewDFA(startState).Minimize()
}
