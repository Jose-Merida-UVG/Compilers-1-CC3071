package lr0

import (
	"fmt"
	"sort"
)

// LR0GraphData is the JSON-serializable representation of the LR(0) automaton,
// written to workspace/parsers/<name>.lr0.json for visualization.
type LR0GraphData struct {
	Nodes []LR0Node `json:"nodes"`
	Edges []LR0Edge `json:"edges"`
}

// LR0Node represents one state in the automaton.
type LR0Node struct {
	ID       string   `json:"id"`
	Label    string   `json:"label"`   // "State N"
	Items    []string `json:"items"`   // human-readable item strings, e.g. "E -> T • + E"
	IsStart  bool     `json:"isStart"` // true for state 0
	IsAccept bool     `json:"isAccept"`// true if it contains the complete augmented item S' -> S •
}

// LR0Edge is a transition between two states on a given symbol.
type LR0Edge struct {
	ID     string `json:"id"`
	Source string `json:"source"`
	Target string `json:"target"`
	Symbol string `json:"symbol"`
}

// Serialize converts the Automaton into its JSON-serializable form.
func (a *Automaton) Serialize() *LR0GraphData {
	productions := a.Productions
	nodes := make([]LR0Node, 0, len(a.States))

	for _, state := range a.States {
		// Build human-readable item strings for this state.
		items := make([]string, len(state.Items))
		for i, item := range state.Items {
			items[i] = a.ItemString(item)
		}

		// A state is an accept state if it contains the complete augmented item:
		// S' → S • (production index 0, dot at end).
		isAccept := false
		for _, item := range state.Items {
			if item.ProdIndex == 0 && item.Dot >= len(productions[0].Body) {
				isAccept = true
				break
			}
		}

		nodes = append(nodes, LR0Node{
			ID:       fmt.Sprintf("%d", state.ID),
			Label:    fmt.Sprintf("State %d", state.ID),
			Items:    items,
			IsStart:  state.ID == 0,
			IsAccept: isAccept,
		})
	}

	// Collect all state IDs in sorted order so edge IDs are deterministic.
	stateIDs := make([]int, 0, len(a.Transitions))
	for id := range a.Transitions {
		stateIDs = append(stateIDs, id)
	}
	sort.Ints(stateIDs)

	var edges []LR0Edge
	edgeID := 0

	for _, sourceID := range stateIDs {
		// Sort symbols for deterministic output.
		symbols := make([]string, 0, len(a.Transitions[sourceID]))
		for sym := range a.Transitions[sourceID] {
			symbols = append(symbols, sym)
		}
		sort.Strings(symbols)

		for _, symbol := range symbols {
			targetID := a.Transitions[sourceID][symbol]
			edgeID++
			edges = append(edges, LR0Edge{
				ID:     fmt.Sprintf("e%d", edgeID),
				Source: fmt.Sprintf("%d", sourceID),
				Target: fmt.Sprintf("%d", targetID),
				Symbol: symbol,
			})
		}
	}

	return &LR0GraphData{Nodes: nodes, Edges: edges}
}
