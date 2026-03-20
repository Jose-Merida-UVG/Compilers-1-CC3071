package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
)

func main() {
	// Workspace is the directory the server reads/writes user files from.
	workspace := filepath.Join(".", "workspace")
	if err := os.MkdirAll(workspace, 0o755); err != nil {
		fmt.Fprintf(os.Stderr, "cannot create workspace: %v\n", err)
		os.Exit(1)
	}

	mux := http.NewServeMux()
	registerHandlers(mux, workspace)

	// Serve the built React frontend from frontend/dist.
	// API routes registered above take precedence because they are more specific.
	dist := filepath.Join(".", "frontend", "dist")
	if _, err := os.Stat(dist); err == nil {
		mux.Handle("/", http.FileServer(http.Dir(dist)))
	}

	addr := ":8080"
	fmt.Printf("YALex server listening on http://localhost%s\n", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		fmt.Fprintf(os.Stderr, "server error: %v\n", err)
		os.Exit(1)
	}
}
