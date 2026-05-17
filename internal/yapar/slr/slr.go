// Package slr builds an SLR(1) parse table from a fully-built LR(0) automaton.
// The table has two parts: an Action table (what to do on a terminal) and a Goto
// table (which state to enter after reducing). Conflicts are recorded and surfaced
// as an error — if any exist, the grammar is not SLR(1).
package slr

import (
	"fmt"

	"github.com/Jose-Merida-UVG/Compilers-1-CC3071/internal/yapar/lr0"
)

// ActionKind distinguishes the four possible parser actions.
type ActionKind int

const (
	ActionShift  ActionKind = iota // push a new state onto the stack
	ActionReduce                   // pop and apply a production
	ActionAccept                   // input fully parsed
	ActionError                    // no valid action (blank cell)
)

// Action is one cell in the Action table.
type Action struct {
	Kind    ActionKind
	State   int // for Shift: the target state to push
	ProdIdx int // for Reduce: which production to apply
}

func (a Action) String() string {
	switch a.Kind {
	case ActionShift:
		return fmt.Sprintf("s%d", a.State)
	case ActionReduce:
		return fmt.Sprintf("r%d", a.ProdIdx)
	case ActionAccept:
		return "acc"
	default:
		return ""
	}
}

// Conflict is a shift/reduce or reduce/reduce collision — two different actions
// competing for the same (state, symbol) cell in the Action table.
type Conflict struct {
	State    int
	Symbol   string
	Existing Action
	Incoming Action
}

func (c Conflict) String() string {
	return fmt.Sprintf("state %d on '%s': %s vs %s", c.State, c.Symbol, c.Existing, c.Incoming)
}

// Table is the complete SLR(1) parse table.
// Action[stateID][terminal] → what to do; Goto[stateID][nonTerminal] → next state after reduce.
type Table struct {
	Action    map[int]map[string]Action
	Goto      map[int]map[string]int
	Conflicts []Conflict
}

// Build constructs the SLR(1) parse table from a fully-built LR(0) automaton.
// For each state and item it fills the Action and Goto tables:
//   - dot before a terminal  → Shift to the transition target
//   - dot before a non-terminal → Goto entry to the transition target
//   - dot at the end → Reduce on every terminal in FOLLOW(head), or Accept for S' → S •
//
// Any (state, symbol) cell that would be written twice is a conflict — the grammar
// is not SLR(1). Conflicts are recorded in Table.Conflicts; call HasConflicts() to check.
func Build(a *lr0.Automaton) *Table {
	productions := a.Productions
	follow := a.Grammar.Follow

	table := &Table{
		Action: make(map[int]map[string]Action),
		Goto:   make(map[int]map[string]int),
	}

	// setAction writes an Action cell, recording a conflict if the cell is already filled.
	setAction := func(state int, symbol string, incoming Action) {
		if table.Action[state] == nil {
			table.Action[state] = make(map[string]Action)
		}
		existing, exists := table.Action[state][symbol]
		if exists && existing != incoming {
			table.Conflicts = append(table.Conflicts, Conflict{
				State:    state,
				Symbol:   symbol,
				Existing: existing,
				Incoming: incoming,
			})
			return
		}
		table.Action[state][symbol] = incoming
	}

	for _, state := range a.States {
		for _, item := range state.Items {
			prod := productions[item.ProdIndex]

			if item.Dot < len(prod.Body) {
				sym := prod.Body[item.Dot]

				if sym.IsTerminal {
					// Dot before a terminal: shift to the transition target.
					if target, ok := a.Transitions[state.ID][sym.Name]; ok {
						setAction(state.ID, sym.Name, Action{Kind: ActionShift, State: target})
					}
				} else {
					// Dot before a non-terminal: fill the Goto table.
					if target, ok := a.Transitions[state.ID][sym.Name]; ok {
						if table.Goto[state.ID] == nil {
							table.Goto[state.ID] = make(map[string]int)
						}
						table.Goto[state.ID][sym.Name] = target
					}
				}

			} else {
				// Complete item — dot is at the end.

				// Augmented production S' → S •: accept on end-of-input.
				if item.ProdIndex == 0 {
					setAction(state.ID, "$", Action{Kind: ActionAccept})
					continue
				}

				// Reduce on every terminal in FOLLOW(head).
				for terminal := range follow[prod.Head] {
					setAction(state.ID, terminal, Action{Kind: ActionReduce, ProdIdx: item.ProdIndex})
				}
			}
		}
	}

	return table
}

// HasConflicts reports whether the table has any shift/reduce or reduce/reduce conflicts.
func (t *Table) HasConflicts() bool {
	return len(t.Conflicts) > 0
}
