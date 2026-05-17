package lr0

// Goto computes the state reached after reading a given symbol from a given state.
// It collects all items where the dot sits right before that symbol, advances the dot
// past it, then runs Closure to expand into the full target state.
// Returns nil if no transition on symbol exists from this state.
func (a *Automaton) Goto(state State, symbol string) []Item {
	productions := a.Productions

	var advancedItems []Item

	for _, currentItem := range state.Items {
		currentProd := productions[currentItem.ProdIndex]

		// Skip complete items, nothing left to read.
		if currentItem.Dot >= len(currentProd.Body) {
			continue
		}

		// Only advance items where the next symbol matches.
		if currentProd.Body[currentItem.Dot].Name != symbol {
			continue
		}

		advancedItems = append(advancedItems, Item{ProdIndex: currentItem.ProdIndex, Dot: currentItem.Dot + 1})
	}

	if len(advancedItems) == 0 {
		return nil
	}

	// Expand the advanced items into a full state via Closure.
	return a.Closure(advancedItems)
}
