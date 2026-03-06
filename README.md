# Compiler Design Project: Direct Regex to DFA

Hey everyone! Welcome to the repository. 

We are currently shifting our focus to the newest objective for this project: **Direct Conversion of Regular Expressions to Deterministic Finite Automata (DFA)**. 

This README serves as a quick guide to help you navigate the codebase and understand where to focus your efforts for this new implementation.

## 🎯 Current Objective: Direct DFA Construction

Instead of the traditional route (Regex -> NFA via Thompson -> DFA via Subset Construction), we are now building the DFA directly from the regex's syntax tree using `firstpos`, `lastpos`, `nullable`, and `followpos` functions.

### 🏗️ Architecture & How It Works

The core idea is to compute a `followpos` table from the augmented syntax tree (appending `#` to the regex). 

Conceptually, you can imagine this table (or the state tracking structure) having a `toDFA()` method. 
- This method would start by taking the `firstpos` of the root node of our syntax tree as the initial state of the new DFA.
- It would then iterate through the active positions, looking up transitions in the `followpos` table for each input symbol.
- For every new set of positions encountered, it creates a new DFA state.
- Ultimately, `toDFA()` returns a fully constructed, direct DFA object without ever building an intermediate NFA!

## 📂 Where to Look (Important Files)

If you're jumping in to help with the direct conversion, these are the files you need to care about:

- **`internal/automata/direct.go`**: This is where the magic should happen. The main logic for the direct construction algorithm lives here.
- **`internal/automata/dfa.go`**: Contains the struct definitions and methods for our DFA object.
- **`internal/regex/syntaxtree.go`**: Crucial for the first step. This parses the regex into a syntax tree and is where the calculations for `nullable`, `firstpos`, `lastpos`, and `followpos` need to be handled.

## 🗃️ Everything Else (The Useless/Extra Stuff for Now)

Don't worry too much about the rest of the repository for this specific task. The following are just supporting structures or old implementations:
- `internal/automata/thompson.go` & `internal/automata/subsetconversion.go`: The old NFA/DFA conversion pipeline.
- `internal/ds/`: Basic data structures (Stacks, Trees) used across the project.
- `internal/graph/`: Code for rendering the graph visuals (Graphviz stuff).
- `data/` & `out/`: Just input strings and output PDFs.

Let's get this direct conversion working!