# YALex — Lexer Generator

You write a `.yal` spec, click **Build** in the UI, and get a visual DFA you can
pan/zoom and a generated Go lexer file ready for the parser phase.

---

## What's Already Done

- Full `.yal` file parser — header, `let` definitions, rules, trailer, nested comments
- `let` identifier expansion (transitive, so `alnum = letter | digit` works)
- Regex pipeline: char classes, escapes, wildcards, `eof`, `+`, `?`, `*`, `|`, grouping
- Direct DFA construction + Hopcroft minimization
- Multi-pattern merge with first-rule-wins priority and per-token IDs
- Interactive DFA viewer in the browser (Graphviz WASM, pan/zoom, hover labels)
- HTTP server with file workspace (read, write, create, delete, rename)
- "Run Lexer" endpoint — greedy longest-match scan, outputs token list

## What's Left to Implement

### Code generation (`internal/codegen/codegen.go`)

This is **the main thing that needs to be built.** The function signature is:

```go
func GenerateLexer(name, sourcePath string, dfa *automata.DFA, actions []string) string
```

It returns the full text of a Go source file. The server already calls it and
writes the result to `workspace/lexers/<name>.go` — you just need to fill in the function.

Everything you need is in `dfa`:

```go
dfa.GetAllStates()         // []*DFAState
dfa.StartState             // the start state

// Each state:
state.ID                   // int
state.Accept               // bool
state.TokenID              // 1-based index into actions[] (0 = not accepting)
state.Transitions          // map[rune]*DFAState

actions[tokenID - 1]       // the raw action string from the .yal file
                           // e.g. "return PLUS", "return int(lxm)"
```

The generated file must:
- Be in `package lexers`
- Export `func <Name>Lex(input string) []Token`
- Track **line and column** as it scans
- Use **longest match** (maximal munch)
- Emit an error token for unrecognized characters rather than panicking

A `Token` should look roughly like this — you can adjust:

```go
type Token struct {
    Kind   int    // TokenID from the DFA (1-based). -1 for errors.
    Lexeme string // the matched text  (yytext)
    Value  any    // semantic value    (yyval) — set by action logic, can be nil
    Line   int    // 1-based line where the token starts
    Col    int    // 1-based column where the token starts
}
```

---

## Project Layout

```
internal/
  yalex/       .yal parsing + let-def expansion + compilation orchestration
  regex/       pattern normalization, shunting-yard, AST construction
  automata/    direct DFA construction, merge, Hopcroft minimization
  codegen/     ← your work goes here
  graph/       DFA → JSON for the frontend viewer
  ds/          generic stack and tree node

frontend/      React + Vite UI (DFA viewer, editor, file tree)
workspace/     runtime user files
  specs/       .yal files
  lexers/      generated .go and .dfa.json files (output of Build)
  input/       test input files
```

---

## Running Locally

```bash
make frontend-install   # first time only
make dev                # Go API on :8080, Vite hot-reload on :5173
```

Open `http://localhost:5173`. To do a production build: `make build && ./yalex`.

---

## .yal File Format

```
(* comments use (* ... *) and can be nested *)
{
  optional Go preamble — copied into generated file header
}

let digit  = ['0'-'9']
let letter = ['a'-'z'] | ['A'-'Z']
let alnum  = letter | digit        (* transitive — this works *)

rule gettoken =
    [' ' '\t']   { return lexbuf }
  | ['\n']       { return EOL }
  | digit+       { return int(lxm) }
  | '+'          { return PLUS }
  | eof          { raise('End of input') }

{
  optional Go trailer — appended to generated file
}
```

Supported pattern syntax: `'c'`, `'\n'`/`'\t'`/`'\r'`/`'\\'`/`'\''`, `"string"`,
`['a'-'z']`, `['a' 'b']`, `[^set]`, `[A] # [B]` (set diff), `_` (any printable ASCII),
`eof`, `x*`, `x+`, `x?`, `x | y`, `(x)`.

First matching rule wins.

---

## API Quick Reference

The frontend talks to the Go server over `/api/`. All file paths are relative to `workspace/`.

| Method | Path | What it does |
|--------|------|--------------|
| GET | `/api/workspace/tree` | Directory listing |
| GET | `/api/file?path=...` | Read a file |
| POST | `/api/file` `{path, content}` | Write a file |
| PUT | `/api/file` `{path}` | Create empty file |
| DELETE | `/api/file?path=...` | Delete file/dir |
| POST | `/api/file/rename` `{oldPath, newPath}` | Rename/move |
| PUT | `/api/directory` `{path}` | Create directory |
| POST | `/api/dfa` `{path}` | Build DFA from .yal → returns graph JSON + writes lexer |
| POST | `/api/lexer` `{inputPath}` | Run lexer on a file → returns token list |

`POST /api/lexer` infers the spec from the input file's extension:
`input/test.arithmetic` → `specs/arithmetic.yal`.
