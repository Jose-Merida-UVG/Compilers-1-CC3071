package grammar

func (g *Grammar) ComputeFollow() {
	// Inicializar sets vacíos para cada no terminal
	g.Follow = make(map[string]map[string]bool)
	for nt := range g.NonTerminals {
		g.Follow[nt] = make(map[string]bool)
	}

	// Regla 1: $ siempre está en FOLLOW del símbolo inicial
	g.Follow[g.StartSymbol]["$"] = true

	changed := true

	for changed {
		changed = false

		// Iterar sobre todas las producciones A → X1 X2 ... Xn
		for i := range g.Productions {
			prod := &g.Productions[i]
			A := prod.Head
			body := prod.Body

			for j := range body {
				sym := &body[j]

				// Solo nos interesan los no terminales del cuerpo
				if sym.IsTerminal {
					continue
				}

				B := sym.Name // el no terminal que estamos procesando

				// Regla 2: agregar FIRST(β) - {ε} a FOLLOW(B)
				// donde β = lo que viene después de B en esta producción
				// β = body[j+1 ... n]
				betaNullable := true // asumimos que β puede derivar ε

				for k := j + 1; k < len(body); k++ {
					next := &body[k]

					if next.IsTerminal {
						// terminal: se agrega directamente y se corta
						if !g.Follow[B][next.Name] {
							g.Follow[B][next.Name] = true
							changed = true
						}
						betaNullable = false
						break
					}

					// no terminal: agregar FIRST(next) - {ε} a FOLLOW(B)
					for t := range g.First[next.Name] {
						if t == Epsilon {
							continue
						}
						if !g.Follow[B][t] {
							g.Follow[B][t] = true
							changed = true
						}
					}

					// si next NO es nullable, β ya no puede derivar ε
					if !g.First[next.Name][Epsilon] {
						betaNullable = false
						break
					}
					// si next ES nullable, seguimos mirando lo que sigue
				}

				// Regla 3: si β puede derivar ε (o B está al final),
				// agregar FOLLOW(A) a FOLLOW(B)
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