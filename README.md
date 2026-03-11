# Direct Regex to DFA Compiler

A Go implementation that converts regular expressions directly to Deterministic Finite Automata (DFA) using syntax tree algorithms—no intermediate NFA required.

## What It Does

This compiler implements the **direct DFA construction algorithm** from formal language theory:
1. Parse a regex string and convert it to postfix notation
2. Build an augmented syntax tree with position attributes (nullable, firstPos, lastPos)
3. Compute the followPos table from the tree structure
4. Construct a DFA directly using position sets as states
5. Simulate and visualize the DFA against test strings

## Quick Start

### Build
```bash
go build -o bin/main cmd/main.go
```

### Run
```bash
./bin/main
```

Reads regexes from `data/regex.txt`, constructs their DFAs, runs test simulations, and generates visualization graphs in `out/`.

## Test Examples

The project includes 3 test regexes:

**Regex 1: `a+b|c*`** (one or more a's + b, OR zero or more c's)
- ✓ `aaaab` → matches
- ✓ `ccc` → matches  
- ✗ `aaabb` → no match
- ✗ `acb` → no match

**Regex 2: `(a|b)*cd?`** (zero or more a/b, then c, optional d)
- ✓ `bababacd` → matches
- ✓ `c` → matches
- ✗ `ababab` → no match
- ✗ `abcdc` → no match

**Regex 3: `a?(b|c)*d*e`** (optional a, zero or more b/c, zero or more d, then e)
- ✓ `abcbcddde` → matches
- ✓ `bcce` → matches
- ✗ `abcbcd` → no match
- ✗ `aabce` → no match

## Where to Find Everything

| Directory | Purpose |
|-----------|---------|
| `cmd/main.go` | Entry point |
| `internal/regex/` | Regex parsing (shunting yard, syntax tree) |
| `internal/automata/` | DFA construction & simulation |
| `internal/ds/` | Data structures (stack, tree) |
| `internal/graph/` | Graphviz visualization |
| `data/regex.txt` | Input test regexes |
| `out/` | Generated DFA & AST graphs (PDF) |

## Output

- **Console**: Transition tables and test results
- **Files** (`out/`): PDF diagrams of syntax trees and DFA state machines

## Requirements

- Go 1.23.0+
- Graphviz (for graph visualization)