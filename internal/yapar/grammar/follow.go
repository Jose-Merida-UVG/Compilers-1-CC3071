package grammar

// ComputeFollow populates g.Follow for every non-terminal using fixed-point iteration.
//
// FOLLOW(A) is the set of terminals that can appear immediately after A in any valid
// derivation. Required by SLR(1) to decide when to reduce. Three rules are applied
// repeatedly until no set changes (β = the suffix of symbols after B in the body,
// scanned left to right. If a symbol in β is nullable we keep scanning past it):
//
//   - $ always follows the start symbol S.
//   - For A → ... B β: scan β left to right, adding each symbol's FIRST − {ε} to FOLLOW(B);
//     stop as soon as a non-nullable symbol is reached.
//   - For A → ... B β where β fully derives ε (or is empty): add FOLLOW(A) to FOLLOW(B).

func (g *Grammar) ComputeFollow() {
	g.Follow = make(map[string]map[string]bool)
	for nt := range g.NonTerminals {
		g.Follow[nt] = make(map[string]bool)
	}

	// End-of-input always follows the start symbol.
	g.Follow[g.StartSymbol]["$"] = true

	changed := true

	// Repeat until no changes
	for changed {
		changed = false

		// Scan all productions
		for i := range g.Productions {
			prod := &g.Productions[i]
			A := prod.Head
			body := prod.Body

			// Go through production's body
			for j := range body {
				sym := &body[j]

				if sym.IsTerminal {
					continue
				}

				// B is the current non-terminal; β is everything after it in the body.
				B := sym.Name
				betaNullable := true

				for k := j + 1; k < len(body); k++ {
					next := &body[k]

					if next.IsTerminal {
						// Terminal in β: add it to FOLLOW(B) and stop since β isn't nullable past here.
						if !g.Follow[B][next.Name] {
							g.Follow[B][next.Name] = true
							changed = true
						}
						betaNullable = false
						break
					}

					// Non-terminal in β: propagate its FIRST − {ε} into FOLLOW(B).
					for t := range g.First[next.Name] {
						if t == Epsilon {
							continue
						}
						if !g.Follow[B][t] {
							g.Follow[B][t] = true
							changed = true
						}
					}

					// If this symbol in β isn't nullable, β can't derive ε. We stop scanning.
					if !g.First[next.Name][Epsilon] {
						betaNullable = false
						break
					}
				}

				// β fully derives ε (or B is at the tail): whatever follows A also follows B.
				if betaNullable {
					for t := range g.Follow[A] {
						if !g.Follow[B][t] {
							g.Follow[B][t] = true
							changed = true
						}
					}
				}
			}
		}
	}
}
