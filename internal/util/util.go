package util

import (
	"fmt"
	"log"

	"github.com/TonitoMC/TDLCGo/internal/automata"
	"github.com/TonitoMC/TDLCGo/internal/graph"
	"github.com/TonitoMC/TDLCGo/internal/regex"
)

var RegularExpressionCounter int

type RegularExpression struct {
	Regex        string
	NFA          *automata.NFA
	DFA          *automata.DFA
	MinimizedDFA *automata.DFA
}

func NewRegularExpression(input string, id int) *RegularExpression {
	outputDir := fmt.Sprintf("out/regex%d", id)
	fmt.Printf("Expresion Regular %d: %s\n", id, input)
	rs := regex.NewRegexString(input)
	if !regex.Balanced(*rs) {
		log.Fatal("Expresion regular no balanceada, invalida")
	}
	rs.HandleCharClasses()
	rs.HandleSpecialOperators()
	rs.HandleExplicitConcatenation()
	rs.ShuntingYard()
	fmt.Printf("Postfix Utilizando Shunting Yard: %s\n", rs.String())

	nfa, err := automata.ThompsonConstruction(rs)
	if err != nil {
		log.Fatal(err)
	}
	graph.BuildGraph(nfa, outputDir, "NFA")

	dfa := nfa.ToDFA()
	graph.BuildDFA(dfa, outputDir, "DFA")

	minimizedDFA := dfa.Minimize()
	graph.BuildDFA(minimizedDFA, outputDir, "Minimized DFA")

	return &RegularExpression{
		Regex:        input,
		NFA:          nfa,
		DFA:          dfa,
		MinimizedDFA: minimizedDFA,
	}
}

func (rx *RegularExpression) Simulate(input string) {
	fmt.Printf("Simulando: %s\n", input)
	fmt.Printf("Simulacion Utilizando NFA: %t\n", rx.NFA.Sim(input))
	fmt.Printf("Simulacion Utilizando DFA: %t\n", rx.DFA.Sim(input))
	fmt.Printf("Simulacion Utilizando DFA Minimizado: %t\n", rx.MinimizedDFA.Sim(input))
}
