package matcher

import (
	"path/filepath"
	"strings"
)

type Matcher struct {
	masks []string
}

func New(masks []string) *Matcher {
	clean := make([]string, 0, len(masks))
	for _, mask := range masks {
		mask = strings.TrimSpace(mask)
		if mask == "" {
			continue
		}
		clean = append(clean, filepath.ToSlash(mask))
	}
	return &Matcher{masks: clean}
}

func (m *Matcher) Match(value string) bool {
	if len(m.masks) == 0 {
		return true
	}

	norm := filepath.ToSlash(value)

	for _, mask := range m.masks {
		if ok, _ := filepath.Match(mask, norm); ok {
			return true
		}
		if strings.Contains(norm, mask) {
			return true
		}
	}

	return false
}
