package graph

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	//"github.com/goccy/go-graphviz"
)

// Custom dot file writer, only implementing functions needed with
// the specific format to graph Finite State Machines

func CreateHeader() string {
	dotContent := "digraph \"graph\" {\n\trankdir=LR size = \"8,5\"\n\tnode[shape=circle]\n"
	return dotContent
}

func AddNode(content string, nodeName string, label string) string {
	nodeText := fmt.Sprintf("\t%s [label=%s]\n", nodeName, label)
	content += nodeText
	return content
}

func AddAcceptNode(content string, nodeName string, label string) string {
	nodeText := fmt.Sprintf("\t%s [label=%s][shape=doublecircle]\n", nodeName, label)
	content += nodeText
	return content
}

func AddEdge(content string, fromNode string, toNode string, label string) string {
	transitionText := fmt.Sprintf("\t%s -> %s [label=\"%s\"]\n", fromNode, toNode, label)
	content += transitionText
	return content
}

// Ensure proper file and directory handling
func WriteFile(content string, dir string, fileName string) error {
	// Ensure the directory exists with full permissions
	err := os.MkdirAll(dir, 0777) // Full read/write/execute permissions for everyone
	if err != nil {
		return err
	}

	absDir, err := filepath.Abs(dir)
	if err != nil {
		return err
	}

	dotFilePath := filepath.Join(absDir, fileName+".DOT")
	content += "}"

	// Create the .DOT file
	file, err := os.Create(dotFilePath)
	if err != nil {
		return err
	}

	_, err = file.WriteString(content)
	if err != nil {
		return err
	}

	err = file.Close()
	if err != nil {
		return err
	}

	pdfFilePath := filepath.Join(absDir, fileName+".pdf")
	cmd := exec.Command("dot", "-Tpdf", dotFilePath, "-o", pdfFilePath)

	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Command execution failed: %s\n", err)
		fmt.Printf("Output: %s\n", string(output))
		return err
	}

	return nil
}
