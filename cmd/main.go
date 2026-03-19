package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/TonitoMC/TDLCGo/internal/automata"
	"github.com/TonitoMC/TDLCGo/internal/graph"
	"github.com/TonitoMC/TDLCGo/internal/regex"
)

func main() {
	file, err := os.Open("data/regex.txt")
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	regexNum := 1

	for scanner.Scan() {
		inputStr := strings.TrimSpace(scanner.Text())
		if inputStr == "" {
			continue
		}

		fmt.Printf("\n========================================\n")
		fmt.Printf("Regex #%d: %s\n", regexNum, inputStr)
		fmt.Printf("========================================\n")

		// Process regex
		rs := regex.NewRegexString(inputStr)
		rs.HandleSpecialOperators()
		rs.HandleCharClasses()
		rs.HandleExplicitConcatenation()
		rs.ShuntingYard()
		rs.AppendEndMarker()
		fmt.Printf("Postfix with '#': %s\n", rs.String())

		// Build AST and direct table
		astRoot := regex.BuildDirectAST(rs)
		table := automata.BuildDirectTable(&astRoot)
		dfa := table.ToDFA()
		minDFA := dfa.Minimize() //minimización

		// Render DFA graphs
		outputDir := "out"
		dfaFile := fmt.Sprintf("DFA_regex%d", regexNum)
		minDfaFile := fmt.Sprintf("MinDFA_regex%d", regexNum)
		graph.BuildDFA(dfa, outputDir, dfaFile)
		graph.BuildDFA(minDFA, outputDir, minDfaFile)
		fmt.Printf("DFA     → %s/%s.pdf\n", outputDir, dfaFile)
		fmt.Printf("MinDFA  → %s/%s.pdf\n", outputDir, minDfaFile)

		fmt.Println("\n--- DFA Directo ---")
		fmt.Printf("States: %d\n", len(dfa.GetAllStates()))
		fmt.Printf("Transitions: %d\n", dfa.CountTransitions())
		dfa.PrintTransitionTable()

		fmt.Println("\n--- DFA Minimizado ---")
		fmt.Printf("States: %d\n", len(minDFA.GetAllStates()))
		fmt.Printf("Transitions: %d\n", minDFA.CountTransitions())
		minDFA.PrintTransitionTable()

		regexNum++
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading file: %v\n", err)
	}
}
