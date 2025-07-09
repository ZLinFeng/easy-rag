package documentloader

type JsonLoader struct {
	JsonPath    string
	FilePath    string
	IsJsonLines bool
}

func (l *JsonLoader) Load() ([]Document, error) {
	return nil, nil
}
