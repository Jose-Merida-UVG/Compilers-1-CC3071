# YALex / YAPar — Pendientes de implementación

---

## Codegen del parser

`GenerateCombined` en `internal/yalex/codegen/codegen.go` ya existe y está lista. El handler `buildParser` en `handlers.go` escribe un stub en lugar de llamarla.

### Convención de nombres

El `.yal` y el `.yalp` deben tener el mismo nombre base:

```
specs/arithmetic.yal   +   specs/arithmetic.yalp   →   programs/arithmetic/
```

Esto es el equivalente frontend al flag `-l` del CLI:

```
yapar parser.yalp -l lexer.yal -o theparser
```

### Flujo del frontend

**Build Parser** (botón activo al abrir un `.yalp`) debe hacer todo en una sola acción — no requiere "Build Lexer" previo:

```
specs/<nombre>.yalp  +  specs/<nombre>.yal  (mismo nombre base)
  → POST /api/yapar
  → parsear ambos archivos
  → compilar .yal → DFA + acciones
  → compilar .yalp → autómata LR(0) + tablas SLR(1)
  → GenerateCombined() → dos archivos:
      programs/<nombre>/lexer/lexer.go   ← lexer sin main()
      programs/<nombre>/parser/parser.go ← parser + main()
  → programs/<nombre>/docs/lr0.json
  → programs/<nombre>/docs/slr.json
```

**Build Lexer** sigue existiendo como herramienta standalone para visualizar el DFA sin necesitar gramática.

### Pasos concretos en `handlers.go: buildParser`

1. Derivar `yalPath = "specs/" + specBase + ".yal"`
2. Leer y parsear el `.yal` con `yalex.ParseYalContent()`
3. Compilar con `yalFile.Compile()` para obtener DFA y acciones
4. Construir `[]codegen.TokenDef` desde los tokens del `.yalp`
5. Llamar `codegen.GenerateCombined()` con ambos
6. Separar el output en dos archivos: lexer sin `main()`, parser con `main()`
7. Reemplazar el stub actual

### Conflicto de `main()`

`GenerateCombined` actualmente genera un solo archivo con `main()`. Al dividir en dos archivos, hay que asegurarse de que solo uno tenga `main()`. El run handler recolecta todos los `.go` de `lexer/` + `parser/` y los pasa juntos a `go run` — funciona siempre que haya exactamente un `main()`.

---

## Reglas del contrato pendientes

### Regla 1 — Advertencia: token declarado pero nunca retornado

Después de parsear el `.yal` y el `.yalp`, cruzar:
- tokens declarados en `%token`
- identificadores en `return X` dentro de las acciones del `.yal`

Los tokens en `%token` que no aparecen en ningún `return` deben emitirse como advertencia en la terminal. No bloquea el build.

### Regla 3 — Error: token ignorado referenciado en gramática

Los tokens bajo `IGNORE` en el `.yal` no deben aparecer en ninguna producción del `.yalp`. Si aparecen, retornar error antes de generar código.

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
