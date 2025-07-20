package core

import (
	"fmt"
)

type TokenType int

type Token struct {
	Type     TokenType
	Value    string
	StartPos int
	EndPos   int
}

type Lexer struct {
	input       string
	postion     int // 当前位置
	readPostion int // 下一个读取的位置
	ch          byte
}

const (
	EOF TokenType = iota
	ROOT
	RECURSIVE_DESCENT
	IDENTIFIER
	WILDCARD
	STRING
	NUMBER
	COLON
	COMMA
	DOT
	LBRACKET
	RBRACKET
)

func (l *Lexer) readChar() {
	if l.readPostion >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPostion]
	}
	l.postion = l.readPostion
	l.readPostion++
}

// 预读取下一个字符
func (l *Lexer) peekChar() byte {
	if l.readPostion >= len(l.input) {
		return 0
	} else {
		return l.input[l.readPostion]
	}
}

func (l *Lexer) readIdentifier() (string, error) {
	startPos := l.postion
	for {
		if !isLetter(l.ch) && !isDigit(l.ch) && l.ch != '_' {
			return "", fmt.Errorf("%c: invalid identifier", l.ch)
		}
		preCh := l.peekChar()
		if preCh != 0 && preCh != '.' && preCh != ',' && preCh != ':' && preCh != '[' && preCh != ']' && preCh != '*' && preCh != '$' && preCh != ' ' {
			l.readChar()
			continue
		} else {
			return l.input[startPos:l.readPostion], nil
		}
	}
}

func (l *Lexer) NextToken() (*Token, error) {
	tok := &Token{}
	l.readChar()
	startPos := l.postion
	switch l.ch {
	case 0:
		tok.Type = EOF
		tok.Value = ""
	case '$':
		tok.Type = ROOT
		tok.Value = "$"
	case '.':
		if l.peekChar() == '.' {
			l.readChar()
			tok.Type = RECURSIVE_DESCENT
			tok.Value = ".."
		} else {
			tok.Type = DOT
			tok.Value = "."
		}
	case '*':
		tok.Type = WILDCARD
		tok.Value = "*"
	case '[':
		tok.Type = LBRACKET
		tok.Value = "["
	case ']':
		tok.Type = RBRACKET
		tok.Value = "]"
	case ',':
		tok.Type = COMMA
		tok.Value = ","
	case ':':
		tok.Type = COLON
		tok.Value = ":"
	default:
		value, err := l.readIdentifier()
		if err != nil {
			return nil, err
		}
		tok.Type = IDENTIFIER
		tok.Value = value
	}
	tok.StartPos = startPos
	tok.EndPos = startPos + len(tok.Value)
	return tok, nil
}

func isLetter(ch byte) bool {
	return ('a' <= ch && ch <= 'z') || ('A' <= ch && ch <= 'Z')
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}
