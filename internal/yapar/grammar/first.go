package grammar

func (g *Grammar) ComputeFirst() {
    //inicializar mapa
    g.First = make(map[string]map[string]bool)

    //crear sets vacíos para cada no terminal
    for nt := range g.NonTerminals {
        g.First[nt] = make(map[string]bool)
    }

    changed := true

}
