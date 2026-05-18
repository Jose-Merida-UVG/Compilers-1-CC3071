# Calc Language Type Rules

## Types

### Number Type
Includes: `INT`, `FLOAT`, `HEX`

Examples:
- `42` (INT)
- `3.14` (FLOAT)
- `0xFF` (HEX = 255)

### Boolean Type
Includes: `true`, `false`, comparison results

Examples:
- `true`
- `false`
- `5 > 3` (evaluates to boolean)
- `x == y` (evaluates to boolean)

## Valid Operations

### Arithmetic Operations (number → number)
These take numbers and produce numbers:
- `+` addition: `5 + 3` → `8`
- `-` subtraction: `10 - 4` → `6`
- `*` multiplication: `3 * 4` → `12`
- `/` division: `10 / 2` → `5`
- `%` modulo: `10 % 3` → `1`
- `**` power: `2 ** 3` → `8`

**Type Rule:** Both operands must be numbers, result is a number.

### Comparison Operations (number → boolean)
These take numbers and produce booleans:
- `<` less than: `3 < 5` → `true`
- `>` greater than: `5 > 3` → `true`
- `<=` less or equal: `3 <= 3` → `true`
- `>=` greater or equal: `5 >= 3` → `true`

**Type Rule:** Both operands must be numbers, result is a boolean.

### Equality Operations (any → boolean)
These work on any type:
- `==` equal: `5 == 5` → `true`, `true == true` → `true`
- `!=` not equal: `5 != 3` → `true`, `false != true` → `true`

**Type Rule:** Both operands must be the same type, result is a boolean.

### Logical Operations (boolean → boolean)
These take booleans and produce booleans:
- `&&` and: `true && false` → `false`
- `||` or: `true || false` → `true`
- `!` not: `!true` → `false`

**Type Rule:** All operands must be booleans, result is a boolean.

## Type Errors (Invalid Operations)

The parser accepts these syntactically, but they are **semantic errors**:

### Arithmetic on Booleans
```calc
true + false;        // ERROR: can't add booleans
true * 5;            // ERROR: can't multiply boolean and number
!42;                 // ERROR: ! requires boolean, not number
```

### Boolean Context Expecting Number
```calc
let x = if 5 ...     // ERROR: condition must be boolean, not number
while 42 ...         // ERROR: loop condition must be boolean
```

### Comparison on Booleans
```calc
true < false;        // ERROR: <, >, <=, >= only work on numbers
```

### Type Mismatches
```calc
5 && 3;              // ERROR: && requires booleans
true + 1;            // ERROR: + requires numbers
```

## Control Flow Type Requirements

### If Statements
```calc
if expr stmtblock
if expr stmtblock else stmtblock
```
**Type Rule:** `expr` must evaluate to boolean.

Examples:
```calc
if x > 5 (          // VALID: x > 5 is boolean
    y = 10;
)

if x (              // ERROR: x must be boolean, not number
    y = 10;
)
```

### While Loops
```calc
while expr stmtblock
```
**Type Rule:** `expr` must evaluate to boolean.

Examples:
```calc
while x > 0 (       // VALID: x > 0 is boolean
    x -= 1;
)

while x (           // ERROR: x must be boolean
    x -= 1;
)
```

## Variable Assignment Type Rules

### Let Bindings
```calc
let x = expr;       // x takes the type of expr
```

### Reassignment
```calc
x = expr;           // expr must match x's type
```

### Compound Assignment
```calc
x += expr;          // x and expr must both be numbers
x -= expr;          // x and expr must both be numbers
x *= expr;          // x and expr must both be numbers
x /= expr;          // x and expr must both be numbers
```

**Type Rule:** Compound assignments only work on numbers.

Error example:
```calc
let flag = true;
flag += 1;          // ERROR: can't add to boolean
```

## Function Calls

```calc
functionName(arg1, arg2, ...)
```

Functions can take any types and return any type. The type rules depend on the specific function being called.

Example:
```calc
let result = pow(2, 8);     // pow returns number
```

## Mixed Type Example (Showing Errors)

```calc
// VALID program
let x = 10;                 // x: number
let y = 20;                 // y: number
let isPositive = x > 0;     // isPositive: boolean (x > 0 is comparison)
let canProceed = isPositive && y < 100;  // canProceed: boolean

if canProceed (             // VALID: canProceed is boolean
    x += 5;                 // VALID: x is number
)

// INVALID program
let a = 5;                  // a: number
let b = true;               // b: boolean
let c = a + b;              // ERROR: can't add number and boolean
let d = !a;                 // ERROR: ! requires boolean, a is number

if a (                      // ERROR: condition must be boolean, a is number
    a += 1;
)

while b + 1 (               // ERROR: condition must be boolean, b + 1 is invalid
    b = false;
)
```

## Summary Table

| Operation | Left Type | Right Type | Result Type | Example |
|-----------|-----------|------------|-------------|---------|
| + - * / % ** | number | number | number | `5 + 3` → `8` |
| < > <= >= | number | number | boolean | `5 > 3` → `true` |
| == != | any | same as left | boolean | `5 == 5` → `true` |
| && \|\| | boolean | boolean | boolean | `true && false` → `false` |
| ! | boolean | - | boolean | `!true` → `false` |
| unary - | number | - | number | `-5` → `-5` |

## Implementation Note

The **parser does not enforce these type rules** - it only checks syntax. Type checking would be done in a separate semantic analysis phase after parsing. However, programs violating these rules are considered semantically invalid even if they parse successfully.
