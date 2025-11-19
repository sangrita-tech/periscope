package transformer

type Transformer interface {
	Transform(path, content string) (string, error)
}

type FuncTransformer func(path, content string) (string, error)

func (f FuncTransformer) Transform(path, content string) (string, error) {
	return f(path, content)
}
