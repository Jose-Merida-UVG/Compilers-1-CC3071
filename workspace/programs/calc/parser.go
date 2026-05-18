package main

import (
	"fmt"
	"os"
	"strings"
)

// parserAction encodes one cell of the SLR(1) ACTION table.
// kind: 1=shift 2=reduce 3=accept; arg: target state (shift) or prod index (reduce).
type parserAction struct{ kind, arg int }

// parserActionTable[state][terminal] → action
var parserActionTable = map[int]map[string]parserAction{
	0: {
		"FLOAT": {1, 3},
		"HEX": {1, 4},
		"IDENT": {1, 10},
		"INT": {1, 2},
		"KWIF": {1, 7},
		"KWLET": {1, 6},
		"LPAREN": {1, 5},
		"MINUS": {1, 18},
		"NOT": {1, 16},
	},
	1: {
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
	2: {
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
	3: {
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
	4: {
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
	5: {
		"FLOAT": {1, 3},
		"HEX": {1, 4},
		"IDENT": {1, 10},
		"INT": {1, 2},
		"KWIF": {1, 7},
		"KWLET": {1, 6},
		"LPAREN": {1, 5},
		"MINUS": {1, 18},
		"NOT": {1, 16},
	},
	6: {
		"IDENT": {1, 20},
	},
	7: {
		"FLOAT": {1, 3},
		"HEX": {1, 4},
		"IDENT": {1, 10},
		"INT": {1, 2},
		"LPAREN": {1, 5},
		"MINUS": {1, 18},
		"NOT": {1, 16},
	},
	8: {
		"$": {2, 8},
		"AND": {2, 8},
		"COMMA": {2, 8},
		"KWELSE": {2, 8},
		"OR": {2, 8},
		"RPAREN": {2, 8},
	},
	9: {
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
	10: {
		"$": {2, 29},
		"AND": {2, 29},
		"COMMA": {2, 29},
		"DIV": {2, 29},
		"EQ": {2, 29},
		"GEQ": {2, 29},
		"GT": {2, 29},
		"KWELSE": {2, 29},
		"LEQ": {2, 29},
		"LPAREN": {1, 22},
		"LT": {2, 29},
		"MINUS": {2, 29},
		"MOD": {2, 29},
		"NEQ": {2, 29},
		"OR": {2, 29},
		"PLUS": {2, 29},
		"RPAREN": {2, 29},
		"TIMES": {2, 29},
	},
	11: {
		"$": {2, 1},
	},
	12: {
		"$": {2, 4},
		"COMMA": {2, 4},
		"OR": {1, 23},
		"RPAREN": {2, 4},
	},
	13: {
		"$": {2, 6},
		"AND": {1, 24},
		"COMMA": {2, 6},
		"KWELSE": {2, 6},
		"OR": {2, 6},
		"RPAREN": {2, 6},
	},
	14: {
		"$": {2, 15},
		"AND": {2, 15},
		"COMMA": {2, 15},
		"EQ": {1, 31},
		"GEQ": {1, 28},
		"GT": {1, 27},
		"KWELSE": {2, 15},
		"LEQ": {1, 26},
		"LT": {1, 25},
		"MINUS": {1, 30},
		"NEQ": {1, 32},
		"OR": {2, 15},
		"PLUS": {1, 29},
		"RPAREN": {2, 15},
	},
	15: {
		"$": {2, 18},
		"AND": {2, 18},
		"COMMA": {2, 18},
		"DIV": {1, 34},
		"EQ": {2, 18},
		"GEQ": {2, 18},
		"GT": {2, 18},
		"KWELSE": {2, 18},
		"LEQ": {2, 18},
		"LT": {2, 18},
		"MINUS": {2, 18},
		"MOD": {1, 35},
		"NEQ": {2, 18},
		"OR": {2, 18},
		"PLUS": {2, 18},
		"RPAREN": {2, 18},
		"TIMES": {1, 33},
	},
	16: {
		"FLOAT": {1, 3},
		"HEX": {1, 4},
		"IDENT": {1, 10},
		"INT": {1, 2},
		"LPAREN": {1, 5},
		"MINUS": {1, 18},
		"NOT": {1, 16},
	},
	17: {
		"$": {3, 0},
	},
	18: {
		"FLOAT": {1, 3},
		"HEX": {1, 4},
		"IDENT": {1, 10},
		"INT": {1, 2},
		"LPAREN": {1, 5},
		"MINUS": {1, 18},
		"NOT": {1, 16},
	},
	19: {
		"RPAREN": {1, 38},
	},
	20: {
		"ASSN": {1, 39},
	},
	21: {
		"KWELSE": {1, 40},
		"OR": {1, 23},
	},
	22: {
		"FLOAT": {1, 3},
		"HEX": {1, 4},
		"IDENT": {1, 10},
		"INT": {1, 2},
		"KWIF": {1, 7},
		"KWLET": {1, 6},
		"LPAREN": {1, 5},
		"MINUS": {1, 18},
		"NOT": {1, 16},
		"RPAREN": {2, 33},
	},
	23: {
		"FLOAT": {1, 3},
		"HEX": {1, 4},
		"IDENT": {1, 10},
		"INT": {1, 2},
		"LPAREN": {1, 5},
		"MINUS": {1, 18},
		"NOT": {1, 16},
	},
	24: {
		"FLOAT": {1, 3},
		"HEX": {1, 4},
		"IDENT": {1, 10},
		"INT": {1, 2},
		"LPAREN": {1, 5},
		"MINUS": {1, 18},
		"NOT": {1, 16},
	},
	25: {
		"FLOAT": {1, 3},
		"HEX": {1, 4},
		"IDENT": {1, 10},
		"INT": {1, 2},
		"LPAREN": {1, 5},
		"MINUS": {1, 18},
		"NOT": {1, 16},
	},
	26: {
		"FLOAT": {1, 3},
		"HEX": {1, 4},
		"IDENT": {1, 10},
		"INT": {1, 2},
		"LPAREN": {1, 5},
		"MINUS": {1, 18},
		"NOT": {1, 16},
	},
	27: {
		"FLOAT": {1, 3},
		"HEX": {1, 4},
		"IDENT": {1, 10},
		"INT": {1, 2},
		"LPAREN": {1, 5},
		"MINUS": {1, 18},
		"NOT": {1, 16},
	},
	28: {
		"FLOAT": {1, 3},
		"HEX": {1, 4},
		"IDENT": {1, 10},
		"INT": {1, 2},
		"LPAREN": {1, 5},
		"MINUS": {1, 18},
		"NOT": {1, 16},
	},
	29: {
		"FLOAT": {1, 3},
		"HEX": {1, 4},
		"IDENT": {1, 10},
		"INT": {1, 2},
		"LPAREN": {1, 5},
		"MINUS": {1, 18},
		"NOT": {1, 16},
	},
	30: {
		"FLOAT": {1, 3},
		"HEX": {1, 4},
		"IDENT": {1, 10},
		"INT": {1, 2},
		"LPAREN": {1, 5},
		"MINUS": {1, 18},
		"NOT": {1, 16},
	},
	31: {
		"FLOAT": {1, 3},
		"HEX": {1, 4},
		"IDENT": {1, 10},
		"INT": {1, 2},
		"LPAREN": {1, 5},
		"MINUS": {1, 18},
		"NOT": {1, 16},
	},
	32: {
		"FLOAT": {1, 3},
		"HEX": {1, 4},
		"IDENT": {1, 10},
		"INT": {1, 2},
		"LPAREN": {1, 5},
		"MINUS": {1, 18},
		"NOT": {1, 16},
	},
	33: {
		"FLOAT": {1, 3},
		"HEX": {1, 4},
		"IDENT": {1, 10},
		"INT": {1, 2},
		"LPAREN": {1, 5},
		"MINUS": {1, 18},
		"NOT": {1, 16},
	},
	34: {
		"FLOAT": {1, 3},
		"HEX": {1, 4},
		"IDENT": {1, 10},
		"INT": {1, 2},
		"LPAREN": {1, 5},
		"MINUS": {1, 18},
		"NOT": {1, 16},
	},
	35: {
		"FLOAT": {1, 3},
		"HEX": {1, 4},
		"IDENT": {1, 10},
		"INT": {1, 2},
		"LPAREN": {1, 5},
		"MINUS": {1, 18},
		"NOT": {1, 16},
	},
	36: {
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
	37: {
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
	38: {
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
	39: {
		"FLOAT": {1, 3},
		"HEX": {1, 4},
		"IDENT": {1, 10},
		"INT": {1, 2},
		"KWIF": {1, 7},
		"KWLET": {1, 6},
		"LPAREN": {1, 5},
		"MINUS": {1, 18},
		"NOT": {1, 16},
	},
	40: {
		"FLOAT": {1, 3},
		"HEX": {1, 4},
		"IDENT": {1, 10},
		"INT": {1, 2},
		"LPAREN": {1, 5},
		"MINUS": {1, 18},
		"NOT": {1, 16},
	},
	41: {
		"RPAREN": {1, 59},
	},
	42: {
		"COMMA": {1, 60},
		"RPAREN": {2, 32},
	},
	43: {
		"COMMA": {2, 34},
		"RPAREN": {2, 34},
	},
	44: {
		"$": {2, 5},
		"AND": {1, 24},
		"COMMA": {2, 5},
		"KWELSE": {2, 5},
		"OR": {2, 5},
		"RPAREN": {2, 5},
	},
	45: {
		"$": {2, 7},
		"AND": {2, 7},
		"COMMA": {2, 7},
		"KWELSE": {2, 7},
		"OR": {2, 7},
		"RPAREN": {2, 7},
	},
	46: {
		"$": {2, 11},
		"AND": {2, 11},
		"COMMA": {2, 11},
		"KWELSE": {2, 11},
		"MINUS": {1, 30},
		"OR": {2, 11},
		"PLUS": {1, 29},
		"RPAREN": {2, 11},
	},
	47: {
		"$": {2, 12},
		"AND": {2, 12},
		"COMMA": {2, 12},
		"KWELSE": {2, 12},
		"MINUS": {1, 30},
		"OR": {2, 12},
		"PLUS": {1, 29},
		"RPAREN": {2, 12},
	},
	48: {
		"$": {2, 13},
		"AND": {2, 13},
		"COMMA": {2, 13},
		"KWELSE": {2, 13},
		"MINUS": {1, 30},
		"OR": {2, 13},
		"PLUS": {1, 29},
		"RPAREN": {2, 13},
	},
	49: {
		"$": {2, 14},
		"AND": {2, 14},
		"COMMA": {2, 14},
		"KWELSE": {2, 14},
		"MINUS": {1, 30},
		"OR": {2, 14},
		"PLUS": {1, 29},
		"RPAREN": {2, 14},
	},
	50: {
		"$": {2, 16},
		"AND": {2, 16},
		"COMMA": {2, 16},
		"DIV": {1, 34},
		"EQ": {2, 16},
		"GEQ": {2, 16},
		"GT": {2, 16},
		"KWELSE": {2, 16},
		"LEQ": {2, 16},
		"LT": {2, 16},
		"MINUS": {2, 16},
		"MOD": {1, 35},
		"NEQ": {2, 16},
		"OR": {2, 16},
		"PLUS": {2, 16},
		"RPAREN": {2, 16},
		"TIMES": {1, 33},
	},
	51: {
		"$": {2, 17},
		"AND": {2, 17},
		"COMMA": {2, 17},
		"DIV": {1, 34},
		"EQ": {2, 17},
		"GEQ": {2, 17},
		"GT": {2, 17},
		"KWELSE": {2, 17},
		"LEQ": {2, 17},
		"LT": {2, 17},
		"MINUS": {2, 17},
		"MOD": {1, 35},
		"NEQ": {2, 17},
		"OR": {2, 17},
		"PLUS": {2, 17},
		"RPAREN": {2, 17},
		"TIMES": {1, 33},
	},
	52: {
		"$": {2, 9},
		"AND": {2, 9},
		"COMMA": {2, 9},
		"KWELSE": {2, 9},
		"MINUS": {1, 30},
		"OR": {2, 9},
		"PLUS": {1, 29},
		"RPAREN": {2, 9},
	},
	53: {
		"$": {2, 10},
		"AND": {2, 10},
		"COMMA": {2, 10},
		"KWELSE": {2, 10},
		"MINUS": {1, 30},
		"OR": {2, 10},
		"PLUS": {1, 29},
		"RPAREN": {2, 10},
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
		"$": {2, 2},
		"COMMA": {2, 2},
		"RPAREN": {2, 2},
	},
	58: {
		"$": {2, 3},
		"COMMA": {2, 3},
		"OR": {1, 23},
		"RPAREN": {2, 3},
	},
	59: {
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
	60: {
		"FLOAT": {1, 3},
		"HEX": {1, 4},
		"IDENT": {1, 10},
		"INT": {1, 2},
		"KWIF": {1, 7},
		"KWLET": {1, 6},
		"LPAREN": {1, 5},
		"MINUS": {1, 18},
		"NOT": {1, 16},
	},
	61: {
		"COMMA": {2, 35},
		"RPAREN": {2, 35},
	},
}

// parserGotoTable[state][nonTerminal] → next state
var parserGotoTable = map[int]map[string]int{
	0: {
		"addexpr": 14,
		"andexpr": 13,
		"cmpexpr": 8,
		"expr": 11,
		"mulexpr": 15,
		"orexpr": 12,
		"primary": 1,
		"program": 17,
		"unary": 9,
	},
	5: {
		"addexpr": 14,
		"andexpr": 13,
		"cmpexpr": 8,
		"expr": 19,
		"mulexpr": 15,
		"orexpr": 12,
		"primary": 1,
		"unary": 9,
	},
	7: {
		"addexpr": 14,
		"andexpr": 13,
		"cmpexpr": 8,
		"mulexpr": 15,
		"orexpr": 21,
		"primary": 1,
		"unary": 9,
	},
	16: {
		"primary": 1,
		"unary": 36,
	},
	18: {
		"primary": 1,
		"unary": 37,
	},
	22: {
		"addexpr": 14,
		"andexpr": 13,
		"arglist": 42,
		"args": 41,
		"cmpexpr": 8,
		"expr": 43,
		"mulexpr": 15,
		"orexpr": 12,
		"primary": 1,
		"unary": 9,
	},
	23: {
		"addexpr": 14,
		"andexpr": 44,
		"cmpexpr": 8,
		"mulexpr": 15,
		"primary": 1,
		"unary": 9,
	},
	24: {
		"addexpr": 14,
		"cmpexpr": 45,
		"mulexpr": 15,
		"primary": 1,
		"unary": 9,
	},
	25: {
		"addexpr": 46,
		"mulexpr": 15,
		"primary": 1,
		"unary": 9,
	},
	26: {
		"addexpr": 47,
		"mulexpr": 15,
		"primary": 1,
		"unary": 9,
	},
	27: {
		"addexpr": 48,
		"mulexpr": 15,
		"primary": 1,
		"unary": 9,
	},
	28: {
		"addexpr": 49,
		"mulexpr": 15,
		"primary": 1,
		"unary": 9,
	},
	29: {
		"mulexpr": 50,
		"primary": 1,
		"unary": 9,
	},
	30: {
		"mulexpr": 51,
		"primary": 1,
		"unary": 9,
	},
	31: {
		"addexpr": 52,
		"mulexpr": 15,
		"primary": 1,
		"unary": 9,
	},
	32: {
		"addexpr": 53,
		"mulexpr": 15,
		"primary": 1,
		"unary": 9,
	},
	33: {
		"primary": 1,
		"unary": 54,
	},
	34: {
		"primary": 1,
		"unary": 55,
	},
	35: {
		"primary": 1,
		"unary": 56,
	},
	39: {
		"addexpr": 14,
		"andexpr": 13,
		"cmpexpr": 8,
		"expr": 57,
		"mulexpr": 15,
		"orexpr": 12,
		"primary": 1,
		"unary": 9,
	},
	40: {
		"addexpr": 14,
		"andexpr": 13,
		"cmpexpr": 8,
		"mulexpr": 15,
		"orexpr": 58,
		"primary": 1,
		"unary": 9,
	},
	60: {
		"addexpr": 14,
		"andexpr": 13,
		"cmpexpr": 8,
		"expr": 61,
		"mulexpr": 15,
		"orexpr": 12,
		"primary": 1,
		"unary": 9,
	},
}

// parserProd describes one production: its head symbol, body symbols, and body length.
type parserProd struct{ head string; body string; bodyLen int }

var parserProds = []parserProd{
	0: {"program'", "program", 1},
	1: {"program", "expr", 1},
	2: {"expr", "KWLET IDENT ASSN expr", 4},
	3: {"expr", "KWIF orexpr KWELSE orexpr", 4},
	4: {"expr", "orexpr", 1},
	5: {"orexpr", "orexpr OR andexpr", 3},
	6: {"orexpr", "andexpr", 1},
	7: {"andexpr", "andexpr AND cmpexpr", 3},
	8: {"andexpr", "cmpexpr", 1},
	9: {"cmpexpr", "addexpr EQ addexpr", 3},
	10: {"cmpexpr", "addexpr NEQ addexpr", 3},
	11: {"cmpexpr", "addexpr LT addexpr", 3},
	12: {"cmpexpr", "addexpr LEQ addexpr", 3},
	13: {"cmpexpr", "addexpr GT addexpr", 3},
	14: {"cmpexpr", "addexpr GEQ addexpr", 3},
	15: {"cmpexpr", "addexpr", 1},
	16: {"addexpr", "addexpr PLUS mulexpr", 3},
	17: {"addexpr", "addexpr MINUS mulexpr", 3},
	18: {"addexpr", "mulexpr", 1},
	19: {"mulexpr", "mulexpr TIMES unary", 3},
	20: {"mulexpr", "mulexpr DIV unary", 3},
	21: {"mulexpr", "mulexpr MOD unary", 3},
	22: {"mulexpr", "unary", 1},
	23: {"unary", "NOT unary", 2},
	24: {"unary", "MINUS unary", 2},
	25: {"unary", "primary", 1},
	26: {"primary", "INT", 1},
	27: {"primary", "FLOAT", 1},
	28: {"primary", "HEX", 1},
	29: {"primary", "IDENT", 1},
	30: {"primary", "LPAREN expr RPAREN", 3},
	31: {"primary", "IDENT LPAREN args RPAREN", 4},
	32: {"args", "arglist", 1},
	33: {"args", "ε", 0},
	34: {"arglist", "expr", 1},
	35: {"arglist", "arglist COMMA expr", 3},
}

var parserIgnore = map[int]bool{}

// Parse runs the SLR(1) parse loop over the token stream produced by lexer l.
// It returns nil on a successful parse, or a descriptive error on failure.
// Tokens whose IDs appear in parserIgnore are silently skipped.
func Parse(l *Lexer) error {
	// State stack — start in state 0.
	stk := []int{0}
	peek := func() int { return stk[len(stk)-1] }

	// Symbol stack — tracks the sentential form for derivation display.
	var symStk []string

	// Fetch the first non-ignored token.
	var cur Lexeme

	var symbolTable []Lexeme
	var sententialForms []string

	nextToken := func() {
		for {
			cur = l.NextToken()

			if cur.Token != EOF && cur.Token != ERROR {
				symbolTable = append(symbolTable, cur)
			}

			if !parserIgnore[cur.Token] {
				break
			}
		}
	}
	nextToken()

	// Map token ID → terminal name for table look-ups.
	tokName := tokenIDToName()

	fmt.Println("── Parse Actions ──")

	for {
		state := peek()

		// If we hit an ERROR token from the lexer, immediately fail.
		if cur.Token == ERROR {
			return fmt.Errorf(
				"lexical error: unrecognized input '%s' at line %d, col %d",
				cur.Value,
				cur.Line,
				cur.Col,
			)
		}

		sym := "$"
		if cur.Token != EOF {
			if name, ok := tokName[cur.Token]; ok {
				sym = name
			} else {
				return fmt.Errorf(
					"syntax error: unexpected token %d at line %d, col %d",
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
			symName := tokenName(cur, tokName)
			fmt.Printf("Shift:  %s\n", symName)
			symStk = append(symStk, symName)
			stk = append(stk, act.arg)
			nextToken()

		case 2: // reduce
			prod := parserProds[act.arg]
			bodyStr := prodBody(act.arg)
			fmt.Printf("Reduce: %s → %s\n", prod.head, bodyStr)

			// Pop body symbols and push head.
			if prod.bodyLen > 0 {
				symStk = symStk[:len(symStk)-prod.bodyLen]
			}
			symStk = append(symStk, prod.head)

			// Record the current sentential form.
			sententialForms = append(sententialForms, joinSymbols(symStk))

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
			fmt.Println("\n\n── Sentential Forms (Rightmost Derivation in Reverse) ──")
			for i, form := range sententialForms {
				fmt.Printf("%d. %s\n", i+1, form)
			}

			fmt.Println("\n── Symbol Table ──")
			fmt.Printf("%-20s %-20s %-10s\n", "LEXEME", "TOKEN", "LINE:COL")

			for _, lex := range symbolTable {
				tokenStr := tokenName(lex, tokName)
				fmt.Printf(
					"%-20s %-20s %d:%d\n",
					lex.Value,
					tokenStr,
					lex.Line,
					lex.Col,
				)
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

func tokenName(lex Lexeme, tokName map[int]string) string {
	if lex.Token == EOF {
		return "$"
	}
	if name, ok := tokName[lex.Token]; ok {
		return name
	}
	return fmt.Sprintf("TOKEN%d", lex.Token)
}

func prodBody(idx int) string {
	return parserProds[idx].body
}

func joinSymbols(syms []string) string {
	if len(syms) == 0 {
		return "ε"
	}
	return strings.Join(syms, " ")
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
		9: "POW",
		10: "EQ",
		11: "NEQ",
		12: "LT",
		13: "LEQ",
		14: "GT",
		15: "GEQ",
		16: "AND",
		17: "OR",
		18: "NOT",
		19: "LPAREN",
		20: "RPAREN",
		21: "COMMA",
		22: "ASSN",
		23: "KWLET",
		24: "KWIF",
		25: "KWELSE",
		26: "WHILE",
		27: "PLEQ",
		28: "MIEQ",
		29: "TIEQ",
		30: "DIEQ",
		31: "IDENT",
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
