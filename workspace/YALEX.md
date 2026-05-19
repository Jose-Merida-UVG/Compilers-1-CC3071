# YALex — Lexical Analyzer Generator

A Go-based lexer generator inspired by OCamllex. You write a `.yal` spec file describing your token patterns; the tool builds a minimized DFA and generates a Go lexer embedded in the combined parser program.

---

## File Structure

A `.yal` file has four sections, in order:

```
(* comment *)

let ident = regex   ← zero or more named regex definitions

rule name =         ← one rule block
  pattern { action }
| pattern { action }
```

The header/trailer sections (verbatim Go blocks wrapped in `{ }`) are supported but not needed when using YALex alongside YAPar — token constants are injected by YAPar automatically.

### Let Definitions

Named regex bindings substituted wherever their name appears in patterns.

```
let digit = ['0'-'9']
let alpha = ['a'-'z'] | ['A'-'Z']
let alnum = alpha | digit
let ident = alpha alnum*
```

Definitions can reference previously defined names but not themselves (no recursion).

### Rule Block

One rule block per `.yal` file. Each arm is `pattern { action }`, optionally prefixed with `|`.

```
rule gettoken =
  ws+   {  }
| int   { return INT }
| ident { return IDENT }
```

Actions are verbatim Go. A `return` exits the scan and delivers the token to the caller. No `return` means the scanner loops and picks up the next token — the standard way to skip whitespace and comments.

---

## Regex Syntax

| Syntax              | Meaning                              |
|---------------------|--------------------------------------|
| `'a'`               | literal character                    |
| `"abc"`             | literal string                       |
| `['a'-'z']`         | character range                      |
| `['a' 'b' 'c']`     | character set                        |
| `['a'-'z' '0'-'9']` | combined range and set               |
| `p q`               | concatenation (implicit)             |
| `p \| q`            | alternation                          |
| `p*`                | zero or more                         |
| `p+`                | one or more                          |
| `p?`                | optional                             |
| `(p)`               | grouping                             |
| `name`              | substitutes the named `let`          |

---

## Input Constraints

- **ASCII only.** Characters outside printable ASCII (0x20–0x7E) plus `\t \n \r` are rejected.
- **Comments** use OCamllex style: `(* ... *)`. They are nestable and may span multiple lines.

---

## Using with YAPar (combined mode)

When paired with a `.yalp` file of the same base name, YALex is compiled together with YAPar into a single program under `workspace/programs/<name>/`. In this mode:

- **Do not define token constants in the `.yal` header.** YAPar generates them from `%token` declarations.
- Every `%token` in the `.yalp` — including `IGNORE` tokens — must have a `return TOKEN_NAME` in the lexer. Token validation requires this or the build fails.
- Tokens declared `IGNORE` are returned by the lexer normally; the parser discards them internally via its ignore set.
- Whitespace and other purely lexer-level skips (no `%token` declaration) use an empty action `{  }` with no return — these never reach the parser at all.

```
(* arithmetic.yal — paired with arithmetic.yalp *)

let digit = ['0'-'9']
let int   = digit+
let float = digit+ '.' digit+
let ws    = [' ' '\t' '\n' '\r']

rule gettoken =
  ws+   {  }
| float { return FLOAT }
| int   { return INT }
| '+'   { return PLUS }
```

---

## Generated Lexer API

The lexer produced lives in `workspace/programs/<name>/lexer.go` as part of `package main`.

### Types

```go
type Lexeme struct {
    Token int    // token ID (EOF=0, ERROR=-1, or a declared constant)
    Value string // matched text
    Line  int    // 1-based line
    Col   int    // 1-based column
}

type Lexer struct { /* internal */ }
```

### Constructor

```go
func New<Name>Lexer(input string) *Lexer
```

### NextToken

```go
func (l *Lexer) NextToken() Lexeme
```

Returns the next token as a `Lexeme`. The underlying scan method named after the rule entrypoint (e.g. `gettoken()`) is also available directly.

### Sentinel values

| Constant | Value | Meaning                             |
|----------|-------|-------------------------------------|
| `EOF`    | 0     | end of input                        |
| `ERROR`  | -1    | one or more unrecognised characters |

### Error recovery

Consecutive unrecognised characters are grouped into a single `ERROR` token. `Lxm` contains the full unrecognised substring. The lexer never stops on bad input.

---

## Scanning Rules

1. **Maximal munch.** Always matches the longest possible string. `<=` is one LEQ token, not `<` then `=`.
2. **First-rule-wins.** When two patterns match at the same length, the earlier one wins. List keywords before the general identifier pattern.
3. **Order matters for operators.** Always list longer operators before their prefixes:
   ```
   | leq  { return LEQ }   (* <= before < *)
   | '<'  { return LT }
   ```

---

## How to Use

### 1. Write your spec

Create `workspace/specs/<name>.yal`. If using standalone (no parser), you can define your own token constants in a header block.

### 2. Build the DFA

Open the `.yal` file in the editor and click **◎ Build**. This compiles the DFA and opens the DFA visualizer tab.

### 3. Run against input (standalone)

Create a test file in `workspace/input/` with an extension matching your spec name (e.g. `test.arithmetic` → `arithmetic.yal`). Open it and click **▶ Run**.

For a combined parser+lexer, use the `.yalp` workflow instead — see YAPAR.md.
