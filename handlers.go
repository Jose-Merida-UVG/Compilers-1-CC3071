package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	"github.com/Jose-Merida-UVG/Compilers-1-CC3071/internal/yalex"
	"github.com/Jose-Merida-UVG/Compilers-1-CC3071/internal/yalex/codegen"
	"github.com/Jose-Merida-UVG/Compilers-1-CC3071/internal/yalex/graph"
	"github.com/Jose-Merida-UVG/Compilers-1-CC3071/internal/yapar"
)


type apiHandler struct {
	workspace string
}

func registerHandlers(mux *http.ServeMux, workspace string) {
	h := &apiHandler{workspace: workspace}

	mux.HandleFunc("GET /api/health", h.health)
	mux.HandleFunc("GET /api/workspace/tree", h.listDirectory)
	mux.HandleFunc("GET /api/file", h.readFile)
	mux.HandleFunc("POST /api/file", h.writeFile)
	mux.HandleFunc("PUT /api/file", h.createFile)
	mux.HandleFunc("DELETE /api/file", h.deleteFile)
	mux.HandleFunc("POST /api/file/rename", h.renameFile)
	mux.HandleFunc("PUT /api/directory", h.createDirectory)
	mux.HandleFunc("POST /api/dfa", h.getDFA)
	mux.HandleFunc("POST /api/yapar", h.buildParser)
	mux.HandleFunc("POST /api/run", h.runFile)
}

// ─── Helpers ──────────────────────────────────────────────────────────────────

func (h *apiHandler) resolve(relPath string) (string, error) {
	clean := filepath.Join(h.workspace, filepath.Clean("/"+relPath))
	if !strings.HasPrefix(clean, h.workspace) {
		return "", fmt.Errorf("path escapes workspace")
	}
	return clean, nil
}

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, code int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": msg})
}

func decode[T any](r *http.Request) (T, error) {
	var v T
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		return v, err
	}
	return v, nil
}

// ─── File-system handlers ─────────────────────────────────────────────────────

type FileNode struct {
	Name     string     `json:"name"`
	Path     string     `json:"path"`
	IsDir    bool       `json:"isDir"`
	Children []FileNode `json:"children,omitempty"`
}

func (h *apiHandler) health(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, map[string]string{"status": "ok"})
}

func (h *apiHandler) listDirectory(w http.ResponseWriter, r *http.Request) {
	nodes, err := h.buildTree(h.workspace, "")
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, nodes)
}

func (h *apiHandler) buildTree(base, relBase string) ([]FileNode, error) {
	entries, err := os.ReadDir(base)
	if err != nil {
		return nil, err
	}
	var nodes []FileNode
	for _, entry := range entries {
		rel := entry.Name()
		if relBase != "" {
			rel = relBase + "/" + entry.Name()
		}
		node := FileNode{Name: entry.Name(), Path: rel, IsDir: entry.IsDir()}
		if entry.IsDir() {
			node.Children, _ = h.buildTree(filepath.Join(base, entry.Name()), rel)
		}
		nodes = append(nodes, node)
	}
	return nodes, nil
}

func (h *apiHandler) readFile(w http.ResponseWriter, r *http.Request) {
	relPath := r.URL.Query().Get("path")
	full, err := h.resolve(relPath)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	data, err := os.ReadFile(full)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Write(data)
}

func (h *apiHandler) writeFile(w http.ResponseWriter, r *http.Request) {
	body, err := decode[struct {
		Path    string `json:"path"`
		Content string `json:"content"`
	}](r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	full, err := h.resolve(body.Path)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	os.MkdirAll(filepath.Dir(full), 0o755)
	if err := os.WriteFile(full, []byte(body.Content), 0o644); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, map[string]bool{"ok": true})
}

func (h *apiHandler) createFile(w http.ResponseWriter, r *http.Request) {
	body, err := decode[struct {
		Path string `json:"path"`
	}](r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	full, err := h.resolve(body.Path)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	os.MkdirAll(filepath.Dir(full), 0o755)
	f, err := os.Create(full)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	f.Close()
	writeJSON(w, map[string]bool{"ok": true})
}

func (h *apiHandler) deleteFile(w http.ResponseWriter, r *http.Request) {
	relPath := r.URL.Query().Get("path")
	full, err := h.resolve(relPath)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := os.RemoveAll(full); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, map[string]bool{"ok": true})
}

func (h *apiHandler) renameFile(w http.ResponseWriter, r *http.Request) {
	body, err := decode[struct {
		OldPath string `json:"oldPath"`
		NewPath string `json:"newPath"`
	}](r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	oldFull, err := h.resolve(body.OldPath)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	newFull, err := h.resolve(body.NewPath)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	os.MkdirAll(filepath.Dir(newFull), 0o755)
	if err := os.Rename(oldFull, newFull); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, map[string]bool{"ok": true})
}

func (h *apiHandler) createDirectory(w http.ResponseWriter, r *http.Request) {
	body, err := decode[struct {
		Path string `json:"path"`
	}](r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	full, err := h.resolve(body.Path)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := os.MkdirAll(full, 0o755); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, map[string]bool{"ok": true})
}

// ─── YALex / DFA handlers ─────────────────────────────────────────────────────

func (h *apiHandler) getDFA(w http.ResponseWriter, r *http.Request) {
	body, err := decode[struct {
		Path string `json:"path"`
	}](r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	full, err := h.resolve(body.Path)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	content, err := os.ReadFile(full)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	yalFile, err := yalex.ParseYalContent(string(content))
	if err != nil {
		writeError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}
	compiled, err := yalFile.Compile()
	if err != nil {
		writeError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}
	serialized := graph.SerializeDFA(compiled.DFA)

	specBase := strings.TrimSuffix(filepath.Base(body.Path), ".yal")
	lexerDir := filepath.Join(h.workspace, "programs", specBase, "lexer")
	docsDir  := filepath.Join(h.workspace, "programs", specBase, "docs")
	os.MkdirAll(lexerDir, 0o755)
	os.MkdirAll(docsDir, 0o755)

	// Write DFA graph JSON to docs/ for the frontend visualizer.
	if dfaJSON, err := json.Marshal(serialized); err == nil {
		os.WriteFile(filepath.Join(docsDir, "dfa.json"), dfaJSON, 0o644)
	}

	// Write lexer.go to lexer/.
	goSrc := codegen.GenerateLexer(specBase, yalFile, compiled.DFA, compiled.Actions)
	os.WriteFile(filepath.Join(lexerDir, "lexer.go"), []byte(goSrc), 0o644)

	writeJSON(w, serialized)
}

// ─── YAPar handlers ──────────────────────────────────────────────────────────

func (h *apiHandler) buildParser(w http.ResponseWriter, r *http.Request) {
	body, err := decode[struct {
		YalpPath string `json:"yalpPath"`
	}](r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	yalpFull, err := h.resolve(body.YalpPath)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	yalpContent, err := os.ReadFile(yalpFull)
	if err != nil {
		writeError(w, http.StatusNotFound, fmt.Sprintf("cannot read %s: %v", body.YalpPath, err))
		return
	}

	yalpFile, err := yapar.ParseYalpContent(string(yalpContent))
	if err != nil {
		writeError(w, http.StatusUnprocessableEntity, fmt.Sprintf("yapar: %v", err))
		return
	}

	specBase := strings.TrimSuffix(filepath.Base(body.YalpPath), ".yalp")

	// Require lexer/lexer.go to exist before building the parser.
	lexerPath  := filepath.Join(h.workspace, "programs", specBase, "lexer", "lexer.go")
	parserDir  := filepath.Join(h.workspace, "programs", specBase, "parser")
	docsDir    := filepath.Join(h.workspace, "programs", specBase, "docs")
	if _, err := os.Stat(lexerPath); err != nil {
		writeError(w, http.StatusUnprocessableEntity, fmt.Sprintf("lexer not found for %q — build the lexer first", specBase))
		return
	}
	os.MkdirAll(parserDir, 0o755)
	os.MkdirAll(docsDir, 0o755)

	compiled, err := yalpFile.Compile(specBase)
	if err != nil {
		writeError(w, http.StatusUnprocessableEntity, fmt.Sprintf("yapar compile: %v", err))
		return
	}

	// Write LR(0) and SLR visualizer data to docs/.
	if lr0JSON, err := json.Marshal(compiled.Automaton.Serialize()); err == nil {
		os.WriteFile(filepath.Join(docsDir, "lr0.json"), lr0JSON, 0o644)
	}

	// Write parser.go stub to parser/.
	os.WriteFile(filepath.Join(parserDir, "parser.go"), []byte("package main\n\n// TODO: parser codegen\n"), 0o644)



	// Build structured SLR table data for the frontend viewer.
	terminals := sortedStringKeys(compiled.Grammar.Terminals)
	terminals = append(terminals, "$")
	nonTerminals := sortedBoolKeys(compiled.Grammar.NonTerminals)

	actions := make(map[int]map[string]string)
	for sid, row := range compiled.Table.Action {
		actions[sid] = make(map[string]string)
		for sym, act := range row {
			actions[sid][sym] = act.String()
		}
	}

	type conflictDTO struct {
		State    int    `json:"state"`
		Symbol   string `json:"symbol"`
		Existing string `json:"existing"`
		Incoming string `json:"incoming"`
	}
	conflicts := make([]conflictDTO, len(compiled.Table.Conflicts))
	for i, c := range compiled.Table.Conflicts {
		conflicts[i] = conflictDTO{c.State, c.Symbol, c.Existing.String(), c.Incoming.String()}
	}

	// Build FIRST/FOLLOW sets.
	first := make(map[string][]string)
	follow := make(map[string][]string)
	for _, nt := range nonTerminals {
		first[nt] = sortedBoolKeys(compiled.Grammar.First[nt])
		follow[nt] = sortedBoolKeys(compiled.Grammar.Follow[nt])
	}

	// Build production list.
	type prodDTO struct {
		Head string `json:"head"`
		Body string `json:"body"`
	}
	prods := make([]prodDTO, len(compiled.Grammar.Productions))
	for i, p := range compiled.Grammar.Productions {
		syms := make([]string, len(p.Body))
		for j, s := range p.Body {
			syms[j] = s.Name
		}
		body := "ε"
		if len(syms) > 0 {
			body = strings.Join(syms, " ")
		}
		prods[i] = prodDTO{p.Head, body}
	}

	slrPayload := map[string]any{
		"terminals":    terminals,
		"nonTerminals": nonTerminals,
		"actions":      actions,
		"gotos":        compiled.Table.Goto,
		"conflicts":    conflicts,
		"first":        first,
		"follow":       follow,
		"productions":  prods,
		"startSymbol":  compiled.Grammar.StartSymbol,
		"stateCount":   len(compiled.Automaton.States),
	}

	// Save slr.json to docs/ for the frontend viewer.
	if slrJSON, err := json.Marshal(slrPayload); err == nil {
		os.WriteFile(filepath.Join(docsDir, "slr.json"), slrJSON, 0o644)
	}

	writeJSON(w, slrPayload)
}

// runFile serves POST /api/run.
// Runs programs/<spec>/*.go against the input file.
func (h *apiHandler) runFile(w http.ResponseWriter, r *http.Request) {
	body, err := decode[struct {
		InputPath string `json:"inputPath"`
	}](r)
	if err != nil || body.InputPath == "" {
		writeError(w, http.StatusBadRequest, "inputPath required")
		return
	}

	inputFull, err := h.resolve(body.InputPath)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	base := filepath.Base(body.InputPath)
	ext := strings.TrimPrefix(filepath.Ext(base), ".")
	if ext == "" {
		writeError(w, http.StatusBadRequest, "cannot infer spec: file has no extension")
		return
	}

	// Collect .go files from lexer/ and parser/ subdirectories.
	var goFiles []string
	for _, sub := range []string{"lexer", "parser"} {
		dir := filepath.Join(h.workspace, "programs", ext, sub)
		entries, err := os.ReadDir(dir)
		if err != nil {
			continue // subdirectory may not exist yet
		}
		for _, e := range entries {
			if !e.IsDir() && strings.HasSuffix(e.Name(), ".go") {
				goFiles = append(goFiles, filepath.Join(dir, e.Name()))
			}
		}
	}
	if len(goFiles) == 0 {
		writeError(w, http.StatusNotFound, fmt.Sprintf("no program found for %q — build the lexer first", ext))
		return
	}

	args := append([]string{"run"}, goFiles...)
	args = append(args, inputFull)
	out, err := exec.Command("go", args...).Output()
	if err != nil {
		if ee, ok := err.(*exec.ExitError); ok {
			writeError(w, http.StatusUnprocessableEntity, string(ee.Stderr))
		} else {
			writeError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	lines := strings.Split(strings.TrimRight(string(out), "\n"), "\n")

	os.MkdirAll(filepath.Join(h.workspace, "output"), 0o755)
	os.WriteFile(filepath.Join(h.workspace, "output", base+".out"), out, 0o644)

	writeJSON(w, map[string]any{"lines": lines})
}

func sortedStringKeys(m map[string]int) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func sortedBoolKeys(m map[string]bool) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
