# YALex — Lexer Generator

You write a `.yal` spec, click **Build** in the UI, and get a visual DFA you can
pan/zoom and a generated Go lexer file ready for the parser phase.

---

## What's Done

- Full `.yal` file parser — header, `let` definitions, rules, trailer, nested comments
- `let` identifier expansion (transitive, so `alnum = letter | digit` works)
- Regex pipeline: char classes, escapes, wildcards, `eof`, `+`, `?`, `*`, `|`, grouping
- Direct DFA construction + Hopcroft minimization
- Multi-pattern merge with first-rule-wins priority and per-token IDs
- Interactive DFA viewer in the browser (Graphviz WASM, pan/zoom, hover labels)
- HTTP server with file workspace (read, write, create, delete, rename)
- "Run Lexer" endpoint — greedy longest-match scan, outputs token list
- **Code generation** (`internal/codegen/codegen.go`) — produces a standalone Go lexer file from the minimized DFA

---

## Code Generation

`GenerateLexer` in `internal/codegen/codegen.go` takes the minimized DFA and produces a
complete, self-contained Go source file under `workspace/lexers/<name>.go`.

```go
func GenerateLexer(name, sourcePath string, dfa *automata.DFA, actions []string) string
```

The server calls this automatically on every **Build** — you never invoke it directly.

### What the generated file contains

| Section | Description |
|---|---|
| `Token` struct | `Kind`, `Lexeme`, `Value`, `Line`, `Col` |
| `<name>Trans` | Static DFA transition table: `map[int]map[rune]int` |
| `<name>Accept` | Static accept table: `map[int]int` (stateID → 1-based TokenID) |
| `yalexInt` | Helper that wraps `strconv.Atoi` — only emitted when `int(lxm)` is used |
| `<Name>Lex` | Exported scanner function: `func <Name>Lex(input string) []Token` |

### How the scanner works

The generated `<Name>Lex` uses **maximal munch** (greedy longest match):

1. From the current position, speculatively advance through the DFA following transitions.
2. Every time an accepting state is reached, snapshot `(pos, line, col)` and the `TokenID`.
3. When no transition exists, backtrack to the last snapshot and commit that match.
4. Run the action for that `TokenID`, advance `pos` to the snapshot, repeat.
5. If no accepting state was ever reached, emit an **error token** (`Kind: -1`) for the
   single unrecognised character and keep scanning — never panics.

Line and column are tracked incrementally during the speculative scan and restored from
the snapshot on backtrack, so positions are always correct even after lookahead.

### Action translation

The `{action}` block in a `.yal` rule is translated to Go as follows:

| `.yal` action | Generated Go |
|---|---|
| `return lexbuf` | `// skip` — lexeme consumed, no token emitted |
| `return IDENT` | `Token{Kind: N, Lexeme: lxm, Value: IDENT, ...}` |
| `return int(lxm)` | `Token{Kind: N, Lexeme: lxm, Value: yalexInt(lxm), ...}` |
| `raise('...')` | Emits `Token{Kind: -2, ...}` then `return tokens` |
| *(empty or comment)* | Falls through to default — emits a basic token, raw action as comment |

`lxm` is a `string` variable in scope at every action site holding the matched text.
Identifiers like `PLUS`, `EOL`, `FLOAT` must be defined in the `.yal` header block
(or an imported package) so the generated file compiles.

### Token Kind conventions

| Kind | Meaning |
|---|---|
| `1..N` | 1-based index of the matching rule (first rule = 1) |
| `-1` | Unrecognised character (scanner keeps going) |
| `-2` | `raise(...)` — EOF or fatal condition, scanner stops |

### Writing a `.yal` file that generates clean Go

```
{
const (
    FLOAT  = 1   (* must match the rule order below — first rule = 1 *)
    INT    = 2
    IDENT  = 3
    PLUS   = 4
    (* ... *)
)
}

let digit  = ['0'-'9']
let letter = ['a'-'z'] | ['A'-'Z']

rule gettoken =
    [' ' '\t' '\n']  { return lexbuf }      (* skip — no token emitted *)
  | digit+ '.' digit+ { return FLOAT }      (* float before int — longest match *)
  | digit+            { return INT }
  | letter (letter | digit)* { return IDENT }
  | '+'               { return PLUS }
  | eof               { raise('End of input') }
```

Things to watch out for:

- **Comments as actions** — `{ (* skip *) }` is stripped to an empty action before
  codegen sees it. Use `{ return lexbuf }` to explicitly skip whitespace.
- **`raise` syntax** — must be `raise('...')` with parentheses and a quoted string.
  Bare `raise EOF` is not recognised.
- **Constant names** — `PLUS`, `EOL`, etc. referenced in actions must be defined in the
  header block. If they are missing the generated file will not compile.
- **Rule order for overlapping patterns** — put longer/more-specific rules first
  (e.g. `float` before `int`) since first-rule-wins applies when two patterns match
  the same length.
- **ASCII only** — the `.yal` parser rejects any character outside `0x20–0x7E`
  (plus `\t`, `\n`, `\r`). Watch for em dashes (`—`) or curly quotes copied from
  rich-text editors.
- **Reserved Go identifiers as let names** — avoid names like `int`, `string`, `type`;
  they are valid YALex identifiers but can confuse the word-boundary substitution in
  `expandDefinitions`.

---

## Project Layout

```
internal/
  yalex/       .yal parsing + let-def expansion + compilation orchestration
  regex/       pattern normalization, shunting-yard, AST construction
  automata/    direct DFA construction, merge, Hopcroft minimization
  codegen/     GenerateLexer — DFA + actions → Go source file
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
|---|---|---|
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