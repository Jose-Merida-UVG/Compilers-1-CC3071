package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/TonitoMC/TDLCGo/internal/automata"
	"github.com/TonitoMC/TDLCGo/internal/ds"
	"github.com/TonitoMC/TDLCGo/internal/graph"
	"github.com/TonitoMC/TDLCGo/internal/regex"
)

func main() {
	// Read regular expressions from data/regex.txt
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

		// Skip empty lines
		if inputStr == "" {
			continue
		}

		fmt.Printf("\n========================================\n")
		fmt.Printf("Regex #%d: %s\n", regexNum, inputStr)
		fmt.Printf("========================================\n")

		// Create RegexString and process it
		rs := regex.NewRegexString(inputStr)
		rs.HandleSpecialOperators()
		rs.HandleCharClasses()
		rs.HandleExplicitConcatenation()
		rs.ShuntingYard()
		rs.AppendEndMarker()

		fmt.Printf("Postfix with '#': %s\n", rs.String())

		// Build AST
		astRoot := regex.BuildDirectAST(rs)

		// Compute attributes and build the DirectTable context
		table := automata.BuildDirectTable(&astRoot)

		dfa := table.ToDFA()

		// Render annotated syntax tree to out/AST_regex<N>.pdf
		outputDir := "out"
		outputFile := fmt.Sprintf("AST_regex%d", regexNum)
		graph.BuildASTGraph(&astRoot, outputDir, outputFile)
		fmt.Printf("Annotated AST written to %s/%s.pdf\n", outputDir, outputFile)

		// Print Tree
		fmt.Println("\n--- AST Attributes (Post-Order) ---")
		printTree(&astRoot, "", true)

		// Print NextPos Table
		fmt.Println("\n--- NextPos Table ---")
		for id, nextSet := range table.NextPos {
			var nextList []int
			for nextID := range nextSet {
				nextList = append(nextList, nextID)
			}
			// Sort for consistent output, useful for debugging
			// slices.Sort(nextList) // Only if Go 1.21+ and slices package is imported
			char := table.PosToChar[id]
			fmt.Printf("ID %d ('%c'): %v\n", id, char, nextList)
		}
		fmt.Printf("\nStart State (FirstPos of root): %v\n", table.StartState)
		fmt.Printf("Accept ID ('#'): %d\n", table.AcceptID)

		// Print transition table
		dfa.PrintTransitionTable()
		fmt.Println("\n--- Simulación DFA Directo, se genera un pdf ---")

		regexNum++
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading file: %v\n", err)
	}
}

func printTree(node *ds.Node[regex.ASTNodeData], prefix string, isTail bool) {
	if node == nil {
		return
	}

	// Determine branch prefix
	branch := "├── "
	if isTail {
		branch = "└── "
	}

	// Format node info
	val := string(node.Value.Value)
	if node.Value.Value == 'ε' {
		val = "ε"
	}
	idStr := ""
	if node.Value.ID != 0 {
		idStr = fmt.Sprintf(" [ID: %d]", node.Value.ID)
	}

	info := fmt.Sprintf("%s%s Null: %t, First: %v, Last: %v",
		val, idStr, node.Value.Nullable, node.Value.FirstPos, node.Value.LastPos)

	fmt.Printf("%s%s%s\n", prefix, branch, info)

	// Adjust prefix for children
	childPrefix := prefix
	if isTail {
		childPrefix += "    "
	} else {
		childPrefix += "│   "
	}

	if node.Left != nil || node.Right != nil {
		printTree(node.Left, childPrefix, node.Right == nil)
		if node.Right != nil {
			printTree(node.Right, childPrefix, true)
		}
	}
}
