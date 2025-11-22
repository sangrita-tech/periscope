package preprocess

import "strings"

type preprocessStripComments struct{}

func newPreprocessStripComments() Step {
	return &preprocessStripComments{}
}

func (s *preprocessStripComments) Name() string {
	return "strip_comments"
}

func (s *preprocessStripComments) Apply(path, content string) (string, map[string]string, error) {
	hadTrailingNewline := strings.HasSuffix(content, "\n")
	lines := strings.Split(content, "\n")

	var b strings.Builder
	b.Grow(len(content))

	wroteAny := false

	for i, line := range lines {
		if i == len(lines)-1 && line == "" {
			continue
		}

		trimmed := strings.TrimLeft(line, " \t")
		if strings.HasPrefix(trimmed, "//") ||
			strings.HasPrefix(trimmed, "#") ||
			strings.HasPrefix(trimmed, "--") ||
			strings.HasPrefix(trimmed, ";") {
			continue
		}

		if wroteAny {
			b.WriteByte('\n')
		}
		wroteAny = true
		b.WriteString(line)
	}

	if hadTrailingNewline && wroteAny {
		b.WriteByte('\n')
	}

	return b.String(), nil, nil
}
