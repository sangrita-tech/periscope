package transformer

type Pipeline struct {
	transformers []Transformer
}

func NewPipeline() *Pipeline {
	return &Pipeline{
		transformers: make([]Transformer, 0),
	}
}

func (p *Pipeline) Use(t Transformer) {
	if t == nil {
		return
	}
	p.transformers = append(p.transformers, t)
}

func (p *Pipeline) Process(path, content string) (string, error) {
	var err error
	out := content

	for _, t := range p.transformers {
		out, err = t.Transform(path, out)
		if err != nil {
			return "", err
		}
	}

	return out, nil
}
