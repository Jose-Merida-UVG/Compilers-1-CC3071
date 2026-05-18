package main

import (
	"fmt"
	"os"
)

// parserAction encodes one cell of the SLR(1) ACTION table.
// kind: 1=shift 2=reduce 3=accept; arg: target state (shift) or prod index (reduce).
type parserAction struct{ kind, arg int }

// parserActionTable[state][terminal] → action
var parserActionTable = map[int]map[string]parserAction{
	0: {
		"FLOAT": {1, 12},
		"HEX": {1, 5},
		"IDENT": {1, 6},
		"INT": {1, 8},
		"KWIF": {1, 9},
		"KWLET": {1, 15},
		"LPAREN": {1, 13},
		"MINUS": {1, 4},
		"NOT": {1, 3},
	},
	1: {
		"$": {3, 0},
	},
	2: {
		"$": {2, 15},
		"AND": {2, 15},
		"COMMA": {2, 15},
		"EQ": {1, 26},
		"GEQ": {1, 23},
		"GT": {1, 22},
		"KWELSE": {2, 15},
		"LEQ": {1, 21},
		"LT": {1, 20},
		"MINUS": {1, 25},
		"NEQ": {1, 19},
		"OR": {2, 15},
		"PLUS": {1, 24},
		"RPAREN": {2, 15},
	},
	3: {
		"FLOAT": {1, 12},
		"HEX": {1, 5},
		"IDENT": {1, 6},
		"INT": {1, 8},
		"LPAREN": {1, 13},
		"MINUS": {1, 4},
		"NOT": {1, 3},
	},
	4: {
		"FLOAT": {1, 12},
		"HEX": {1, 5},
		"IDENT": {1, 6},
		"INT": {1, 8},
		"LPAREN": {1, 13},
		"MINUS": {1, 4},
		"NOT": {1, 3},
	},
	5: {
		"$": {2, 28},
		"AND": {2, 28},
		"COMMA": {2, 28},
		"DIV": {2, 28},
		"EQ": {2, 28},
		"GEQ": {2, 28},
		"GT": {2, 28},
		"KWELSE": {2, 28},
		"LEQ": {2, 28},
		"LT": {2, 28},
		"MINUS": {2, 28},
		"MOD": {2, 28},
		"NEQ": {2, 28},
		"OR": {2, 28},
		"PLUS": {2, 28},
		"RPAREN": {2, 28},
		"TIMES": {2, 28},
	},
	6: {
		"$": {2, 29},
		"AND": {2, 29},
		"COMMA": {2, 29},
		"DIV": {2, 29},
		"EQ": {2, 29},
		"GEQ": {2, 29},
		"GT": {2, 29},
		"KWELSE": {2, 29},
		"LEQ": {2, 29},
		"LPAREN": {1, 29},
		"LT": {2, 29},
		"MINUS": {2, 29},
		"MOD": {2, 29},
		"NEQ": {2, 29},
		"OR": {2, 29},
		"PLUS": {2, 29},
		"RPAREN": {2, 29},
		"TIMES": {2, 29},
	},
	7: {
		"$": {2, 4},
		"COMMA": {2, 4},
		"OR": {1, 30},
		"RPAREN": {2, 4},
	},
	8: {
		"$": {2, 26},
		"AND": {2, 26},
		"COMMA": {2, 26},
		"DIV": {2, 26},
		"EQ": {2, 26},
		"GEQ": {2, 26},
		"GT": {2, 26},
		"KWELSE": {2, 26},
		"LEQ": {2, 26},
		"LT": {2, 26},
		"MINUS": {2, 26},
		"MOD": {2, 26},
		"NEQ": {2, 26},
		"OR": {2, 26},
		"PLUS": {2, 26},
		"RPAREN": {2, 26},
		"TIMES": {2, 26},
	},
	9: {
		"FLOAT": {1, 12},
		"HEX": {1, 5},
		"IDENT": {1, 6},
		"INT": {1, 8},
		"LPAREN": {1, 13},
		"MINUS": {1, 4},
		"NOT": {1, 3},
	},
	10: {
		"$": {2, 6},
		"AND": {1, 32},
		"COMMA": {2, 6},
		"KWELSE": {2, 6},
		"OR": {2, 6},
		"RPAREN": {2, 6},
	},
	11: {
		"$": {2, 25},
		"AND": {2, 25},
		"COMMA": {2, 25},
		"DIV": {2, 25},
		"EQ": {2, 25},
		"GEQ": {2, 25},
		"GT": {2, 25},
		"KWELSE": {2, 25},
		"LEQ": {2, 25},
		"LT": {2, 25},
		"MINUS": {2, 25},
		"MOD": {2, 25},
		"NEQ": {2, 25},
		"OR": {2, 25},
		"PLUS": {2, 25},
		"RPAREN": {2, 25},
		"TIMES": {2, 25},
	},
	12: {
		"$": {2, 27},
		"AND": {2, 27},
		"COMMA": {2, 27},
		"DIV": {2, 27},
		"EQ": {2, 27},
		"GEQ": {2, 27},
		"GT": {2, 27},
		"KWELSE": {2, 27},
		"LEQ": {2, 27},
		"LT": {2, 27},
		"MINUS": {2, 27},
		"MOD": {2, 27},
		"NEQ": {2, 27},
		"OR": {2, 27},
		"PLUS": {2, 27},
		"RPAREN": {2, 27},
		"TIMES": {2, 27},
	},
	13: {
		"FLOAT": {1, 12},
		"HEX": {1, 5},
		"IDENT": {1, 6},
		"INT": {1, 8},
		"KWIF": {1, 9},
		"KWLET": {1, 15},
		"LPAREN": {1, 13},
		"MINUS": {1, 4},
		"NOT": {1, 3},
	},
	14: {
		"$": {2, 1},
	},
	15: {
		"IDENT": {1, 34},
	},
	16: {
		"$": {2, 8},
		"AND": {2, 8},
		"COMMA": {2, 8},
		"KWELSE": {2, 8},
		"OR": {2, 8},
		"RPAREN": {2, 8},
	},
	17: {
		"$": {2, 18},
		"AND": {2, 18},
		"COMMA": {2, 18},
		"DIV": {1, 36},
		"EQ": {2, 18},
		"GEQ": {2, 18},
		"GT": {2, 18},
		"KWELSE": {2, 18},
		"LEQ": {2, 18},
		"LT": {2, 18},
		"MINUS": {2, 18},
		"MOD": {1, 37},
		"NEQ": {2, 18},
		"OR": {2, 18},
		"PLUS": {2, 18},
		"RPAREN": {2, 18},
		"TIMES": {1, 35},
	},
	18: {
		"$": {2, 22},
		"AND": {2, 22},
		"COMMA": {2, 22},
		"DIV": {2, 22},
		"EQ": {2, 22},
		"GEQ": {2, 22},
		"GT": {2, 22},
		"KWELSE": {2, 22},
		"LEQ": {2, 22},
		"LT": {2, 22},
		"MINUS": {2, 22},
		"MOD": {2, 22},
		"NEQ": {2, 22},
		"OR": {2, 22},
		"PLUS": {2, 22},
		"RPAREN": {2, 22},
		"TIMES": {2, 22},
	},
	19: {
		"FLOAT": {1, 12},
		"HEX": {1, 5},
		"IDENT": {1, 6},
		"INT": {1, 8},
		"LPAREN": {1, 13},
		"MINUS": {1, 4},
		"NOT": {1, 3},
	},
	20: {
		"FLOAT": {1, 12},
		"HEX": {1, 5},
		"IDENT": {1, 6},
		"INT": {1, 8},
		"LPAREN": {1, 13},
		"MINUS": {1, 4},
		"NOT": {1, 3},
	},
	21: {
		"FLOAT": {1, 12},
		"HEX": {1, 5},
		"IDENT": {1, 6},
		"INT": {1, 8},
		"LPAREN": {1, 13},
		"MINUS": {1, 4},
		"NOT": {1, 3},
	},
	22: {
		"FLOAT": {1, 12},
		"HEX": {1, 5},
		"IDENT": {1, 6},
		"INT": {1, 8},
		"LPAREN": {1, 13},
		"MINUS": {1, 4},
		"NOT": {1, 3},
	},
	23: {
		"FLOAT": {1, 12},
		"HEX": {1, 5},
		"IDENT": {1, 6},
		"INT": {1, 8},
		"LPAREN": {1, 13},
		"MINUS": {1, 4},
		"NOT": {1, 3},
	},
	24: {
		"FLOAT": {1, 12},
		"HEX": {1, 5},
		"IDENT": {1, 6},
		"INT": {1, 8},
		"LPAREN": {1, 13},
		"MINUS": {1, 4},
		"NOT": {1, 3},
	},
	25: {
		"FLOAT": {1, 12},
		"HEX": {1, 5},
		"IDENT": {1, 6},
		"INT": {1, 8},
		"LPAREN": {1, 13},
		"MINUS": {1, 4},
		"NOT": {1, 3},
	},
	26: {
		"FLOAT": {1, 12},
		"HEX": {1, 5},
		"IDENT": {1, 6},
		"INT": {1, 8},
		"LPAREN": {1, 13},
		"MINUS": {1, 4},
		"NOT": {1, 3},
	},
	27: {
		"$": {2, 23},
		"AND": {2, 23},
		"COMMA": {2, 23},
		"DIV": {2, 23},
		"EQ": {2, 23},
		"GEQ": {2, 23},
		"GT": {2, 23},
		"KWELSE": {2, 23},
		"LEQ": {2, 23},
		"LT": {2, 23},
		"MINUS": {2, 23},
		"MOD": {2, 23},
		"NEQ": {2, 23},
		"OR": {2, 23},
		"PLUS": {2, 23},
		"RPAREN": {2, 23},
		"TIMES": {2, 23},
	},
	28: {
		"$": {2, 24},
		"AND": {2, 24},
		"COMMA": {2, 24},
		"DIV": {2, 24},
		"EQ": {2, 24},
		"GEQ": {2, 24},
		"GT": {2, 24},
		"KWELSE": {2, 24},
		"LEQ": {2, 24},
		"LT": {2, 24},
		"MINUS": {2, 24},
		"MOD": {2, 24},
		"NEQ": {2, 24},
		"OR": {2, 24},
		"PLUS": {2, 24},
		"RPAREN": {2, 24},
		"TIMES": {2, 24},
	},
	29: {
		"FLOAT": {1, 12},
		"HEX": {1, 5},
		"IDENT": {1, 6},
		"INT": {1, 8},
		"KWIF": {1, 9},
		"KWLET": {1, 15},
		"LPAREN": {1, 13},
		"MINUS": {1, 4},
		"NOT": {1, 3},
		"RPAREN": {2, 33},
	},
	30: {
		"FLOAT": {1, 12},
		"HEX": {1, 5},
		"IDENT": {1, 6},
		"INT": {1, 8},
		"LPAREN": {1, 13},
		"MINUS": {1, 4},
		"NOT": {1, 3},
	},
	31: {
		"KWELSE": {1, 50},
		"OR": {1, 30},
	},
	32: {
		"FLOAT": {1, 12},
		"HEX": {1, 5},
		"IDENT": {1, 6},
		"INT": {1, 8},
		"LPAREN": {1, 13},
		"MINUS": {1, 4},
		"NOT": {1, 3},
	},
	33: {
		"RPAREN": {1, 52},
	},
	34: {
		"ASSN": {1, 53},
	},
	35: {
		"FLOAT": {1, 12},
		"HEX": {1, 5},
		"IDENT": {1, 6},
		"INT": {1, 8},
		"LPAREN": {1, 13},
		"MINUS": {1, 4},
		"NOT": {1, 3},
	},
	36: {
		"FLOAT": {1, 12},
		"HEX": {1, 5},
		"IDENT": {1, 6},
		"INT": {1, 8},
		"LPAREN": {1, 13},
		"MINUS": {1, 4},
		"NOT": {1, 3},
	},
	37: {
		"FLOAT": {1, 12},
		"HEX": {1, 5},
		"IDENT": {1, 6},
		"INT": {1, 8},
		"LPAREN": {1, 13},
		"MINUS": {1, 4},
		"NOT": {1, 3},
	},
	38: {
		"$": {2, 10},
		"AND": {2, 10},
		"COMMA": {2, 10},
		"KWELSE": {2, 10},
		"MINUS": {1, 25},
		"OR": {2, 10},
		"PLUS": {1, 24},
		"RPAREN": {2, 10},
	},
	39: {
		"$": {2, 11},
		"AND": {2, 11},
		"COMMA": {2, 11},
		"KWELSE": {2, 11},
		"MINUS": {1, 25},
		"OR": {2, 11},
		"PLUS": {1, 24},
		"RPAREN": {2, 11},
	},
	40: {
		"$": {2, 12},
		"AND": {2, 12},
		"COMMA": {2, 12},
		"KWELSE": {2, 12},
		"MINUS": {1, 25},
		"OR": {2, 12},
		"PLUS": {1, 24},
		"RPAREN": {2, 12},
	},
	41: {
		"$": {2, 13},
		"AND": {2, 13},
		"COMMA": {2, 13},
		"KWELSE": {2, 13},
		"MINUS": {1, 25},
		"OR": {2, 13},
		"PLUS": {1, 24},
		"RPAREN": {2, 13},
	},
	42: {
		"$": {2, 14},
		"AND": {2, 14},
		"COMMA": {2, 14},
		"KWELSE": {2, 14},
		"MINUS": {1, 25},
		"OR": {2, 14},
		"PLUS": {1, 24},
		"RPAREN": {2, 14},
	},
	43: {
		"$": {2, 16},
		"AND": {2, 16},
		"COMMA": {2, 16},
		"DIV": {1, 36},
		"EQ": {2, 16},
		"GEQ": {2, 16},
		"GT": {2, 16},
		"KWELSE": {2, 16},
		"LEQ": {2, 16},
		"LT": {2, 16},
		"MINUS": {2, 16},
		"MOD": {1, 37},
		"NEQ": {2, 16},
		"OR": {2, 16},
		"PLUS": {2, 16},
		"RPAREN": {2, 16},
		"TIMES": {1, 35},
	},
	44: {
		"$": {2, 17},
		"AND": {2, 17},
		"COMMA": {2, 17},
		"DIV": {1, 36},
		"EQ": {2, 17},
		"GEQ": {2, 17},
		"GT": {2, 17},
		"KWELSE": {2, 17},
		"LEQ": {2, 17},
		"LT": {2, 17},
		"MINUS": {2, 17},
		"MOD": {1, 37},
		"NEQ": {2, 17},
		"OR": {2, 17},
		"PLUS": {2, 17},
		"RPAREN": {2, 17},
		"TIMES": {1, 35},
	},
	45: {
		"$": {2, 9},
		"AND": {2, 9},
		"COMMA": {2, 9},
		"KWELSE": {2, 9},
		"MINUS": {1, 25},
		"OR": {2, 9},
		"PLUS": {1, 24},
		"RPAREN": {2, 9},
	},
	46: {
		"RPAREN": {1, 57},
	},
	47: {
		"COMMA": {2, 34},
		"RPAREN": {2, 34},
	},
	48: {
		"COMMA": {1, 58},
		"RPAREN": {2, 32},
	},
	49: {
		"$": {2, 5},
		"AND": {1, 32},
		"COMMA": {2, 5},
		"KWELSE": {2, 5},
		"OR": {2, 5},
		"RPAREN": {2, 5},
	},
	50: {
		"FLOAT": {1, 12},
		"HEX": {1, 5},
		"IDENT": {1, 6},
		"INT": {1, 8},
		"LPAREN": {1, 13},
		"MINUS": {1, 4},
		"NOT": {1, 3},
	},
	51: {
		"$": {2, 7},
		"AND": {2, 7},
		"COMMA": {2, 7},
		"KWELSE": {2, 7},
		"OR": {2, 7},
		"RPAREN": {2, 7},
	},
	52: {
		"$": {2, 30},
		"AND": {2, 30},
		"COMMA": {2, 30},
		"DIV": {2, 30},
		"EQ": {2, 30},
		"GEQ": {2, 30},
		"GT": {2, 30},
		"KWELSE": {2, 30},
		"LEQ": {2, 30},
		"LT": {2, 30},
		"MINUS": {2, 30},
		"MOD": {2, 30},
		"NEQ": {2, 30},
		"OR": {2, 30},
		"PLUS": {2, 30},
		"RPAREN": {2, 30},
		"TIMES": {2, 30},
	},
	53: {
		"FLOAT": {1, 12},
		"HEX": {1, 5},
		"IDENT": {1, 6},
		"INT": {1, 8},
		"KWIF": {1, 9},
		"KWLET": {1, 15},
		"LPAREN": {1, 13},
		"MINUS": {1, 4},
		"NOT": {1, 3},
	},
	54: {
		"$": {2, 19},
		"AND": {2, 19},
		"COMMA": {2, 19},
		"DIV": {2, 19},
		"EQ": {2, 19},
		"GEQ": {2, 19},
		"GT": {2, 19},
		"KWELSE": {2, 19},
		"LEQ": {2, 19},
		"LT": {2, 19},
		"MINUS": {2, 19},
		"MOD": {2, 19},
		"NEQ": {2, 19},
		"OR": {2, 19},
		"PLUS": {2, 19},
		"RPAREN": {2, 19},
		"TIMES": {2, 19},
	},
	55: {
		"$": {2, 20},
		"AND": {2, 20},
		"COMMA": {2, 20},
		"DIV": {2, 20},
		"EQ": {2, 20},
		"GEQ": {2, 20},
		"GT": {2, 20},
		"KWELSE": {2, 20},
		"LEQ": {2, 20},
		"LT": {2, 20},
		"MINUS": {2, 20},
		"MOD": {2, 20},
		"NEQ": {2, 20},
		"OR": {2, 20},
		"PLUS": {2, 20},
		"RPAREN": {2, 20},
		"TIMES": {2, 20},
	},
	56: {
		"$": {2, 21},
		"AND": {2, 21},
		"COMMA": {2, 21},
		"DIV": {2, 21},
		"EQ": {2, 21},
		"GEQ": {2, 21},
		"GT": {2, 21},
		"KWELSE": {2, 21},
		"LEQ": {2, 21},
		"LT": {2, 21},
		"MINUS": {2, 21},
		"MOD": {2, 21},
		"NEQ": {2, 21},
		"OR": {2, 21},
		"PLUS": {2, 21},
		"RPAREN": {2, 21},
		"TIMES": {2, 21},
	},
	57: {
		"$": {2, 31},
		"AND": {2, 31},
		"COMMA": {2, 31},
		"DIV": {2, 31},
		"EQ": {2, 31},
		"GEQ": {2, 31},
		"GT": {2, 31},
		"KWELSE": {2, 31},
		"LEQ": {2, 31},
		"LT": {2, 31},
		"MINUS": {2, 31},
		"MOD": {2, 31},
		"NEQ": {2, 31},
		"OR": {2, 31},
		"PLUS": {2, 31},
		"RPAREN": {2, 31},
		"TIMES": {2, 31},
	},
	58: {
		"FLOAT": {1, 12},
		"HEX": {1, 5},
		"IDENT": {1, 6},
		"INT": {1, 8},
		"KWIF": {1, 9},
		"KWLET": {1, 15},
		"LPAREN": {1, 13},
		"MINUS": {1, 4},
		"NOT": {1, 3},
	},
	59: {
		"$": {2, 3},
		"COMMA": {2, 3},
		"OR": {1, 30},
		"RPAREN": {2, 3},
	},
	60: {
		"$": {2, 2},
		"COMMA": {2, 2},
		"RPAREN": {2, 2},
	},
	61: {
		"COMMA": {2, 35},
		"RPAREN": {2, 35},
	},
}

// parserGotoTable[state][nonTerminal] → next state
var parserGotoTable = map[int]map[string]int{
	0: {
		"addexpr": 2,
		"andexpr": 10,
		"cmpexpr": 16,
		"expr": 14,
		"mulexpr": 17,
		"orexpr": 7,
		"primary": 11,
		"program": 1,
		"unary": 18,
	},
	3: {
		"primary": 11,
		"unary": 27,
	},
	4: {
		"primary": 11,
		"unary": 28,
	},
	9: {
		"addexpr": 2,
		"andexpr": 10,
		"cmpexpr": 16,
		"mulexpr": 17,
		"orexpr": 31,
		"primary": 11,
		"unary": 18,
	},
	13: {
		"addexpr": 2,
		"andexpr": 10,
		"cmpexpr": 16,
		"expr": 33,
		"mulexpr": 17,
		"orexpr": 7,
		"primary": 11,
		"unary": 18,
	},
	19: {
		"addexpr": 38,
		"mulexpr": 17,
		"primary": 11,
		"unary": 18,
	},
	20: {
		"addexpr": 39,
		"mulexpr": 17,
		"primary": 11,
		"unary": 18,
	},
	21: {
		"addexpr": 40,
		"mulexpr": 17,
		"primary": 11,
		"unary": 18,
	},
	22: {
		"addexpr": 41,
		"mulexpr": 17,
		"primary": 11,
		"unary": 18,
	},
	23: {
		"addexpr": 42,
		"mulexpr": 17,
		"primary": 11,
		"unary": 18,
	},
	24: {
		"mulexpr": 43,
		"primary": 11,
		"unary": 18,
	},
	25: {
		"mulexpr": 44,
		"primary": 11,
		"unary": 18,
	},
	26: {
		"addexpr": 45,
		"mulexpr": 17,
		"primary": 11,
		"unary": 18,
	},
	29: {
		"addexpr": 2,
		"andexpr": 10,
		"arglist": 48,
		"args": 46,
		"cmpexpr": 16,
		"expr": 47,
		"mulexpr": 17,
		"orexpr": 7,
		"primary": 11,
		"unary": 18,
	},
	30: {
		"addexpr": 2,
		"andexpr": 49,
		"cmpexpr": 16,
		"mulexpr": 17,
		"primary": 11,
		"unary": 18,
	},
	32: {
		"addexpr": 2,
		"cmpexpr": 51,
		"mulexpr": 17,
		"primary": 11,
		"unary": 18,
	},
	35: {
		"primary": 11,
		"unary": 54,
	},
	36: {
		"primary": 11,
		"unary": 55,
	},
	37: {
		"primary": 11,
		"unary": 56,
	},
	50: {
		"addexpr": 2,
		"andexpr": 10,
		"cmpexpr": 16,
		"mulexpr": 17,
		"orexpr": 59,
		"primary": 11,
		"unary": 18,
	},
	53: {
		"addexpr": 2,
		"andexpr": 10,
		"cmpexpr": 16,
		"expr": 60,
		"mulexpr": 17,
		"orexpr": 7,
		"primary": 11,
		"unary": 18,
	},
	58: {
		"addexpr": 2,
		"andexpr": 10,
		"cmpexpr": 16,
		"expr": 61,
		"mulexpr": 17,
		"orexpr": 7,
		"primary": 11,
		"unary": 18,
	},
}

// parserProd describes one production: its head symbol and body length.
type parserProd struct{ head string; bodyLen int }

var parserProds = []parserProd{
	0: {"program'", 1},
	1: {"program", 1},
	2: {"expr", 4},
	3: {"expr", 4},
	4: {"expr", 1},
	5: {"orexpr", 3},
	6: {"orexpr", 1},
	7: {"andexpr", 3},
	8: {"andexpr", 1},
	9: {"cmpexpr", 3},
	10: {"cmpexpr", 3},
	11: {"cmpexpr", 3},
	12: {"cmpexpr", 3},
	13: {"cmpexpr", 3},
	14: {"cmpexpr", 3},
	15: {"cmpexpr", 1},
	16: {"addexpr", 3},
	17: {"addexpr", 3},
	18: {"addexpr", 1},
	19: {"mulexpr", 3},
	20: {"mulexpr", 3},
	21: {"mulexpr", 3},
	22: {"mulexpr", 1},
	23: {"unary", 2},
	24: {"unary", 2},
	25: {"unary", 1},
	26: {"primary", 1},
	27: {"primary", 1},
	28: {"primary", 1},
	29: {"primary", 1},
	30: {"primary", 3},
	31: {"primary", 4},
	32: {"args", 1},
	33: {"args", 0},
	34: {"arglist", 1},
	35: {"arglist", 3},
}

var parserIgnore = map[int]bool{}

// Parse runs the SLR(1) parse loop over the token stream produced by lexer l.
// It returns nil on a successful parse, or a descriptive error on failure.
// Tokens whose IDs appear in parserIgnore are silently skipped.
func Parse(l *Lexer) error {
	// State stack — start in state 0.
	stk := []int{0}
	peek := func() int { return stk[len(stk)-1] }

	// Fetch the first non-ignored token.
	var cur Lexeme

	symbolTable := map[string][]Lexeme{}

	nextToken := func() {
		for {
			cur = l.NextToken()

			if cur.Token != EOF && cur.Token != ERROR {
				symbolTable[cur.Value] = append(
					symbolTable[cur.Value],
					cur,
				)
			}

			if !parserIgnore[cur.Token] {
				break
			}
		}
	}
	nextToken()

	// Map token ID → terminal name for table look-ups.
	tokName := tokenIDToName()

	for {
		state := peek()

		sym := "$"
		if cur.Token != EOF {
			if name, ok := tokName[cur.Token]; ok {
				sym = name
			} else {
				return fmt.Errorf(
					"syntax error: unexpected %d at line %d, col %d",
					cur.Token,
					cur.Line,
					cur.Col,
				)
			}
		}

		row, ok := parserActionTable[state]
		if !ok {
			return fmt.Errorf(
				"line %d col %d: no actions in state %d (token %q)",
				cur.Line,
				cur.Col,
				state,
				sym,
			)
		}

		act, ok := row[sym]
		if !ok {
			return fmt.Errorf(
				"syntax error: unexpected %s at line %d, col %d",
				sym,
				cur.Line,
				cur.Col,
			)
		}

		switch act.kind {

		case 1: // shift
			stk = append(stk, act.arg)
			nextToken()

		case 2: // reduce
			prod := parserProds[act.arg]

			// Pop |body| states off the stack.
			stk = stk[:len(stk)-prod.bodyLen]

			// Look up Goto[top][head] to find the new state.
			top := peek()

			gotoRow, ok := parserGotoTable[top]
			if !ok {
				return fmt.Errorf(
					"state %d: no goto row (reducing by %q)",
					top,
					prod.head,
				)
			}

			next, ok := gotoRow[prod.head]
			if !ok {
				return fmt.Errorf(
					"state %d: no goto for %q",
					top,
					prod.head,
				)
			}

			stk = append(stk, next)

		case 3: // accept
			fmt.Println("\n── Tabla de símbolos ──")
			fmt.Printf("%-20s %-10s %s\n", "LEXEMA", "TOKEN", "LÍNEA:COL")

			for lexeme, occs := range symbolTable {
				for _, occ := range occs {
					fmt.Printf(
						"%-20s %-10d %d:%d\n",
						lexeme,
						occ.Token,
						occ.Line,
						occ.Col,
					)
				}
			}

			return nil

		default:
			return fmt.Errorf(
				"line %d col %d: parse error (state %d, token %q)",
				cur.Line,
				cur.Col,
				state,
				sym,
			)
		}
	}
}

// tokenIDToName returns a map from token integer ID to its grammar name.
// This is the inverse of the constants emitted by GenerateCombined.
func tokenIDToName() map[int]string {
	return map[int]string{
		1: "FLOAT",
		2: "INT",
		3: "HEX",
		4: "PLUS",
		5: "MINUS",
		6: "TIMES",
		7: "DIV",
		8: "MOD",
		9: "EQ",
		10: "NEQ",
		11: "LT",
		12: "LEQ",
		13: "GT",
		14: "GEQ",
		15: "AND",
		16: "OR",
		17: "NOT",
		18: "LPAREN",
		19: "RPAREN",
		20: "COMMA",
		21: "ASSN",
		22: "KWLET",
		23: "KWIF",
		24: "KWELSE",
		25: "IDENT",
	}
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
	if err := Parse(l); err != nil {
		fmt.Fprintln(os.Stderr, "parse error:", err)
		os.Exit(1)
	}
	fmt.Println("OK — input accepted by the grammar.")
}
