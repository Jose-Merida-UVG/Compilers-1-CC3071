package lr0

// Closure expands a seed set of items into a complete LR(0) state.
// For each item where the dot sits before a non-terminal B, every production of B
// is added at dot=0. Newly added items are expanded the same way until nothing
// new can be added.
func (a *Automaton) Closure(seedItems []Item) []Item {
	productions := a.Productions

	// allItems is the growing set, start with whatever is passed in
	allItems := make([]Item, len(seedItems))
	copy(allItems, seedItems)

	// Track which productions have been added at dot=0 to avoid duplicates, we
	// can match production index alone since closure only adds these at dot=0.
	expandedProds := make(map[int]bool)

	// Pre-mark seeds already at dot=0 so we don't re-add them.
	for _, seedItem := range allItems {
		if seedItem.Dot == 0 {
			expandedProds[seedItem.ProdIndex] = true
		}
	}

	// Go through all items
	for i := 0; i < len(allItems); i++ {
		currentItem := allItems[i]
		currentProd := productions[currentItem.ProdIndex]

		// Complete item, dot is at the end so we have nothing to expand.
		if currentItem.Dot >= len(currentProd.Body) {
			continue
		}

		symbolAfterDot := currentProd.Body[currentItem.Dot]

		// Terminals have no productions to expand into.
		if symbolAfterDot.IsTerminal {
			continue
		}

		// Non-terminal after dot: add all its productions at dot=0.
		for _, prodIndex := range a.ProdsByHead[symbolAfterDot.Name] {
			if expandedProds[prodIndex] {
				continue
			}
			expandedProds[prodIndex] = true
			allItems = append(allItems, Item{ProdIndex: prodIndex, Dot: 0})
		}
	}

	return allItems
}
