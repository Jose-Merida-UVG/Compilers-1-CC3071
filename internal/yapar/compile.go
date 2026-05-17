package yapar

import (
	"fmt"

	"github.com/Jose-Merida-UVG/Compilers-1-CC3071/internal/yapar/grammar"
	"github.com/Jose-Merida-UVG/Compilers-1-CC3071/internal/yapar/lr0"
	"github.com/Jose-Merida-UVG/Compilers-1-CC3071/internal/yapar/slr"
)

// CompiledParser is the result of compiling a YalpFile: the augmented grammar,
// the canonical LR(0) automaton, and the SLR(1) parse table derived from it.
type CompiledParser struct {
	Grammar   *grammar.Grammar
	Automaton *lr0.Automaton
	Table     *slr.Table
}

// Compile builds the augmented grammar, computes FIRST/FOLLOW sets, constructs
// the LR(0) automaton, and derives the SLR(1) parse table. This is the main
// orchestration method, calling methods from grammar/, slr/ and lr0/ packages.
func (f *YalpFile) Compile(specName string) (*CompiledParser, error) {
	// Grammar package uses RawProductions (identical to the normal ones) except
	// we declare them to avoid circular dependencies
	prods := make([]grammar.RawProduction, len(f.Productions))
	for i, p := range f.Productions {
		prods[i] = grammar.RawProduction{Name: p.Name, Rules: p.Rules}
	}

	// Build grammar
	g := grammar.Build(grammar.YalpSpec{
		Name:         specName,
		StartSymbol:  f.StartSymbol,
		TokenMap:     f.TokenMap,
		NonTerminals: f.NonTerminals,
		Productions:  prods,
		IgnoreList:   f.IgnoreList,
	})

	// Compute first & follow
	g.ComputeFirst()
	g.ComputeFollow()

	// Augment + build lr0 automaton
	augmented := g.Augment()
	automaton := lr0.New(augmented)
	automaton.Build()

	// Build SLR table and fail fast if the grammar has conflicts.
	table := slr.Build(automaton)
	if len(table.Conflicts) > 0 {
		return nil, fmt.Errorf("grammar is not SLR(1): %d conflict(s) — first: %s", len(table.Conflicts), table.Conflicts[0])
	}

	return &CompiledParser{
		Grammar:   g,
		Automaton: automaton,
		Table:     table,
	}, nil
}
