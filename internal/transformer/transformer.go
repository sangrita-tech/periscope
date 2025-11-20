package transformer

type Result struct {
	Name    string
	Mapping map[string]string
}

type Transformer interface {
	Transform(path, content string) (string, Result, error)
}

type FuncTransformer func(path, content string) (string, Result, error)

func (f FuncTransformer) Transform(path, content string) (string, Result, error) {
	return f(path, content)
}
