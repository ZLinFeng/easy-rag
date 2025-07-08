package documentloader

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/saintfish/chardet"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/encoding/traditionalchinese"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
	"io"
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

	text, err := Decode(content, t.Encoding)
	if err != nil {
		return nil, err
	}
	res := make([]Document, 1)
	res[0] = Document{
		PageContent: text,
		MetaData:    make(map[string]interface{}),
	}
	return res, nil
}

func Decode(bs []byte, enc Encoding) (string, error) {
	var decoder *encoding.Decoder
	switch enc {
	case EncodingUTF8:
		return string(bs), nil
	case EncodingUTF8BOM:
		if len(bs) >= 3 && bs[0] == 0xef && bs[1] == 0xbb && bs[2] == 0xbf {
			bs = bs[3:]
			return string(bs), nil
		} else {
			return "", errors.New("document loader: utf8bom is not supported")
		}
	case EncodingUTF16BE:
		decoder = unicode.UTF16(unicode.BigEndian, unicode.UseBOM).NewDecoder()
	case EncodingUTF16LE:
		decoder = unicode.UTF16(unicode.LittleEndian, unicode.UseBOM).NewDecoder()
	case EncodingBIG5:
		decoder = traditionalchinese.Big5.NewDecoder()
	case EncodingGBK:
		decoder = simplifiedchinese.GBK.NewDecoder()
	case EncodingISO88591:
		decoder = charmap.ISO8859_1.NewDecoder()
	case EncodingISO88595:
		decoder = charmap.ISO8859_5.NewDecoder()
	case EncodingKOI8R:
		decoder = charmap.KOI8R.NewDecoder()
	case EncodingWindows1251:
		decoder = charmap.Windows1251.NewDecoder()
	case EncodingWindows1252:
		decoder = charmap.Windows1252.NewDecoder()
	case EncodingAuto:
		detector := chardet.NewTextDetector()
		codingRes, err := detector.DetectBest(bs)
		if err != nil {
			return "", errors.New("auto detect error")
		}
		mapCode, ok := encodingMap[codingRes.Charset]
		if !ok {
			return "", fmt.Errorf("%s is not support", codingRes.Charset)
		}
		return Decode(bs, mapCode)
	default:
		return string(bs), nil
	}

	if decoder == nil {
		return "", errors.New("document loader: unsupported encoding")
	}
	reader := transform.NewReader(bytes.NewReader(bs), decoder)
	res, err := io.ReadAll(reader)
	if err != nil {
		return "", err
	}
	return string(res), nil
}
