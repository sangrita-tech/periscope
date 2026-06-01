package output

import "io"

type TerminalWriter struct {
	out io.Writer
}

func NewTerminalWriter(out io.Writer) *TerminalWriter {
	return &TerminalWriter{
		out: out,
	}
}

func (w *TerminalWriter) Write(content string) error {
	_, err := io.WriteString(w.out, content)
	return err
}
