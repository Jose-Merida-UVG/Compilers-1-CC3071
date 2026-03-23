# YALex — Lexical Analyzer Generator

A Go-based lexer generator inspired by OCamllex. You write a `.yal` spec file,
click **◎ Build** in the UI, and get a minimized DFA visualization plus a
self-contained Go lexer that scans input using maximal munch.

---

## Software Architecture

### Overview

```
┌─────────────────────────────────────────────────────────────┐
│                     Browser (React + Vite)                   │
│                                                              │
│  FileTree  ──→  EditorPane  ──→  DFAViewer / MarkdownViewer │
│                      │                                       │
│              toolbar: ◎ Build, ▶ Run                         │
└──────────────────────┬───────────────────────────────────────┘
                       │ HTTP /api/*
┌──────────────────────▼───────────────────────────────────────┐
│                   Go HTTP Server (main.go)                    │
│                                                              │
│  apiHandler                                                  │
│    getDFA   ──→ yalex.Compile ──→ codegen.GenerateLexer      │
│    runLexer ──→ go run workspace/lexers/<name>.go <input>    │
│    file CRUD, workspace tree                                 │
└──────────────────────────────────────────────────────────────┘
```

### Pipeline: `.yal` → tokens

```
.yal file
   │
   ▼ internal/yalex
   Parse header / let defs / rule blocks / trailer
   Expand let references (transitive substitution)
   │
   ▼ internal/regex
   Normalize pattern strings → postfix (shunting-yard)
   Build NFA via Thompson construction
   │
   ▼ internal/automata
   Direct DFA construction (subset construction)
   Merge per-rule DFAs with first-rule-wins priority
   Hopcroft minimization
   │
   ├──▶ internal/graph   → JSON for DFA viewer (Graphviz WASM)
   │
   └──▶ internal/codegen → workspace/lexers/<name>.go
                           (package main, Scan() int, func main())
```

### Go Packages

| Package | Responsibility |
|---|---|
| `internal/yalex` | Parse `.yal` files; expand `let` definitions; orchestrate compilation |
| `internal/regex` | Pattern normalization, shunting-yard to postfix, AST node types |
| `internal/automata` | NFA→DFA (direct method), multi-DFA merge, Hopcroft minimization |
| `internal/codegen` | Emit complete `package main` Go source from minimized DFA + actions |
| `internal/graph` | Serialize DFA to Graphviz/JSON for the frontend visualizer |
| `internal/ds` | Generic stack and tree-node used by the regex/automata phases |

### Frontend Components

| Component | Role |
|---|---|
| `App.tsx` | State root — open tabs, active file, file-save debounce (200 ms) |
| `FileTree` | Workspace directory browser; create/rename/delete files & folders |
| `EditorPane` | Tab bar + Monaco editor; routes `.md` files to MarkdownViewer, DFA tabs to DFAViewer |
| `DFAViewer` | Interactive Graphviz WASM renderer with pan/zoom and hover labels |
| `MarkdownViewer` | Renders `.md` files as styled HTML via `react-markdown` + `remark-gfm` |
| `Terminal` | Displays output lines from the **▶ Run** action |

### HTTP API

| Method | Path | Description |
|---|---|---|
| `GET` | `/api/workspace/tree` | Directory listing |
| `GET` | `/api/file?path=...` | Read file contents |
| `POST` | `/api/file` `{path, content}` | Write file |
| `PUT` | `/api/file` `{path}` | Create empty file |
| `DELETE` | `/api/file?path=...` | Delete file or directory |
| `POST` | `/api/file/rename` `{oldPath, newPath}` | Rename / move |
| `PUT` | `/api/directory` `{path}` | Create directory |
| `POST` | `/api/dfa` `{path}` | Build DFA from `.yal` → returns graph JSON + writes lexer |
| `POST` | `/api/lexer` `{inputPath}` | Run generated lexer on input → returns token lines |

`POST /api/lexer` infers the spec from the input file extension:
`input/test.arithmetic` → `specs/arithmetic.yal` → `lexers/arithmetic.go`.

---

## Project Layout

```
.
├── main.go               entry point — HTTP server + handler wiring
├── handlers.go           getDFA, runLexer, file CRUD handlers
├── Makefile
├── workspace/
│   ├── specs/            .yal spec files (source)
│   ├── lexers/           generated .go lexers + .dfa.json (output of Build)
│   ├── input/            test input files
│   └── YALEX.md          YALex spec reference (opens in MarkdownViewer)
├── internal/
│   ├── yalex/            parser + compiler orchestration
│   ├── regex/            pattern → postfix → AST
│   ├── automata/         DFA construction + minimization
│   ├── codegen/          DFA + actions → Go source
│   ├── graph/            DFA → JSON for viewer
│   └── ds/               generic stack / tree node
└── frontend/
    ├── src/
    │   ├── App.tsx
    │   ├── components/
    │   │   ├── Editor/
    │   │   ├── FileTree/
    │   │   ├── DFAViewer/
    │   │   ├── MarkdownViewer/
    │   │   └── Terminal/
    │   └── lib/
    │       └── monaco-yal.ts   YALex syntax highlighting for Monaco
    └── vite.config.ts
```

---

## Running Locally

```bash
make frontend-install   # first time only
make dev                # Go API :8080 + Vite hot-reload :5173
```

Open `http://localhost:5173`. Production build: `make build && ./yalex`.

---

## Tech Stack

| Layer | Technology |
|---|---|
| Backend | Go 1.22, stdlib `net/http` |
| Frontend | React 18, TypeScript, Vite |
| Editor | Monaco Editor (`@monaco-editor/react`) |
| DFA visualization | Graphviz WASM (`@hpcc-js/wasm-graphviz`), React Flow |
| Markdown rendering | `react-markdown` + `remark-gfm` |
| Generated lexers | Pure Go — no runtime dependencies |

---

## How the Generated Lexer Works

The generated `workspace/lexers/<name>.go` is a complete `package main` program.

```go
type Lexer struct {
    Lxm string   // matched lexeme
    Ln  int       // 1-based line
    Col int       // 1-based column
    // internal fields
}

func New<Name>Lexer(input string) *Lexer
func (l *Lexer) Scan() int   // returns token ID; 0 = EOF, -1 = ERROR
```

**Scanning rules:**

1. **Maximal munch** — always matches the longest possible string.
2. **First-rule-wins** — when two patterns match the same length, the earlier rule in the `.yal` file wins.
3. **Error recovery** — consecutive unrecognised characters are grouped into a single `ERROR` token; scanning never stops.

Actions in the `.yal` file are verbatim Go. An action that contains `return` exits `Scan()` and delivers the token; an action without `return` loops and picks up the next token (standard way to skip whitespace).

For the full `.yal` syntax reference, open `workspace/YALEX.md` in the editor.
