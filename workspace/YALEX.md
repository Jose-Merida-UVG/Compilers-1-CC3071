# YALex — Lexical Analyzer Generator

A Go-based lexer generator inspired by OCamllex. You write a `.yal` spec file
describing your token patterns; the tool builds a minimized DFA and generates
a self-contained Go program that scans input using maximal munch.

---

## File Structure

A `.yal` file has four sections, in order:

```
{ header }          ← optional: verbatim Go (constants, imports)

let ident = regex   ← zero or more named regex definitions

rule name =         ← one or more rule blocks
  pattern { action }
| pattern { action }

{ trailer }         ← optional: verbatim Go (helpers, extra code)
```

### Header / Trailer
Go source copied verbatim into the generated file before and after the scanner.
Use the header to define your token constants and any imports your actions need.
Use the trailer for helper functions referenced in actions.

```
{
    const (
        INT   = 1
        PLUS  = 2
        WS    = 3
    )
}
```

### Let Definitions
Named regex bindings. They are substituted by name wherever they appear in
patterns. Names must be plain ASCII identifiers — no underscores.

```
let digit = ['0'-'9']
let alpha = ['a'-'z'] | ['A'-'Z']
let ident = alpha (alpha | digit)*
```

Definitions can reference previously defined names:

```
let alnum = alpha | digit
let ident = alpha alnum*
```

### Rule Blocks
Each arm is `pattern { action }`, optionally prefixed with `|`.

```
rule tokens =
  ws+   {  }
| int   { return INT }
| ident { return IDENT }
```

Actions are verbatim Go code. If an action contains `return`, it exits `Scan()`
and delivers the token to the caller. If it does not `return`, the scanner
loops and picks up the next token — the standard way to skip whitespace.

---

## Regex Syntax

| Syntax             | Meaning                              |
|--------------------|--------------------------------------|
| `'a'`              | literal character                    |
| `['a'-'z']`        | character range                      |
| `['a' 'b' 'c']`    | character set                        |
| `['a'-'z' '0'-'9']`| combined range and set               |
| `p q`              | concatenation (implicit)             |
| `p \| q`           | alternation                          |
| `p*`               | Kleene star (zero or more)           |
| `p+`               | one or more                          |
| `p?`               | optional (zero or one)               |
| `(p)`              | grouping                             |
| `name`             | substitutes the named let definition |

Multi-character tokens are written as concatenated char literals:

```
let leq = '<' '='     (* matches <= *)
let inc = '+' '+'     (* matches ++ *)
```

---

## Input Constraints

- **ASCII only.** The validator rejects any character outside printable ASCII
  (0x20–0x7E) plus the whitespace control codes `\t`, `\n`, `\r`.
  This prevents unexpected behavior from Unicode look-alikes in patterns.

- **Comments** use OCamllex style: `(* ... *)`. They are nestable and may
  span multiple lines. The sequence `(*` inside a string literal inside an
  action block is not supported — it will be treated as a comment opener.

---

## Generated Lexer

Building a `.yal` file produces `lexers/<name>.go` — a self-contained
`package main` program. Its public API:

### Lexer struct

```go
type Lexer struct {
    Lxm string // matched lexeme (yytext equivalent)
    Ln  int    // 1-based line where the token starts
    Col int    // 1-based column where the token starts
    // internal fields omitted
}
```

### Constructor

```go
func New<Name>Lexer(input string) *Lexer
```

### Scan

```go
func (l *Lexer) Scan() int
```

Advances to the next token and returns its integer ID. Sets `Lxm`, `Ln`, and
`Col` before running the action. Typical usage:

```go
l := NewArithmeticLexer(input)
for {
    tok := l.Scan()
    if tok == EOF { break }
    fmt.Printf("%d  %q  ln=%d col=%d\n", tok, l.Lxm, l.Ln, l.Col)
}
```

### Sentinel return values

| Constant | Value | Meaning                              |
|----------|-------|--------------------------------------|
| `EOF`    | 0     | end of input                         |
| `ERROR`  | -1    | one or more unrecognised characters  |

### Error recovery

Consecutive unrecognised characters are grouped into a single `ERROR` token.
`l.Lxm` contains the full unrecognised substring. Scanning always continues
after an error — the lexer never panics or stops on bad input.

### Token ID constants

The generator also emits positional constants (`TOKEN_1`, `TOKEN_2`, …) in
the order patterns appear in the spec. These are for reference; you will
typically use your own named constants from the header instead.

---

## Scanning Rules

1. **Maximal munch.** The scanner always matches the longest possible string.
   `<=` is one token (LEQ), not `<` followed by `=`.

2. **First-rule-wins.** When two patterns match the same string at the same
   length, the one listed first in the rule wins. Use this to give keywords
   priority over identifiers for exact matches — but because of maximal munch,
   `iftrue` will still match as IDENT (longer than `if`).

3. **Order matters for operators.** Always list longer operators before their
   single-character prefixes:
   ```
   | leq  { return LEQ }   (* <= before < *)
   | '<'  { return LT }
   ```

---

## How to Use

### 1. Write your spec

Create a file in `workspace/specs/` with the `.yal` extension.

### 2. Build the DFA

Open the `.yal` file in the editor and click **◎ Build**. This:
- Parses the spec and compiles the minimized DFA
- Opens a DFA visualizer tab
- Generates `lexers/<name>.go`

### 3. Run against input

Create a test file in `workspace/input/` with an extension matching your spec
name (e.g. `test.arithmetic` uses `specs/arithmetic.yal`). Open it and click
**▶ Run**. Output appears in the terminal and is saved to
`workspace/output/<filename>.out`.

### 4. Use the generated file standalone

The generated `lexers/<name>.go` is a complete Go program:

```
go run workspace/lexers/arithmetic.go workspace/input/test.arithmetic
```

---

## Example: arithmetic.yal

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
  ws+   {  }
| float { return FLOAT }
| int   { return INT }
| '+'   { return PLUS }
| '-'   { return MINUS }
| '*'   { return TIMES }
| '/'   { return DIV }
| '%'   { return MOD }
| '('   { return LPAREN }
| ')'   { return RPAREN }
```
