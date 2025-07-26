package core

import (
	"strings"
	"testing"
)

// TestJsonParser 测试 JSONPath 解析器的功能
func TestJsonParser(t *testing.T) {
	// 定义测试用例
	tests := []struct {
		name     string
		input    string
		expected string // 期望的 AST 字符串表示（仅对有效用例）
		isValid  bool   // 是否为有效 JSONPath
		errMsg   string // 期望的错误信息（仅对无效用例）
	}{
		// 有效 JSONPath 测试用例
		{
			name:     "Simple root path",
			input:    "$.store",
			expected: "Root($) -> Dot(.) -> Identifier(store)",
			isValid:  true,
		},
		{
			name:     "Nested path with dot",
			input:    "$.store.book",
			expected: "Root($) -> Dot(.) -> Identifier(book)",
			isValid:  true,
		},
		{
			name:     "Array index",
			input:    "$.book[0]",
			expected: "Root($) -> Dot(.) -> Identifier(book) -> ArrayIndex(0)",
			isValid:  true,
		},
		{
			name:     "Array wildcard",
			input:    "$.store.book[*]",
			expected: "Root($) -> Dot(.) -> Identifier(book) -> ArrayWildcard([*])",
			isValid:  true,
		},
		{
			name:     "Recursive descent",
			input:    "$..book",
			expected: "Root($) -> RecursiveDescent(..) -> Identifier(book)",
			isValid:  true,
		},
		{
			name:     "Wildcard property",
			input:    "$.store.*",
			expected: "Root($) -> Dot(.) -> Wildcard(*)",
			isValid:  true,
		},
		{
			name:     "Complex path",
			input:    "$.store.book[0].author",
			expected: "Root($) -> Dot(.) -> Identifier(book) -> ArrayIndex(0) -> Dot(.) -> Identifier(author)",
			isValid:  true,
		},
		{
			name:     "Recursive descent with array",
			input:    "$..book[*].title",
			expected: "Root($) -> RecursiveDescent(..) -> Identifier(book) -> ArrayWildcard([*]) -> Dot(.) -> Identifier(title)",
			isValid:  true,
		},

		// 无效 JSONPath 测试用例
		{
			name:    "Empty input",
			input:   "",
			isValid: false,
			errMsg:  "JSONPath must start with $",
		},
		{
			name:    "Missing root",
			input:   ".store",
			isValid: false,
			errMsg:  "JSONPath must start with $",
		},
		{
			name:    "Invalid character in identifier",
			input:   "$.store@book",
			isValid: false,
			errMsg:  "@: invalid identifier",
		},
		{
			name:    "Unexpected token after identifier",
			input:   "$.store:book",
			isValid: false,
			errMsg:  "expected Lbracket or Dot or RecursiveDescent or EOF after Identifier, got Colon",
		},
		{
			name:    "Unclosed bracket",
			input:   "$.store.book[0",
			isValid: false,
			errMsg:  "expected Rbracket after Lbracket and Number, got EOF",
		},
		{
			name:    "Invalid token after Lbracket",
			input:   "$.store.book[abc]",
			isValid: false,
			errMsg:  "expected Wildcard or Number after Lbracket, got Identifier",
		},
		{
			name:    "Invalid token after dot",
			input:   "$.store.[0]",
			isValid: false,
			errMsg:  "expected Identifier or Wildcard after Dot, got Lbracket",
		},
		{
			name:    "Invalid token after recursive descent",
			input:   "$..[0]",
			isValid: false,
			errMsg:  "expected Identifier or Wildcard after RecursiveDescent, got Lbracket",
		},
		{
			name:    "Unexpected token after array",
			input:   "$.store.book[0],",
			isValid: false,
			errMsg:  "expected Lbracket or EOF or Dot or RecursiveDescent after Array, got Comma",
		},
		{
			name:    "Unexpected EOF after root",
			input:   "$",
			isValid: false,
			errMsg:  "unexpected EOF after root",
		},
	}

	// 运行测试用例
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建新的解析器
			parser := NewParser(tt.input)
			// 解析输入
			node, err := parser.Parse()

			if tt.isValid {
				// 验证有效 JSONPath
				if err != nil {
					t.Errorf("Expected no error, got: %v", err)
					return
				}
				if node == nil {
					t.Error("Expected non-nil AST node")
					return
				}
				// 验证 AST 的字符串表示
				got := node.String()
				if got != tt.expected {
					t.Errorf("Expected AST:\n%s\nGot:\n%s", tt.expected, got)
				}
			} else {
				// 验证无效 JSONPath
				if err == nil {
					t.Error("Expected error, got nil")
					return
				}
				if !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("Expected error message containing %q, got: %v", tt.errMsg, err)
				}
			}
		})
	}
}
