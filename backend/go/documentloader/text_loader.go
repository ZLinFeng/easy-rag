package documentloader

import (
	"errors"
	"golang.org/x/net/html/charset"
	"golang.org/x/text/encoding"
	"os"
)

type TextLoader struct {
	FilePath string
	Encoding Encoding
}

func DefaultTextLoader(filePath string) *TextLoader {
	return &TextLoader{
		FilePath: filePath,
		Encoding: EncodingUTF8,
	}
}

func (t *TextLoader) Load() ([]Document, error) {
	// 是否为可读文件
	fileInfo, err := os.Stat(t.FilePath)
	if err != nil {
		return nil, err
	}
	if fileInfo.IsDir() {
		return nil, errors.New("document loader: not a file")
	}

	// 内容是否为空
	content, err := os.ReadFile(t.FilePath)
	if err != nil {
		return nil, err
	}
	if len(content) == 0 {
		return nil, errors.New("document loader: empty file")
	}

	return []Document{}, nil
}
