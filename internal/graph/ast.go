// ast.go renders an annotated syntax tree (ASTNodeData) as a Graphviz PDF.
// Each node shows {firstpos} on the left, the symbol in the centre, {lastpos}
// on the right, and T / F (nullable) in a second row below.
package graph

import (
	"fmt"
	"sort"
	"strings"

	"github.com/TonitoMC/TDLCGo/internal/ds"
	"github.com/TonitoMC/TDLCGo/internal/regex"
)

// BuildASTGraph produces a top-down annotated tree diagram and writes it to
// <dir>/<filename>.pdf via Graphviz.  Errors from the DOT pipeline are printed
// to stdout (matching the behaviour of BuildGraph / BuildDFA).
func BuildASTGraph(root *ds.Node[regex.ASTNodeData], dir string, filename string) {
	if root == nil {
		return
	}
	content := astHeader()
	counter := 0
	content, _ = addASTNode(root, &counter, content)
	WriteFile(content, dir, filename)
}

// astHeader returns the DOT preamble for a top-down tree layout.
func astHeader() string {
	return "digraph \"ast\" {\n\trankdir=TB\n\tnode [shape=plaintext fontname=\"Courier\"]\n"
}

// addASTNode recursively emits the node declaration and outgoing edges for
// node, then for its children.  It returns the updated DOT content and the
// integer ID that was assigned to this node so the caller can draw an edge.
// Node IDs are assigned in pre-order via the shared counter.
func addASTNode(node *ds.Node[regex.ASTNodeData], counter *int, content string) (string, int) {
	myID := *counter
	*counter++

	label := astNodeLabel(node.Value)
	content += fmt.Sprintf("\tn%d [label=<%s>]\n", myID, label)

	if node.Left != nil {
		var leftID int
		content, leftID = addASTNode(node.Left, counter, content)
		content += fmt.Sprintf("\tn%d -> n%d\n", myID, leftID)
	}

	if node.Right != nil {
		var rightID int
		content, rightID = addASTNode(node.Right, counter, content)
		content += fmt.Sprintf("\tn%d -> n%d\n", myID, rightID)
	}

	return content, myID
}

// astNodeLabel builds an HTML-table label for a single AST node.
//
// Visual layout inside the node box:
//
//	┌──────────────────────────────┐
//	│  {firstpos}  │  {lastpos}   │
//	│  sym         │  T / F       │
//	└──────────────────────────────┘
func astNodeLabel(data regex.ASTNodeData) string {
	firstStr := formatPosSet(data.FirstPos)
	lastStr := formatPosSet(data.LastPos)
	sym := htmlEscape(displaySymbol(data.Value))
	nullable := "F"
	if data.Nullable {
		nullable = "T"
	}

	return fmt.Sprintf(
		`<TABLE BORDER="1" CELLBORDER="0" CELLSPACING="3">`+
			`<TR>`+
			`<TD ALIGN="CENTER"><FONT COLOR="navy">%s</FONT></TD>`+
			`<TD ALIGN="CENTER"><FONT COLOR="darkgreen">%s</FONT></TD>`+
			`</TR>`+
			`<TR>`+
			`<TD COLSPAN="2" ALIGN="CENTER"><B>%s</B> &nbsp; <I>%s</I></TD>`+
			`</TR>`+
			`</TABLE>`,
		firstStr, lastStr, sym, nullable,
	)
}

// formatPosSet converts a slice of ints into "{1,2,3}" notation.
// The slice is sorted before formatting so the output is deterministic.
func formatPosSet(positions []int) string {
	if len(positions) == 0 {
		return "{}"
	}
	sorted := make([]int, len(positions))
	copy(sorted, positions)
	sort.Ints(sorted)

	parts := make([]string, len(sorted))
	for i, p := range sorted {
		parts[i] = fmt.Sprintf("%d", p)
	}
	return "{" + strings.Join(parts, ",") + "}"
}

// displaySymbol maps internal operator runes to human-readable symbols.
// The concatenation operator is stored as '~' internally but is conventionally
// displayed as '·' in syntax-tree diagrams.
func displaySymbol(r rune) string {
	switch r {
	case '~':
		return "\u00b7" // middle dot  ·
	default:
		return string(r)
	}
}

// htmlEscape escapes the characters that are special inside Graphviz HTML labels.
func htmlEscape(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	s = strings.ReplaceAll(s, "\"", "&quot;")
	return s
}
