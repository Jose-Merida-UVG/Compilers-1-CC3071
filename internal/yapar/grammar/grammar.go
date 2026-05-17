// Package grammar builds a structured representation of a context-free grammar
// from a parsed YalpFile and exposes the data needed to compute FIRST/FOLLOW sets
// and construct SLR parse tables.
package grammar


// Epsilon is the sentinel string representing ε in FIRST sets and production display.
// An empty Body slice in a Production means the rule derives ε.
const Epsilon = "ε"

// Symbol is a grammar symbol, either a terminal (token) or a non-terminal.
type Symbol struct {
	Name       string
	IsTerminal bool
	TokenID    int // only meaningful for terminals; 0 for non-terminals
}

// Production is a single grammar rule: Head → Body.
type Production struct {
	Head string
	Body []Symbol
}

// Grammar is the fully-structured CFG ready for SLR table construction.
type Grammar struct {
	Name         string
	StartSymbol  string
	Terminals    map[string]int  // name → token ID
	NonTerminals map[string]bool // set of non-terminal names
	Productions  []Production
	IgnoreList   []string

	// Use Epsilon ("ε") as a key in First to indicate the symbol is nullable.
	First  map[string]map[string]bool
	Follow map[string]map[string]bool
}

// YalpSpec holds the subset of a parsed .yalp file that Build needs.
// It matches the fields of yapar.YalpFile so callers can pass one directly
// without grammar importing the yapar package.
type YalpSpec struct {
	Name         string
	StartSymbol  string
	TokenMap     map[string]int
	NonTerminals map[string]bool
	Productions  []RawProduction
	IgnoreList   []string
}

// RawProduction is a single production from the .yalp file before symbols are resolved.
type RawProduction struct {
	Name  string
	Rules [][]string
}

// Build constructs a Grammar from a YalpSpec.
func Build(spec YalpSpec) *Grammar {
	prods := make([]Production, 0)
	for _, p := range spec.Productions {
		for _, rule := range p.Rules {
			body := make([]Symbol, 0, len(rule))
			for _, sym := range rule {
				if id, ok := spec.TokenMap[sym]; ok {
					body = append(body, Symbol{Name: sym, IsTerminal: true, TokenID: id})
				} else {
					body = append(body, Symbol{Name: sym, IsTerminal: false})
				}
			}
			prods = append(prods, Production{Head: p.Name, Body: body})
		}
	}

	return &Grammar{
		Name:         spec.Name,
		StartSymbol:  spec.StartSymbol,
		Terminals:    spec.TokenMap,
		NonTerminals: spec.NonTerminals,
		Productions:  prods,
		IgnoreList:   spec.IgnoreList,
		First:        make(map[string]map[string]bool),
		Follow:       make(map[string]map[string]bool),
	}
}

// ProductionsFor returns all productions with the given head symbol.
func (g *Grammar) ProductionsFor(head string) []Production {
	var out []Production
	for _, p := range g.Productions {
		if p.Head == head {
			out = append(out, p)
		}
	}
	return out
}

// IsTerminal reports whether name is a declared token.
func (g *Grammar) IsTerminal(name string) bool {
	_, ok := g.Terminals[name]
	return ok
}

// IsNonTerminal reports whether name is a production head.
func (g *Grammar) IsNonTerminal(name string) bool {
	return g.NonTerminals[name]
}

// IsNullable reports whether a non-terminal can derive ε.
// Only valid after FIRST sets have been computed.
func (g *Grammar) IsNullable(name string) bool {
	return g.First[name][Epsilon]
}

// Augment returns a new Grammar with an augmented start production S' → S
// prepended at index 0. The original grammar is not modified.
func (g *Grammar) Augment() *Grammar {
	augHead := g.StartSymbol + "'"
	augProd := Production{
		Head: augHead,
		Body: []Symbol{{Name: g.StartSymbol, IsTerminal: false}},
	}

	prods := make([]Production, 0, len(g.Productions)+1)
	prods = append(prods, augProd)
	prods = append(prods, g.Productions...)

	nonTerminals := make(map[string]bool, len(g.NonTerminals)+1)
	for k, v := range g.NonTerminals {
		nonTerminals[k] = v
	}
	nonTerminals[augHead] = true

	return &Grammar{
		Name:         g.Name,
		StartSymbol:  augHead,
		Terminals:    g.Terminals,
		NonTerminals: nonTerminals,
		Productions:  prods,
		IgnoreList:   g.IgnoreList,
		First:        g.First,
		Follow:       g.Follow,
	}
}
