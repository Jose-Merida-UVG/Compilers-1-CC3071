// Package lr0 builds the canonical LR(0) automaton for a given (augmented) grammar.
// The main entry point is New, which initializes the automaton, followed by Build,
// which computes all reachable states via Closure and Goto. The resulting Automaton
// is consumed by the slr package to construct the SLR(1) parse table.
// Serialize produces a JSON-friendly form used by the frontend graph visualizer.
package lr0

import (
	"fmt"
	"sort"
	"strings"

	"github.com/Jose-Merida-UVG/Compilers-1-CC3071/internal/yapar/grammar"
)

// Item is a single LR(0) item.
type Item struct {
	ProdIndex int // index into Grammar.Productions
	Dot       int // position of the dot (0-indexed)
}

// State is one node in the LR(0) automaton.
// It represents every parsing position the parser could simultaneously be in
// after having read the same sequence of symbols to get here.
type State struct {
	ID    int    // unique state number, assigned in discovery order
	Items []Item // full set of LR(0) items that define this state
}

// Automaton is the complete LR(0) automaton
// Transitions[stateID][symbol] → the ID of the next state after reading symbol.
type Automaton struct {
	Grammar     *grammar.Grammar
	Productions []grammar.Production // local ref to Grammar.Productions to avoid deep chaining
	ProdsByHead map[string][]int     // non-terminal name → indices of its productions in Productions, for practicality
	States      []State
	Transitions map[int]map[string]int
}

// New builds the initial Automaton from an already-augmented Grammar.
// The augmented grammar must have S' → S as its first production (index 0).
// Must call g.Augment() before passing the grammar here.
func New(g *grammar.Grammar) *Automaton {
	// Build the head → production indices index up front so Closure and Goto
	// can look up productions for a non-terminal.
	prodsByHead := make(map[string][]int, len(g.NonTerminals))
	for idx, prod := range g.Productions {
		prodsByHead[prod.Head] = append(prodsByHead[prod.Head], idx)
	}

	// Initialize automata
	automaton := &Automaton{
		Grammar:     g,
		Productions: g.Productions,
		ProdsByHead: prodsByHead,
		Transitions: make(map[int]map[string]int),
	}

	// State 0 starts from the single item S' → • S, closure then expands it into the full initial state.
	initialItem := Item{ProdIndex: 0, Dot: 0}
	state0 := State{ID: 0, Items: automaton.Closure([]Item{initialItem})}
	automaton.States = append(automaton.States, state0)

	return automaton
}

// Build computes the full canonical LR(0) collection starting from state 0.
// It repeatedly applies Goto over every symbol that appears after a dot in
// each state, collecting new states until no new ones are found.
func (a *Automaton) Build() {
	productions := a.Productions

	// stateKeys maps a state's canonical string key to its ID.
	// Used to check if a Goto result is a state we've already seen.
	stateKeys := make(map[string]int)
	stateKeys[stateKey(a.States[0].Items)] = 0

	// Process states as they are discovered
	for i := 0; i < len(a.States); i++ {
		currentState := a.States[i]

		// Collect every symbol that appears after a dot in this state.
		symbolsSeen := make(map[string]bool)
		for _, item := range currentState.Items {
			prod := productions[item.ProdIndex]
			if item.Dot < len(prod.Body) {
				symbolsSeen[prod.Body[item.Dot].Name] = true
			}
		}

		// Apply Goto for each symbol and record the transition.
		for symbol := range symbolsSeen {
			gotoItems := a.Goto(currentState, symbol)
			if len(gotoItems) == 0 {
				continue
			}

			key := stateKey(gotoItems)

			// If this state hasn't been seen before, add it.
			targetID, exists := stateKeys[key]
			if !exists {
				targetID = len(a.States)
				stateKeys[key] = targetID
				a.States = append(a.States, State{ID: targetID, Items: gotoItems})
			}

			// Record the transition: currentState w/ symbol--> targetID
			if a.Transitions[currentState.ID] == nil {
				a.Transitions[currentState.ID] = make(map[string]int)
			}
			a.Transitions[currentState.ID][symbol] = targetID
		}
	}
}

// ItemString returns a human-readable representation of an item.
// Example:  E -> T • + E
func (a *Automaton) ItemString(item Item) string {
	prod := a.Productions[item.ProdIndex]

	var builder strings.Builder
	builder.WriteString(prod.Head)
	builder.WriteString(" ->")

	// Write each body symbol, placing the dot • before the symbol at position Dot.
	for i, sym := range prod.Body {
		if i == item.Dot {
			builder.WriteString(" •")
		}
		builder.WriteString(" ")
		builder.WriteString(sym.Name)
	}

	// If the dot is past the last symbol, append it at the end.
	if item.Dot >= len(prod.Body) {
		builder.WriteString(" •")
	}

	return builder.String()
}

// stateKey returns a canonical string that uniquely identifies a set of items.
// Items are sorted so that two states with the same items in different order
// produce the same key. Uses brackets to prevent any possible collision between
// adjacent numbers (e.g., distinguishes (1,23)|(12,3) from (12,3)|(1,23)).
func stateKey(items []Item) string {
	keys := make([]string, len(items))
	for i, item := range items {
		keys[i] = fmt.Sprintf("[%d:%d]", item.ProdIndex, item.Dot)
	}
	sort.Strings(keys)
	return strings.Join(keys, "|")
}
