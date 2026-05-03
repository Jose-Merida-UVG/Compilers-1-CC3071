package graph

import (
	"fmt"
	"sort"

	"github.com/Jose-Merida-UVG/Compilers-1-CC3071/internal/yalex/automata"
)

// SerializeDFA converts a DFA into the JSON-serialisable form consumed by the frontend.
// Accept states are labelled with their TokenID (e.g. "q3/t2").
func SerializeDFA(dfa *automata.DFA) *DFAGraphData {
	states := dfa.GetAllStates()

	nodes := make([]DFAGraphNode, 0, len(states))
	for _, s := range states {
		label := fmt.Sprintf("q%d", s.ID)
		if s.Accept && s.TokenID > 0 {
			label = fmt.Sprintf("q%d/t%d", s.ID, s.TokenID)
		}
		nodes = append(nodes, DFAGraphNode{
			ID:     fmt.Sprintf("%d", s.ID),
			Label:  label,
			Accept: s.Accept,
			Start:  s.ID == dfa.StartState.ID,
		})
	}

	var edges []DFAGraphEdge
	eID := 0
	for _, s := range states {
		// Collapse parallel edges by grouping symbols that go to the same target.
		byTarget := make(map[int][]rune)
		for sym, next := range s.Transitions {
			byTarget[next.ID] = append(byTarget[next.ID], sym)
		}
		for targetID, syms := range byTarget {
			sort.Slice(syms, func(i, j int) bool { return syms[i] < syms[j] })
			parts := make([]string, len(syms))
			for i, sym := range syms {
				parts[i] = string(sym)
			}
			eID++
			edges = append(edges, DFAGraphEdge{
				ID:      fmt.Sprintf("e%d", eID),
				Source:  fmt.Sprintf("%d", s.ID),
				Target:  fmt.Sprintf("%d", targetID),
				Symbols: parts,
			})
		}
	}

	return &DFAGraphData{Nodes: nodes, Edges: edges}
}

// DFAGraphData is the JSON-serialisable representation of a DFA sent to the frontend.
type DFAGraphData struct {
	Nodes []DFAGraphNode `json:"nodes"`
	Edges []DFAGraphEdge `json:"edges"`
}

type DFAGraphNode struct {
	ID     string `json:"id"`
	Label  string `json:"label"`
	Accept bool   `json:"accept"`
	Start  bool   `json:"start"`
}

type DFAGraphEdge struct {
	ID      string   `json:"id"`
	Source  string   `json:"source"`
	Target  string   `json:"target"`
	Symbols []string `json:"symbols"`
}
