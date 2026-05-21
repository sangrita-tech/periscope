package clipboard

import (
	"fmt"

	systemclipboard "github.com/atotto/clipboard"
	"github.com/sangrita-tech/periscope/internal/domain"
)

type clipboardWriter struct{}

func NewClipboardWriter() domain.ClipboardWriter {
	return &clipboardWriter{}
}

func (w *clipboardWriter) Write(result domain.InspectionResult) error {
	if err := systemclipboard.WriteAll(result.Text); err != nil {
		return fmt.Errorf("copy to clipboard: %w", err)
	}
	return nil
}
