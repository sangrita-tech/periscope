package replacement

import (
	"strings"

	"github.com/sangrita-tech/periscope/internal/model"
)

func Apply(text string, rules []model.Replacement) string {
	for _, rule := range rules {
		if rule.Pattern == "" {
			continue
		}
		text = strings.ReplaceAll(text, rule.Pattern, rule.Value)
	}
	return text
}
