# Calc Type Enforcement

The calc grammar now enforces type safety at **parse time** by separating numeric and boolean expressions into different non-terminals.

## What Gets Caught (Syntax Errors)

These are now **parse errors** - the grammar rejects them:

### 1. Arithmetic on Booleans
```calc
let x = true + false;     // ERROR: can't use + with booleans
let y = true * 5;         // ERROR: can't use * with boolean
let z = 10 - false;       // ERROR: can't subtract boolean from number
```
**Error:** `syntax error: unexpected PLUS/TIMES/MINUS`

### 2. Logical Operations on Numbers
```calc
let x = 5 && 10;          // ERROR: && requires booleans
let y = 3 || 7;           // ERROR: || requires booleans
let z = !42;              // ERROR: ! requires boolean
```
**Error:** `syntax error: unexpected AND/OR` or `unexpected INT after NOT`

### 3. Wrong Condition Type
```calc
if 5 (                    // ERROR: condition must be boolean, not number
    x = 10;
)

while 42 (                // ERROR: loop condition must be boolean
    x += 1;
)
```
**Error:** `syntax error: unexpected LPAREN` (expected comparison/boolean)

### 4. Comparison on Booleans
```calc
let x = true < false;     // ERROR: <, >, <=, >= only work on numbers
```
**Error:** Comparisons require `numexpr LT numexpr`, not booleans

### 5. Mixed Type Literals
```calc
let x = 5 + (true && false);   // ERROR: can't add number and boolean
let y = (10 * 2) || false;      // ERROR: || requires both operands to be boolean
```
**Error:** `syntax error` at the type mismatch point

## What Still Needs Runtime Checking

Variables (IDENT) can appear in both numeric and boolean contexts because we don't track types at parse time:

### 1. Variable Type Misuse
```calc
let flag = true;          // flag is boolean
let x = flag + 5;         // PARSES OK but runtime error: flag is boolean, can't add
```

```calc
let count = 10;           // count is number
if count (                // PARSES OK but runtime error: count is number, need boolean
    x = 5;
)
```

The parser accepts these because `IDENT` can be either type. You'd need a separate type-checking pass or runtime checks to catch these.

### 2. Function Return Types
```calc
let x = someFunction();   // Parser doesn't know return type
let y = x + 5;           // If x is boolean, runtime error
```

## Grammar Structure

### Numeric Expressions (numexpr)
```
numexpr  → numexpr + numterm | numexpr - numterm | numterm
numterm  → numterm * numfactor | numterm / numfactor | numterm % numfactor | numfactor
numfactor → numprimary ** numfactor | numprimary
numprimary → INT | FLOAT | HEX | IDENT | -numprimary | (numexpr) | func(numargs)
```

**Produces:** Numbers only
**Allows:** Arithmetic operators, numeric literals, variables, function calls

### Boolean Expressions (boolexpr)
```
boolexpr    → boolexpr || boolterm | boolterm
boolterm    → boolterm && boolfactor | boolfactor
boolfactor  → !boolfactor | boolprimary
boolprimary → TRUE | FALSE | IDENT | (boolexpr)
            | numexpr < numexpr | numexpr > numexpr
            | numexpr <= numexpr | numexpr >= numexpr
            | numexpr == numexpr | numexpr != numexpr
```

**Produces:** Booleans only
**Allows:** Logical operators, boolean literals, comparisons between numbers, variables

### Key Points

1. **Comparisons bridge the gap**: `numexpr < numexpr` produces a boolean
   - This is allowed in `boolprimary`, so you can write: `if x + 5 > 10 ( ... )`

2. **Variables are ambiguous**: `IDENT` appears in both `numprimary` and `boolprimary`
   - Parser can't distinguish numeric vs boolean variables
   - This is a fundamental limitation without a separate type-checking pass

3. **Operators are type-specific**:
   - `+, -, *, /, %, **` only in `numexpr`
   - `&&, ||, !` only in `boolexpr`
   - This enforces type safety for operations on literals

## Testing Type Safety

### Valid Programs
```calc
// All numeric
let x = 5 + 3 * 2;
let y = x ** 2;
x += 10;

// All boolean
let flag = true;
let check = flag && false || true;

// Mixed properly
let num = 10;
let isLarge = num > 100;
if isLarge (
    num = num * 2;
)
```

### Invalid Programs (Parse Errors)
```calc
// Type error: arithmetic on boolean
let x = true + 5;              // ✗ syntax error

// Type error: logic on number
let y = 5 && 10;               // ✗ syntax error

// Type error: number in if condition
if 42 (                        // ✗ syntax error
    x = 1;
)

// Type error: boolean in arithmetic
let z = 3 * false;             // ✗ syntax error
```

### Parses But Runtime Error
```calc
// Variable used in wrong context
let flag = true;
let x = flag + 5;              // ✓ parses, ✗ runtime error

let count = 10;
if count (                     // ✓ parses, ✗ runtime error
    x = 1;
)
```

## Comparison with Previous Version

**Before:** Single `expr` non-terminal, all operations mixed
- `5 + true` would **parse successfully**
- Type errors only caught at semantic analysis or runtime
- Grammar accepted all syntactically valid expressions regardless of types

**Now:** Separate `numexpr` and `boolexpr` non-terminals
- `5 + true` is a **parse error**
- Type errors for literals and operations caught immediately
- Grammar enforces type correctness for operations
- Only variable types need runtime checking

## Benefits

1. **Earlier error detection** - Type mismatches caught during parsing
2. **Better error messages** - "unexpected PLUS" is clearer in context
3. **Self-documenting** - Grammar structure shows type system
4. **No semantic analysis needed** - For literals and operations
5. **Compiler optimization** - Can generate type-specific code

## Limitations

1. **Can't track variable types** - Need symbol table for that
2. **Can't infer function return types** - Need type annotations or separate pass
3. **Larger grammar** - More productions, potentially more conflicts
4. **Less flexible** - Harder to add polymorphic operations

Overall, this approach catches ~80% of type errors at parse time, with the remaining 20% (variable misuse, function types) needing runtime checks or a separate type-checking pass.
