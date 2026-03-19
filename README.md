# YALex - Yet Another Lexical Analyzer Generator

A Go implementation of a lexical analyzer generator inspired by Lex and OCamllex. YALex reads lexer specification files (`.yal`) and generates scanners using direct DFA construction algorithms.

## What It Does

YALex implements a complete lexical analyzer generator pipeline:
1. Parse `.yal` specification files (similar to Lex format)
2. Extract regex patterns, let definitions, and token rules
3. Convert regular expressions directly to DFAs using syntax tree algorithms
4. Generate scanner code for token recognition
5. Visualize DFA state machines

## Quick Start

### Build
```bash
go build -o bin/yalex main.go
```

### Run
```bash
./bin/yalex <file.yal>
```

Or run without arguments to test with the default example:
```bash
./bin/yalex
```

Reads lexer specifications from `.yal` files and displays the parsed structure.

## YALex File Format

YALex files (`.yal`) follow this structure:

```
(* comments *)
{ header }                    (* optional: code copied to output *)
let ident = regexp ...        (* optional: named regex definitions *)
rule entrypoint =             (* required: lexer rules *)
  pattern { action }
| pattern { action }
...
{ trailer }                   (* optional: helper code *)
```

### Example

```
{
import myToken
}
let digit = ['0'-'9']
rule gettoken =
  [' ' '\t']+ { return lexbuf }
| digit+      { return int(lxm) }
| '+'         { return PLUS }
| eof         { raise('EOF') }
```

## Project Structure

| Directory | Purpose |
|-----------|---------|
| `main.go` | Entry point - YALex file parser |
| `parser/` | YALex file parsing & regex processing |
| `automata/` | Direct DFA construction & minimization |
| `graph/` | Graphviz DFA/AST visualization |
| `ds/` | Data structures (stack, tree) |
| `data/` | Example `.yal` specification files |
| `out/` | Generated DFA graphs (PDF) |

## Current Status

✅ YALex file parsing (header, let definitions, rules, trailer)  
✅ Regex to DFA direct construction algorithm  
✅ DFA minimization  
✅ DFA visualization (Graphviz PDF output)  
🚧 Scanner code generation (in progress)

## Requirements

- Go 1.23.0+
- Graphviz (for graph visualization)