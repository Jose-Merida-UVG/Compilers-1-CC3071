# CONTRACT

This document defines the integration contract between YALex and YAPar.
Violating any of these rules will produce either a compile error in the
generated Go file or undefined runtime behavior. Both outcomes are intentional.

---

## FILES

| File | Role |
|------|------|
| `<name>.yal`  | Lexer specification. Defines patterns and actions. |
| `<name>.yalp` | Parser specification. Defines tokens and grammar. |

Both files must be **strictly ASCII**. Non-ASCII bytes will be rejected before
any processing begins.

---

## TOKEN OWNERSHIP

**YAPar owns the token constants.** The `%token` declarations in the `.yalp`
file are the single source of truth for token names and their integer IDs.

When the parser is built, YAPar injects these constants as the header of the
generated lexer file. Any header already present in the `.yal` file that
redefines token constants will conflict and cause a compile error.

---

## LEXER ACTIONS

Actions in the `.yal` file are verbatim Go code. They can do anything:

- **Return a token** — `return PLUS` — emits that token to the parser.
- **Skip** — no `return` statement — the lexer loops and scans the next token.
  Use this for whitespace, comments, and anything the parser should not see.
- **Side effects** — update counters, log, etc. — then either return or skip.

There is no requirement that every action returns. Skipping is valid and expected.

---

## THE CONTRACT

### Rule 1 — Every token in `%token` must be returnable by the lexer.

Every token declared in `.yalp` must appear as `return TOKENNAME` in at least
one `.yal` action. A token that is declared but never returned means the parser
can never shift it — the grammar is effectively broken for that token.

This is a **warning**, not a hard error. The build will proceed but the
terminal will flag it.

### Rule 2 — Every `return X` in a `.yal` action must be a declared token.

If a `.yal` action returns an identifier that is not declared in `%token`, the
generated Go file will not compile because the constant is undefined.

This is a **hard compile error**. The error message from the Go compiler will
be the diagnostic. No special handling is needed — it is intentionally loud.

### Rule 3 — Ignored tokens must not appear in the grammar.

Tokens listed under `IGNORE` in the `.yal` file must not appear in any
production rule in the `.yalp` file. The parser will never see them, so any
production that references them is unreachable.

This is a **hard error** caught at build time by YAPar.

### Rule 4 — The `.yal` and `.yalp` files must be built together.

The `-l` flag in the YAPar invocation specifies the lexer spec:

```
yapar parser.yalp -l lexer.yal -o theparser
```

Building the parser without a lexer spec is not supported. The generated parser
file assumes the lexer is present and exposes `NextToken() Lexeme`.

---

## LEXEME

The lexer exposes one function to the parser:

```go
func (l *Lexer) NextToken() Lexeme
```

Where `Lexeme` is:

```go
type Lexeme struct {
    Token int    // integer ID matching the %token declaration
    Value string // the raw matched string from the input
    Line  int    // 1-based line number where the token starts
    Col   int    // 1-based column number where the token starts
}
```

The parser calls `NextToken()` exclusively. It does not access any internal
lexer fields directly. The lexer's internal scan function (`gettoken`) is not
part of the public interface.

**From the user's perspective**, actions only ever write `return PLUS` (or
whichever token name). That `return` exits `gettoken()` with the integer ID.
`NextToken()` then wraps the result automatically:

```go
func (l *Lexer) NextToken() Lexeme {
    tok := l.gettoken()  // "return PLUS" in the action lands here
    return Lexeme{Token: tok, Value: l.Lxm, Line: l.Ln, Col: l.Col}
}
```

The user never constructs a `Lexeme` directly. The struct assembly is
transparent and handled entirely by the generated code.

---

## WHAT WILL CRASH AND WHY

| Situation | Result |
|-----------|--------|
| `return PEEPEEBIG` in `.yal`, `PEEPEEBIG` not in `%token` | Go compile error — undefined constant |
| Token in `%token` never returned by any action | Parser hangs waiting for a token it never receives |
| Non-ASCII character in either file | Hard error before processing |
| `.yalp` references a token in `IGNORE` | Hard error at build time |
| Parser built without a lexer spec | Not supported |
