# YALex — Lexer Generator

A web-based lexer generator. You write a `.yal` spec file, click **Build** in the
UI, and get back a visual DFA diagram plus a generated Go lexer ready for the
parser phase.

---

## Project Layout

```
.
├── main.go                      Entry point — HTTP server on :8080
├── handlers.go                  API route handlers (file I/O, DFA build, lexer run)
├── Makefile                     Build / dev targets
│
├── internal/
│   ├── yalex/
│   │   ├── yalex.go             Types: YalFile, Rule, LetDefinition, RulePattern
│   │   ├── scanner.go           Comment stripping, ASCII validation, section detection
│   │   └── compile.go           let-def expansion + DFA compilation orchestration
│   │
│   ├── regex/
│   │   ├── regexstring.go       Token types; sentinel runes (RuneEOF, RuneEndMarker)
│   │   ├── normalize.go         YALex syntax  →  flat token stream
│   │   │                          'c'  "str"  ['a'-'z']  [^set]  _ (wildcard)  eof
│   │   ├── shuntingyard.go      Infix → postfix (RPN) conversion
│   │   ├── syntaxtree.go        AST node construction from postfix
│   │   └── regex.go             Preprocess() pipeline entry point
│   │
│   ├── automata/
│   │   ├── dfa.go               DFAState / DFA types; SetToken(), Alphabet(), Sim()
│   │   ├── direct.go            Direct DFA construction from regex AST
│   │   ├── compile.go           Compile() (one regex → DFA) + Merge() (N DFAs → one)
│   │   └── minimization.go      Hopcroft-style DFA minimization
│   │
│   ├── codegen/
│   │   └── codegen.go           *** YOUR CODE GOES HERE — see section below ***
│   │
│   ├── graph/
│   │   └── graph.go             SerializeDFA — DFA → JSON for the frontend viewer
│   │
│   └── ds/
│       ├── stack.go             Generic stack
│       └── tree.go              Generic tree node (used for the regex AST)
│
├── frontend/                    React + Vite UI
│   ├── src/
│   │   ├── components/
│   │   │   └── DFAViewer/       Interactive DFA graph (Graphviz WASM, pan/zoom)
│   │   └── types/index.ts       TypeScript mirrors of the Go API types
│   └── vite.config.ts
│
└── workspace/                   Runtime user files (not committed)
    ├── specs/                   .yal spec files go here
    ├── lexers/                  Generated output: <name>.go + <name>.dfa.json
    └── input/                   Test input files for the "Run Lexer" feature
```

---

## How the Pipeline Works

```
.yal file
   │
   ▼
yalex.ParseYalContent()          Parse: header, let-defs, rules, trailer
   │
   ▼
YalFile.Compile()                For each rule pattern:
   │  ├─ expandDefinitions()       Substitute let-identifiers (transitive, word-boundary)
   │  ├─ regex.Preprocess()
   │  │     ├─ NormalizePattern    YALex syntax → token stream
   │  │     ├─ HandleSpecialOps    +  →  xx*     ?  →  (x|ε)
   │  │     ├─ ExplicitConcat      insert ~ concatenation tokens
   │  │     ├─ ShuntingYard        infix → postfix
   │  │     └─ AppendEndMarker     append # end-marker
   │  ├─ automata.Compile()        direct DFA construction + minimize  (one per pattern)
   │  └─ automata.Merge()          merge all pattern DFAs into one minimized DFA
   │                               first-rule-wins priority; TokenID = pattern index (1-based)
   ▼
CompiledLexer { DFA, Actions[] }
   │
   ├──▶  graph.SerializeDFA()      → JSON sent to browser → DFA Viewer
   │
   └──▶  codegen.GenerateLexer()   → workspace/lexers/<name>.go   ← implement this
```

---

## What You Need to Implement

**File:** `internal/codegen/codegen.go`

The stub is already there with a full doc comment. The function signature is:

```go
func GenerateLexer(name, sourcePath string, dfa *automata.DFA, actions []string) string
```

Return the complete text of a Go source file. The handler writes it to
`workspace/lexers/<name>.go` — you don't need to touch anything else.

### DFA data you get

```go
states := dfa.GetAllStates()   // []*automata.DFAState

state.ID                        // unique int
state.Accept                    // true on accepting states
state.TokenID                   // 1-based token index (0 on non-accept states)
state.Transitions               // map[rune]*DFAState

dfa.StartState                  // *DFAState — the initial state

actions[tokenID - 1]            // raw action string for that token
                                // e.g. "return PLUS", "return int(lxm)"
```

### Output file contract

The generated file must be in `package lexers` and export:

```go
func <Name>Lex(input string) []Token
```

`Token` should be defined in the generated file (or shared in `workspace/lexers/tokens.go`):

```go
type Token struct {
    Kind   int    // TokenID (1-based). Use -1 for error tokens.
    Lexeme string // matched text         — yytext equivalent
    Value  any    // semantic value       — yyval equivalent; nil until action sets it
    Line   int    // 1-based line number where this token starts
    Col    int    // 1-based column       where this token starts
}
```

### Lexer behavior

| Requirement | Details |
|-------------|---------|
| **Longest match** | Follow DFA transitions as far as possible; commit to the last accepting state seen (maximal munch). |
| **Line tracking** | Increment `Line` on every `'\n'` consumed; reset `Col` to 1. |
| **Column tracking** | Increment `Col` for every non-`'\n'` rune consumed. |
| **Token position** | Record `Line`/`Col` at the *start* of the lexeme, before consuming it. |
| **Error recovery** | If no accepting state was seen, emit `Token{Kind: -1, Lexeme: <current rune>}` and advance one rune. |
| **Actions** | `actions[tokenID-1]` is user-supplied Go code. Embed it in a `switch` arm or inline it so the caller can dispatch on token kind. |

### Lex utility checklist

| Standard name | Maps to     | Notes                                             |
|---------------|-------------|---------------------------------------------------|
| `yytext`      | `Token.Lexeme` | The raw matched string                         |
| `yyval`       | `Token.Value`  | Any semantic value; typed by the action        |
| `yyline`      | `Token.Line`   | 1-based; updated on `'\n'`                     |
| `yycol`       | `Token.Col`    | 1-based; reset to 1 after `'\n'`               |

### Skeleton of a generated file

```go
// Code generated by YALex — DO NOT EDIT.
// Source: specs/arithmetic.yal
package lexers

type Token struct {
    Kind   int
    Lexeme string
    Value  any
    Line   int
    Col    int
}

// arithmeticTransitions[state][char] = nextState
var arithmeticTransitions = map[int]map[rune]int{ /* ... populated by GenerateLexer */ }

// arithmeticTokenID[state] = tokenID for accept states (0 = non-accept)
var arithmeticTokenID = map[int]int{ /* ... */ }

const arithmeticStart = 0

func ArithmeticLex(input string) []Token {
    var tokens []Token
    runes := []rune(input)
    pos, line, col := 0, 1, 1

    for pos < len(runes) {
        state := arithmeticStart
        lastAccept, lastAcceptPos := -1, pos
        startLine, startCol := line, col
        cur := pos

        for cur < len(runes) {
            next, ok := arithmeticTransitions[state][runes[cur]]
            if !ok { break }
            state = next
            if arithmeticTokenID[state] > 0 {
                lastAccept, lastAcceptPos = arithmeticTokenID[state], cur+1
            }
            if runes[cur] == '\n' { line++; col = 1 } else { col++ }
            cur++
        }

        if lastAccept < 0 {
            tokens = append(tokens, Token{Kind: -1, Lexeme: string(runes[pos]),
                Line: startLine, Col: startCol})
            if runes[pos] == '\n' { line++; col = 1 } else { col++ }
            pos++
            continue
        }

        tok := Token{Kind: lastAccept, Lexeme: string(runes[pos:lastAcceptPos]),
            Line: startLine, Col: startCol}

        switch lastAccept {
        case 1: /* action for token 1 */
        case 2: /* action for token 2 */
        // ...
        }

        tokens = append(tokens, tok)
        pos = lastAcceptPos
    }
    return tokens
}
```

---

## Running Locally

```bash
# First time: install frontend dependencies
make frontend-install

# Development mode (Go API on :8080, Vite dev server on :5173 with hot-reload)
make dev

# Production build (bundles frontend, compiles Go binary)
make build
./yalex
```

Open `http://localhost:5173` in dev mode, or `http://localhost:8080` after a production build.

---

## .yal File Format

```
(* optional comment *)
{
  Go import / preamble code  — copied verbatim into the generated file header
}

let digit  = ['0'-'9']
let letter = ['a'-'z'] | ['A'-'Z']
let alnum  = letter | digit

rule gettoken =
    [' ' '\t']    { return lexbuf }
  | ['\n']        { return EOL }
  | digit+        { return int(lxm) }
  | '+'           { return PLUS }
  | eof           { raise('End of input') }

{
  Go trailer code  — appended verbatim to the generated file
}
```

**Pattern syntax:**

| Syntax           | Meaning                                        |
|------------------|------------------------------------------------|
| `'c'`            | Literal character                              |
| `'\n'` `'\t'` …  | Escape sequences: `n t r \\ ' "`              |
| `"hello"`        | String — characters matched in sequence        |
| `['a'-'z']`      | Character range                                |
| `['a' 'b' 'c']`  | Character list                                 |
| `[^'0'-'9']`     | Negated class (all printable ASCII minus set)  |
| `[A] # [B]`      | Set difference A − B                           |
| `_`              | Wildcard — any printable ASCII character       |
| `eof`            | End-of-buffer sentinel (rune U+E002)           |
| `x*` `x+` `x?`  | Kleene / one-or-more / optional                |
| `x \| y`         | Alternation                                    |
| `(x)`            | Grouping                                       |

**Priority:** the first rule whose pattern matches wins (standard lex behaviour).

---

## API Reference

All paths are relative to the `workspace/` directory.

| Method | Path                  | Body / Query         | Description                                      |
|--------|-----------------------|----------------------|--------------------------------------------------|
| GET    | `/api/health`         | —                    | Liveness check                                   |
| GET    | `/api/workspace/tree` | —                    | Recursive directory listing of workspace         |
| GET    | `/api/file`           | `?path=rel/path`     | Read a file                                      |
| POST   | `/api/file`           | `{path, content}`    | Write / overwrite a file                         |
| PUT    | `/api/file`           | `{path}`             | Create an empty file                             |
| DELETE | `/api/file`           | `?path=rel/path`     | Delete file or directory                         |
| POST   | `/api/file/rename`    | `{oldPath, newPath}` | Rename / move                                    |
| PUT    | `/api/directory`      | `{path}`             | Create directory                                 |
| POST   | `/api/dfa`            | `{path}`             | Parse .yal → build DFA → return graph JSON + write lexer .go |
| POST   | `/api/lexer`          | `{inputPath}`        | Run lexer on an input file, return token list    |

`POST /api/lexer` infers the spec from the input file extension:
`input/test.arithmetic` → looks for `specs/arithmetic.yal`.
