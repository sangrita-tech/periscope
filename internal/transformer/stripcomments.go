package transformer

import (
	"bufio"
	"strings"
)

type stripCommentsTransformer struct{}

func StripComments() Transformer {
	return stripCommentsTransformer{}
}

func (stripCommentsTransformer) Transform(path, content string) (string, error) {
	var b strings.Builder
	scanner := bufio.NewScanner(strings.NewReader(content))

	for scanner.Scan() {
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)

		if strings.HasPrefix(trimmed, "//") ||
			strings.HasPrefix(trimmed, "#") ||
			strings.HasPrefix(trimmed, "--") ||
			strings.HasPrefix(trimmed, ";") {
			continue
		}

		b.WriteString(line)
		b.WriteString("\n")
	}

	if err := scanner.Err(); err != nil {
		return content, nil
	}

	return b.String(), nil
}
