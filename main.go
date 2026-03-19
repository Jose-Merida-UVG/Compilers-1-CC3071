package main

import (
	"fmt"
	"os"

	"github.com/Jose-Merida-UVG/Compilers-1-CC3071/parser"
)

func main() {
	// Check if a .yal file was provided as argument
	if len(os.Args) < 2 {
		fmt.Println("Usage: yalex <file.yal>")
		fmt.Println("\nNo file provided. Testing with default example file...")
		testDefaultFile()
		return
	}

	filename := os.Args[1]

	// Parse the YALex file
	fmt.Printf("Reading YALex file: %s\n\n", filename)
	yalFile, err := parser.ParseYALexFile(filename)
	if err != nil {
		fmt.Printf("Error parsing file: %v\n", err)
		os.Exit(1)
	}

	// Print the parsed structure
	yalFile.PrintStructure()
}

func testDefaultFile() {
	// Test with the example file
	filename := "data/ejemplo.yal"
	fmt.Printf("Reading YALex file: %s\n\n", filename)

	yalFile, err := parser.ParseYALexFile(filename)
	if err != nil {
		fmt.Printf("Error parsing file: %v\n", err)
		os.Exit(1)
	}

	// Print the parsed structure
	yalFile.PrintStructure()
}
