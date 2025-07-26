package core

import (
	"errors"
	"fmt"
	"strings"
)

type JsonTokenType int

type JsonToken struct {
	Type     JsonTokenType
	Value    string
	StartPos int
	EndPos   int
}

type JsonLexer struct {
	input        string
	position     int // 当前位置
	readPosition int // 下一个读取的位置
	ch           byte
}

const (
	EOF JsonTokenType = iota
	Root
	RecursiveDescent
	Identifier
	Number
	Wildcard
	Colon
	Comma
	Dot
	Lbracket
	Rbracket
)

func (l *JsonLexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition++
}

// 预读取下一个字符
func (l *JsonLexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	} else {
		return l.input[l.readPosition]
	}
}

func (l *JsonLexer) readIdentifier() (string, error) {
	startPos := l.position
	for {
		if !isLetter(l.ch) && !isDigit(l.ch) && l.ch != '_' {
			return "", fmt.Errorf("%c: invalid identifier", l.ch)
		}
		preCh := l.peekChar()
		if preCh != 0 && preCh != '.' && preCh != ',' && preCh != ':' && preCh != '[' && preCh != ']' && preCh != '*' && preCh != '$' && preCh != ' ' {
			l.readChar()
			continue
		} else {
			return l.input[startPos:l.readPosition], nil
		}
	}
}

func isLetter(ch byte) bool {
	return ('a' <= ch && ch <= 'z') || ('A' <= ch && ch <= 'Z')
}

func isNumber(s string) bool {
	return strings.IndexFunc(s, func(r rune) bool {
		return !isDigit(byte(r))
	}) == -1
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func (l *JsonLexer) NextToken() (*JsonToken, error) {
	tok := &JsonToken{}
	l.readChar()
	startPos := l.position
	switch l.ch {
	case 0:
		tok.Type = EOF
		tok.Value = ""
	case '$':
		tok.Type = Root
		tok.Value = "$"
	case '.':
		if l.peekChar() == '.' {
			l.readChar()
			tok.Type = RecursiveDescent
			tok.Value = ".."
		} else {
			tok.Type = Dot
			tok.Value = "."
		}
	case '*':
		tok.Type = Wildcard
		tok.Value = "*"
	case '[':
		tok.Type = Lbracket
		tok.Value = "["
	case ']':
		tok.Type = Rbracket
		tok.Value = "]"
	case ',':
		tok.Type = Comma
		tok.Value = ","
	case ':':
		tok.Type = Colon
		tok.Value = ":"
	default:
		value, err := l.readIdentifier()
		if err != nil {
			return nil, err
		}
		if isNumber(value) {
			tok.Type = Number
		} else {
			tok.Type = Identifier
		}
		tok.Value = value
	}
	tok.StartPos = startPos
	tok.EndPos = startPos + len(tok.Value)
	return tok, nil
}

type JsonAstNodeType int

const (
	NodeRoot JsonAstNodeType = iota
	NodeIdentifier
	NodeDot
	NodeRecursiveDescent
	NodeWildcard
	NodeArrayIndex
	NodeArrayWildcard
)

type JsonAstNode struct {
	Type  JsonAstNodeType
	Value string
	Child *JsonAstNode
}

// 为 JsonAstNode 添加 String 方法用于直观打印
func (n *JsonAstNode) String() string {
	if n == nil {
		return ""
	}

	// 构建当前节点的字符串表示
	nodeStr := fmt.Sprintf("%s(%s)", n.getTypeName(), n.Value)

	// 如果没有子节点，则直接返回当前节点信息
	if n.Child == nil {
		return nodeStr
	}

	// 递归处理子节点
	return fmt.Sprintf("%s -> %s", nodeStr, n.Child.String())
}

// getTypeName 返回类型的可读名称
func (n *JsonAstNode) getTypeName() string {
	switch n.Type {
	case NodeRoot:
		return "Root"
	case NodeIdentifier:
		return "Identifier"
	case NodeDot:
		return "Dot"
	case NodeRecursiveDescent:
		return "RecursiveDescent"
	case NodeWildcard:
		return "Wildcard"
	case NodeArrayIndex:
		return "ArrayIndex"
	case NodeArrayWildcard:
		return "ArrayWildcard"
	default:
		return "Unknown"
	}
}

type JsonParser struct {
	lexer        *JsonLexer
	currentToken *JsonToken
}

func (jp *JsonParser) nextToken() error {
	if tok, err := jp.lexer.NextToken(); err != nil {
		return err
	} else {
		jp.currentToken = tok
	}
	return nil
}

func (jp *JsonParser) parseRest() (*JsonAstNode, error) {
	var node *JsonAstNode

	switch jp.currentToken.Type {
	case Dot:
		node = &JsonAstNode{Type: NodeDot, Value: jp.currentToken.Value, Child: &JsonAstNode{}}
		if err := jp.nextToken(); err != nil {
			return nil, err
		}
		if jp.currentToken.Type != Identifier && jp.currentToken.Type != Wildcard {
			return nil, fmt.Errorf("expected Identifier or Wildcard after Dot, got %v", jp.currentToken.Type)
		}
	case Identifier:
		node = &JsonAstNode{
			Type:  NodeIdentifier,
			Value: jp.currentToken.Value,
		}
		if err := jp.nextToken(); err != nil {
			return nil, err
		}
		if jp.currentToken.Type != Lbracket && jp.currentToken.Type != Dot && jp.currentToken.Type != RecursiveDescent && jp.currentToken.Type != EOF {
			return nil, fmt.Errorf("expected Lbracket or Dot or RecursiveDescent or EOF after Identifier, got %v", jp.currentToken.Type)
		}
	case RecursiveDescent:
		node = &JsonAstNode{
			Type:  NodeRecursiveDescent,
			Value: "..",
		}
		if err := jp.nextToken(); err != nil {
			return nil, err
		}
		if jp.currentToken.Type != Identifier && jp.currentToken.Type != Wildcard {
			return nil, fmt.Errorf("expected Identifier or Wildcard after RecursiveDescent, got %v", jp.currentToken.Type)
		}
	case Lbracket:
		if err := jp.nextToken(); err != nil {
			return nil, err
		}
		if jp.currentToken.Type == Wildcard {
			if err := jp.nextToken(); err != nil {
				return nil, err
			}
			if jp.currentToken.Type != Rbracket {
				return nil, fmt.Errorf("expected Rbracket after Lbracket and Wildcard, got %v", jp.currentToken.Type)
			}
			node = &JsonAstNode{
				Type:  NodeArrayWildcard,
				Value: "[*]",
			}
		} else if jp.currentToken.Type == Number {
			numValue := jp.currentToken.Value
			if err := jp.nextToken(); err != nil {
				return nil, err
			}
			if jp.currentToken.Type != Rbracket {
				return nil, fmt.Errorf("expected Rbracket after Lbracket and Number, got %v", jp.currentToken.Type)
			}
			node = &JsonAstNode{
				Type:  NodeArrayIndex,
				Value: numValue,
			}
		} else {
			return nil, fmt.Errorf("expected Wildcard or Number after Lbracket, got %v", jp.currentToken.Type)
		}
		if err := jp.nextToken(); err != nil {
			return nil, err
		}
		if jp.currentToken.Type != EOF && jp.currentToken.Type != Dot && jp.currentToken.Type != Lbracket && jp.currentToken.Type != RecursiveDescent {
			return nil, fmt.Errorf("expected Lbracket or EOF or Dot or RecursiveDescent after Array, got %v", jp.currentToken.Type)
		}
	case Wildcard:
		node = &JsonAstNode{
			Type:  NodeWildcard,
			Value: "*",
		}
		if err := jp.nextToken(); err != nil {
			return nil, err
		}
		if jp.currentToken.Type != EOF && jp.currentToken.Type != Dot && jp.currentToken.Type != Lbracket && jp.currentToken.Type != RecursiveDescent {
			return nil, fmt.Errorf("expected Lbracket or EOF or Dot or RecursiveDescent after Wildcard, got %v", jp.currentToken.Type)
		}
	case EOF:
		return nil, nil
	default:
		return nil, fmt.Errorf("expected token %v at position %d", jp.currentToken.Type, jp.currentToken.StartPos)
	}

	if jp.currentToken.Type != EOF {
		if child, err := jp.parseRest(); err != nil {
			return nil, err
		} else {
			node.Child = child
		}
	}

	return node, nil
}

func NewParser(input string) *JsonParser {
	lexer := &JsonLexer{input: input}
	return &JsonParser{lexer: lexer}
}

func (jp *JsonParser) Parse() (*JsonAstNode, error) {

	if err := jp.nextToken(); err != nil {
		return nil, err
	}
	if jp.currentToken.Type != Root {
		return nil, errors.New("JSONPath must start with $")
	}

	rootNode := &JsonAstNode{
		Type:  NodeRoot,
		Value: "$",
		Child: &JsonAstNode{},
	}

	if err := jp.nextToken(); err != nil {
		return nil, err
	}
	if jp.currentToken.Type == EOF {
		return nil, errors.New("unexpected EOF after root")
	}
	if jp.currentToken.Type != Dot && jp.currentToken.Type != Lbracket && jp.currentToken.Type != RecursiveDescent {
		return nil, fmt.Errorf("expected '.' or '[' or '..', got %v", jp.currentToken.Type)
	}

	rest, err := jp.parseRest()
	if err != nil {
		return nil, err
	}
	rootNode.Child = rest

	if jp.currentToken.Type != EOF {
		return nil, fmt.Errorf("expected EOF, got %v", jp.currentToken.Type)
	}

	return rootNode, nil
}
