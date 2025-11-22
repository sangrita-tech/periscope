package contentbuilder

import (
	"bytes"
	"strings"
)

type ContentBuilder struct {
	buf   bytes.Buffer
	first bool
}

func New() *ContentBuilder {
	return &ContentBuilder{first: true}
}

func (b *ContentBuilder) AddBlock(label, content string) error {
	if strings.TrimSpace(content) == "" {
		return nil
	}

	header := label + "\n\n"

	if !b.first {
		if err := b.writeString("\n"); err != nil {
			return err
		}
	}

	if err := b.writeString(header); err != nil {
		return err
	}

	if err := b.writeString(content); err != nil {
		return err
	}

	if !strings.HasSuffix(content, "\n") {
		if err := b.writeString("\n"); err != nil {
			return err
		}
	}

	b.first = false
	return nil
}

func (b *ContentBuilder) Result() string {
	return b.buf.String()
}

func (b *ContentBuilder) writeString(s string) error {
	_, err := b.buf.WriteString(s)
	return err
}
