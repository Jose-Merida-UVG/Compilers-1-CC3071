# YALex / YAPar — Generador de Analizadores Léxicos y Sintácticos

Proyecto de CC3071 - Diseño de Lenguajes de Programación, Universidad del Valle de Guatemala.

El sistema tiene dos componentes principales:
- **YALex** — genera un analizador léxico (lexer) a partir de un archivo `.yal`
- **YAPar** — construye y analiza la gramática de un parser a partir de un archivo `.yalp`

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

Para hacer un build de producción (genera binario + frontend compilado):

```bash
make build
./yalex
```

El binario sirve el frontend directamente en `http://localhost:8080`.

Para limpiar artefactos generados:

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

Todos los archivos de trabajo del usuario viven en `workspace/`:

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

1. Creá un archivo `.yal` en `workspace/specs/` desde el explorador de archivos de la UI
2. Escribí tu especificación léxica en el editor
3. Hacé click en **◎ Build Lexer** — genera el DFA y el lexer en `lexers/`
4. Creá un archivo de prueba en `workspace/input/` con extensión igual al nombre de tu spec (ej. `test.arithmetic` para `arithmetic.yal`)
5. Abrí el archivo de prueba y hacé click en **▶ Run** — el output aparece en la terminal de la UI

---

## Usar el Parser (YAPar)

### Archivos necesarios

Cada parser necesita un par de archivos con el **mismo nombre base**:

| Archivo | Propósito |
|---|---|
| `specs/nombre.yal` | Especificación léxica (patrones y acciones) |
| `specs/nombre.yalp` | Especificación de la gramática (tokens y producciones) |

### Estructura de un archivo `.yalp`

```
/* comentario estilo C */

%token TOKEN_A
%token TOKEN_B TOKEN_C    ← podés declarar varios en una línea
%token WS
IGNORE WS                 ← el parser ignora este token aunque el lexer lo produzca

%%

produccion1:
    produccion1 TOKEN_A produccion2
  | produccion2
;

produccion2:
    TOKEN_B
  | TOKEN_C
  | /* empty */            ← producción épsilon
;
```

Reglas del formato:
- Los **tokens** (terminales) van en **MAYÚSCULAS** — deben coincidir con los nombres de reglas en el `.yal`
- Las **producciones** (no terminales) van en **minúsculas**
- Las producciones vacías (ε) se escriben como `| /* empty */`
- Las dos secciones se separan con `%%`
- Los comentarios son `/* ... */`

### Cómo usar Build Parser

1. Abrí el archivo `.yalp` en el editor de la UI
2. Hacé click en **◎ Build Parser**
3. En la terminal de la UI vas a ver:
   - Resumen de la gramática (terminales, no terminales, producciones)
   - Sets FIRST de cada no terminal
   - Sets FOLLOW de cada no terminal

---

## Para el equipo — dónde implementar FIRST y FOLLOW

Todo el trabajo del parser vive en:

```
internal/yapar/grammar/
  grammar.go   ← struct Grammar, Build(), Summary() y helpers
  first.go     ← ComputeFirst()   ← implementar acá
  follow.go    ← ComputeFollow()  ← implementar acá
```

### Flujo completo de comunicación

```
Click "Build Parser" en la UI
         ↓
frontend llama POST /api/yapar { yalpPath }
         ↓
handlers.go → ParseYalpContent() → grammar.Build() → g.Summary()
         ↓
responde { summary: []string }
         ↓
cada string del slice se imprime como una línea en la terminal de la UI
```

**El equipo solo toca `first.go` y `follow.go`.** El handler y el frontend no cambian.

### El struct Grammar

```go
type Grammar struct {
    Name         string
    StartSymbol  string
    Terminals    map[string]int   // nombre → token ID
    NonTerminals map[string]bool  // conjunto de no terminales
    Productions  []Production     // todas las reglas, una por alternativa
    IgnoreList   []string

    // Llenar estos en ComputeFirst y ComputeFollow:
    First  map[string]map[string]bool
    Follow map[string]map[string]bool
}
```

### Convenciones

| Concepto | Cómo representarlo |
|---|---|
| ε en FIRST | `g.First["X"][grammar.Epsilon] = true` |
| Símbolo nullable | `sym.Nullable = true` + la línea de arriba |
| Fin de input en FOLLOW | usar el string `"$"` como clave |
| Producción ε | `len(production.Body) == 0` |

### Helpers disponibles en Grammar

```go
g.ProductionsFor("expr")  // todas las producciones con ese head
g.IsTerminal("PLUS")      // true si es token declarado
g.IsNonTerminal("expr")   // true si es cabeza de producción
g.IsNullable("expr")      // true si ε ∈ FIRST(expr), válido post-ComputeFirst
```

### Cómo agregar output a la UI

Dentro de `Summary()` en `grammar.go`, los resultados de FIRST y FOLLOW ya están conectados. Si querés agregar secciones extra, simplemente agreguen líneas al slice antes de retornarlo:

```go
lines = append(lines, "── Mi sección ──")
lines = append(lines, fmt.Sprintf("  resultado: %v", valor))
```

---

## Gramáticas de ejemplo

En `workspace/specs/` hay dos gramáticas listas para probar:

### `dragon.yalp` — Dragon Book §4.5

Gramática clásica de expresiones aritméticas, sin epsilon. Útil para verificar FIRST/FOLLOW contra los resultados del libro.

```
FIRST(e) = { LPAREN, ID }
FIRST(t) = { LPAREN, ID }
FIRST(f) = { LPAREN, ID }
```

### `epsilon_test.yalp` — Gramática con múltiples épsilon

Basada en ejercicio de clase:
```
S → a B D h
B → c C
C → b C | ε
D → E F
E → g | ε
F → f | ε
```

Los sets FIRST y FOLLOW esperados están documentados dentro del archivo.

---

## Estructura interna del proyecto

```
internal/
  yalex/              ← pipeline del lexer (autónomo)
    regex/            ← normalización de patrones → postfix
    automata/         ← construcción directa de DFA, minimización Hopcroft
    graph/            ← DFA → JSON para visualización
    codegen/          ← DFA + acciones → código Go
    yalex.go / scanner.go / compile.go
  yapar/              ← pipeline del parser (autónomo)
    yapar.go          ← parseo de archivos .yalp
    grammar/
      grammar.go      ← struct Grammar, Build(), Summary()
      first.go        ← ComputeFirst()
      follow.go        ← ComputeFollow()
  ds/                 ← estructuras de datos compartidas
```
