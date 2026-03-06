package main

import (
	"fmt"

	"github.com/TonitoMC/TDLCGo/internal/automata"
	"github.com/TonitoMC/TDLCGo/internal/ds"
	"github.com/TonitoMC/TDLCGo/internal/regex"
	"github.com/TonitoMC/TDLCGo/internal/graph"
)

func main() {
	inputStr := "(a|b)*abb"
	fmt.Printf("Testing regex: %s\n", inputStr)

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

	graph.BuildDFA(dfa, "out/direct", "DirectDFA")

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

	dfa.PrintTransitionTable()
	fmt.Println("\n--- Simulación DFA Directo, se genera un pdf ---")

	//obtener un string de consola
	var input string
	fmt.Print("\nIngrese una cadena: ")
	fmt.Scanln(&input)

	if dfa.Sim(input) {
		fmt.Println("La cadena PERTENECE al lenguaje.")
	} else {
		fmt.Println("La cadena NO pertenece al lenguaje.")
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
