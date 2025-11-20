package transformer

type Results map[string]Result

type Pipeline struct {
	transformers []Transformer
}

func New() *Pipeline {
	return &Pipeline{
		transformers: make([]Transformer, 0),
	}
}

func (p *Pipeline) Add(t Transformer) *Pipeline {
	if t == nil {
		return p
	}
	p.transformers = append(p.transformers, t)
	return p
}

func (p *Pipeline) Process(path, content string) (string, Results, error) {
	out := content
	results := make(Results, len(p.transformers))

	for _, t := range p.transformers {
		r := Result{}
		var err error

		out, r, err = t.Transform(path, out)
		if err != nil {
			return "", nil, err
		}

		if r.Name != "" {
			results[r.Name] = r
		}
	}

	return out, results, nil
}
