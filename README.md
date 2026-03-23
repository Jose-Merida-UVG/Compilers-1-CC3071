# YALex

A lexical analyzer generator written in Go. You write a `.yal` spec describing
your token patterns, click **◎ Build**, and get a minimized DFA you can explore
visually plus a standalone Go lexer file you can run against any input.

---

## Project Layout

```
.
├── main.go               HTTP server entry point
├── handlers.go           API handlers (build DFA, run lexer, file CRUD)
├── internal/
│   ├── yalex/            .yal parser + let-def expansion + compilation
│   ├── regex/            pattern normalization → postfix → token list
│   ├── automata/         direct DFA construction, merge, Hopcroft minimization
│   ├── codegen/          DFA + actions → self-contained Go source file
│   ├── graph/            DFA → JSON for the frontend visualizer
│   └── ds/               generic stack and tree node
├── frontend/             React + Vite
│   └── src/components/
│       ├── Editor/       Monaco editor with YALex syntax highlighting
│       ├── FileTree/     workspace file browser
│       ├── DFAViewer/    Graphviz WASM renderer (pan/zoom, hover labels)
│       ├── MarkdownViewer/  renders .md files in-editor
│       └── Terminal/     output panel for ▶ Run
└── workspace/
    ├── specs/            your .yal files go here
    ├── lexers/           generated .go files land here after Build
    └── input/            test input files
```

---

## How It Works

### 1. Parsing the spec

`internal/yalex` reads a `.yal` file in four sections:

```
{ header }          optional verbatim Go — constants, imports
let ident = regex   named pattern definitions (transitively expanded)
rule name =         one or more pattern → action arms
  pattern { action }
{ trailer }         optional verbatim Go — helper functions
```

Let definitions are expanded transitively before any DFA is built. If you write:

```
let digit = ['0'-'9']
let int   = digit+
```

then every occurrence of `int` in rule patterns is substituted with `((['0'-'9'])+)`
before the regex engine sees it.

### 2. Regex → postfix → DFA (per pattern)

Each expanded pattern goes through `internal/regex`:

- **Normalize** — converts the YALex pattern syntax (`['a'-'z']`, char literals,
  `+`, `?`, `|`, grouping) into a flat token list (`RegexString`).
- **Explicit concatenation** — inserts `~` operators between adjacent tokens so
  the grammar is fully explicit.
- **Shunting-yard** — converts infix token list to postfix.
- **End-marker append** — appends a unique sentinel rune (`\uE000`) so the DFA
  knows where each pattern ends.

The postfix `RegexString` is then handed to `internal/automata`:

- **Direct method** — builds the DFA directly from the syntax tree without going
  through an NFA. Computes `nullable`, `firstpos`, `lastpos`, and `followpos` on
  the AST, then runs subset construction on position sets.
- **Hopcroft minimization** — reduces the DFA to its minimal equivalent.

### 3. Merging all patterns

`automata.Merge` takes one minimized DFA per pattern and builds a single combined
DFA using **parallel simulation**: each state in the merged DFA is a tuple of
per-pattern state IDs. When a tuple is accepting, the lowest-indexed pattern wins
(**first-rule-wins**). The merged DFA is minimized again.

### 4. Code generation

`internal/codegen.GenerateLexer` emits a complete `package main` Go file:

- The user's header block verbatim (your token constants go here).
- Static transition table (`map[int]map[rune]int`) and accept table
  (`map[int]int`) derived from the minimized DFA.
- A `Lexer` struct with `Lxm string`, `Ln int`, `Col int` fields.
- `New<Name>Lexer(input string) *Lexer` constructor.
- `Scan() int` — advances to the next token, returns its ID. Actions are
  embedded verbatim; a `return` in an action exits `Scan()`, no `return`
  keeps scanning (how you skip whitespace).
- `func main()` so the file can be run with `go run`.

### 5. Running the lexer

When you click **▶ Run**, the backend does:

```
go run workspace/lexers/<name>.go workspace/input/<your-file>
```

Output lines come back to the terminal panel in the UI.

---

## Pipeline at a glance

```
arithmetic.yal
│
│  let digit = ['0'-'9']
│  let float = digit+ '.' digit+
│  rule gettoken =
│    float  { return FLOAT }
│  | '+'    { return PLUS }
│
▼ yalex.Parse + expandDefinitions
│
│  pattern 1: ((['0'-'9'])+)'.'((['0'-'9'])+)
│  pattern 2: '+'
│
▼ regex.Preprocess (normalize → explicit concat → shunting-yard → postfix)
▼ automata.Compile (direct method → minimize)   — one DFA per pattern
▼ automata.Merge (parallel simulation → minimize) — one combined DFA
▼ codegen.GenerateLexer
│
└─▶ workspace/lexers/arithmetic.go   (package main, Scan() int)
└─▶ DFA graph JSON                   (rendered in browser)
```

---

## Example

`workspace/specs/arithmetic.yal`:

```
{
    const (
        FLOAT  = 1
        INT    = 2
        PLUS   = 3
        MINUS  = 4
        TIMES  = 5
        DIV    = 6
        MOD    = 7
        LPAREN = 8
        RPAREN = 9
    )
}

let digit = ['0'-'9']
let int   = digit+
let float = digit+ '.' digit+
let ws    = [' ' '\t' '\n' '\r']

rule gettoken =
  ws+    {  }
| float  { return FLOAT }
| int    { return INT }
| '+'    { return PLUS }
| '-'    { return MINUS }
| '*'    { return TIMES }
| '/'    { return DIV }
| '%'    { return MOD }
| '('    { return LPAREN }
| ')'    { return RPAREN }
```

Running `workspace/input/test_a.arithmetic` through the generated lexer:

```
FLOAT  "3.14"   ln=1 col=1
PLUS   "+"      ln=1 col=5
INT    "2"      ln=1 col=6
TIMES  "*"      ln=1 col=8
LPAREN "("      ln=1 col=9
...
```

---

## Running Locally

```bash
make frontend-install   # once
make dev                # Go API :8080 + Vite :5173
```

Open `http://localhost:5173`. Production: `make build && ./yalex`.
