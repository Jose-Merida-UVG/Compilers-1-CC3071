# YAPar

A parser generator written in Go, inspired by Yacc and ocamlyacc. You write a `.yalp` grammar spec and a paired `.yal` lexer spec, click **◎ Build Parser**, and get a working SLR(1) parser you can run against input right in the browser — along with visualizers for the LR(0) automaton, SLR(1) table, and lexer DFA.

YALex (the lexer generator) is the foundation YAPar is built on. It is also available standalone for building and visualizing DFAs without a grammar.

---

## Project Layout

```
.
├── main.go                     HTTP server entry point + workspace dir setup
├── handlers.go                 All API handlers (build, run, file CRUD)
├── internal/
│   ├── yalex/                  .yal parser, let-def expansion, top-level compile
│   │   ├── yalex.go            YalFile struct, ParseYalContent
│   │   ├── scanner.go          .yal tokenizer
│   │   ├── compile.go          orchestrates per-pattern DFA build + merge
│   │   ├── regex/              pattern → postfix → syntax tree
│   │   │   ├── normalize.go    pattern syntax → flat token list (RegexString)
│   │   │   ├── regexstring.go  RegexString type + explicit concatenation pass
│   │   │   ├── shuntingyard.go infix → postfix (Dijkstra's shunting-yard)
│   │   │   └── syntaxtree.go   postfix → AST with nullable/firstpos/lastpos
│   │   ├── automata/           DFA construction and minimization
│   │   │   ├── direct.go       direct DFA from syntax tree (followpos method)
│   │   │   ├── compile.go      per-pattern compile + merge of all patterns
│   │   │   ├── dfa.go          DFA / DFAState types
│   │   │   └── minimization.go Hopcroft minimization
│   │   ├── codegen/            Go source generation
│   │   │   ├── codegen.go      GenerateCombined — lexer without main()
│   │   │   └── parsergen.go    GenerateParser — ACTION/GOTO tables + Parse()
│   │   └── graph/              DFA → JSON for frontend visualizer
│   ├── yapar/                  .yalp parser, grammar, LR(0), SLR(1)
│   │   ├── yapar.go            YalpFile struct, ParseYalpContent
│   │   ├── compile.go          orchestrates grammar → automaton → table
│   │   ├── validate_tokens.go  cross-validates .yal and .yalp token sets
│   │   ├── grammar/
│   │   │   ├── grammar.go      Grammar struct (productions, terminals, non-terminals)
│   │   │   ├── first.go        FIRST set computation
│   │   │   └── follow.go       FOLLOW set computation
│   │   ├── lr0/
│   │   │   ├── automata.go     LR(0) automaton + state/item types
│   │   │   ├── closure.go      closure(I) computation
│   │   │   ├── goto.go         goto(I, X) computation
│   │   │   └── graph.go        automaton → JSON for frontend visualizer
│   │   └── slr/
│   │       └── slr.go          SLR(1) ACTION/GOTO table construction + conflict detection
│   └── ds/
│       ├── stack.go            generic stack
│       └── tree.go             generic tree node
├── cmd/
│   └── checkgrammars/          standalone CLI for batch grammar checking
├── frontend/                   React + Vite
│   └── src/
│       ├── api/                typed wrappers around the REST API
│       ├── components/
│       │   ├── Editor/         Monaco editor with .yal syntax highlighting
│       │   ├── Sidebar/        workspace file browser
│       │   ├── DFAViewer/      Graphviz WASM renderer (pan/zoom)
│       │   ├── LR0Viewer/      Graphviz WASM renderer for LR(0) automaton
│       │   ├── SLRViewer/      ACTION/GOTO table + FIRST/FOLLOW + conflict list
│       │   ├── MarkdownViewer/ renders .md files in-editor
│       │   ├── StatusBar/      bottom bar (current file, build status)
│       │   └── Terminal/       output panel for ▶ Run
│       └── types/              shared TypeScript types
└── workspace/
    ├── specs/                  .yal + .yalp pairs for real grammars
    ├── demo/                   .yal + .yalp pairs for small test/demo grammars
    ├── input/                  test input files (extension = spec name)
    ├── output/                 .out files written after each Run
    └── programs/
        └── <name>/
            ├── lexer.go        generated lexer (token constants, DFA scan loop)
            ├── parser.go       generated parser (tables, Parse(), main())
            └── docs/           dfa.json, lr0.json, slr.json for visualizers
```

---

## How It Works

### Full pipeline

```
specs/arithmetic.yal + specs/arithmetic.yalp
│
├─▶ yalex.ParseYalContent        parse .yal → YalFile (rules, let-defs, actions)
├─▶ yapar.ParseYalpContent       parse .yalp → YalpFile (tokens, productions)
├─▶ yapar.ValidateTokenCoverage  cross-check: every %token ↔ a lexer return
│
├─▶ yalFile.Compile              .yal → minimized combined DFA + action list
│     regex.Normalize            pattern syntax → RegexString token list
│     regex.ExplicitConcat       insert ~ between adjacent tokens
│     regex.ShuntingYard         infix → postfix
│     automata.DirectDFA         syntax tree → DFA (followpos method)
│     automata.Minimize          Hopcroft minimization per pattern
│     automata.Merge             parallel simulation → single combined DFA
│     automata.Minimize          minimize the merged DFA
│
├─▶ yalpFile.Compile             grammar → LR(0) automaton → SLR(1) table
│     grammar.Build              extract terminals/non-terminals, augment grammar
│     grammar.ComputeFirst       FIRST sets (nullable-aware)
│     grammar.ComputeFollow      FOLLOW sets
│     lr0.BuildAutomaton         closure + goto → full LR(0) item set automaton
│     slr.BuildTable             ACTION/GOTO from automaton + FOLLOW sets
│
├─▶ codegen.GenerateCombined     → lexer.go  (DFA tables, NextToken, Lexeme)
├─▶ codegen.GenerateParser       → parser.go (tables, Parse(), main())
│
└─▶ workspace/programs/arithmetic/
        lexer.go
        parser.go
        docs/dfa.json, lr0.json, slr.json
```

---

## YAPar: Parser Generator

### Grammar construction (`internal/yapar/grammar/`)

`grammar.Build` reads the productions from `YalpFile`, augments the grammar with a new start production `S' → S`, and classifies every symbol as terminal or non-terminal.

`ComputeFirst` computes FIRST sets bottom-up, propagating through nullable symbols (those that can derive ε). `ComputeFollow` uses the FIRST sets and the augmented start production to compute FOLLOW sets.

### LR(0) automaton (`internal/yapar/lr0/`)

`closure(I)` expands a set of LR(0) items by repeatedly adding `B → • γ` for every item `A → α • B β` where the dot precedes a non-terminal.

`goto(I, X)` computes the state reached from item set `I` on symbol `X` — shifts the dot past `X` in every applicable item and takes the closure of the result.

`BuildAutomaton` runs the standard worklist algorithm: starting from the closure of `{S' → • S}`, it applies goto for every grammar symbol reachable from each state and collects all states until no new ones are discovered.

### SLR(1) table (`internal/yapar/slr/slr.go`)

For each LR(0) state, fills the ACTION and GOTO tables:

- **Shift** — if `goto(s, a)` is defined for terminal `a`, add `shift goto(s, a)` to `ACTION[s][a]`.
- **Reduce** — if item `A → α •` is in state `s`, add `reduce A → α` to `ACTION[s][t]` for every `t` in `FOLLOW(A)`.
- **Accept** — if `S' → S •` is in state `s`, add `accept` to `ACTION[s][$]`.

Any cell written twice is a conflict. Conflicts are recorded (not fatal) and surfaced in the SLR viewer tab. A grammar with conflicts will still generate a parser, but it will behave unpredictably in the conflicting states.

### Generated parser (`internal/yalex/codegen/parsergen.go`)

`GenerateParser` emits `parser.go` as part of `package main` alongside the lexer:

- `parserActionTable` and `parserGotoTable` — full SLR(1) tables as Go map literals
- `parserProds` — head name, body string, and body length for each production
- `parserIgnore` — token IDs to silently discard before the grammar sees them
- `Parse(l *Lexer) error` — the shift/reduce loop; prints each action + current sentential form, then a top-down derivation and symbol table on accept
- `tokenIDToName()` — reverse map from int ID to grammar name, for display
- `main()` — reads the input file, runs lexer + parser, exits with error on failure

---

## YALex: Lexer Generator

### Regex → DFA

**Normalization (`internal/yalex/regex/normalize.go`)**
Converts YALex pattern syntax — character classes (`['a'-'z' '0'-'9']`), string literals, char literals, `+`/`?`/`*`, named `let` substitutions — into a flat `RegexString` token list.

**Explicit concatenation (`internal/yalex/regex/regexstring.go`)**
Inserts explicit `~` (concat) operators between adjacent tokens so the grammar is fully explicit before shunting-yard runs.

**Shunting-yard (`internal/yalex/regex/shuntingyard.go`)**
Dijkstra's algorithm converts the infix token stream to postfix (RPN) respecting operator precedence: `* + ?` > `~` > `|`.

**Syntax tree + position functions (`internal/yalex/regex/syntaxtree.go`)**
Builds an AST from the postfix stream and computes `nullable`, `firstpos`, `lastpos`, and `followpos` on every node — the inputs to direct DFA construction.

**Direct DFA construction (`internal/yalex/automata/direct.go`)**
Builds the DFA directly from the position sets without going through an NFA. Each DFA state is a set of positions; transitions follow `followpos`. A unique end-marker position appended to each pattern determines the accepting state's token ID.

**Hopcroft minimization (`internal/yalex/automata/minimization.go`)**
Reduces the DFA to its minimal equivalent by iteratively refining partitions of indistinguishable states.

**Merge (`internal/yalex/automata/compile.go`)**
All per-pattern DFAs are merged into one using parallel simulation: each merged state is a tuple of per-pattern states. When a tuple accepts for multiple patterns, the lowest-indexed pattern wins (first-rule-wins). The merged DFA is minimized again.

### Generated lexer (`internal/yalex/codegen/codegen.go`)

`GenerateCombined` emits `lexer.go` without a `main()` (the parser owns that):

- Token ID constants from `%token` declarations
- `Lexeme` struct (`Token int`, `Value string`, `Line int`, `Col int`)
- `Lexer` struct with internal scan state (`input`, `pos`, `line`, `col`)
- Static transition table and accept table derived from the minimized DFA
- `New<Name>Lexer(input string) *Lexer`
- `<rulename>() int` — raw scan method with maximal munch and ERROR grouping
- `NextToken() Lexeme` — public wrapper used by `Parse()`

---

## API

All paths are sandboxed to `workspace/` — any path escaping it is rejected.

| Method | Path | Body / Query | What it does |
|--------|------|-------------|--------------|
| `GET` | `/api/health` | — | Liveness check |
| `GET` | `/api/workspace/tree` | — | Recursive file tree |
| `GET` | `/api/file` | `?path=` | Read a file |
| `POST` | `/api/file` | `{path, content}` | Write a file |
| `PUT` | `/api/file` | `{path}` | Create empty file |
| `DELETE` | `/api/file` | `?path=` | Delete file or directory |
| `POST` | `/api/file/rename` | `{oldPath, newPath}` | Rename |
| `PUT` | `/api/directory` | `{path}` | Create directory |
| `POST` | `/api/dfa` | `{path}` | Build DFA from `.yal`, return serialized DFA |
| `POST` | `/api/yapar` | `{yalpPath}` | Build parser from `.yalp` + paired `.yal`, return SLR payload |
| `POST` | `/api/run` | `{inputPath}` | `go run programs/<ext>/*.go <input>`, return output lines |

`POST /api/run` infers the spec name from the input file's extension (`test.arithmetic` → `programs/arithmetic/`).

---

## Running Locally

```bash
# Install frontend dependencies (once)
cd frontend && npm install && cd ..

# Development
cd frontend && npm run dev &   # Vite dev server on :5173
go run .                       # API server on :8080

# Production
cd frontend && npm run build && cd ..
go run .    # serves built frontend + API on :8080
```
