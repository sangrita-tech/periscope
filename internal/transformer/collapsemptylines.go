package transformer

import (
	"bufio"
	"strings"
)

type collapseEmptyLinesTransformer struct{}

func CollapseEmptyLines() Transformer {
	return collapseEmptyLinesTransformer{}
}

func (collapseEmptyLinesTransformer) Transform(path, content string) (string, error) {
	var b strings.Builder

	scanner := bufio.NewScanner(strings.NewReader(content))
	prevEmpty := false

	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) == "" {
			if prevEmpty {
				continue
			}
			prevEmpty = true
			b.WriteString("\n")
			continue
		}

		prevEmpty = false
		b.WriteString(line)
		b.WriteString("\n")
	}

	if err := scanner.Err(); err != nil {
		return content, nil
	}

	out := b.String()

	out = strings.TrimRight(out, "\n")
	if out == "" {
		return "", nil
	}
	out += "\n"

	return out, nil
}
