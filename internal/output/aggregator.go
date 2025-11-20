package output

import (
	"bytes"
	"fmt"
	"io"
	"path/filepath"
	"strings"
)

type Aggregator struct {
	buf             bytes.Buffer
	first           bool
	copyToClipboard bool
	writer          io.Writer
}

func NewAggregator(copyToClipboard bool, writer io.Writer) *Aggregator {
	return &Aggregator{
		copyToClipboard: copyToClipboard,
		writer:          writer,
		first:           true,
	}
}

func (a *Aggregator) HandleFile(path, content string) error {
	if strings.TrimSpace(content) == "" {
		return nil
	}

	absPath, err := filepath.Abs(path)
	if err != nil {
		absPath = path
	}

	header := fmt.Sprintf("[FILE] %s\n\n", absPath)

	if !a.copyToClipboard && a.writer != nil {
		if !a.first {
			fmt.Fprintln(a.writer)
		}

		fmt.Fprint(a.writer, header)
		fmt.Fprint(a.writer, content)
		if !strings.HasSuffix(content, "\n") {
			fmt.Fprintln(a.writer)
		}
	}

	if !a.first {
		a.buf.WriteString("\n")
	}
	a.buf.WriteString(header)
	a.buf.WriteString(content)
	if !strings.HasSuffix(content, "\n") {
		a.buf.WriteString("\n")
	}

	a.first = false
	return nil
}

func (a *Aggregator) Result() string {
	return a.buf.String()
}
