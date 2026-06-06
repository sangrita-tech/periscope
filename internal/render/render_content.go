package render

import (
	"bytes"
	"fmt"
	"path"
	"sort"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/sangrita-tech/periscope/internal/model"
)

func RenderContent(src model.Source, entries []model.Entry, generatedAt time.Time) string {
	var buffer bytes.Buffer

	writeHeader(&buffer, src.Name, generatedAt)

	for _, entry := range sortEntries(entries) {
		if isBinary(entry.Data) {
			continue
		}
		writeContentEntry(&buffer, src.Root, entry)
	}

	return limitNewLines(buffer.String())
}

func isBinary(data []byte) bool {
	if len(data) == 0 {
		return false
	}

	if bytes.Contains(data, []byte{0}) {
		return true
	}

	return !utf8.Valid(data)
}

func writeContentEntry(buffer *bytes.Buffer, root string, entry model.Entry) {
	filePath := path.Join(root, entry.RelPath)
	fence := fenceCode(entry.Data)

	fmt.Fprintf(buffer, "## %s\n\n", filePath)
	fmt.Fprintf(buffer, "%s\n", fence)

	buffer.Write(entry.Data)
	fmt.Fprintf(buffer, "%s\n\n", fence)
}

func fenceCode(data []byte) string {
	maxRun := 2
	currentRun := 0

	for _, char := range data {
		if char == '`' {
			currentRun++
			maxRun = max(maxRun, currentRun)
			continue
		}

		currentRun = 0
	}

	return strings.Repeat("`", maxRun+1)
}

func sortEntries(entries []model.Entry) []model.Entry {
	result := append([]model.Entry(nil), entries...)

	sort.SliceStable(result, func(i, j int) bool {
		return result[i].RelPath < result[j].RelPath
	})

	return result
}

func limitNewLines(value string) string {
	value = strings.ReplaceAll(value, "\r\n", "\n")
	value = strings.ReplaceAll(value, "\r", "\n")

	for strings.Contains(value, "\n\n\n") {
		value = strings.ReplaceAll(value, "\n\n\n", "\n\n")
	}

	return value
}
