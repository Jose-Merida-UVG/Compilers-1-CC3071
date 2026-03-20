# CONNECTION.md — Backend ↔ Frontend Contract

This document is the authoritative reference for how the Go HTTP server and the React frontend communicate. Read this before touching either side.

---

## Big Picture

```
Browser (React + Vite)
        │
        │  HTTP/JSON over localhost (port 8080)
        │  All routes prefixed /api
        ▼
Go net/http server  (main.go → handlers.go)
        │
        ├── File system ops  (workspace/ directory)
        ├── YALex parser     (internal/parser)
        ├── Automata engine  (internal/automata)   ← TODO: real DFA
        └── Code generator   (codegen.go)           ← TODO: real output
```

The frontend **never** touches the file system directly. Every read, write, create, delete goes through the API. The Go server keeps every path sandboxed inside `workspace/`.

---

## Workspace Directory Layout

```
workspace/
  specs/      .yal lexer specification files  (user-authored)
  input/      test input files                (user-authored)
  lexers/     generated .dfa and .go files    (written by Build)
  output/     .out result files               (written by Run)
```

**All paths sent from the frontend are relative to `workspace/`.**
Never send an absolute path. The server rejects any path that tries to escape the workspace with `..`.

---

## API Endpoints

### GET /api/health
Health check. Returns `{"status":"ok"}`. No request body.

---

### GET /api/workspace/tree
Returns the full workspace directory tree as a JSON array.

**Response** — `FileNode[]`
```json
[
  {
    "name": "specs",
    "path": "specs",
    "isDir": true,
    "children": [
      { "name": "arithmetic.yal", "path": "specs/arithmetic.yal", "isDir": false }
    ]
  }
]
```
The frontend calls this on startup and after every mutating operation to keep the sidebar in sync.

---

### GET /api/file?path=<relpath>
Read a file. `path` is URL-encoded, relative to workspace root.

**Response** — raw `text/plain` content (not JSON).

Used to populate Monaco editor tabs and to read `.dfa` files for the DFA viewer (parsed as JSON by the frontend).

---

### POST /api/file
Save (overwrite) a file.

**Request body**
```json
{ "path": "specs/arithmetic.yal", "content": "let digit = ..." }
```
**Response** — `{"ok": true}`

Creates parent directories automatically. This is what Ctrl+S triggers.

---

### PUT /api/file
Create a new empty file (touch).

**Request body**
```json
{ "path": "input/test.arithmetic" }
```
**Response** — `{"ok": true}`

---

### DELETE /api/file?path=<relpath>
Delete a file or directory (recursive).

**Response** — `{"ok": true}`

---

### POST /api/file/rename
Rename or move a file/directory.

**Request body**
```json
{ "oldPath": "specs/foo.yal", "newPath": "specs/bar.yal" }
```
**Response** — `{"ok": true}`

---

### PUT /api/directory
Create a directory (and all parents).

**Request body**
```json
{ "path": "input/subdir" }
```
**Response** — `{"ok": true}`

---

### POST /api/dfa  ← THE KEY ONE TO IMPLEMENT
Triggered when the user clicks **◎ Build** on a `.yal` file.

**Request body**
```json
{ "path": "specs/arithmetic.yal" }
```

**What it must do (currently placeholder):**
1. Read and parse the `.yal` file via `parser.ParseYALexContent`.
2. Build a combined DFA from **all** patterns across all rules (not just the first one).
3. Serialize that DFA as `DFAGraphData` (nodes + edges).
4. Write the serialized graph to `lexers/<name>.dfa` (JSON).
5. Generate a Go state machine and write it to `lexers/<name>.go`.
6. Return the `DFAGraphData` as the HTTP response.

**Response** — `DFAGraphData`
```json
{
  "nodes": [
    { "id": "q0", "label": "q0", "accept": false, "start": true },
    { "id": "q1", "label": "q1", "accept": true,  "start": false }
  ],
  "edges": [
    { "id": "e1", "source": "q0", "target": "q1", "label": "a" }
  ]
}
```

The frontend immediately renders this in the DFA viewer tab (React Flow graph). It also reads `lexers/<name>.go` as a second tab.

**Error responses:**
- `422 Unprocessable Entity` — `.yal` file failed to parse (returns `{"error": "..."}`)
- `404 Not Found` — `.yal` file does not exist
- `400 Bad Request` — bad request body

---

### POST /api/lexer  ← THE OTHER KEY ONE
Triggered when the user clicks **▶ Run** on an `input/` file.

**Request body**
```json
{ "inputPath": "input/test.arithmetic" }
```

**What it does:**
1. Reads the input file from `input/test.arithmetic`.
2. Infers the spec from the file extension: `.arithmetic` → loads `specs/arithmetic.yal`.
3. Parses the `.yal` spec.
4. For every rule pattern, builds a minimal DFA via the full regex pipeline.
5. Runs a greedy longest-match scan over the input.
6. Writes results to `output/test.arithmetic.out`.
7. Returns the lines as JSON.

**Response** — `LexerOutput`
```json
{
  "lines": [
    "MATCH  \"123\"               action=return INT",
    "MATCH  \"+\"                 action=return PLUS",
    "ERROR  unrecognized \"@\" at position 4"
  ]
}
```

Each line is either:
- `MATCH  <lexeme>  action=<action-body>` — successful token
- `ERROR  unrecognized <char> at position <n>` — no pattern matched

**Error responses:**
- `400` — no `inputPath`, or the file has no extension (can't infer spec)
- `404` — input file or spec file not found
- `422` — spec failed to parse

---

## How the Frontend Calls the API

All calls go through `frontend/src/api/index.ts`. It wraps `fetch`, prefixes `/api`, and throws on non-2xx.

```ts
api.listDirectory()                   // GET /api/workspace/tree
api.readFile(path)                    // GET /api/file?path=...
api.writeFile(path, content)          // POST /api/file
api.createFile(path)                  // PUT /api/file
api.deleteFile(path)                  // DELETE /api/file?path=...
api.renameFile(oldPath, newPath)      // POST /api/file/rename
api.createDirectory(path)             // PUT /api/directory
api.getDFA(yalPath)                   // POST /api/dfa  → DFAGraphData
api.runLexer(inputPath)               // POST /api/lexer → LexerOutput
```

---

## Data Flow: Build DFA

```
User clicks ◎ Build on specs/arithmetic.yal
        │
        ▼
App.tsx: buildDFA("specs/arithmetic.yal")
        │
        ├── api.getDFA("specs/arithmetic.yal")
        │         │
        │         ▼  POST /api/dfa {"path":"specs/arithmetic.yal"}
        │         Go: parse .yal → build DFA → write lexers/arithmetic.dfa
        │                                     → write lexers/arithmetic.go
        │         Returns: DFAGraphData
        │
        ├── refreshTree()  ← tree now shows lexers/arithmetic.dfa and .go
        │
        ├── Opens tab "arithmetic" with dfaData set → renders DFA viewer
        └── Opens tab "arithmetic.go" with Go source → renders Monaco
```

---

## Data Flow: Run Lexer

```
User clicks ▶ Run on input/test.arithmetic
        │
        ▼
App.tsx: runFile("input/test.arithmetic")
        │
        ├── If tab is dirty: api.writeFile(...) first
        │
        ├── api.runLexer("input/test.arithmetic")
        │         │
        │         ▼  POST /api/lexer {"inputPath":"input/test.arithmetic"}
        │         Go: read input → load specs/arithmetic.yal
        │             → build per-pattern DFAs → greedy scan
        │             → write output/test.arithmetic.out
        │         Returns: { lines: [...] }
        │
        ├── Each line appended to terminal pane
        └── refreshTree()  ← output/ now shows the .out file
```

---

## The YALex File Format

Specs live in `specs/*.yal`. The parser (`internal/parser/yalexfile.go`) produces a `YALexFile`:

```
(* Optional header in { } *)

let digit  = ['0'-'9']         ← named regex macro
let letter = ['a'-'z'|'A'-'Z']

rule tokens =                  ← rule entrypoint name
  digit+           { return INT   }   ← pattern  { action }
| letter digit*    { return IDENT }
| '+'              { return PLUS  }
| eof              { raise EOF    }
```

**Parsed into:**
```go
type YALexFile struct {
    Header      string          // Go code between first { }
    Definitions []LetDefinition // let name = regex
    Rules       []Rule          // rule name = patterns
    Trailer     string          // Go code between last { }
}
type LetDefinition struct { Identifier, Regex string }
type Rule          struct { Entrypoint string; Patterns []RulePattern }
type RulePattern   struct { Pattern, Action string }
```

When building DFAs, `let` macro names are expanded inline:
`digit+` → `(['0'-'9'])+` before feeding into the regex pipeline.

---

## Regex → DFA Pipeline

Lives entirely in `internal/`. The steps in order:

```
RegexString (raw pattern string)
   │
   ├── HandleCharClasses()         ['a'-'z'] → (a|b|c|...|z)
   ├── HandleSpecialOperators()    ? → (x|ε),  + → (xx*)
   ├── HandleExplicitConcatenation()  insert ~ between adjacent tokens
   ├── ShuntingYard()              infix → postfix (Dijkstra)
   ├── AppendEndMarker()           append # at end
   │
   ▼
BuildDirectAST(postfix)          annotated syntax tree (Nullable/FirstPos/LastPos)
   │
   ▼
BuildDirectTable(ast)            compute NextPos table  →  DirectTable
   │
   ▼
DirectTable.ToDFA()              subset construction  →  *DFA
   │
   ▼
DFA.Minimize()                   Hopcroft partition refinement  →  *DFA (minimal)
```

**Internal operator:** concatenation is represented as `~` (tilde).
**Epsilon:** stored as rune `'ε'` (U+03B5).

---

## DFA Struct → Graph JSON

`internal/graph/dfa.go: SerializeDFA(*automata.DFA) *DFAGraphData`

Walks the DFA state graph (BFS via `Transitions` map on each `DFAState`) and produces the `DFAGraphData` struct that is both returned from `/api/dfa` and written to `lexers/<name>.dfa`.

The frontend reads this and renders it with **React Flow** (nodes as circles, edges as smoothstep curves, self-loops as a custom arc). Layout is computed by **Dagre** in LR (left-to-right) mode before rendering.

---

## What Still Needs Real Implementation

| Location | Current state | What it should do |
|---|---|---|
| `handlers.go: getDFA` | Returns `placeholderDFA()` — a hardcoded 7-state graph | Build a **combined DFA** from all patterns in the `.yal`, serialize it via `graph.SerializeDFA`, write files |
| `handlers.go: getDFA` (codegen) | Writes a stub comment `.go` file | Call `generateGoLexer(specBase, path, realDFA)` with the real minimized DFA |
| `runLexer` | Fully implemented | Already builds per-pattern DFAs and runs greedy scan — this works today |

### Implementing the real getDFA

The `runLexer` handler already shows the correct pattern. Adapt it:

```go
// 1. Parse the spec
yal, err := parser.ParseYALexContent(string(content))

// 2. For each pattern, expand macros then run the full pipeline
for _, rule := range yal.Rules {
    for _, pat := range rule.Patterns {
        regexStr := pat.Pattern
        for _, def := range yal.Definitions {
            regexStr = strings.ReplaceAll(regexStr, def.Identifier, "("+def.Regex+")")
        }
        rs := parser.NewRegexString(regexStr)
        rs.HandleCharClasses()
        rs.HandleSpecialOperators()
        rs.HandleExplicitConcatenation()
        rs.ShuntingYard()
        rs.AppendEndMarker()
        root := parser.BuildDirectAST(rs)
        table := automata.BuildDirectTable(&root)
        dfa := table.ToDFA().Minimize()
        // dfa is now ready for serialization or combination
    }
}

// 3. Serialize and return
serialized := graph.SerializeDFA(combinedDFA)
```

To build a **single combined DFA** from multiple patterns, wrap all patterns in an alternation before running the pipeline, or merge the individual DFAs after construction. The simplest correct approach is to prefix each pattern with a tag position and build one big regex: `(pat1)|(pat2)|(pat3)` — the accept states in the result each correspond to a token class.
