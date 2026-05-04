package grammar

func (g *Grammar) ComputeFirst() {
    //inicializar mapa
    g.First = make(map[string]map[string]bool)

    //crear sets vacíos para cada no terminal para evitar crashes
    for nt := range g.NonTerminals {
        g.First[nt] = make(map[string]bool)
    }

    changed := true //bandera para saber si algo cambió

    for changed {

        changed = false //repetimos el loop hasta que ya no hayan cambios

        //iterar sobre las producciones
        for i := range g.Productions {
            prod := &g.Productions[i] //obtenemos la producción actual

            A := prod.Head //lado izquierdo
            body := prod.Body //lado derecho

            //caso ε
            if len(body) == 0 {
                if !g.First[A][Epsilon] {
                    g.First[A][Epsilon] = true //si aún no se agregó epsilon a First(A) lo agregamos
                    changed = true //y marcamos el cambio realizado
                }
                continue
            }

            //procesar α = X1 X2 X3 ...
            allNullable := true //se asume que todos pueden producir epsilon

            for j := range body { //recorre simbolos
                sym := &body[j]
                name := sym.Name

                //si es terminal
                if sym.IsTerminal {
                    if !g.First[A][name] {
                        g.First[A][name] = true //se agrega directamente el terminal
                        /*si el lado derecho empieza con un terminal, ese es el FIRST*/
                        changed = true //se marca el cambio
                    }
                    allNullable = false
                    break //se termina, porque el terminal bloquea todo lo demás
                }

                //si es no terminal
                if g.IsNonTerminal(name) {
                    //agregar FIRST del símbolo sin ε
                    for t := range g.First[name] {
                        if t == Epsilon {
                            continue
                        }
                        if !g.First[A][t] { //agregamos cada terminal
                            g.First[A][t] = true
                            changed = true
                        }
                    }

                    //se verifica si sym es nullable y si no lo es, se detiene
                    if !g.First[name][Epsilon] {
                        allNullable = false
                        break
                    }
                }
            }

            //si todos eran nullable se agrega epsilon
            if allNullable {
                if !g.First[A][Epsilon] {
                    g.First[A][Epsilon] = true
                    changed = true
                }
            }
        }
    }

    //marcar Nullable en los símbolos, importante para FOLLOW
    for i := range g.Productions {
        prod := &g.Productions[i]
        for j := range prod.Body {
            sym := &prod.Body[j]
            if !sym.IsTerminal {
                sym.Nullable = g.First[sym.Name][Epsilon]
            }
        }
    }
}
