package matcher

import (
	"regexp"
	"strings"
)

type Matcher struct {
	raw  []string
	res  []*regexp.Regexp
}

func New(masks []string) *Matcher {
	clean := make([]string, 0, len(masks))
	res := make([]*regexp.Regexp, 0, len(masks))

	for _, mask := range masks {
		mask = strings.TrimSpace(mask)
		if mask == "" {
			continue
		}

		clean = append(clean, mask)
		res = append(res, compileMask(mask))
	}

	return &Matcher{raw: clean, res: res}
}

func (m *Matcher) Match(value string) bool {
	if len(m.res) == 0 {
		return true
	}

	for _, re := range m.res {
		if re.MatchString(value) {
			return true
		}
	}

	return false
}

func compileMask(mask string) *regexp.Regexp {
	parts := strings.Split(mask, "*")
	var b strings.Builder
	b.Grow(len(mask) * 2)

	b.WriteByte('^')
	for i, p := range parts {
		if i > 0 {
			b.WriteString(".*")
		}
		b.WriteString(regexp.QuoteMeta(p))
	}
	b.WriteByte('$')

	return regexp.MustCompile(b.String())
}
