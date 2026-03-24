# Edge Cases para Testing del Lexer

## 1. Maximal Munch / Backtracking

Estos prueban que el scanner siempre tome el match más largo y haga backtrack correctamente.

| Input | Esperado (arithmetic) | Por qué |
|-------|----------------------|---------|
| `3.14.15` | `FLOAT(3.14)` `ERROR(.)` `INT(15)` | El segundo `.` no puede iniciar float, backtrack |
| `**=` | `POW(**)` `ERROR(=)` | `**` es más largo que `*`, `=` no tiene patrón |
| `+-*/**%` | `PLUS` `MINUS` `TIMES` `DIV` `POW` `MOD` | Sin espacios, cada operador debe resolverse |

| Input | Esperado (gol) | Por qué |
|-------|----------------|---------|
| `++=` | `INC(++)` `ASSN(=)` | Maximal munch agarra `++` primero, no `+` y `+=` |
| `<<=` | `LT(<)` `LEQ(<=)` | `<<` no es un patrón, así que solo toma `<`, luego `<=` |
| `===` | `EQ(==)` `ASSN(=)` | No `ASSN ASSN ASSN` |

## 2. First-Rule-Wins / Keyword vs Ident

| Input | Esperado | Por qué |
|-------|----------|---------|
| `po` | `IDENT(po)` | No es `pow`, maximal munch no hace match parcial de keyword |
| `pow` | `KWPOW(pow)` | Keyword exacto, first-rule-wins (kwpow antes de ident) |
| `power` | `IDENT(power)` | Maximal munch: `power` > `pow`, y `power` matchea `alpha alnum*` |
| `iffy` | `IDENT(iffy)` | `iffy` > `if`, maximal munch gana |
| `if` | `KWIF(if)` | Keyword exacto |
| `forloop` | `IDENT(forloop)` | `forloop` > `for` |

## 3. Error Recovery / Agrupación

El scanner agrupa caracteres inválidos consecutivos en un solo `ERROR`.

| Input | Esperado | Por qué |
|-------|----------|---------|
| `@@@hello` | `ERROR(@@@)` `IDENT(hello)` | 3 chars inválidos agrupados |
| `3@4` | `INT(3)` `ERROR(@)` `INT(4)` | Error entre tokens válidos |
| `hello@world` | `IDENT(hello)` `ERROR(@)` `IDENT(world)` | Error en medio de idents |
| `@#$%` | Un solo `ERROR(@#$%)` | Todos son inválidos, se agrupan |
| `` ~` `` | `ERROR` de ambos | Backtick y tilde no están en ningún patrón |

## 4. Input Vacío / Degenerado

| Input | Esperado | Por qué |
|-------|----------|---------|
| (archivo vacío) | Nada, solo `EOF` | Sin tokens |
| Solo whitespace/newlines | Nada, solo `EOF` | Whitespace se skipea |
| `@@@` (solo inválidos) | Un solo `ERROR(@@@)` | Todo agrupado |
| `+` (un solo char válido) | `PLUS(+)` | Token mínimo |
| `@` (un solo char inválido) | `ERROR(@)` | Error mínimo |

## 5. Whitespace y Line Tracking

| Input | Verificar |
|-------|-----------|
| `3\t\t+\t4` | Tabs entre tokens, ¿cols correctos? (tab = 1 col en el scanner) |
| `\r\n` entre tokens | `\r` no incrementa línea, solo `\n`. ¿Line tracking correcto en Windows-style? |
| Whitespace al inicio y final del archivo | No produce tokens, solo `EOF` |
| Líneas vacías entre tokens | Line numbers saltan correctamente |

## 6. Hex Literals (gol/calc)

| Input | Esperado (gol) | Por qué |
|-------|----------------|---------|
| `0xFF` | `HEX(0xFF)` | Caso normal |
| `0x` | `INT(0)` `IDENT(x)` | `hexdig+` requiere al menos un dígito hex. DFA no matchea, backtrack a `0` como INT |
| `0xGG` | `INT(0)` `IDENT(xGG)` | `G` no es hexdig, mismo backtrack |
| `0X1a2B` | `HEX(0X1a2B)` | Mayúscula `X`, mezcla de case en digits |

## 7. Bad Literals (gol)

`badlit = digit+ alpha alnum*` matchea cosas como `10abc`.

| Input | Esperado | Por qué |
|-------|----------|---------|
| `10abc` | `BADLIT(10abc)` | Dígitos seguidos de letras |
| `10abc20def` | `BADLIT(10abc20def)` | `alnum*` consume todo |
| `123_` | `BADLIT(123_)` | `_` es parte de `alpha` |
| `0abc` | `BADLIT(0abc)` | Empieza con dígito, no es hex (no `0x`) |

## 8. Float Edge Cases

| Input | Esperado | Por qué |
|-------|----------|---------|
| `3.` | `INT(3)` `ERROR(.)` | Float requiere dígitos después del punto |
| `.5` | `ERROR(.)` `INT(5)` | Float requiere dígitos antes del punto |
| `3.14` | `FLOAT(3.14)` | Normal |
| `3.14.15` | `FLOAT(3.14)` `ERROR(.)` `INT(15)` | Segundo punto no puede ser float |
| `0.0` | `FLOAT(0.0)` | Ceros |
| `00.00` | `FLOAT(00.00)` | Leading zeros (¿es válido? el pattern lo permite) |

## 9. Multi-char Operator Overlap (calc)

Calc tiene muchos operadores de 2 chars que comparten prefijo.

| Input | Esperado | Por qué |
|-------|----------|---------|
| `+= -=` | `PLEQ(+=)` `MIEQ(-=)` | Compound assignment |
| `== !=` | `EQ(==)` `NEQ(!=)` | Comparison |
| `= ==` | `ASSN(=)` `EQ(==)` | Simple vs double |
| `!! ` | `NOT(!)` `NOT(!)` | No existe `!!` como operador |
| `&&||` | `AND(&&)` `OR(\|\|)` | Sin espacio entre operadores lógicos |

## 10. Spec Parsing (no end-to-end, sino el parser del .yal)

- Nested comments: `(* outer (* inner *) still outer *)`
- Spec sin header ni trailer (solo let-defs y rules)
- Spec sin let-defs (rules directamente)
- Cadena de let-defs referenciándose: `let a = ...`, `let b = a ...`, `let c = b ...`
- Acción con braces en strings: `{ fmt.Printf("{%d}", x); return TOK }`
- Patrón sin acción (sin `{}`)

## 11. `#` Set Difference

`gol.yal` usa `_ # ['a'-'z' 'A'-'Z' '0'-'9' '_' ' ' '\t' '\n' '\r']` para punct.

| Input | Verificar |
|-------|-----------|
| Todos los printable ASCII uno por uno | Cada uno clasificado correctamente |
| `!@#$%^&*()-+=[]{}:;<>,./?~\|` | Todos deberían ser tokens de punct o sus operadores específicos |

## Test Files Sugeridos

### `workspace/input/edge.arithmetic`
```
3.14.15
**=
po
pow
power
@@@hello
3@4
@#$%
+-*/**%
((((1))))
99999999999999999999



+
@
```

### `workspace/input/edge.gol`
```
++=
<<=
===
0xGG
0x
hello@world
@@@func
10abc20def
forloop
~`
.5
00.00
&&||
!!
```

### `workspace/input/edge.calc`
```
+=-=*=/=
===
!==
&&||
0xFF
0x
0xGG
let x = 3.14.15
iffy
if
@@@hello@@@
```
