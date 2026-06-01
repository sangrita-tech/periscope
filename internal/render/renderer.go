package render

import (
	"github.com/sangrita-tech/periscope/internal/domain"
)

type Renderer interface {
	Render(source domain.Source, entries []domain.Entry) string
}
