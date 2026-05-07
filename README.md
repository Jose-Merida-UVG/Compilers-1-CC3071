# YALex / YAPar — Generador de Analizadores Léxicos y Sintácticos

Proyecto de CC3071 - Diseño de Lenguajes de Programación, Universidad del Valle de Guatemala.

---

## ¿Qué hace el proyecto?

Este proyecto implementa un generador de analizadores léxicos y sintácticos con interfaz web:

- **YALex**: toma un archivo `.yal` con patrones de expresiones regulares y acciones, construye un DFA minimizado y genera un lexer en Go.
- **YAPar**: toma un archivo `.yalp` con la gramática en formato BNF, parsea las producciones, y computa los conjuntos **FIRST** y **FOLLOW** para cada no terminal.

Todo se controla desde una UI web: abrís un archivo en el editor, hacés click en **◎ Build Lexer** o **◎ Build Parser**, y el output aparece en la terminal integrada.

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
%token TOKEN_B TOKEN_C
%token WS
IGNORE WS

%%

produccion1:
    produccion1 TOKEN_A produccion2
  | produccion2
;

produccion2:
    TOKEN_B
  | /* empty */
;
```

- Terminales en **MAYÚSCULAS**, no terminales en **minúsculas**
- Producciones ε: `| /* empty */`
- Secciones separadas por `%%`

Abrí el `.yalp` en el editor y hacé click en **◎ Build Parser** — los conjuntos FIRST y FOLLOW aparecen en la terminal de la UI.

---

## Gramáticas de ejemplo

En `workspace/specs/` hay dos gramáticas listas.

### `dragon.yalp` — Dragon Book §4.5

```
e → e PLUS t | t
t → t TIMES f | f
f → LPAREN e RPAREN | ID
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

---

## Estructura interna

```
internal/
  yalex/              ← pipeline del lexer
    regex/            ← normalización de patrones → postfix
    automata/         ← construcción directa de DFA, minimización Hopcroft
    graph/            ← DFA → JSON para visualización
    codegen/          ← DFA + acciones → código Go
  yapar/              ← pipeline del parser
    yapar.go          ← parseo de archivos .yalp
    grammar/
      grammar.go      ← struct Grammar, Build(), Summary()
      first.go        ← ComputeFirst()
      follow.go       ← ComputeFollow()
  ds/                 ← estructuras de datos compartidas

workspace/
  specs/    ← archivos .yal y .yalp
  lexers/   ← lexers generados (.go)
  parsers/  ← parsers generados (.go)
  input/    ← archivos de prueba
  output/   ← resultados de corridas (.out)
```
