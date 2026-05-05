package grammar

// ComputeFollow calcula el conjunto FOLLOW para todos los no terminales de la gramática.
//
// El conjunto FOLLOW(A) contiene todos los terminales que pueden aparecer inmediatamente
// a la derecha de A en alguna forma sentencial derivada del símbolo inicial.
// FOLLOW se utiliza en el análisis sintáctico (parsing) para construir tablas de parsing
// LL(1) y también es necesario para parsers SLR(1).
//
// El algoritmo implementa las tres reglas clásicas del cálculo de FOLLOW:
//
//	Regla 1: $ ∈ FOLLOW(S), donde S es el símbolo inicial de la gramática.
//	Regla 2: Si existe una producción A → αBβ, entonces FIRST(β) - {ε} ⊆ FOLLOW(B).
//	Regla 3: Si existe una producción A → αBβ y ε ∈ FIRST(β) (o β es vacío),
//	          entonces FOLLOW(A) ⊆ FOLLOW(B).
//
// Se usa un enfoque de punto fijo (fixed-point iteration): se repiten las reglas
// hasta que ningún conjunto FOLLOW cambie, ya que la Regla 3 puede propagar
// elementos de FOLLOW entre no terminales de forma transitiva.
func (g *Grammar) ComputeFollow() {

	// ──────────────────────────────────────────────
	// Paso 0: Inicialización
	// Se crea un mapa vacío de FOLLOW para cada no terminal de la gramática.
	// Cada FOLLOW(nt) es un conjunto (set) representado como map[string]bool.
	// ──────────────────────────────────────────────
	g.Follow = make(map[string]map[string]bool)
	for nt := range g.NonTerminals {
		g.Follow[nt] = make(map[string]bool)
	}

	// ──────────────────────────────────────────────
	// Regla 1: El marcador de fin de entrada "$" siempre pertenece
	// al conjunto FOLLOW del símbolo inicial de la gramática.
	// Esto indica que al terminar de leer toda la entrada,
	// el símbolo inicial debe poder ser seguido por "$".
	// ──────────────────────────────────────────────
	g.Follow[g.StartSymbol]["$"] = true

	// Variable centinela para el ciclo de punto fijo.
	// Se mantiene en true mientras algún conjunto FOLLOW haya cambiado.
	changed := true

	// ──────────────────────────────────────────────
	// Ciclo de punto fijo (fixed-point iteration):
	// Se repite la aplicación de las reglas 2 y 3 hasta que
	// no se agregue ningún nuevo terminal a ningún conjunto FOLLOW.
	// Esto es necesario porque la Regla 3 puede causar propagación
	// en cadena: FOLLOW(A) → FOLLOW(B) → FOLLOW(C) → ...
	// ──────────────────────────────────────────────
	for changed {
		changed = false

		// Recorrer todas las producciones de la gramática.
		// Cada producción tiene la forma A → X1 X2 ... Xn
		// donde A es la cabeza (Head) y X1...Xn es el cuerpo (Body).
		for i := range g.Productions {
			prod := &g.Productions[i] // Puntero a la producción actual para evitar copias
			A := prod.Head            // No terminal de la cabeza de la producción
			body := prod.Body         // Cuerpo de la producción (slice de símbolos)

			// Recorrer cada símbolo del cuerpo de la producción.
			for j := range body {
				sym := &body[j] // Puntero al símbolo actual

				// Solo nos interesan los no terminales del cuerpo,
				// ya que FOLLOW solo se define para no terminales.
				// Si el símbolo es terminal, lo saltamos.
				if sym.IsTerminal {
					continue
				}

				// B es el no terminal actual del cuerpo para el que
				// vamos a calcular/actualizar FOLLOW(B).
				B := sym.Name

				// ──────────────────────────────────────────────
				// Regla 2: Para la producción A → αBβ, agregar
				// FIRST(β) - {ε} a FOLLOW(B).
				//
				// β es la subcadena que viene después de B en el cuerpo,
				// es decir, body[j+1], body[j+2], ..., body[n-1].
				//
				// betaNullable indica si toda la subcadena β puede
				// derivar ε (la cadena vacía). Inicialmente asumimos true
				// porque si B está al final de la producción (β es vacío),
				// entonces β trivialmente deriva ε.
				// ──────────────────────────────────────────────
				betaNullable := true

				// Recorrer cada símbolo de β (lo que sigue después de B).
				for k := j + 1; k < len(body); k++ {
					next := &body[k] // Siguiente símbolo en β

					if next.IsTerminal {
						// Caso: el siguiente símbolo es un terminal.
						// Un terminal siempre está en su propio FIRST,
						// así que lo agregamos directamente a FOLLOW(B).
						// Como un terminal no puede derivar ε,
						// β ya no es nullable y dejamos de recorrer.
						if !g.Follow[B][next.Name] {
							g.Follow[B][next.Name] = true
							changed = true // Se modificó FOLLOW(B)
						}
						betaNullable = false // β no puede derivar ε
						break                // No hay más que revisar en β
					}

					// Caso: el siguiente símbolo es un no terminal.
					// Agregamos FIRST(next) - {ε} a FOLLOW(B).
					// Es decir, todos los terminales que pueden iniciar
					// una derivación de 'next', excepto ε.
					for t := range g.First[next.Name] {
						if t == Epsilon {
							continue // No agregar ε al conjunto FOLLOW
						}
						if !g.Follow[B][t] {
							g.Follow[B][t] = true
							changed = true // Se modificó FOLLOW(B)
						}
					}

					// Si este no terminal NO puede derivar ε,
					// entonces β tampoco puede derivar ε.
					// Dejamos de recorrer β.
					if !g.First[next.Name][Epsilon] {
						betaNullable = false
						break
					}
					// Si este no terminal SÍ puede derivar ε,
					// seguimos revisando los símbolos que le siguen,
					// porque necesitamos considerar su FIRST también.
				}

				// ──────────────────────────────────────────────
				// Regla 3: Si β puede derivar ε (betaNullable == true),
				// o si B está al final de la producción (β es vacío),
				// entonces todo lo que puede seguir a A también
				// puede seguir a B. Es decir:
				//     FOLLOW(A) ⊆ FOLLOW(B)
				//
				// Esto se debe a que si β desaparece (deriva ε),
				// lo que sigue a la producción completa A → αBβ
				// será lo mismo que lo que sigue a A.
				// ──────────────────────────────────────────────
				if betaNullable {
					for t := range g.Follow[A] {
						if !g.Follow[B][t] {
							g.Follow[B][t] = true
							changed = true // Se modificó FOLLOW(B)
						}
					}
				}
			}
		}
	}
}