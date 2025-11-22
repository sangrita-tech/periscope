package preprocess

type Step interface {
	Name() string
	Apply(path, content string) (string, map[string]string, error)
}

type StepFunc struct {
	name string
	fn   func(path, content string) (string, map[string]string, error)
}

func NewStep(name string, fn func(path, content string) (string, map[string]string, error)) Step {
	return StepFunc{name: name, fn: fn}
}

func (s StepFunc) Name() string {
	return s.name
}

func (s StepFunc) Apply(path, content string) (string, map[string]string, error) {
	return s.fn(path, content)
}
