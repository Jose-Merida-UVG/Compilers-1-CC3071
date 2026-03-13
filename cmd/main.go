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
// Regex 4: (a|b)*a(a|b)*a(a|b)* — needs at least two a's
// Regex 5: ((a|b)*)* — nested Kleene star
// Regex 6: (ab|ac)b — shared prefix
// Regex 7: (a*b*)*c — nested loops
// Regex 8: (a|b)*a(a|b)*b — needs at least one a and one b with a before b
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
	4: {
		{"aa", true}, // minimal: two a's
		{"a", false}, // only one a
	},
	5: {
		{"abab", true}, // nested star accepts
		{"", true},     // empty string accepted
	},
	6: {
		{"abb", true}, // ab then b
		{"acb", true}, // ac then b
	},
	7: {
		{"aabbbc", true}, // loops then c
		{"c", true},      // just c
	},
	8: {
		{"aab", true}, // a before b
		{"ab", true},  // minimal
	},
	9: {
		{"aaa", true}, // three a's
		{"bbb", true}, // three b's
	},
	10: {
		{"aaac", true}, // a* then c
		{"bbbc", true}, // b* then c
	},
	11: {
		{"aaa", true}, // a*a
		{"bb", true},  // b*b
	},
	12: {
		{"a", true},   // single a
		{"aaa", true}, // multiple a's
	},
	13: {
		{"ab", true},   // a then b
		{"aabb", true}, // aa then bb
	},
	14: {
		{"abb", true},    // contains abb
		{"aabbaa", true}, // contains abb in middle
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
		minDFA := dfa.Minimize() //minimización

		// Render DFA graphs
		outputDir := "out"
		dfaFile := fmt.Sprintf("DFA_regex%d", regexNum)
		minDfaFile := fmt.Sprintf("MinDFA_regex%d", regexNum)
		graph.BuildDFA(dfa, outputDir, dfaFile)
		graph.BuildDFA(minDFA, outputDir, minDfaFile)
		fmt.Printf("DFA     → %s/%s.pdf\n", outputDir, dfaFile)
		fmt.Printf("MinDFA  → %s/%s.pdf\n", outputDir, minDfaFile)

		fmt.Println("\n--- Comparación ---")
		fmt.Printf("DFA directo:     %d estados, %d transiciones\n",
			len(dfa.GetAllStates()), dfa.CountTransitions())

		fmt.Printf("DFA minimizado:  %d estados, %d transiciones\n",
			len(minDFA.GetAllStates()), minDFA.CountTransitions())

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

			fmt.Println("\n--- Simulación con DFA Minimizado ---")
			for _, tc := range cases {
				display := tc.input
				if display == "" {
					display = "ε"
				}

				result := minDFA.Sim(tc.input)
				status := "✓"
				if result != tc.expected {
					status = "✗ MISMATCH"
				}

				fmt.Printf("  %-16s → accepted=%-5t (expected %-5t) %s\n",
					display, result, tc.expected, status)
			}
		}

		regexNum++
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading file: %v\n", err)
	}
}
