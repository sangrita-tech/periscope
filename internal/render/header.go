package render

import (
	"bytes"
	"time"
)

const headerTimeLayout = "2006-01-02 15:04:05"

func writeHeader(buffer *bytes.Buffer, name string, generatedAt time.Time) {
	buffer.WriteString("# Periscoped project " + name + " " + generatedAt.Format(headerTimeLayout) + "\n\n")
}
