package ignore

import (
	"fmt"
	"path"
	"regexp"
	"strings"
)

type Matcher struct {
	rules []rule
}

type rule struct {
	raw           string
	pathPattern   string
	directoryOnly bool
	glob          *regexp.Regexp
	regex         *regexp.Regexp
}

func NewMatcher(patterns []string) (*Matcher, error) {
	rules := make([]rule, 0)

	for _, pattern := range patterns {
		r, ok, err := newRule(pattern)
		if err != nil {
			return nil, err
		}

		if !ok {
			continue
		}

		rules = append(rules, r)
	}

	return &Matcher{rules: rules}, nil
}

func (m *Matcher) Match(relPath string, isDir bool) bool {
	if m == nil || len(m.rules) == 0 {
		return false
	}

	relPath = cleanPath(relPath)
	if relPath == "" || relPath == "." {
		return false
	}

	base := path.Base(relPath)

	for _, r := range m.rules {
		if r.matches(relPath, base, isDir) {
			return true
		}
	}

	return false
}

func newRule(pattern string) (rule, bool, error) {
	pattern = strings.TrimSpace(pattern)
	if pattern == "" {
		return rule{}, false, nil
	}

	r := rule{raw: pattern}

	r.directoryOnly = strings.HasSuffix(pattern, "/") || strings.HasSuffix(pattern, "\\")
	r.pathPattern = cleanPatternPath(pattern)

	if hasGlobMeta(r.pathPattern) {
		glob, err := compileGlob(r.pathPattern)
		if err != nil {
			return rule{}, false, fmt.Errorf("compile ignore glob %q: %w", pattern, err)
		}

		r.glob = glob
	}

	if shouldCompileRegex(pattern) {
		reText := strings.ReplaceAll(pattern, `\|`, "|")
		re, err := regexp.Compile(reText)
		if err != nil {
			return rule{}, false, fmt.Errorf("compile ignore regexp %q: %w", pattern, err)
		}

		r.regex = re
	}

	return r, true, nil
}

func (r rule) matches(relPath, base string, isDir bool) bool {
	if r.directoryOnly && !isDir && !isDescendantOf(relPath, r.pathPattern) {
		return false
	}

	if r.matchesPath(relPath, base) {
		return true
	}

	if r.glob != nil && (r.glob.MatchString(relPath) || r.glob.MatchString(base)) {
		return true
	}

	if r.regex != nil && (r.regex.MatchString(relPath) || r.regex.MatchString(base)) {
		return true
	}

	return false
}

func (r rule) matchesPath(relPath, base string) bool {
	pattern := r.pathPattern
	if pattern == "" || pattern == "." {
		return false
	}

	if relPath == pattern || base == pattern {
		return true
	}

	return isDescendantOf(relPath, pattern)
}

func isDescendantOf(relPath, dirPath string) bool {
	dirPath = strings.TrimSuffix(cleanPath(dirPath), "/")
	if dirPath == "" || dirPath == "." {
		return false
	}

	return strings.HasPrefix(relPath, dirPath+"/")
}

func cleanPath(value string) string {
	value = strings.TrimSpace(value)
	value = strings.ReplaceAll(value, "\\", "/")
	value = strings.TrimPrefix(value, "./")
	value = strings.TrimPrefix(value, "/")
	value = path.Clean(value)

	if value == "." {
		return ""
	}

	return value
}

func cleanPatternPath(value string) string {
	value = strings.TrimSpace(value)
	value = strings.TrimSuffix(value, "/")
	value = strings.TrimSuffix(value, "\\")

	value = strings.ReplaceAll(value, `\|`, "|")
	value = strings.ReplaceAll(value, "\\", "/")

	value = strings.TrimPrefix(value, "./")
	value = strings.TrimPrefix(value, "/")
	value = path.Clean(value)

	if value == "." {
		return ""
	}

	return value
}

func hasGlobMeta(pattern string) bool {
	return strings.ContainsAny(pattern, "*?[")
}

func shouldCompileRegex(pattern string) bool {
	return strings.Contains(pattern, `\|`) || strings.ContainsAny(pattern, "|()[]^$+") || strings.Contains(pattern, `\`)
}

func compileGlob(pattern string) (*regexp.Regexp, error) {
	return regexp.Compile("^" + globToRegexp(pattern) + "$")
}

func globToRegexp(pattern string) string {
	var b strings.Builder

	for i := 0; i < len(pattern); i++ {
		ch := pattern[i]

		switch ch {
		case '*':
			if i+1 < len(pattern) && pattern[i+1] == '*' {
				if i+2 < len(pattern) && pattern[i+2] == '/' {
					b.WriteString("(?:.*/)?")
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
