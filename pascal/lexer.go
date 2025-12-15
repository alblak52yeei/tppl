package main

import (
	"fmt"
	"unicode"
	"unicode/utf8"
)

// TokenType представляет тип токена
type TokenType int

const (
	TokenEOF TokenType = iota
	TokenBEGIN
	TokenEND
	TokenDOT
	TokenSEMICOLON
	TokenASSIGN
	TokenPLUS
	TokenMINUS
	TokenMULTIPLY
	TokenDIVIDE
	TokenLPAREN
	TokenRPAREN
	TokenIDENTIFIER
	TokenNUMBER
)

// Token представляет токен с типом и значением
type Token struct {
	Type  TokenType
	Value string
	Pos   int
}

// Lexer представляет лексер для Pascal
type Lexer struct {
	input  string
	pos    int
	start  int
	tokens []Token
}

// NewLexer создает новый лексер
func NewLexer(input string) *Lexer {
	return &Lexer{
		input:  input,
		pos:    0,
		start:  0,
		tokens: []Token{},
	}
}

// Tokenize разбивает входную строку на токены
func (l *Lexer) Tokenize() ([]Token, error) {
	for l.pos < len(l.input) {
		l.skipWhitespace()
		if l.pos >= len(l.input) {
			break
		}

		r, size := l.peekRune()
		if size == 0 {
			break
		}

		switch {
		case r == '.':
			l.emit(TokenDOT)
			l.advance()
		case r == ';':
			l.emit(TokenSEMICOLON)
			l.advance()
		case r == ':':
			l.advance() // пропускаем ':'
			// Пропускаем пробелы между ':' и '='
			for l.pos < len(l.input) {
				nextR, _ := l.peekRune()
				if !unicode.IsSpace(nextR) {
					break
				}
				l.advance()
			}
			if l.pos < len(l.input) && l.input[l.pos] == '=' {
				l.advance() // пропускаем '='
				l.emit(TokenASSIGN)
			} else {
				return nil, fmt.Errorf("ожидался '=' после ':' на позиции %d", l.pos)
			}
		case r == '+':
			l.emit(TokenPLUS)
			l.advance()
		case r == '-':
			l.emit(TokenMINUS)
			l.advance()
		case r == '*':
			l.emit(TokenMULTIPLY)
			l.advance()
		case r == '/':
			l.emit(TokenDIVIDE)
			l.advance()
		case r == '(':
			l.emit(TokenLPAREN)
			l.advance()
		case r == ')':
			l.emit(TokenRPAREN)
			l.advance()
		case unicode.IsDigit(r):
			l.readNumber()
		case unicode.IsLetter(r):
			l.readIdentifier()
		default:
			return nil, fmt.Errorf("неожиданный символ '%c' на позиции %d", r, l.pos)
		}
	}

	l.emit(TokenEOF)
	return l.tokens, nil
}

func (l *Lexer) skipWhitespace() {
	for l.pos < len(l.input) {
		r, size := l.peekRune()
		if size == 0 {
			break
		}
		if !unicode.IsSpace(r) {
			break
		}
		l.advance()
	}
	l.start = l.pos
}

func (l *Lexer) readNumber() {
	for l.pos < len(l.input) {
		r, size := l.peekRune()
		if size == 0 || !unicode.IsDigit(r) {
			break
		}
		l.advance()
	}
	l.emit(TokenNUMBER)
}

func (l *Lexer) readIdentifier() {
	for l.pos < len(l.input) {
		r, size := l.peekRune()
		if size == 0 || (!unicode.IsLetter(r) && !unicode.IsDigit(r)) {
			break
		}
		l.advance()
	}
	value := l.input[l.start:l.pos]
	
	// Проверяем ключевые слова
	switch value {
	case "BEGIN":
		l.emit(TokenBEGIN)
	case "END":
		l.emit(TokenEND)
	default:
		l.emit(TokenIDENTIFIER)
	}
}

func (l *Lexer) peekRune() (rune, int) {
	if l.pos >= len(l.input) {
		return 0, 0
	}
	return utf8.DecodeRuneInString(l.input[l.pos:])
}

func (l *Lexer) peekNext() rune {
	if l.pos+1 >= len(l.input) {
		return 0
	}
	r, _ := utf8.DecodeRuneInString(l.input[l.pos+1:])
	return r
}

func (l *Lexer) advance() {
	if l.pos < len(l.input) {
		_, size := utf8.DecodeRuneInString(l.input[l.pos:])
		l.pos += size
	}
}

func (l *Lexer) emit(t TokenType) {
	value := l.input[l.start:l.pos]
	l.tokens = append(l.tokens, Token{
		Type:  t,
		Value: value,
		Pos:   l.start,
	})
	l.start = l.pos
}

