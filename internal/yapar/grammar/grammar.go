// Package grammar builds a structured representation of a context-free grammar
// from a parsed YalpFile and exposes the data needed to compute FIRST/FOLLOW sets
// and construct SLR parse tables.
package grammar

import (
	"fmt"
	"sort"
	"strings"

	"github.com/Jose-Merida-UVG/Compilers-1-CC3071/internal/yapar"
)

// Epsilon is the sentinel string representing ε in FIRST sets and production display.
// An empty Body slice in a Production means the rule derives ε.
const Epsilon = "ε"

// Symbol is a grammar symbol, either a terminal (token) or a non-terminal.
type Symbol struct {
	Name       string
	IsTerminal bool
	TokenID    int  // only meaningful for terminals; 0 for non-terminals
	Nullable   bool // true if this symbol can derive ε; set by ComputeFirst
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

// Build constructs a Grammar from a parsed YalpFile.
func Build(name string, yf *yapar.YalpFile) *Grammar {
	prods := make([]Production, 0)
	for _, p := range yf.Productions {
		for _, rule := range p.Rules {
			body := make([]Symbol, 0, len(rule))
			for _, sym := range rule {
				if id, ok := yf.TokenMap[sym]; ok {
					body = append(body, Symbol{Name: sym, IsTerminal: true, TokenID: id})
				} else {
					body = append(body, Symbol{Name: sym, IsTerminal: false})
				}
			}
			prods = append(prods, Production{Head: p.Name, Body: body})
		}
	}

	return &Grammar{
		Name:         name,
		StartSymbol:  yf.StartSymbol,
		Terminals:    yf.TokenMap,
		NonTerminals: yf.NonTerminals,
		Productions:  prods,
		IgnoreList:   yf.IgnoreList,
		First:        make(map[string]map[string]bool),
		Follow:       make(map[string]map[string]bool),
	}
}

// Summary returns human-readable lines describing the grammar.
// These are printed to the frontend terminal on "Build Parser".
func (g *Grammar) Summary() []string {
	var lines []string

	lines = append(lines, fmt.Sprintf("── Grammar: %s ──", g.Name))
	lines = append(lines, fmt.Sprintf("Start symbol : %s", g.StartSymbol))

	terms := sortedKeys(g.Terminals)
	lines = append(lines, fmt.Sprintf("Terminals    (%d): %s", len(terms), strings.Join(terms, "  ")))

	nonTerms := sortedBoolKeys(g.NonTerminals)
	lines = append(lines, fmt.Sprintf("Non-terminals(%d): %s", len(nonTerms), strings.Join(nonTerms, "  ")))

	if len(g.IgnoreList) > 0 {
		lines = append(lines, fmt.Sprintf("Ignored      : %s", strings.Join(g.IgnoreList, "  ")))
	}

	lines = append(lines, fmt.Sprintf("Productions  (%d):", len(g.Productions)))
	for _, p := range g.Productions {
		body := Epsilon // display ε for empty bodies
		if len(p.Body) > 0 {
			syms := make([]string, len(p.Body))
			for i, s := range p.Body {
				syms[i] = s.Name
			}
			body = strings.Join(syms, " ")
		}
		lines = append(lines, fmt.Sprintf("  %s → %s", p.Head, body))
	}

	g.ComputeFirst()
	g.ComputeFollow()

	lines = append(lines, "── FIRST sets ──")
	for _, nt := range nonTerms {
		if len(g.First[nt]) == 0 {
			lines = append(lines, fmt.Sprintf("  FIRST(%s) = { }", nt))
			continue
		}
		lines = append(lines, fmt.Sprintf("  FIRST(%s) = { %s }", nt, strings.Join(sortedBoolKeys(g.First[nt]), ", ")))
	}

	lines = append(lines, "── FOLLOW sets ──")
	for _, nt := range nonTerms {
		if len(g.Follow[nt]) == 0 {
			lines = append(lines, fmt.Sprintf("  FOLLOW(%s) = { }", nt))
			continue
		}
		lines = append(lines, fmt.Sprintf("  FOLLOW(%s) = { %s }", nt, strings.Join(sortedBoolKeys(g.Follow[nt]), ", ")))
	}

	return lines
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

func sortedKeys(m map[string]int) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func sortedBoolKeys(m map[string]bool) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
