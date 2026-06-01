package render

import (
	"bytes"
	"slices"
	"strings"
	"time"

	"github.com/sangrita-tech/periscope/internal/domain"
)

type Renderer struct {
}

func NewRenderer() *Renderer {
	return &Renderer{}
}

func (r *Renderer) Render(files []domain.Entry) string {
	r.sortEntries(files)

	var buffer bytes.Buffer

	buffer.WriteString("# Periscoped project " + time.Now().Format("2006-01-02 15:04:05") + "\n\n")

	for _, file := range files {
		buffer.WriteString("## " + file.RelPath + "\n\n")
		buffer.WriteString("```\n")
		buffer.Write(file.Data)
		buffer.WriteString("```\n\n")
	}

	bufferString := buffer.String()

	return r.limitNewLines(bufferString)
}

func (r *Renderer) sortEntries(entries []domain.Entry) {
	slices.SortFunc(entries, func(a, b domain.Entry) int {
		if a.RelPath < b.RelPath {
			return -1
		}

		if a.RelPath > b.RelPath {
			return 1
		}

		return 0
	})
}

func (r *Renderer) limitNewLines(value string) string {
	value = strings.ReplaceAll(value, "\r\n", "\n")
	value = strings.ReplaceAll(value, "\r", "\n")

	for strings.Contains(value, "\n\n\n") {
		value = strings.ReplaceAll(value, "\n\n\n", "\n\n")
	}

	return value
}
