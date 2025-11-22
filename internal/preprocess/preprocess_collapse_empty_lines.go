package preprocess

import "strings"

type preprocessCollapseEmptyLines struct{}

func newPreprocessCollapseEmptyLines() Step {
	return &preprocessCollapseEmptyLines{}
}

func (s *preprocessCollapseEmptyLines) Name() string {
	return "collapse_empty_lines"
}

func (s *preprocessCollapseEmptyLines) Apply(path, content string) (string, map[string]string, error) {
	hadTrailingNewline := strings.HasSuffix(content, "\n")
	lines := strings.Split(content, "\n")

	var b strings.Builder
	b.Grow(len(content))

	firstOut := true
	prevEmpty := false

	for i, line := range lines {
		if i == len(lines)-1 && line == "" {
			continue
		}

		if strings.TrimSpace(line) == "" {
			if prevEmpty {
				continue
			}
			prevEmpty = true

			if !firstOut {
				b.WriteByte('\n')
			}
			firstOut = false
			continue
		}

		prevEmpty = false

		if !firstOut {
			b.WriteByte('\n')
		}
		firstOut = false

		b.WriteString(line)
	}

	if hadTrailingNewline && !firstOut {
		b.WriteByte('\n')
	}

	return b.String(), nil, nil
}
