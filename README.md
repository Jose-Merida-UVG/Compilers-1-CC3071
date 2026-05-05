# YALex / YAPar — Generador de Analizadores Léxicos y Sintácticos

Proyecto de CC3071 - Diseño de Lenguajes de Programación, Universidad del Valle de Guatemala.
---

## Para el equipo — FIRST y FOLLOW

Todo el trabajo vive en:

```
internal/yapar/grammar/
  grammar.go   ← struct Grammar, Build(), Summary() y helpers
  first.go     ← ComputeFirst()   ← implementar acá
  follow.go    ← ComputeFollow()  ← implementar acá
```

Implementás `ComputeFirst()` y `ComputeFollow()`, llenás `g.First` y `g.Follow`, y los resultados aparecen automáticamente en la terminal de la UI al hacer click en **◎ Build Parser**. No hay que tocar nada más.

### El struct Grammar

```go
type Grammar struct {
    Name         string
    StartSymbol  string
    Terminals    map[string]int   // nombre → token ID
    NonTerminals map[string]bool  // conjunto de no terminales
    Productions  []Production     // todas las reglas, una por alternativa
    IgnoreList   []string

    First  map[string]map[string]bool  // llenar en ComputeFirst
    Follow map[string]map[string]bool  // llenar en ComputeFollow
}
```

### Convenciones

| Concepto | Cómo representarlo |
|---|---|
| ε en FIRST | `g.First["X"][grammar.Epsilon] = true` |
| Símbolo nullable | `sym.Nullable = true` + la línea de arriba |
| Fin de input en FOLLOW | usar el string `"$"` como clave |
| Producción ε | `len(production.Body) == 0` |

### Helpers disponibles

```go
g.ProductionsFor("expr")  // todas las producciones con ese head
g.IsTerminal("PLUS")      // true si es token declarado
g.IsNonTerminal("expr")   // true si es cabeza de producción
g.IsNullable("expr")      // true si ε ∈ FIRST(expr), válido post-ComputeFirst
```

### Flujo de comunicación con la UI

```
Click "Build Parser"
         ↓
POST /api/yapar { yalpPath }
         ↓
ParseYalpContent() → grammar.Build() → g.Summary()
         ↓
{ summary: []string }
         ↓
cada string = una línea en la terminal de la UI
```

Si querés agregar output extra al terminal, agregá líneas en `Summary()`:

```go
lines = append(lines, "── Mi sección ──")
lines = append(lines, fmt.Sprintf("  resultado: %v", valor))
```

---

## Gramáticas de ejemplo

En `workspace/specs/` hay dos gramáticas listas — abrí cualquiera en la UI y hacé click en **◎ Build Parser** para ver el output.

### `dragon.yalp` — Dragon Book §4.5

```
e → e PLUS t | t
t → t TIMES f | f
f → LPAREN e RPAREN | ID
```

Output esperado:
```
FIRST(e) = { ID, LPAREN }
FIRST(t) = { ID, LPAREN }
FIRST(f) = { ID, LPAREN }

FOLLOW(e) = { $, PLUS, RPAREN }
FOLLOW(t) = { $, PLUS, TIMES, RPAREN }
FOLLOW(f) = { $, PLUS, TIMES, RPAREN }
```

### `epsilon_test.yalp` — Gramática con múltiples ε

```
S → a B D h
B → c C
C → b C | ε
D → E F
E → g | ε
F → f | ε
```

Output esperado:
```
FIRST(s)     = { A }
FIRST(bprod) = { C }
FIRST(cprod) = { B, ε }
FIRST(dprod) = { G, F, ε }
FIRST(eprod) = { G, ε }
FIRST(fprod) = { F, ε }

FOLLOW(s)     = { $ }
FOLLOW(bprod) = { G, F, H }
FOLLOW(cprod) = { G, F, H }
FOLLOW(dprod) = { H }
FOLLOW(eprod) = { F, H }
FOLLOW(fprod) = { H }
```

---

## Requisitos

- [Go 1.23+](https://go.dev/dl/)
- [Node.js 18+](https://nodejs.org/) y npm

---

## Cómo correr el proyecto

### Opción A — Makefile (recomendado)

```bash
# Primera vez: instalar dependencias del frontend
make frontend-install

# Modo desarrollo: levanta backend y frontend en paralelo
make dev
```

`make dev` inicia dos servidores al mismo tiempo:
- Backend Go en `http://localhost:8080`
- Frontend Vite en `http://localhost:5173`

Abrí el navegador en **http://localhost:5173**.

Para build de producción:

```bash
make build
./yalex
```

Para limpiar artefactos:

```bash
make clean
```

---

### Opción B — Manual

**Terminal 1 — Backend:**
```bash
go run .
```

**Terminal 2 — Frontend:**
```bash
cd frontend
npm install      # solo la primera vez
npm run dev
```

Abrí el navegador en **http://localhost:5173**.

---

## Estructura del workspace

```
workspace/
  specs/    ← tus archivos .yal y .yalp van acá
  lexers/   ← lexers generados automáticamente (.go)
  parsers/  ← parsers generados automáticamente (.go)
  input/    ← archivos de prueba para correr el lexer/parser
  output/   ← resultados de las corridas (.out)
```

---

## Usar el Lexer (YALex)

1. Creá un archivo `.yal` en `workspace/specs/`
2. Escribí tu especificación léxica en el editor
3. Hacé click en **◎ Build Lexer** — genera el DFA y el lexer en `lexers/`
4. Creá un archivo de prueba en `workspace/input/` con extensión igual al nombre de tu spec (ej. `test.arithmetic` para `arithmetic.yal`)
5. Abrí el archivo de prueba y hacé click en **▶ Run**

---

## Usar el Parser (YAPar)

Cada parser necesita un par de archivos con el **mismo nombre base**:

| Archivo | Propósito |
|---|---|
| `specs/nombre.yal` | Especificación léxica (patrones y acciones) |
| `specs/nombre.yalp` | Especificación de la gramática (tokens y producciones) |

### Formato de un archivo `.yalp`

```
/* comentario */

%token TOKEN_A
%token TOKEN_B TOKEN_C    ← varios en una línea
%token WS
IGNORE WS                 ← el parser ignora este token

%%

produccion1:
    produccion1 TOKEN_A produccion2
  | produccion2
;

produccion2:
    TOKEN_B
  | /* empty */            ← producción épsilon
;
```

- Terminales en **MAYÚSCULAS**, no terminales en **minúsculas**
- Producciones ε: `| /* empty */`
- Secciones separadas por `%%`

Abrí el `.yalp` en el editor y hacé click en **◎ Build Parser** — el output aparece en la terminal de la UI.

---

## Estructura interna

```
internal/
  yalex/              ← pipeline del lexer (autónomo)
    regex/            ← normalización de patrones → postfix
    automata/         ← construcción directa de DFA, minimización Hopcroft
    graph/            ← DFA → JSON para visualización
    codegen/          ← DFA + acciones → código Go
  yapar/              ← pipeline del parser (autónomo)
    yapar.go          ← parseo de archivos .yalp
    grammar/
      grammar.go      ← struct Grammar, Build(), Summary()
      first.go        ← ComputeFirst()
      follow.go       ← ComputeFollow()
  ds/                 ← estructuras de datos compartidas
```
