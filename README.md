# Direct Regex to DFA Compiler

A Go implementation that converts regular expressions directly to Deterministic Finite Automata (DFA) using syntax tree algorithms—no intermediate NFA required. Includes DFA minimization and comprehensive test simulation.

## Video Demos

**Direct Method (DFA Construction)**
> [Add video link here]

**State Minimization**
> [Add video link here]

## What It Does

This compiler implements the **direct DFA construction algorithm** with state minimization:

1. **Parse & Normalize**: Convert regex string to postfix notation with special operator handling
2. **Build Syntax Tree**: Create augmented syntax tree with position attributes (nullable, firstPos, lastPos)
3. **Compute Tables**: Generate followPos table and direct transition table
4. **Construct DFA**: Build DFA directly from position sets as states
5. **Minimize States**: Apply partition refinement algorithm to minimize the DFA
6. **Simulate & Test**: Run test strings against both regular and minimized DFAs
7. **Visualize**: Generate PDF diagrams of DFA state machines

## Quick Start

### Build
```bash
go build -o bin/main cmd/main.go
```

### Run
```bash
./bin/main
```

Reads regexes from `data/regex.txt`, constructs DFAs, minimizes them, runs test simulations, and generates visualization graphs in `out/`.

## Current Functionality

The compiler processes each regex and outputs:

- **Postfix notation** with end marker
- **Regular DFA**: Complete state transition table
- **Minimized DFA**: Reduced state version using partition refinement
- **Test Simulations**: For each regex, 2+ test cases showing:
  - Strings that belong to the language ✓
  - Strings that don't belong to the language ✗
  - Verification that DFA and MinDFA produce identical results
- **Visualization**: PDF graphs for both DFA and MinDFA in `out/`

## Test Examples

### Regex 1: `(a|b)*cd?`
Matches: zero or more a/b, then mandatory 'c', optional 'd'

| Input | Expected | DFA | MinDFA | Status |
|-------|----------|-----|--------|--------|
| `bababac` | ✓ | ✓ | ✓ | Pass |
| `bababacd` | ✓ | ✓ | ✓ | Pass |
| `ababdcd` | ✗ | ✗ | ✗ | Pass |

### Regex 2: `(a|b)*abb(a|b)*`
Matches: any string containing 'abb' as a substring

| Input | Expected | DFA | MinDFA | Status |
|-------|----------|-----|--------|--------|
| `babababbababab` | ✓ | ✓ | ✓ | Pass |
| `abbbbbbbbbbbb` | ✓ | ✓ | ✓ | Pass |
| `bababababab` | ✗ | ✗ | ✗ | Pass |

## Architecture

| Directory | Purpose |
|-----------|---------|
| `cmd/main.go` | Entry point - orchestrates regex processing, DFA construction, minimization, and test simulation |
| `internal/regex/` | Regex parsing (shunting yard, syntax tree, balanced parentheses) |
| `internal/automata/` | DFA construction, minimization, and simulation |
| `internal/ds/` | Data structures (stack, tree node) |
| `internal/graph/` | Graphviz visualization |
| `data/regex.txt` | Input test regexes |
| `out/` | Generated DFA & MinDFA graphs (PDF) |

## Output

### Console Output
- Postfix notation representation
- DFA state count and transition count
- MinDFA state count and transition count (typically fewer states)
- Transition tables for both DFA and MinDFA
- Test simulation results comparing expected vs actual

### Generated Files
- `out/DFA_regex#.pdf` - Visualization of the complete DFA
- `out/MinDFA_regex#.pdf` - Visualization of the minimized DFA

## How It Works

### Direct DFA Construction
The algorithm builds a DFA directly from a regex without constructing an NFA first:
- Position attributes (nullable, firstPos, lastPos) are computed bottom-up on the syntax tree
- The followPos function determines which positions can follow each position
- States in the DFA represent sets of positions from the original regex

### State Minimization
Uses partition refinement algorithm:
- Initial partitions separate accepting and non-accepting states
- Iteratively refines partitions based on transition behavior
- Stops when no further refinement is possible
- Typically reduces state count significantly (see test examples above)

## Requirements

- Go 1.23.0+
- Graphviz (for PDF graph visualization)
