package transformers

import (
	"strings"

	"github.com/sangrita-tech/periscope/internal/transformer"
)

type collapseEmptyLinesTransformer struct{}

func CollapseEmptyLines() transformer.Transformer {
	return &collapseEmptyLinesTransformer{}
}

func (t *collapseEmptyLinesTransformer) Transform(path, content string) (string, transformer.Result, error) {
	lines := strings.Split(content, "\n")

	var b strings.Builder
	b.Grow(len(content))

	prevEmpty := false

	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			if prevEmpty {
				continue
			}
			prevEmpty = true
			b.WriteByte('\n')
			continue
		}

		prevEmpty = false
		b.WriteString(line)
		b.WriteByte('\n')
	}

	res := transformer.Result{
		Name:    "collapse_empty_lines",
		Mapping: nil,
	}

	return b.String(), res, nil
}
