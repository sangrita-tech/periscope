package output

import "github.com/atotto/clipboard"

type ClipboardWriter struct{}

func NewClipboardWriter() *ClipboardWriter {
	return &ClipboardWriter{}
}

func (w *ClipboardWriter) Write(content string) error {
	return clipboard.WriteAll(content)
}
