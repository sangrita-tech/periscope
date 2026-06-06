package ignore

import (
	"fmt"
	"path"
	"regexp"
	"strings"
)

type Matcher struct {
	globs []*regexp.Regexp
}

func NewMatcher(patterns []string) (*Matcher, error) {
	globs := make([]*regexp.Regexp, 0)

	for _, pattern := range patterns {
		glob, err := compilePattern(pattern)
		if err != nil {
			return nil, err
		}
		globs = append(globs, glob)
	}

	return &Matcher{globs: globs}, nil
}

func (m *Matcher) Match(target string) bool {
	if len(m.globs) == 0 {
		return false
	}

	target = path.Clean(target)
	if target == "" || target == "." {
		return false
	}

	for _, glob := range m.globs {
		if glob.MatchString(target) {
			return true
		}
	}

	return false
}

func compilePattern(pattern string) (*regexp.Regexp, error) {
	pattern = path.Clean(pattern)

	if hasGlobMeta(pattern) {
		return compileGlob(pattern)
	}

	return compilePlain(pattern)
}

func compilePlain(pattern string) (*regexp.Regexp, error) {
	pattern = regexp.QuoteMeta(pattern)

	if strings.Contains(pattern, "/") {
		return regexp.Compile("^" + pattern + "(/.*)?$")
	}

	return regexp.Compile("(^|/)" + pattern + "(/|$)")
}

func compileGlob(pattern string) (*regexp.Regexp, error) {
	glob := globToRegexp(pattern)

	if !strings.Contains(pattern, "/") {
		glob = "(.*/)?" + glob
	}

	re, err := regexp.Compile("^" + glob + "$")
	if err != nil {
		return nil, fmt.Errorf("compile ignore glob %q: %w", pattern, err)
	}

	return re, nil
}

func hasGlobMeta(pattern string) bool {
	return strings.ContainsAny(pattern, "*?[")
}

func globToRegexp(pattern string) string {
	var b strings.Builder

	for i := 0; i < len(pattern); i++ {
		ch := pattern[i]

		switch ch {
		case '*':
			if i+1 < len(pattern) && pattern[i+1] == '*' {
				if i+2 < len(pattern) && pattern[i+2] == '/' {
					b.WriteString("(.*/)?")
					i += 2
					continue
				}

				b.WriteString(".*")
				i++
				continue
			}

			b.WriteString("[^/]*")

		case '?':
			b.WriteString("[^/]")

		case '[':
			end := strings.IndexByte(pattern[i+1:], ']')
			if end == -1 {
				b.WriteString(regexp.QuoteMeta(string(ch)))
				continue
			}

			class := pattern[i : i+end+2]
			b.WriteString(class)
			i += end + 1

		default:
			b.WriteString(regexp.QuoteMeta(string(ch)))
		}
	}

	return b.String()
}
