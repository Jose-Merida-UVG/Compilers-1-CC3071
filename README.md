# YALex / YAPar — Pendientes de implementación

---

## Codegen del parser

`GenerateCombined` en `internal/yalex/codegen/codegen.go` ya existe y está lista. El handler `buildParser` en `handlers.go` ya construye el lexer inline pero todavía escribe un stub en `parser/parser.go` en lugar de llamar a `GenerateCombined`.

### Convención de nombres

El `.yal` y el `.yalp` deben tener el mismo nombre base — equivalente al flag `-l` del CLI:

```
yapar parser.yalp -l lexer.yal -o theparser
```

```
specs/arithmetic.yal  +  specs/arithmetic.yalp  →  programs/arithmetic/
```

### Flujo actual de `buildParser` (handlers.go)

**Build Parser** ya no requiere "Build Lexer" previo — hace todo en una sola acción:

```
specs/<nombre>.yalp  +  specs/<nombre>.yal
  → POST /api/yapar
  → parsear .yalp  →  yalpFile
  → parsear .yal   →  yalFile
  → yalFile.Compile()     →  DFA + acciones
  → yalpFile.Compile()    →  autómata LR(0) + tablas SLR(1)
  → GenerateLexer()       →  programs/<nombre>/lexer/lexer.go   ✓ ya se escribe
  → docs/dfa.json, lr0.json, slr.json                           ✓ ya se escriben
  → parser/parser.go      ←  ⚠️ STUB todavía
```

### Lo que falta: llamar a GenerateCombined

En `handlers.go: buildParser`, reemplazar el stub con:

1. Construir `[]codegen.TokenDef` desde los tokens del `.yalp`
2. Llamar `codegen.GenerateCombined(specBase, yalFile, compiledLexer.DFA, compiledLexer.Actions, tokenDefs)`
3. Separar el output en dos archivos: lexer sin `main()` → `lexer/lexer.go`, parser con `main()` → `parser/parser.go`

`GenerateCombined` actualmente genera un solo archivo. Al dividir, debe haber exactamente un `main()` entre los dos — el run handler hace `go run lexer/lexer.go parser/parser.go <input>` y falla si hay dos.

---

## Manejo de errores en el output

### Errores léxicos (ya funcional)

El lexer generado agrupa caracteres no reconocidos en tokens `ERROR` y continúa escaneando. El output en la terminal muestra cada ERROR con su posición:

```
ERROR  "@@@@"               ln=9 col=1-4
ERROR  "0x"                 ln=10 col=1-2
```

Esto ya funciona — el lexer no para ante un carácter inválido, sigue hasta EOF y reporta todos los errores juntos.

### Errores sintácticos (a implementar en el parser generado)

Cuando el parser no encuentra una acción válida para el token actual (celda vacía en la tabla SLR), debe:

1. Emitir un mensaje con el token ofensivo y su posición (disponible en `Lexeme.Line` / `Lexeme.Col`)
2. Intentar recuperación de errores o parar

Output esperado en terminal:

```
syntax error: unexpected PLUS at line 3, col 7
```

La información de posición ya viene en el `Lexeme` — el parser solo necesita usarla al detectar una celda vacía en la tabla Action.

**Highlighting en el editor:** la UI actualmente solo muestra texto plano en la terminal. Para resaltar el error en el editor habría que parsear el mensaje de error desde el frontend y marcar el rango en Monaco. Eso está fuera del scope del codegen pero es viable si el formato del mensaje es consistente (ej. siempre `line X, col Y`).

---

## Reglas del contrato pendientes

### Regla 1 — Advertencia: token declarado pero nunca retornado

En `buildParser`, después de parsear ambos archivos, cruzar:
- tokens declarados en `%token` (del `.yalp`)
- identificadores en `return X` dentro de las acciones del `.yal`

Emitir advertencia en la respuesta JSON para que aparezca en la terminal. No bloquea el build.

### Regla 3 — Error: token ignorado referenciado en gramática

Los tokens bajo `IGNORE` en el `.yal` no deben aparecer en ninguna producción del `.yalp`. Si aparecen, retornar 422 antes de generar código.

---

## Tabla de símbolos

El `main()` del parser generado debe acumular todos los tokens consumidos vía `NextToken()` e imprimir una tabla al final. El parser ya ve cada token — no requiere cambios en el lexer.

```go
symbolTable := map[string][]Lexeme{}
// al consumir cada token:
symbolTable[tok.Value] = append(symbolTable[tok.Value], tok)

// al terminar:
fmt.Println("\n── Tabla de símbolos ──")
fmt.Printf("%-20s %-10s %s\n", "LEXEMA", "TOKEN", "LÍNEA:COL")
for lexeme, occurrences := range symbolTable {
    for _, occ := range occurrences {
        fmt.Printf("%-20s %-10d %d:%d\n", lexeme, occ.Token, occ.Line, occ.Col)
    }
}
```

Output aparece en la terminal de la UI al hacer **▶ Run**.
