# PENDING

## YALex — Add `NextToken()` to codegen

`internal/yalex/codegen/codegen.go` currently only generates `gettoken()` returning a bare `int`.

Need to add to the generated lexer file:

```go
func (l *Lexer) NextToken() Lexeme {
    tok := l.gettoken()
    return Lexeme{Token: tok, Value: l.Lxm, Line: l.Ln, Col: l.Col}
}
```

And the `Lexeme` struct:

```go
type Lexeme struct {
    Token int
    Value string
    Line  int
    Col   int
}
```

This is the public interface the parser will use. See CONTRACT.md for full context.
