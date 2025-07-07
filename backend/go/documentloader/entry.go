package documentloader

type Document struct {
	PageContent string
	MetaData    map[string]interface{}
}

type BaseLoader interface {
	Load() ([]Document, error)
}

type Encoding int

const (
	EncodingAuto Encoding = iota
	EncodingUTF8
	EncodingUTF8BOM
	EncodingUTF16LE
	EncodingUTF16BE
	EncodingGBK
	EncodingBIG5
	EncodingISO88591
	EncodingWindows1251 // 西里尔字母编码
	EncodingWindows1252 // 拉丁字母编码
	EncodingKOI8R       // 俄语编码
	EncodingISO88595
)

func (e Encoding) Value() string {
	switch e {
	case EncodingAuto:
		return "auto"
	case EncodingUTF8:
		return "utf-8"
	case EncodingUTF8BOM:
		return "utf-8"
	case EncodingUTF16LE:
		return "utf-16le"
	case EncodingUTF16BE:
		return "utf-16be"
	case EncodingGBK:
		return "gbk"
	case EncodingBIG5:
		return "big5"
	case EncodingISO88591:
		return "iso-8859-1"
	case EncodingWindows1251:
		return "windows-1251"
	case EncodingWindows1252:
		return "windows-1252"
	case EncodingKOI8R:
		return "koi8-r"
	case EncodingISO88595:
		return "iso-8859-5"
	default:
		return "utf-8"
	}
}
