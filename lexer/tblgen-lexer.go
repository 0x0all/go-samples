package main

import (
	"fmt"
	"unicode/utf8"
)

type TokenName int

// Values for TokenName
const (
	// Special tokens
	ERROR TokenName = iota
	EOF

	COMMENT
	IDENTIFIER
	NUMBER
	QUOTE

	// Operators
	PLUS
	MINUS
	MULTIPLY
	PERIOD
	BACKSLASH
	COLON
	PERCENT
	PIPE
	EXCLAMATION
	QUESTION
	POUND
	AMPERSAND
	SEMI
	COMMA
	L_PAREN
	R_PAREN
	L_ANG
	R_ANG
	L_BRACE
	R_BRACE
	L_BRACKET
	R_BRACKET
	EQUALS
)

var tokenNames = [...]string{
	ERROR:       "ERROR",
	EOF:         "EOF",
	COMMENT:     "COMMENT",
	IDENTIFIER:  "IDENTIFIER",
	NUMBER:      "NUMBER",
	QUOTE:       "QUOTE",
	PLUS:        "PLUS",
	MINUS:       "MINUS",
	MULTIPLY:    "MULTIPLY",
	PERIOD:      "PERIOD",
	BACKSLASH:   "BACKSLASH",
	COLON:       "COLON",
	PERCENT:     "PERCENT",
	PIPE:        "PIPE",
	EXCLAMATION: "EXCLAMATION",
	QUESTION:    "QUESTION",
	POUND:       "POUND",
	AMPERSAND:   "AMPERSAND",
	SEMI:        "SEMI",
	COMMA:       "COMMA",
	L_PAREN:     "L_PAREN",
	R_PAREN:     "R_PAREN",
	L_ANG:       "L_ANG",
	R_ANG:       "R_ANG",
	L_BRACE:     "L_BRACE",
	R_BRACE:     "R_BRACE",
	L_BRACKET:   "L_BRACKET",
	R_BRACKET:   "R_BRACKET",
	EQUALS:      "EQUALS",
}

type Token struct {
	Name TokenName
	Val  string
	Pos  int
}

func (tok Token) String() string {
	return fmt.Sprintf("Token{%s, '%s', %d}", tokenNames[tok.Name], tok.Val, tok.Pos)
}

func makeErrorToken(pos int) Token {
	return Token{ERROR, "", pos}
}

var opTable = [...]TokenName{
	'+':  PLUS,
	'-':  MINUS,
	'*':  MULTIPLY,
	'.':  PERIOD,
	'\\': BACKSLASH,
	':':  COLON,
	'%':  PERCENT,
	'|':  PIPE,
	'!':  EXCLAMATION,
	'?':  QUESTION,
	'#':  POUND,
	'&':  AMPERSAND,
	';':  SEMI,
	',':  COMMA,
	'(':  L_PAREN,
	')':  R_PAREN,
	'<':  L_ANG,
	'>':  R_ANG,
	'{':  L_BRACE,
	'}':  R_BRACE,
	'[':  L_BRACKET,
	']':  R_BRACKET,
	'=':  EQUALS,
}

type Lexer struct {
	buf []byte

	// Current rune.
	r rune

	// Position of the current rune in buf.
	rpos int

	// Position of the next rune in buf.
	nextpos int
}

func NewLexer(buf []byte) *Lexer {
	lex := Lexer{buf, -1, 0, 0}

	// Prime the lexer by calling .next
	lex.next()
	return &lex
}

func (lex *Lexer) NextToken() Token {
	lex.skipNontokens()

	if lex.r < 0 {
		return Token{EOF, "", lex.nextpos}
	}

	if int(lex.r) < len(opTable); opName := opTable[lex.r]; opName != ERROR {
			startpos := lex.rpos
			lex.next()
			return Token{opName, string(lex.buf[startpos:lex.rpos]), startpos}
		}
	} else if isAlpha(lex.r) {
		return lex.scanIdentifier()
	} else if isDigit(lex.r) {
		return lex.scanNumber()
	} else if lex.r == '"' {
		return lex.scanQuote()
	}

	return makeErrorToken(lex.rpos)
}

func (lex *Lexer) next() {
	if lex.nextpos < len(lex.buf) {
		lex.rpos = lex.nextpos

		// r is the current rune, w is its width. We start by assuming the
		// common case - that the current rune is ASCII (and thus has width=1).
		r, w := rune(lex.buf[lex.nextpos]), 1

		if r > utf8.RuneSelf {
			// The current rune is not actually ASCII, so we have to decode it
			// properly.
			r, w = utf8.DecodeRune(lex.buf[lex.nextpos:])
		}

		lex.nextpos += w
		lex.r = r
	} else {
		lex.rpos = len(lex.buf)
		lex.r = -1 // EOF
	}
}

func (lex *Lexer) skipNontokens() {
	for lex.r == ' ' || lex.r == '\t' || lex.r == '\n' || lex.r == '\r' {
		lex.next()
	}
}

func (lex *Lexer) scanIdentifier() Token {
	startpos := lex.rpos
	for isAlpha(lex.r) || isDigit(lex.r) {
		lex.next()
	}
	return Token{IDENTIFIER, string(lex.buf[startpos:lex.rpos]), startpos}
}

func (lex *Lexer) scanNumber() Token {
	startpos := lex.rpos
	for isDigit(lex.r) {
		lex.next()
	}
	return Token{NUMBER, string(lex.buf[startpos:lex.rpos]), startpos}
}

func (lex *Lexer) scanQuote() Token {
	startpos := lex.rpos
	lex.next()
	for lex.r != '"' {
		lex.next()
	}

	if lex.r < 0 {
		return makeErrorToken(startpos)
	} else {
		lex.next()
		return Token{QUOTE, string(lex.buf[startpos:lex.rpos]), startpos}
	}
}

func isAlpha(r rune) bool {
	return 'a' <= r && r <= 'z' || 'A' <= r && r <= 'Z' || r == '_'
}

func isDigit(r rune) bool {
	return '0' <= r && r <= '9'
}

//------------------------------------------------------------------------------

func main() {
	const sample = `foo
3456 baz "本ä" 3 `
	fmt.Println(sample)

	nl := NewLexer([]byte(sample))
	fmt.Println(nl)

	for {
		nt := nl.NextToken()
		fmt.Println(nt)
		if nt.Name == EOF {
			break
		}
	}
}
