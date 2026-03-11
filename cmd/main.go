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

// testCase holds a string to simulate and whether we expect it to be accepted.
type testCase struct {
	input    string
	expected bool
}

// testCases maps each regex (by line order) to its test strings.
// Regex 1: a+b|c*    — accepts (a+b) or (c*)
// Regex 2: (a|b)*cd? — accepts zero or more a/b, then c, then optional d
// Regex 3: a?(b|c)*d*e — accepts optional a, zero or more b/c, zero or more d, then e
var testCases = map[int][]testCase{
	1: {
		{"aaaab", true},  // long a+ chain followed by b
		{"ccc", true},    // c* with multiple c's
		{"aaabb", false}, // two b's after a+ — no match
		{"acb", false},   // interleaved chars, matches neither side
	},
	2: {
		{"bababacd", true}, // long alternating (a|b)*, then cd
		{"c", true},        // empty (a|b)*, just c, no d
		{"ababab", false},  // missing required c entirely
		{"abcdc", false},   // extra c after the optional d
	},
	3: {
		{"abcbcddde", true}, // all parts exercised: a, bcbc, ddd, e
		{"bcce", true},      // no a, mixed b/c repetition, straight to e
		{"abcbcd", false},   // everything present except the required final e
		{"aabce", false},    // two a's — a? only allows one
	},
}

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

		// Render AST and DFA graphs
		outputDir := "out"
		astFile := fmt.Sprintf("AST_regex%d", regexNum)
		dfaFile := fmt.Sprintf("DFA_regex%d", regexNum)
		graph.BuildASTGraph(&astRoot, outputDir, astFile)
		graph.BuildDFA(dfa, outputDir, dfaFile)
		fmt.Printf("AST  → %s/%s.pdf\n", outputDir, astFile)
		fmt.Printf("DFA  → %s/%s.pdf\n", outputDir, dfaFile)

		// Print NextPos table
		fmt.Println("\n--- NextPos Table ---")
		for id, nextSet := range table.NextPos {
			var nextList []int
			for nextID := range nextSet {
				nextList = append(nextList, nextID)
			}
			fmt.Printf("  ID %d ('%c'): %v\n", id, table.PosToChar[id], nextList)
		}
		fmt.Printf("\nStart State: %v\n", table.StartState)
		fmt.Printf("Accept ID:   %d\n", table.AcceptID)

		// Print transition table
		fmt.Println("\n--- DFA Transition Table ---")
		dfa.PrintTransitionTable()

		// Simulate test strings
		if cases, ok := testCases[regexNum]; ok {
			fmt.Println("\n--- String Simulation ---")
			for _, tc := range cases {
				display := tc.input
				if display == "" {
					display = "ε"
				}
				result := dfa.Sim(tc.input)
				status := "✓"
				if result != tc.expected {
					status = "✗ MISMATCH"
				}
				fmt.Printf("  %-16s → accepted=%-5t (expected %-5t) %s\n", display, result, tc.expected, status)
			}
		}

		regexNum++
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading file: %v\n", err)
	}
}
