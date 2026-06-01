package render

import (
	"github.com/sangrita-tech/periscope/internal/domain"
)

type Renderer interface {
	Render(entries []domain.Entry) string
}
