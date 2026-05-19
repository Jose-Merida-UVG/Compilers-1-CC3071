# YAPar — Parser Generator

A Go-based SLR(1) parser generator inspired by Yacc and ocamlyacc. You write a `.yalp` spec declaring your tokens and grammar; the tool computes FIRST/FOLLOW sets, builds the LR(0) automaton, constructs SLR(1) parse tables, and generates a combined Go program together with the paired `.yal` lexer.

The intended invocation (as a CLI tool) would be:
```
yapar parser.yalp -l lexer.yal -o theparser
```
The current implementation is delivered as a web UI — the build and run buttons in the editor are the equivalent of that call.

---

## File Structure

A `.yalp` file has two sections separated by `%%`:

```
/* comment */

%token TOKEN_A
%token TOKEN_B TOKEN_C
IGNORE WS

%%

nonterminal:
    nonterminal TOKEN_A other
  | other
;

other:
    TOKEN_B
  | /* empty */
;
```

### Token declarations

```
%token NAME
%token NAME_A NAME_B NAME_C
```

Declares terminals. Each name gets a unique integer ID (1-based, in declaration order). These constants are generated and injected into the combined `lexer.go` output — do **not** define them in the `.yal` header.

Multiple names on one `%token` line are allowed.

### IGNORE

```
IGNORE WS
IGNORE COMMENT NEWLINE
```

Tokens that exist in the lexer but are silently discarded before reaching the parser. A token listed under `IGNORE` must not appear in any production rule — the parser never sees it, so any such rule would be unreachable. This is a hard error at build time.

### Productions

```
nonterminal:
    body1
  | body2
;
```

- Non-terminals are lowercase identifiers.
- Terminals are UPPERCASE and must be declared with `%token`.
- Empty productions: `| /* empty */`
- Each rule ends with `;`.
- The first non-terminal declared is the start symbol.

---

## Pairing with a `.yal` file

Every `.yalp` must be paired with a `.yal` of the same base name in the same directory:

```
specs/arithmetic.yal    ←  lexer patterns + actions
specs/arithmetic.yalp   ←  tokens + grammar
```

The generator reads both, validates that the token sets match, and produces the combined output.

### Lexer requirements

1. **Do not define token constants in the `.yal` header** — YAPar generates these from `%token`.
2. **Every `%token` — including `IGNORE` tokens — must have a `return TOKEN_NAME` in the lexer.** Token validation checks all `%token` declarations, so a missing return is a build error.
3. **Ignored tokens are filtered by the parser, not the lexer.** The lexer returns them normally; the parser's `nextToken` loop sees the token ID in `parserIgnore` and silently discards it before the grammar ever sees it.
4. **Whitespace skipped at lexer level** (the common case) does not need a `%token` or `IGNORE` declaration at all — just use an empty action `{  }` with no return and don't declare the token.
5. `EOF` (0) and `ERROR` (-1) are always available; do not redefine them.

---

## Input Constraints

- **ASCII only.** Both files are rejected if any non-ASCII byte is found.
- **Comments** use C-style block syntax: `/* ... */`. Single-line `//` is not supported.
- Terminals **must** be UPPERCASE; non-terminals **must** be lowercase.

---

## What the tool produces

Building a `.yalp` generates files under `workspace/programs/<name>/`:

| File | Contents |
|------|----------|
| `lexer.go` | Combined lexer: DFA scan loop, `Lexeme` type, `NextToken()`, token constants |
| `parser.go` | ACTION/GOTO tables, `Parse()` function, `main()` |
| `docs/dfa.json` | Lexer DFA graph for the DFA visualizer tab |
| `docs/lr0.json` | LR(0) automaton for the LR(0) visualizer tab |
| `docs/slr.json` | SLR(1) table for the SLR viewer tab |

The generated program is a complete runnable binary:

```bash
go run workspace/programs/arithmetic/lexer.go \
        workspace/programs/arithmetic/parser.go \
        workspace/input/test.arithmetic
```

---

## Parser Output Format

### Parse Actions

Each step shows the action taken and the resulting sentential form (the symbol stack) directly below it:

```
── Parse Actions ──
Shift  KWPOW(pow)
  KWPOW(pow)
Shift  LPAREN(()
  KWPOW(pow) LPAREN(()
Shift  INT(2)
  KWPOW(pow) LPAREN(() INT(2)
Reduce atom → INT
  KWPOW(pow) LPAREN(() atom
Reduce power → atom
  KWPOW(pow) LPAREN(() power
```

Terminals show their matched lexeme in parentheses. Non-terminals show just their name.

### Derivation

After all actions, the rightmost derivation is printed in reverse (top-down) so you can read it as an expansion from the start symbol down to the input:

```
── Derivation (top-down) ──
1. expr
2. term
3. power
4. atom
5. KWPOW(pow) LPAREN(() atom COMMA(,) atom RPAREN())
...
```

### Symbol Table

```
── Symbol Table ──
LEXEME               TOKEN                LINE:COL
pow                  KWPOW                1:1
(                    LPAREN               1:4
2                    INT                  1:5
```

### Result

- `OK — input accepted by the grammar.` on success
- `syntax error: unexpected TOKEN at line X, col Y` on parse failure
- `lexical error: unrecognized input 'X' at line Y, col Z` on lex failure

---

## Visualizer Tabs

After clicking **◎ Build Parser** three tabs open:

- **SLR table** — ACTION and GOTO tables, FIRST/FOLLOW sets, production list, and any conflicts.
- **LR(0) graph** — full automaton with item sets and transitions. Supports pan and scroll-to-zoom.
- **DFA** — lexer DFA for the paired `.yal`.

---

## SLR(1) Constraints

The grammar must be SLR(1). Common conflict sources:

**Shift/reduce** — ambiguous associativity. Fix with explicit left recursion:
```
expr: expr PLUS term | term   ← left-associative, unambiguous
```

**Dangling else** — `if cond stmt` followed by `else`. Fix by requiring braces or always requiring an else branch.

**Reduce/reduce** — two productions match the same lookahead. Restructure so each lookahead maps to at most one reduce.

**Type-safe grammars** — separating `numexpr` and `boolexpr` into parallel hierarchies causes SLR(1) conflicts whenever a shared token like `IDENT` or `LPAREN` can start either. Type enforcement belongs in a semantic analysis pass, not in the grammar.

Conflicts are reported in the SLR tab and do not block the build, but the parser will behave unpredictably in conflicting states.

---

## How to Use

### 1. Write your specs

Create `workspace/specs/<name>.yal` and `workspace/specs/<name>.yalp` with the same base name. For demo/test grammars that don't need input files, put them in `workspace/demo/`.

### 2. Build

Open the `.yalp` in the editor and click **◎ Build Parser**. The tool reads both files, cross-validates tokens, generates the combined program, and opens the visualizer tabs.

### 3. Run against input

Create a test file in `workspace/input/` with an extension matching the spec name (e.g. `test.arithmetic` → uses `arithmetic.yal` + `arithmetic.yalp`). Open it and click **▶ Run**.

---

## Common Issues

### Tokens out of sync
```
token validation: token FOO declared in parser but never returned by lexer
```
Every `%token` in the `.yalp` must have a corresponding `return FOO` somewhere in the `.yal` actions, and vice versa.

### IGNORE token used in production
```
token WS is IGNORE but used in production "stmt"
```
Remove the ignored token from all production rules.

### No lexer spec found
```
no lexer spec found for "dragon" — expected dragon.yal alongside the .yalp
```
The `.yal` and `.yalp` must share the same base name and be in the same directory.

### Grammar has conflicts
```
grammar is not SLR(1): N conflict(s) — first: state X on 'TOKEN': rY vs sZ
```
Check the SLR visualizer tab for the conflicting states and symbols. Usually requires grammar restructuring.
