package preprocess

type Chain struct {
	steps []Step
}

func New() *Chain {
	return &Chain{steps: make([]Step, 0)}
}

func (c *Chain) Add(step Step) *Chain {
	if step == nil {
		return c
	}
	c.steps = append(c.steps, step)
	return c
}

func (c *Chain) Process(path, content string) (string, Results, error) {
	out := content
	results := make(Results, len(c.steps))

	for _, s := range c.steps {
		var (
			r   map[string]string
			err error
		)

		out, r, err = s.Apply(path, out)
		if err != nil {
			return "", nil, err
		}

		if r != nil {
			results[s.Name()] = r
		}
	}

	return out, results, nil
}
