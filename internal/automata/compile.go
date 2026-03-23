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

// Merge combines multiple per-pattern DFAs into one minimized DFA.
// Each state in the result is a tuple of per-DFA state IDs (-1 = dead/no
// transition). Accepting tuples resolve priority by lowest index — first-rule-wins.
// The merged DFA is minimized again so equivalent merged states are collapsed.
func Merge(dfas []*DFA) *DFA {
	if len(dfas) == 0 {
		return nil
	}

	n := len(dfas)

	// Build ID → *DFAState lookup for each input DFA so we can follow
	// transitions by ID during the BFS without pointer indirection issues.
	stateByID := make([]map[int]*DFAState, n)
	for i, dfa := range dfas {
		stateByID[i] = make(map[int]*DFAState)
		for _, s := range dfa.GetAllStates() {
			stateByID[i][s.ID] = s
		}
	}

	// Union alphabet across all DFAs.
	alphabet := make(map[rune]bool)
	for _, dfa := range dfas {
		for r := range dfa.Alphabet() {
			alphabet[r] = true
		}
	}

	// A multiState is a slice of state IDs, one per input DFA (-1 = dead).
	type multiState = []int

	keyOf := func(ms multiState) string {
		var b strings.Builder
		for _, id := range ms {
			b.WriteString(strconv.Itoa(id))
			b.WriteByte(',')
		}
		return b.String()
	}

	// tokenOf returns the 1-based TokenID for the first (highest-priority)
	// accepting component, or 0 if the multiState is not accepting.
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

	startMS := make(multiState, n)
	for i, dfa := range dfas {
		startMS[i] = dfa.StartState.ID
	}

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
					next[i] = -1
				}
			}
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

	return NewDFA(startState).Minimize()
}
