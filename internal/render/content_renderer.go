package render

import (
	"bytes"
	"slices"
	"strings"
	"time"

	"github.com/sangrita-tech/periscope/internal/domain"
)

type ContentRenderer struct {
}

func NewContentRenderer() *ContentRenderer {
	return &ContentRenderer{}
}

func (r *ContentRenderer) Render(source domain.Source, entries []domain.Entry) string {
	r.sortEntries(entries)

	var buffer bytes.Buffer

	buffer.WriteString("# Periscoped project " + source.Name + " " + time.Now().Format("2006-01-02 15:04:05") + "\n\n")

	for _, entry := range entries {
		buffer.WriteString("## " + source.Name + "/" + entry.RelPath + "\n\n")
		buffer.WriteString("```\n")
		buffer.Write(entry.Data)
		buffer.WriteString("```\n\n")
	}

	bufferString := buffer.String()

	return r.limitNewLines(bufferString)
}

func (r *ContentRenderer) sortEntries(entries []domain.Entry) {
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

func (r *ContentRenderer) limitNewLines(value string) string {
	value = strings.ReplaceAll(value, "\r\n", "\n")
	value = strings.ReplaceAll(value, "\r", "\n")

	for strings.Contains(value, "\n\n\n") {
		value = strings.ReplaceAll(value, "\n\n\n", "\n\n")
	}

	return value
}
