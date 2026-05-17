package main

import (
	"fmt"
	"os"
)

// Sentinel return values for Yylex.
const (
	EOF   = 0  // end of input
	ERROR = -1 // unrecognised character
)

// Token ID constants — one per pattern, in the order they appear in the spec.
const (
	TOKEN_1 = 1
	TOKEN_2 = 2 // fmt.Print("WHILE  "); return WHILE
	TOKEN_3 = 3 // fmt.Print("FLOAT  "); return FLOAT
	TOKEN_4 = 4 // fmt.Print("HEX    "); return HEX
	TOKEN_5 = 5 // fmt.Print("INT    "); return INT
	TOKEN_6 = 6 // fmt.Print("POW    "); return POW
	TOKEN_7 = 7 // fmt.Print("PLEQ   "); return PLEQ
	TOKEN_8 = 8 // fmt.Print("MIEQ   "); return MIEQ
	TOKEN_9 = 9 // fmt.Print("TIEQ   "); return TIEQ
	TOKEN_10 = 10 // fmt.Print("DIEQ   "); return DIEQ
	TOKEN_11 = 11 // fmt.Print("EQ     "); return EQ
	TOKEN_12 = 12 // fmt.Print("NEQ    "); return NEQ
	TOKEN_13 = 13 // fmt.Print("LEQ    "); return LEQ
	TOKEN_14 = 14 // fmt.Print("GEQ    "); return GEQ
	TOKEN_15 = 15 // fmt.Print("AND    "); return AND
	TOKEN_16 = 16 // fmt.Print("OR     "); return OR
	TOKEN_17 = 17 // fmt.Print("PLUS   "); return PLUS
	TOKEN_18 = 18 // fmt.Print("MINUS  "); return MINUS
	TOKEN_19 = 19 // fmt.Print("TIMES  "); return TIMES
	TOKEN_20 = 20 // fmt.Print("DIV    "); return DIV
	TOKEN_21 = 21 // fmt.Print("MOD    "); return MOD
	TOKEN_22 = 22 // fmt.Print("ASSN   "); return ASSN
	TOKEN_23 = 23 // fmt.Print("LT     "); return LT
	TOKEN_24 = 24 // fmt.Print("GT     "); return GT
	TOKEN_25 = 25 // fmt.Print("NOT    "); return NOT
	TOKEN_26 = 26 // fmt.Print("LPAREN "); return LPAREN
	TOKEN_27 = 27 // fmt.Print("RPAREN "); return RPAREN
	TOKEN_28 = 28 // fmt.Print("COMMA  "); return COMMA
	TOKEN_29 = 29 // fmt.Print("KWLET  "); return KWLET
	TOKEN_30 = 30 // fmt.Print("KWIF   "); return KWIF
	TOKEN_31 = 31 // fmt.Print("KWELSE "); return KWELSE
	TOKEN_32 = 32 // fmt.Print("IDENT  "); return IDENT
)

// Lexeme is the unit the parser consumes: a token ID plus the matched text and position.
type Lexeme struct {
	Token int
	Value string
	Line  int
	Col   int
}

// --- header ---
    const (
        FLOAT  = 1
        INT    = 2
        HEX    = 3
        PLUS   = 4
        MINUS  = 5
        TIMES  = 6
        DIV    = 7
        MOD    = 8
        POW    = 9
        PLEQ   = 10
        MIEQ   = 11
        TIEQ   = 12
        DIEQ   = 13
        ASSN   = 14
        EQ     = 15
        NEQ    = 16
        LT     = 17
        LEQ    = 18
        GT     = 19
        GEQ    = 20
        AND    = 21
        OR     = 22
        NOT    = 23
        LPAREN = 24
        RPAREN = 25
        COMMA  = 26
        KWLET  = 27
        KWIF   = 28
        KWELSE = 29
        IDENT  = 30
        WHILE = 31
    )

// Lexer holds the scanning state between calls to Scan.
type Lexer struct {
	input []rune
	pos   int
	line  int
	col   int
	Lxm   string // matched lexeme
	Ln    int    // 1-based line where the current token starts
	Col   int    // 1-based column where the current token starts
}

// NewCalcLexer creates a Lexer ready to scan input.
func NewCalcLexer(input string) *Lexer {
	return &Lexer{input: []rune(input), pos: 0, line: 1, col: 1}
}

// calcTrans is the DFA transition table: stateID → rune → nextStateID.
var calcTrans = map[int]map[rune]int{
	227: {
		'\t': 232,
		'\n': 232,
		'\r': 232,
		' ': 232,
		'!': 257,
		'%': 253,
		'&': 228,
		'(': 258,
		')': 259,
		'*': 251,
		'+': 249,
		',': 260,
		'-': 250,
		'/': 252,
		'0': 236,
		'1': 237,
		'2': 237,
		'3': 237,
		'4': 237,
		'5': 237,
		'6': 237,
		'7': 237,
		'8': 237,
		'9': 237,
		'<': 255,
		'=': 254,
		'>': 256,
		'A': 274,
		'B': 274,
		'C': 274,
		'D': 274,
		'E': 274,
		'F': 274,
		'G': 274,
		'H': 274,
		'I': 274,
		'J': 274,
		'K': 274,
		'L': 274,
		'M': 274,
		'N': 274,
		'O': 274,
		'P': 274,
		'Q': 274,
		'R': 274,
		'S': 274,
		'T': 274,
		'U': 274,
		'V': 274,
		'W': 274,
		'X': 274,
		'Y': 274,
		'Z': 274,
		'_': 274,
		'a': 274,
		'b': 274,
		'c': 274,
		'd': 274,
		'e': 272,
		'f': 274,
		'g': 274,
		'h': 274,
		'i': 266,
		'j': 274,
		'k': 274,
		'l': 268,
		'm': 274,
		'n': 274,
		'o': 274,
		'p': 274,
		'q': 274,
		'r': 274,
		's': 274,
		't': 274,
		'u': 274,
		'v': 274,
		'w': 273,
		'x': 274,
		'y': 274,
		'z': 274,
		'|': 231,
	},
	274: {
		'0': 274,
		'1': 274,
		'2': 274,
		'3': 274,
		'4': 274,
		'5': 274,
		'6': 274,
		'7': 274,
		'8': 274,
		'9': 274,
		'A': 274,
		'B': 274,
		'C': 274,
		'D': 274,
		'E': 274,
		'F': 274,
		'G': 274,
		'H': 274,
		'I': 274,
		'J': 274,
		'K': 274,
		'L': 274,
		'M': 274,
		'N': 274,
		'O': 274,
		'P': 274,
		'Q': 274,
		'R': 274,
		'S': 274,
		'T': 274,
		'U': 274,
		'V': 274,
		'W': 274,
		'X': 274,
		'Y': 274,
		'Z': 274,
		'_': 274,
		'a': 274,
		'b': 274,
		'c': 274,
		'd': 274,
		'e': 274,
		'f': 274,
		'g': 274,
		'h': 274,
		'i': 274,
		'j': 274,
		'k': 274,
		'l': 274,
		'm': 274,
		'n': 274,
		'o': 274,
		'p': 274,
		'q': 274,
		'r': 274,
		's': 274,
		't': 274,
		'u': 274,
		'v': 274,
		'w': 274,
		'x': 274,
		'y': 274,
		'z': 274,
	},
	236: {
		'.': 229,
		'0': 237,
		'1': 237,
		'2': 237,
		'3': 237,
		'4': 237,
		'5': 237,
		'6': 237,
		'7': 237,
		'8': 237,
		'9': 237,
		'X': 230,
		'x': 230,
	},
	237: {
		'.': 229,
		'0': 237,
		'1': 237,
		'2': 237,
		'3': 237,
		'4': 237,
		'5': 237,
		'6': 237,
		'7': 237,
		'8': 237,
		'9': 237,
	},
	229: {
		'0': 234,
		'1': 234,
		'2': 234,
		'3': 234,
		'4': 234,
		'5': 234,
		'6': 234,
		'7': 234,
		'8': 234,
		'9': 234,
	},
	234: {
		'0': 234,
		'1': 234,
		'2': 234,
		'3': 234,
		'4': 234,
		'5': 234,
		'6': 234,
		'7': 234,
		'8': 234,
		'9': 234,
	},
	230: {
		'0': 235,
		'1': 235,
		'2': 235,
		'3': 235,
		'4': 235,
		'5': 235,
		'6': 235,
		'7': 235,
		'8': 235,
		'9': 235,
		'A': 235,
		'B': 235,
		'C': 235,
		'D': 235,
		'E': 235,
		'F': 235,
		'a': 235,
		'b': 235,
		'c': 235,
		'd': 235,
		'e': 235,
		'f': 235,
	},
	235: {
		'0': 235,
		'1': 235,
		'2': 235,
		'3': 235,
		'4': 235,
		'5': 235,
		'6': 235,
		'7': 235,
		'8': 235,
		'9': 235,
		'A': 235,
		'B': 235,
		'C': 235,
		'D': 235,
		'E': 235,
		'F': 235,
		'a': 235,
		'b': 235,
		'c': 235,
		'd': 235,
		'e': 235,
		'f': 235,
	},
	232: {
		'\t': 232,
		'\n': 232,
		'\r': 232,
		' ': 232,
	},
	249: {
		'=': 239,
	},
	252: {
		'=': 242,
	},
	257: {
		'=': 244,
	},
	250: {
		'=': 240,
	},
	228: {
		'&': 247,
	},
	255: {
		'=': 245,
	},
	272: {
		'0': 274,
		'1': 274,
		'2': 274,
		'3': 274,
		'4': 274,
		'5': 274,
		'6': 274,
		'7': 274,
		'8': 274,
		'9': 274,
		'A': 274,
		'B': 274,
		'C': 274,
		'D': 274,
		'E': 274,
		'F': 274,
		'G': 274,
		'H': 274,
		'I': 274,
		'J': 274,
		'K': 274,
		'L': 274,
		'M': 274,
		'N': 274,
		'O': 274,
		'P': 274,
		'Q': 274,
		'R': 274,
		'S': 274,
		'T': 274,
		'U': 274,
		'V': 274,
		'W': 274,
		'X': 274,
		'Y': 274,
		'Z': 274,
		'_': 274,
		'a': 274,
		'b': 274,
		'c': 274,
		'd': 274,
		'e': 274,
		'f': 274,
		'g': 274,
		'h': 274,
		'i': 274,
		'j': 274,
		'k': 274,
		'l': 270,
		'm': 274,
		'n': 274,
		'o': 274,
		'p': 274,
		'q': 274,
		'r': 274,
		's': 274,
		't': 274,
		'u': 274,
		'v': 274,
		'w': 274,
		'x': 274,
		'y': 274,
		'z': 274,
	},
	270: {
		'0': 274,
		'1': 274,
		'2': 274,
		'3': 274,
		'4': 274,
		'5': 274,
		'6': 274,
		'7': 274,
		'8': 274,
		'9': 274,
		'A': 274,
		'B': 274,
		'C': 274,
		'D': 274,
		'E': 274,
		'F': 274,
		'G': 274,
		'H': 274,
		'I': 274,
		'J': 274,
		'K': 274,
		'L': 274,
		'M': 274,
		'N': 274,
		'O': 274,
		'P': 274,
		'Q': 274,
		'R': 274,
		'S': 274,
		'T': 274,
		'U': 274,
		'V': 274,
		'W': 274,
		'X': 274,
		'Y': 274,
		'Z': 274,
		'_': 274,
		'a': 274,
		'b': 274,
		'c': 274,
		'd': 274,
		'e': 274,
		'f': 274,
		'g': 274,
		'h': 274,
		'i': 274,
		'j': 274,
		'k': 274,
		'l': 274,
		'm': 274,
		'n': 274,
		'o': 274,
		'p': 274,
		'q': 274,
		'r': 274,
		's': 265,
		't': 274,
		'u': 274,
		'v': 274,
		'w': 274,
		'x': 274,
		'y': 274,
		'z': 274,
	},
	265: {
		'0': 274,
		'1': 274,
		'2': 274,
		'3': 274,
		'4': 274,
		'5': 274,
		'6': 274,
		'7': 274,
		'8': 274,
		'9': 274,
		'A': 274,
		'B': 274,
		'C': 274,
		'D': 274,
		'E': 274,
		'F': 274,
		'G': 274,
		'H': 274,
		'I': 274,
		'J': 274,
		'K': 274,
		'L': 274,
		'M': 274,
		'N': 274,
		'O': 274,
		'P': 274,
		'Q': 274,
		'R': 274,
		'S': 274,
		'T': 274,
		'U': 274,
		'V': 274,
		'W': 274,
		'X': 274,
		'Y': 274,
		'Z': 274,
		'_': 274,
		'a': 274,
		'b': 274,
		'c': 274,
		'd': 274,
		'e': 263,
		'f': 274,
		'g': 274,
		'h': 274,
		'i': 274,
		'j': 274,
		'k': 274,
		'l': 274,
		'm': 274,
		'n': 274,
		'o': 274,
		'p': 274,
		'q': 274,
		'r': 274,
		's': 274,
		't': 274,
		'u': 274,
		'v': 274,
		'w': 274,
		'x': 274,
		'y': 274,
		'z': 274,
	},
	263: {
		'0': 274,
		'1': 274,
		'2': 274,
		'3': 274,
		'4': 274,
		'5': 274,
		'6': 274,
		'7': 274,
		'8': 274,
		'9': 274,
		'A': 274,
		'B': 274,
		'C': 274,
		'D': 274,
		'E': 274,
		'F': 274,
		'G': 274,
		'H': 274,
		'I': 274,
		'J': 274,
		'K': 274,
		'L': 274,
		'M': 274,
		'N': 274,
		'O': 274,
		'P': 274,
		'Q': 274,
		'R': 274,
		'S': 274,
		'T': 274,
		'U': 274,
		'V': 274,
		'W': 274,
		'X': 274,
		'Y': 274,
		'Z': 274,
		'_': 274,
		'a': 274,
		'b': 274,
		'c': 274,
		'd': 274,
		'e': 274,
		'f': 274,
		'g': 274,
		'h': 274,
		'i': 274,
		'j': 274,
		'k': 274,
		'l': 274,
		'm': 274,
		'n': 274,
		'o': 274,
		'p': 274,
		'q': 274,
		'r': 274,
		's': 274,
		't': 274,
		'u': 274,
		'v': 274,
		'w': 274,
		'x': 274,
		'y': 274,
		'z': 274,
	},
	266: {
		'0': 274,
		'1': 274,
		'2': 274,
		'3': 274,
		'4': 274,
		'5': 274,
		'6': 274,
		'7': 274,
		'8': 274,
		'9': 274,
		'A': 274,
		'B': 274,
		'C': 274,
		'D': 274,
		'E': 274,
		'F': 274,
		'G': 274,
		'H': 274,
		'I': 274,
		'J': 274,
		'K': 274,
		'L': 274,
		'M': 274,
		'N': 274,
		'O': 274,
		'P': 274,
		'Q': 274,
		'R': 274,
		'S': 274,
		'T': 274,
		'U': 274,
		'V': 274,
		'W': 274,
		'X': 274,
		'Y': 274,
		'Z': 274,
		'_': 274,
		'a': 274,
		'b': 274,
		'c': 274,
		'd': 274,
		'e': 274,
		'f': 262,
		'g': 274,
		'h': 274,
		'i': 274,
		'j': 274,
		'k': 274,
		'l': 274,
		'm': 274,
		'n': 274,
		'o': 274,
		'p': 274,
		'q': 274,
		'r': 274,
		's': 274,
		't': 274,
		'u': 274,
		'v': 274,
		'w': 274,
		'x': 274,
		'y': 274,
		'z': 274,
	},
	262: {
		'0': 274,
		'1': 274,
		'2': 274,
		'3': 274,
		'4': 274,
		'5': 274,
		'6': 274,
		'7': 274,
		'8': 274,
		'9': 274,
		'A': 274,
		'B': 274,
		'C': 274,
		'D': 274,
		'E': 274,
		'F': 274,
		'G': 274,
		'H': 274,
		'I': 274,
		'J': 274,
		'K': 274,
		'L': 274,
		'M': 274,
		'N': 274,
		'O': 274,
		'P': 274,
		'Q': 274,
		'R': 274,
		'S': 274,
		'T': 274,
		'U': 274,
		'V': 274,
		'W': 274,
		'X': 274,
		'Y': 274,
		'Z': 274,
		'_': 274,
		'a': 274,
		'b': 274,
		'c': 274,
		'd': 274,
		'e': 274,
		'f': 274,
		'g': 274,
		'h': 274,
		'i': 274,
		'j': 274,
		'k': 274,
		'l': 274,
		'm': 274,
		'n': 274,
		'o': 274,
		'p': 274,
		'q': 274,
		'r': 274,
		's': 274,
		't': 274,
		'u': 274,
		'v': 274,
		'w': 274,
		'x': 274,
		'y': 274,
		'z': 274,
	},
	268: {
		'0': 274,
		'1': 274,
		'2': 274,
		'3': 274,
		'4': 274,
		'5': 274,
		'6': 274,
		'7': 274,
		'8': 274,
		'9': 274,
		'A': 274,
		'B': 274,
		'C': 274,
		'D': 274,
		'E': 274,
		'F': 274,
		'G': 274,
		'H': 274,
		'I': 274,
		'J': 274,
		'K': 274,
		'L': 274,
		'M': 274,
		'N': 274,
		'O': 274,
		'P': 274,
		'Q': 274,
		'R': 274,
		'S': 274,
		'T': 274,
		'U': 274,
		'V': 274,
		'W': 274,
		'X': 274,
		'Y': 274,
		'Z': 274,
		'_': 274,
		'a': 274,
		'b': 274,
		'c': 274,
		'd': 274,
		'e': 267,
		'f': 274,
		'g': 274,
		'h': 274,
		'i': 274,
		'j': 274,
		'k': 274,
		'l': 274,
		'm': 274,
		'n': 274,
		'o': 274,
		'p': 274,
		'q': 274,
		'r': 274,
		's': 274,
		't': 274,
		'u': 274,
		'v': 274,
		'w': 274,
		'x': 274,
		'y': 274,
		'z': 274,
	},
	267: {
		'0': 274,
		'1': 274,
		'2': 274,
		'3': 274,
		'4': 274,
		'5': 274,
		'6': 274,
		'7': 274,
		'8': 274,
		'9': 274,
		'A': 274,
		'B': 274,
		'C': 274,
		'D': 274,
		'E': 274,
		'F': 274,
		'G': 274,
		'H': 274,
		'I': 274,
		'J': 274,
		'K': 274,
		'L': 274,
		'M': 274,
		'N': 274,
		'O': 274,
		'P': 274,
		'Q': 274,
		'R': 274,
		'S': 274,
		'T': 274,
		'U': 274,
		'V': 274,
		'W': 274,
		'X': 274,
		'Y': 274,
		'Z': 274,
		'_': 274,
		'a': 274,
		'b': 274,
		'c': 274,
		'd': 274,
		'e': 274,
		'f': 274,
		'g': 274,
		'h': 274,
		'i': 274,
		'j': 274,
		'k': 274,
		'l': 274,
		'm': 274,
		'n': 274,
		'o': 274,
		'p': 274,
		'q': 274,
		'r': 274,
		's': 274,
		't': 261,
		'u': 274,
		'v': 274,
		'w': 274,
		'x': 274,
		'y': 274,
		'z': 274,
	},
	261: {
		'0': 274,
		'1': 274,
		'2': 274,
		'3': 274,
		'4': 274,
		'5': 274,
		'6': 274,
		'7': 274,
		'8': 274,
		'9': 274,
		'A': 274,
		'B': 274,
		'C': 274,
		'D': 274,
		'E': 274,
		'F': 274,
		'G': 274,
		'H': 274,
		'I': 274,
		'J': 274,
		'K': 274,
		'L': 274,
		'M': 274,
		'N': 274,
		'O': 274,
		'P': 274,
		'Q': 274,
		'R': 274,
		'S': 274,
		'T': 274,
		'U': 274,
		'V': 274,
		'W': 274,
		'X': 274,
		'Y': 274,
		'Z': 274,
		'_': 274,
		'a': 274,
		'b': 274,
		'c': 274,
		'd': 274,
		'e': 274,
		'f': 274,
		'g': 274,
		'h': 274,
		'i': 274,
		'j': 274,
		'k': 274,
		'l': 274,
		'm': 274,
		'n': 274,
		'o': 274,
		'p': 274,
		'q': 274,
		'r': 274,
		's': 274,
		't': 274,
		'u': 274,
		'v': 274,
		'w': 274,
		'x': 274,
		'y': 274,
		'z': 274,
	},
	231: {
		'|': 248,
	},
	254: {
		'=': 243,
	},
	256: {
		'=': 246,
	},
	273: {
		'0': 274,
		'1': 274,
		'2': 274,
		'3': 274,
		'4': 274,
		'5': 274,
		'6': 274,
		'7': 274,
		'8': 274,
		'9': 274,
		'A': 274,
		'B': 274,
		'C': 274,
		'D': 274,
		'E': 274,
		'F': 274,
		'G': 274,
		'H': 274,
		'I': 274,
		'J': 274,
		'K': 274,
		'L': 274,
		'M': 274,
		'N': 274,
		'O': 274,
		'P': 274,
		'Q': 274,
		'R': 274,
		'S': 274,
		'T': 274,
		'U': 274,
		'V': 274,
		'W': 274,
		'X': 274,
		'Y': 274,
		'Z': 274,
		'_': 274,
		'a': 274,
		'b': 274,
		'c': 274,
		'd': 274,
		'e': 274,
		'f': 274,
		'g': 274,
		'h': 271,
		'i': 274,
		'j': 274,
		'k': 274,
		'l': 274,
		'm': 274,
		'n': 274,
		'o': 274,
		'p': 274,
		'q': 274,
		'r': 274,
		's': 274,
		't': 274,
		'u': 274,
		'v': 274,
		'w': 274,
		'x': 274,
		'y': 274,
		'z': 274,
	},
	271: {
		'0': 274,
		'1': 274,
		'2': 274,
		'3': 274,
		'4': 274,
		'5': 274,
		'6': 274,
		'7': 274,
		'8': 274,
		'9': 274,
		'A': 274,
		'B': 274,
		'C': 274,
		'D': 274,
		'E': 274,
		'F': 274,
		'G': 274,
		'H': 274,
		'I': 274,
		'J': 274,
		'K': 274,
		'L': 274,
		'M': 274,
		'N': 274,
		'O': 274,
		'P': 274,
		'Q': 274,
		'R': 274,
		'S': 274,
		'T': 274,
		'U': 274,
		'V': 274,
		'W': 274,
		'X': 274,
		'Y': 274,
		'Z': 274,
		'_': 274,
		'a': 274,
		'b': 274,
		'c': 274,
		'd': 274,
		'e': 274,
		'f': 274,
		'g': 274,
		'h': 274,
		'i': 269,
		'j': 274,
		'k': 274,
		'l': 274,
		'm': 274,
		'n': 274,
		'o': 274,
		'p': 274,
		'q': 274,
		'r': 274,
		's': 274,
		't': 274,
		'u': 274,
		'v': 274,
		'w': 274,
		'x': 274,
		'y': 274,
		'z': 274,
	},
	269: {
		'0': 274,
		'1': 274,
		'2': 274,
		'3': 274,
		'4': 274,
		'5': 274,
		'6': 274,
		'7': 274,
		'8': 274,
		'9': 274,
		'A': 274,
		'B': 274,
		'C': 274,
		'D': 274,
		'E': 274,
		'F': 274,
		'G': 274,
		'H': 274,
		'I': 274,
		'J': 274,
		'K': 274,
		'L': 274,
		'M': 274,
		'N': 274,
		'O': 274,
		'P': 274,
		'Q': 274,
		'R': 274,
		'S': 274,
		'T': 274,
		'U': 274,
		'V': 274,
		'W': 274,
		'X': 274,
		'Y': 274,
		'Z': 274,
		'_': 274,
		'a': 274,
		'b': 274,
		'c': 274,
		'd': 274,
		'e': 274,
		'f': 274,
		'g': 274,
		'h': 274,
		'i': 274,
		'j': 274,
		'k': 274,
		'l': 264,
		'm': 274,
		'n': 274,
		'o': 274,
		'p': 274,
		'q': 274,
		'r': 274,
		's': 274,
		't': 274,
		'u': 274,
		'v': 274,
		'w': 274,
		'x': 274,
		'y': 274,
		'z': 274,
	},
	264: {
		'0': 274,
		'1': 274,
		'2': 274,
		'3': 274,
		'4': 274,
		'5': 274,
		'6': 274,
		'7': 274,
		'8': 274,
		'9': 274,
		'A': 274,
		'B': 274,
		'C': 274,
		'D': 274,
		'E': 274,
		'F': 274,
		'G': 274,
		'H': 274,
		'I': 274,
		'J': 274,
		'K': 274,
		'L': 274,
		'M': 274,
		'N': 274,
		'O': 274,
		'P': 274,
		'Q': 274,
		'R': 274,
		'S': 274,
		'T': 274,
		'U': 274,
		'V': 274,
		'W': 274,
		'X': 274,
		'Y': 274,
		'Z': 274,
		'_': 274,
		'a': 274,
		'b': 274,
		'c': 274,
		'd': 274,
		'e': 233,
		'f': 274,
		'g': 274,
		'h': 274,
		'i': 274,
		'j': 274,
		'k': 274,
		'l': 274,
		'm': 274,
		'n': 274,
		'o': 274,
		'p': 274,
		'q': 274,
		'r': 274,
		's': 274,
		't': 274,
		'u': 274,
		'v': 274,
		'w': 274,
		'x': 274,
		'y': 274,
		'z': 274,
	},
	233: {
		'0': 274,
		'1': 274,
		'2': 274,
		'3': 274,
		'4': 274,
		'5': 274,
		'6': 274,
		'7': 274,
		'8': 274,
		'9': 274,
		'A': 274,
		'B': 274,
		'C': 274,
		'D': 274,
		'E': 274,
		'F': 274,
		'G': 274,
		'H': 274,
		'I': 274,
		'J': 274,
		'K': 274,
		'L': 274,
		'M': 274,
		'N': 274,
		'O': 274,
		'P': 274,
		'Q': 274,
		'R': 274,
		'S': 274,
		'T': 274,
		'U': 274,
		'V': 274,
		'W': 274,
		'X': 274,
		'Y': 274,
		'Z': 274,
		'_': 274,
		'a': 274,
		'b': 274,
		'c': 274,
		'd': 274,
		'e': 274,
		'f': 274,
		'g': 274,
		'h': 274,
		'i': 274,
		'j': 274,
		'k': 274,
		'l': 274,
		'm': 274,
		'n': 274,
		'o': 274,
		'p': 274,
		'q': 274,
		'r': 274,
		's': 274,
		't': 274,
		'u': 274,
		'v': 274,
		'w': 274,
		'x': 274,
		'y': 274,
		'z': 274,
	},
	251: {
		'*': 238,
		'=': 241,
	},
}

// calcAccept maps accepting-state IDs to their 1-based TokenID.
var calcAccept = map[int]int{
	274: 32,
	236: 5,
	237: 5,
	234: 3,
	235: 4,
	232: 1,
	249: 17,
	239: 7,
	252: 20,
	242: 10,
	257: 25,
	244: 12,
	250: 18,
	240: 8,
	247: 15,
	259: 27,
	255: 23,
	245: 13,
	272: 32,
	270: 32,
	265: 32,
	263: 31,
	253: 21,
	266: 32,
	262: 30,
	260: 28,
	268: 32,
	267: 32,
	261: 29,
	248: 16,
	254: 22,
	243: 11,
	256: 24,
	246: 14,
	258: 26,
	273: 32,
	271: 32,
	269: 32,
	264: 32,
	233: 2,
	251: 19,
	238: 6,
	241: 9,
}

// gettoken advances to the next token and returns its ID.
// Lxm, Ln, and Col are set before the action runs.
// Actions that return emit the token to the caller.
// Actions that don't return let the scan loop continue (skip).
// Consecutive unrecognised characters are grouped into one ERROR token.
// Returns EOF (0) at end of input, ERROR (-1) for unrecognised characters.
func (l *Lexer) gettoken() int {
	for l.pos < len(l.input) {
		// Snapshot the position at the start of each token attempt.
		// If the inner loop overshoots, we backtrack here via lastPos/lastLine/lastCol.
		startPos  := l.pos
		startLine := l.line
		startCol  := l.col

		state    := 227  // start state of the DFA
		lastTok  := 0   // TokenID of the last accepting state seen (0 = none yet)
		lastPos  := l.pos
		lastLine := l.line
		lastCol  := l.col
		curLine  := l.line
		curCol   := l.col

		// Inner loop: maximal munch — keep advancing as long as the DFA has a
		// transition. Every time we land on an accepting state, snapshot it.
		// When no transition exists we stop and backtrack to the last snapshot.
		for l.pos < len(l.input) {
			ch := l.input[l.pos]
			row, ok := calcTrans[state]
			if !ok {
				break
			}
			next, ok := row[ch]
			if !ok {
				break
			}
			l.pos++
			if ch == '\n' {
				curLine++
				curCol = 1
			} else {
				curCol++
			}
			state = next
			if tok := calcAccept[state]; tok != 0 {
				// Accepting state — snapshot position and token ID.
				// A later, longer match may overwrite this snapshot.
				lastTok  = tok
				lastPos  = l.pos
				lastLine = curLine
				lastCol  = curCol
			}
		}

		if lastTok == 0 {
			// No accepting state was ever reached from startPos — unrecognised input.
			// Consume characters one-by-one until the DFA's start state can fire again,
			// grouping them all into a single ERROR token rather than one per character.
			errStart := startPos
			errLn    := l.line
			errCol   := l.col
			for l.pos < len(l.input) {
				ch := l.input[l.pos]
				// Stop as soon as the start state has a transition for this char —
				// the next call to Scan will pick up from here cleanly.
				if row, ok := calcTrans[227]; ok {
					if _, ok := row[ch]; ok {
						break
					}
				}
				if ch == '\n' {
					l.line++
					l.col = 1
				} else {
					l.col++
				}
				l.pos++
			}
			l.Lxm = string(l.input[errStart:l.pos])
			l.Ln  = errLn
			l.Col = errCol
			return ERROR
		}

		// Backtrack to the last accepting snapshot and commit that match.
		l.pos = lastPos
		l.line = lastLine
		l.col  = lastCol
		l.Lxm  = string(l.input[startPos:lastPos])
		l.Ln   = startLine
		l.Col  = startCol

		// Dispatch to the verbatim action for the winning token.
		// Actions that execute a return statement exit Scan() with that token ID.
		// Actions without a return fall through and the outer loop retries (skip).
		switch lastTok {
		case 1:
			// no action
		case 2:
			fmt.Print("WHILE  "); return WHILE
		case 3:
			fmt.Print("FLOAT  "); return FLOAT
		case 4:
			fmt.Print("HEX    "); return HEX
		case 5:
			fmt.Print("INT    "); return INT
		case 6:
			fmt.Print("POW    "); return POW
		case 7:
			fmt.Print("PLEQ   "); return PLEQ
		case 8:
			fmt.Print("MIEQ   "); return MIEQ
		case 9:
			fmt.Print("TIEQ   "); return TIEQ
		case 10:
			fmt.Print("DIEQ   "); return DIEQ
		case 11:
			fmt.Print("EQ     "); return EQ
		case 12:
			fmt.Print("NEQ    "); return NEQ
		case 13:
			fmt.Print("LEQ    "); return LEQ
		case 14:
			fmt.Print("GEQ    "); return GEQ
		case 15:
			fmt.Print("AND    "); return AND
		case 16:
			fmt.Print("OR     "); return OR
		case 17:
			fmt.Print("PLUS   "); return PLUS
		case 18:
			fmt.Print("MINUS  "); return MINUS
		case 19:
			fmt.Print("TIMES  "); return TIMES
		case 20:
			fmt.Print("DIV    "); return DIV
		case 21:
			fmt.Print("MOD    "); return MOD
		case 22:
			fmt.Print("ASSN   "); return ASSN
		case 23:
			fmt.Print("LT     "); return LT
		case 24:
			fmt.Print("GT     "); return GT
		case 25:
			fmt.Print("NOT    "); return NOT
		case 26:
			fmt.Print("LPAREN "); return LPAREN
		case 27:
			fmt.Print("RPAREN "); return RPAREN
		case 28:
			fmt.Print("COMMA  "); return COMMA
		case 29:
			fmt.Print("KWLET  "); return KWLET
		case 30:
			fmt.Print("KWIF   "); return KWIF
		case 31:
			fmt.Print("KWELSE "); return KWELSE
		case 32:
			fmt.Print("IDENT  "); return IDENT
		}
	}
	return EOF
}

// NextToken advances to the next token and returns a Lexeme.
func (l *Lexer) NextToken() Lexeme {
	tok := l.gettoken()
	return Lexeme{Token: tok, Value: l.Lxm, Line: l.Ln, Col: l.Col}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: calc <inputfile>")
		os.Exit(1)
	}
	data, err := os.ReadFile(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	l := NewCalcLexer(string(data))
	for {
		tok := l.gettoken()
		if tok == EOF {
			break
		}
		if tok == ERROR {
			fmt.Printf("ERROR  %-20q  ln=%d col=%d-%d\n", l.Lxm, l.Ln, l.Col, l.Col+len([]rune(l.Lxm))-1)
			continue
		}
		fmt.Printf("%d  %-20q  ln=%d col=%d\n", tok, l.Lxm, l.Ln, l.Col)
	}
}
