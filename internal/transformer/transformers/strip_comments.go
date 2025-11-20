package transformers

import (
	"strings"

	"github.com/sangrita-tech/periscope/internal/transformer"
)

type stripCommentsTransformer struct{}

func StripComments() transformer.Transformer {
	return &stripCommentsTransformer{}
}

func (t *stripCommentsTransformer) Transform(path, content string) (string, transformer.Result, error) {
	lines := strings.Split(content, "\n")

	var b strings.Builder
	b.Grow(len(content))

	for _, line := range lines {
		trimmed := strings.TrimLeft(line, " \t")

		if strings.HasPrefix(trimmed, "//") ||
			strings.HasPrefix(trimmed, "#") ||
			strings.HasPrefix(trimmed, "--") ||
			strings.HasPrefix(trimmed, ";") {
			continue
		}

		b.WriteString(line)
		b.WriteByte('\n')
	}

	res := transformer.Result{
		Name:    "strip_comments",
		Mapping: nil,
	}

	return b.String(), res, nil
}
