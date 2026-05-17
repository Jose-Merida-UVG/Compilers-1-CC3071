// Package slr builds an SLR(1) parse table from an LR(0) automaton.
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

// Conflict records a shift/reduce or reduce/reduce conflict.
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
//
//	Action[stateID][terminal]    → what to do
//	Goto[stateID][nonTerminal]   → which state to push after a reduce
type Table struct {
	Action    map[int]map[string]Action
	Goto      map[int]map[string]int
	Conflicts []Conflict
}

// Build constructs the SLR(1) parse table from the LR(0) automaton.
// The automaton must have been fully built (automaton.Build() called).
func Build(a *lr0.Automaton) *Table {
	productions := a.Productions
	follow := a.Grammar.Follow

	table := &Table{
		Action: make(map[int]map[string]Action),
		Goto:   make(map[int]map[string]int),
	}

	// Helper: set an Action cell, recording conflicts if a cell is already filled.
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
				// Dot is NOT at the end — symbol after dot drives a shift or goto.
				sym := prod.Body[item.Dot]

				if sym.IsTerminal {
					// Rule 1: A → α • a β  →  Action[state][a] = Shift(target)
					if target, ok := a.Transitions[state.ID][sym.Name]; ok {
						setAction(state.ID, sym.Name, Action{Kind: ActionShift, State: target})
					}
				} else {
					// Goto table: A → α • B β  →  Goto[state][B] = target
					if target, ok := a.Transitions[state.ID][sym.Name]; ok {
						if table.Goto[state.ID] == nil {
							table.Goto[state.ID] = make(map[string]int)
						}
						table.Goto[state.ID][sym.Name] = target
					}
				}

			} else {
				// Dot is at the end — complete item.

				// Rule 3: augmented production S' → S •  →  Action[state][$] = Accept
				if item.ProdIndex == 0 {
					setAction(state.ID, "$", Action{Kind: ActionAccept})
					continue
				}

				// Rule 2: A → α •  →  for each t in FOLLOW(A), Action[state][t] = Reduce
				for terminal := range follow[prod.Head] {
					setAction(state.ID, terminal, Action{Kind: ActionReduce, ProdIdx: item.ProdIndex})
				}
			}
		}
	}

	return table
}
