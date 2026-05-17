# YAPar — Parser Generator

A Go-based SLR(1) parser generator. You write a `.yalp` spec file declaring
your tokens and grammar; the tool computes FIRST/FOLLOW sets, builds the LR(0)
automaton, and generates SLR(1) parse tables. Always paired with a `.yal` file
of the same name.

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

Declares terminals the parser knows about. Each name gets a unique integer ID
(1-based, in declaration order). These constants are injected into the generated
file — do **not** define them yourself in the `.yal` header.

Multiple names on one `%token` line are allowed.

### IGNORE

```
IGNORE WS
```

Tells the parser that this token exists in the lexer but will never reach the
parser (the lexer skips it). A token listed under `IGNORE` must not appear in
any production rule — the parser would never see it, making any such rule
unreachable. This is a hard error at build time.

### Productions

```
nonterminal:
    body1
  | body2
;
```

- Non-terminals are lowercase identifiers.
- Terminals are uppercase (must be declared with `%token`).
- Empty productions: `| /* empty */`
- Each rule ends with `;`
- The first non-terminal declared is the start symbol.

---

## Pairing with a `.yal` file

Every `.yalp` must be paired with a `.yal` of the same base name:

```
specs/arithmetic.yal   ←  lexer patterns + actions
specs/arithmetic.yalp  ←  tokens + grammar
```

This mirrors the CLI flag `-l`:

```
yapar parser.yalp -l lexer.yal -o theparser
```

**YAPar owns the token constants.** The `%token` declarations are the source of
truth. The lexer actions return these names (`return PLUS`), and YAPar injects
the constants so they compile. Any constant block in the `.yal` header that
redefines the same names will cause a Go compile error.

---

## Input Constraints

- **ASCII only.** Both files are rejected before processing if any non-ASCII
  byte is found.
- **Comments** use C-style block syntax: `/* ... */`. Single-line `//` is not
  supported.
- Terminals **must** be uppercase, non-terminals **must** be lowercase.

---

## What the tool produces

Building a `.yalp` generates files under `programs/<name>/`:

| File | Contents |
|------|----------|
| `lexer/lexer.go` | Generated lexer (from the paired `.yal`) |
| `parser/parser.go` | Generated parser with `main()` |
| `docs/dfa.json` | DFA graph for the visualizer tab |
| `docs/lr0.json` | LR(0) automaton for the visualizer tab |
| `docs/slr.json` | SLR(1) table for the visualizer tab |

---

## Visualizer tabs

After clicking **◎ Build Parser** three tabs open automatically:

- **SLR table** — Action and Goto tables, FIRST/FOLLOW sets, production list,
  and any shift/reduce or reduce/reduce conflicts.
- **LR(0) graph** — the full automaton with item sets and transitions.
- **DFA** — the lexer DFA for the paired `.yal`.

---

## SLR(1) constraints

The grammar must be SLR(1). Common sources of conflicts:

**Shift/reduce** — ambiguous operator associativity. Fix with left recursion for
left-associative operators:
```
expr: expr PLUS term | term   ← left-recursive, unambiguous
```

**Dangling else** — `if cond stmt` followed by `else` is a classic conflict.
Fix by requiring braces, or by always requiring the `else` branch.

**Reduce/reduce** — two productions apply for the same lookahead. Usually
signals an ambiguous grammar. Restructure so each lookahead leads to at most
one reduce action.

Conflicts are reported in the SLR visualizer tab and do not block the build,
but the generated parser will behave unpredictably in conflicting states.

---

## Example: arithmetic.yalp

```
%token FLOAT INT
%token POW KWPOW
%token PLUS MINUS TIMES DIV MOD
%token LPAREN RPAREN COMMA
%token IDENT

%%

expr:
    expr PLUS term
  | expr MINUS term
  | term
;

term:
    term TIMES power
  | term DIV power
  | term MOD power
  | power
;

power:
    atom POW power
  | atom
;

atom:
    INT
  | FLOAT
  | IDENT
  | LPAREN expr RPAREN
  | KWPOW LPAREN expr COMMA expr RPAREN
  | IDENT LPAREN args RPAREN
;

args:
    arglist
  | /* empty */
;

arglist:
    expr
  | arglist COMMA expr
;
```

---

## How to Use

### 1. Write your specs

Create `specs/<name>.yal` and `specs/<name>.yalp` with the same base name.

### 2. Build

Open the `.yalp` in the editor and click **◎ Build Parser**. The tool reads
both files, compiles the lexer and parser together, and opens the visualizer
tabs.

### 3. Run against input

Create a test file in `workspace/input/` with an extension matching the spec
name (e.g. `test.arithmetic` → uses `arithmetic.yal` + `arithmetic.yalp`).
Open it and click **▶ Run**. Output — including the symbol table — appears in
the terminal and is saved to `workspace/output/<filename>.out`.
