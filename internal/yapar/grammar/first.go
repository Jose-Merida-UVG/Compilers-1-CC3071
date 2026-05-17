package grammar

// ComputeFirst populates g.First for every non-terminal using fixed-point iteration,
// there's probably a fancier way of doing this with tracking dependencies, etc. but
// this project was done on limited time so we went for a simpler approach.
//
// FIRST(A) is the set of terminals that can begin a string derived from A,
// plus ε if A is nullable. The algorithm applies three rules until no set changes
// (uppercase = non-terminal, lowercase = terminal, X1..Xn = any body symbols):
//
//  1. A → ε              ⟹  ε ∈ FIRST(A)
//  2. A → a rest         ⟹  a ∈ FIRST(A)
//  3. A → X1 ... Xn      ⟹  add FIRST(Xi) − {ε} to FIRST(A) for each Xi;
//     stop if Xi is not nullable; add ε if all Xi are nullable.
func (g *Grammar) ComputeFirst() {
	// Initialize an empty FIRST set for every non-terminal.
	g.First = make(map[string]map[string]bool)
	for nt := range g.NonTerminals {
		g.First[nt] = make(map[string]bool)
	}

	changed := true

	// Repeat until sets don't change
	for changed {
		changed = false

		// Iterate through every production
		for i := range g.Productions {
			prod := &g.Productions[i]

			// Note: LHS non-terminal will be referred to as A in code & comments, while
			// the current examined symbol in the RHS (body) will be referred to as sym.
			A := prod.Head
			body := prod.Body

			// Add epsilon to FIRST(A) if symbol derives the empty string
			if len(body) == 0 {
				if !g.First[A][Epsilon] {
					g.First[A][Epsilon] = true
					changed = true
				}
				continue
			}

			// Process body = X1 X2 … Xn
			allNullable := true

			for j := range body {
				sym := &body[j]
				name := sym.Name

				// Terminal: add sym to FIRST(A) and stop.
				if sym.IsTerminal {
					if !g.First[A][name] {
						g.First[A][name] = true
						changed = true
					}
					allNullable = false
					break
				}

				// Non-terminal: add FIRST(sym) − {ε} to FIRST(A).
				for t := range g.First[name] {
					if t == Epsilon {
						continue
					}
					if !g.First[A][t] {
						g.First[A][t] = true
						changed = true
					}
				}

				// If sym is not nullable, we cannot derive ε so we stop here.
				if !g.First[name][Epsilon] {
					allNullable = false
					break
				}
			}

			// If every symbol in the body is nullable, A is also nullable.
			if allNullable {
				if !g.First[A][Epsilon] {
					g.First[A][Epsilon] = true
					changed = true
				}
			}
		}
	}
}
