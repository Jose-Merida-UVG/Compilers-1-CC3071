package main

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/Jose-Merida-UVG/Compilers-1-CC3071/internal/yapar"
)

func check(name, path string) {
	content, err := os.ReadFile(path)
	if err != nil { fmt.Printf("ERROR reading %s: %v\n", path, err); return }
	yf, err := yapar.ParseYalpContent(string(content))
	if err != nil { fmt.Printf("PARSE ERROR %s: %v\n", name, err); return }
	compiled, err := yf.Compile(name)
	if err != nil { fmt.Printf("COMPILE ERROR %s: %v\n", name, err); return }
	g := compiled.Grammar
	fmt.Printf("\n══════ %s ══════\n", strings.ToUpper(name))
	fmt.Printf("Start: %s   Productions: %d   States: %d   Conflicts: %d\n",
		g.StartSymbol, len(g.Productions), len(compiled.Automaton.States), len(compiled.Table.Conflicts))
	for _, c := range compiled.Table.Conflicts { fmt.Printf("  ⚠ %s\n", c) }
	nts := sortedKeys(g.NonTerminals)
	fmt.Println("\nFIRST:")
	for _, nt := range nts { fmt.Printf("  %-20s { %s }\n", nt, joinSet(g.First[nt])) }
	fmt.Println("\nFOLLOW:")
	for _, nt := range nts { fmt.Printf("  %-20s { %s }\n", nt, joinSet(g.Follow[nt])) }
}

func sortedKeys(m map[string]bool) []string {
	k := make([]string, 0, len(m))
	for s := range m { k = append(k, s) }
	sort.Strings(k)
	return k
}
func joinSet(m map[string]bool) string { return strings.Join(sortedKeys(m), ", ") }

func main() {
	check("arithmetic", "workspace/specs/arithmetic.yalp")
	check("calc",       "workspace/specs/calc.yalp")
	check("gol",        "workspace/specs/gol.yalp")
}
