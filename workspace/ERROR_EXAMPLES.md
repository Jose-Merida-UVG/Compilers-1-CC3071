# Parser Error Examples

This document shows how to introduce errors and what error messages you'll see.

## Arithmetic Grammar Errors

### 1. Missing Closing Parenthesis
**Input:** `(2 + 3 * 4`
**Error:** `syntax error: unexpected $ at line 1, col X`
**Explanation:** Parser expects RPAREN but gets EOF ($)

### 2. Invalid Token Sequence
**Input:** `2 + + 3`
**Error:** `syntax error: unexpected PLUS at line 1, col 5`
**Explanation:** After PLUS, parser expects a primary expression (number, identifier, lparen), not another operator

### 3. Missing Operator
**Input:** `2 3`
**Error:** `syntax error: unexpected INT at line 1, col 3`
**Explanation:** After reducing `2` to expr, parser expects an operator or EOF, not another number

### 4. Unmatched Parentheses
**Input:** `2 + 3)`
**Error:** `syntax error: unexpected RPAREN at line 1, col 6`
**Explanation:** Extra closing parenthesis with no matching opening

### 5. Function Call Errors
**Input:** `pow(2,)`
**Error:** `syntax error: unexpected RPAREN at line 1, col 7`
**Explanation:** Parser expects an expression after comma, not RPAREN

**Input:** `pow(2 3)`
**Error:** `syntax error: unexpected INT at line 1, col 7`
**Explanation:** Arguments must be separated by commas

## Calc Grammar Errors

### 1. Missing ELSE in Conditional
**Input:** `if x > 5`
**Error:** `syntax error: unexpected $ at line 1, col X`
**Explanation:** `if` expression requires `else` branch in this grammar

### 2. Missing Assignment in LET
**Input:** `let x`
**Error:** `syntax error: unexpected $ at line 1, col X`
**Explanation:** `let` binding requires `= expr`

### 3. Comparison Chaining (Not Allowed)
**Input:** `1 < 2 < 3`
**Error:** `syntax error: unexpected LT at line 1, col 7`
**Explanation:** Comparisons are not associative; use `&&` instead: `1 < 2 && 2 < 3`

### 4. Missing Operand
**Input:** `2 +`
**Error:** `syntax error: unexpected $ at line 1, col X`
**Explanation:** Binary operator requires right operand

### 5. Invalid Function Call
**Input:** `func()`
**Error:** Depends on grammar - if IDENT is in primary but function calls require args

## GoL (Mini-Go) Grammar Errors

### 1. Missing Semicolon
**Input:**
```go
package main
import fmt
func add() int
```
**Error:** `syntax error: unexpected KWFUNC at line 3`
**Explanation:** Import statement must end with SEMI

### 2. Missing Package Clause
**Input:**
```go
func add() int {
    return 1
}
```
**Error:** `syntax error: unexpected KWFUNC at line 1`
**Explanation:** Program must start with `package <name> ;`

### 3. Invalid Variable Declaration
**Input:**
```go
var x
```
**Error:** `syntax error: unexpected SEMI`
**Explanation:** Variable declaration requires type or initializer

### 4. Missing Block Braces
**Input:**
```go
func add() int
    return 1
```
**Error:** `syntax error: unexpected KWRETURN`
**Explanation:** Function body requires `{ ... }`

### 5. Invalid For Loop
**Input:**
```go
for {
    x = x + 1
```
**Error:** `syntax error: unexpected $ at EOF`
**Explanation:** Missing closing brace for block

## Lexical Errors

These occur BEFORE parsing:

### 1. Invalid Character
**Input:** `2 + @3`
**Error:** `lexical error: unrecognized input '@' at line 1, col 5`
**Explanation:** `@` is not in the lexer alphabet

### 2. Malformed Number
**Input:** `0xGHI` (if grammar doesn't support letters G,H,I in hex)
**Error:** `lexical error: unrecognized input 'G' at line 1, col 3`
**Explanation:** Invalid hex digit

### 3. Unterminated String (if strings are in grammar)
**Input:** `"hello`
**Error:** `lexical error: unrecognized input at EOF`
**Explanation:** String literal not closed

## How Errors Are Detected

1. **Lexical Errors**: Lexer encounters input that doesn't match any token pattern
   - Returns ERROR token (-1)
   - Parser immediately fails with "lexical error" message

2. **Syntax Errors**: Parser cannot find valid action in SLR table
   - `parserActionTable[state][symbol]` lookup fails
   - Reports "syntax error: unexpected SYMBOL at line X, col Y"
   - Shows what token was encountered but not expected

3. **State Errors**: Parser state machine gets confused
   - No action row for current state (rare, means table generation bug)
   - No goto entry after reduce (grammar ambiguity or bug)

## Testing Strategy

1. **Start Simple**: Test valid inputs first
2. **One Error at a Time**: Introduce single error to understand message
3. **Boundary Cases**: Empty input, single token, deeply nested expressions
4. **Recovery**: This parser has NO error recovery - it fails on first error
5. **Error Position**: Check that line/column numbers are accurate
