# YALex / YAPar — Pendientes de implementación

El único trabajo pendiente está en **`internal/yalex/codegen/codegen.go`** — la función `GenerateCombined` necesita ser llamada desde el handler y su output dividido en dos archivos.

El handler `buildParser` (`handlers.go`) ya parsea y compila tanto el `.yal` como el `.yalp` juntos. Solo le falta usar ese output para generar código real en lugar del stub actual (`// TODO: parser codegen` en `parser/parser.go`).

---

## Qué hacer

En `handlers.go: buildParser` (línea ~364), reemplazar el stub con:

1. Construir `[]codegen.TokenDef` desde los tokens del `.yalp`
2. Llamar `codegen.GenerateCombined(specBase, yalFile, compiledLexer.DFA, compiledLexer.Actions, tokenDefs)`
3. Separar el output en dos archivos:
   - `programs/<nombre>/lexer/lexer.go` — lexer **sin** `main()`
   - `programs/<nombre>/parser/parser.go` — parser **con** `main()`

El run handler hace `go run lexer/lexer.go parser/parser.go <input>` — necesita exactamente un `main()` entre los dos archivos.

---

## Convención de nombres

`.yal` y `.yalp` deben tener el mismo nombre base (equivalente al flag `-l`):

```
yapar parser.yalp -l lexer.yal -o theparser
```

```
specs/arithmetic.yal  +  specs/arithmetic.yalp  →  programs/arithmetic/
```

---

## Manejo de errores

### Léxicos — ya funcional

El lexer agrupa caracteres no reconocidos en tokens `ERROR` y sigue escaneando hasta EOF. Nunca para ante input inválido.

### Sintácticos — implementar en el `main()` generado

Cuando la tabla SLR tiene celda vacía para el token actual, emitir:

```
syntax error: unexpected PLUS at line 3, col 7
```

La posición ya viene en `Lexeme.Line` / `Lexeme.Col` — solo hay que usarla.

---

## Tabla de símbolos — implementar en el `main()` generado

El parser llama `NextToken()` para cada token, así que acumular y al final imprimir es trivial:

```go
symbolTable := map[string][]Lexeme{}
// al consumir cada token:
symbolTable[tok.Value] = append(symbolTable[tok.Value], tok)

// al terminar:
fmt.Println("\n── Tabla de símbolos ──")
fmt.Printf("%-20s %-10s %s\n", "LEXEMA", "TOKEN", "LÍNEA:COL")
for lexeme, occs := range symbolTable {
    for _, occ := range occs {
        fmt.Printf("%-20s %-10d %d:%d\n", lexeme, occ.Token, occ.Line, occ.Col)
    }
}
```

---

## Reglas del contrato pendientes

**Regla 1** — después de parsear ambos archivos, cruzar los tokens en `%token` contra los `return X` en las acciones del `.yal`. Los que nunca se retornan → advertencia en terminal. No bloquea el build.

**Regla 3** — si un token bajo `IGNORE` en el `.yal` aparece en alguna producción del `.yalp` → error 422 antes de generar código.
