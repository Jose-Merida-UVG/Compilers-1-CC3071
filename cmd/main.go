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

		// Render DFA graph to out/DFA_regex<N>.pdf
		dfaOutputFile := fmt.Sprintf("DFA_regex%d", regexNum)
		graph.BuildDFA(dfa, outputDir, dfaOutputFile)
		fmt.Printf("DFA graph written to %s/%s.pdf\n", outputDir, dfaOutputFile)

		// Print transition table
		fmt.Println("\n--- DFA Transition Table ---")
		dfa.PrintTransitionTable()

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
