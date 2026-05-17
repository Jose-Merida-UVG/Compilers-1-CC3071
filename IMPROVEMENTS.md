# Code Improvements

## Comments

- **first.go / follow.go** — comments are in Spanish, everything else is English. Pick one and stick to it.
- Most comments explain *what* the code does (which the code already shows), not *why*. Trim those. A comment should only exist when the reason behind something is non-obvious.
- `follow.go` has multi-paragraph block comments above every section — way too verbose for what it's doing. One line per step is enough.
- `slr/slr.go` — `setAction` closure has no comment on *why* conflicts are skipped (first writer wins) rather than overwritten.

## Naming

- `yapar.go` — `yf` is too short for the main data structure being built. `yalpFile` everywhere would be cleaner.
- `yapar.go` — `seen` should be `declaredTokens`, `current` should be `currentProd`.
- `handlers.go` — `ext` in `runLexer` is used as a spec name, not a file extension. Call it `specName`.
- Grammar algorithm files use single-letter variables (`A`, `B`, `body`) which is fine and standard for grammar theory — don't change those.

## Duplication

- `sortedKeys` / `sortedBoolKeys` are defined in both `grammar/grammar.go` and `slr/summary.go`. Move them somewhere shared or just inline them.
- `runLexer` and `runFile` in `handlers.go` share the same resolve → stat → exec → write pattern. Extract the exec part into a helper.

## Structure

- `handlers.go` — file permission `0o755` appears 10+ times. One constant at the top.
- `handlers.go` — `writeFile`, `createFile`, `renameFile`, `createDirectory` all repeat the same resolve + error pattern. Could be tightened.
- `main.go` — HTTP server has no read/write timeouts. Fine for local dev but worth noting.
- `yapar.go` — `removeComments` silently ignores unclosed block comments (loop just exits). Should probably return an error.

## Styling

- Blank lines inside functions are inconsistent — some functions are wall-to-wall code, others have breathing room between logical steps. Pick a style.
- `handlers.go` section headers use `// ─── Name ───` which is good — keep that pattern everywhere instead of mixing with plain `//` comments.
