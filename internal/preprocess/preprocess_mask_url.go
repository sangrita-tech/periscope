package preprocess

import (
	"hash/fnv"
	"strings"

	"mvdan.cc/xurls/v2"
)

type preprocessMaskURL struct {
	fakeRootDomain string
	mapping        map[string]string
}

func newPreprocessMaskURL() Step {
	return &preprocessMaskURL{
		fakeRootDomain: "example.com",
		mapping:        make(map[string]string, 32),
	}
}

func (s *preprocessMaskURL) Name() string {
	return "mask_url"
}

func (s *preprocessMaskURL) Apply(path, content string) (string, map[string]string, error) {
	xurlsStrict := xurls.Strict()

	out := xurlsStrict.ReplaceAllStringFunc(content, func(u string) string {
		schemeEnd := strings.Index(u, "://")
		if schemeEnd == -1 {
			return u
		}

		authorityStart := schemeEnd + 3
		authorityEnd := len(u)

		for i, ch := range u[authorityStart:] {
			if ch == '/' || ch == '?' || ch == '#' {
				authorityEnd = authorityStart + i
				break
			}
		}

		authority := u[authorityStart:authorityEnd]
		if authority == "" {
			return u
		}

		domain := authority

		if lastColon := strings.LastIndex(domain, ":"); lastColon != -1 {
			if !strings.Contains(domain, "]") {
				domain = domain[:lastColon]
			}
		}

		if atIndex := strings.LastIndex(domain, "@"); atIndex != -1 {
			domain = domain[atIndex+1:]
		}

		if domain == "" {
			return u
		}

		fake, ok := s.mapping[domain]
		if !ok {
			fake = s.generateFakeDomain(domain)
			s.mapping[domain] = fake
		}

		return u[:authorityStart] + strings.Replace(u[authorityStart:], domain, fake, 1)
	})

	return out, s.mapping, nil
}

func (s *preprocessMaskURL) generateFakeDomain(domain string) string {
	const alphabet = "abcdefghijklmnopqrstuvwxyz0123456789"
	const base = uint32(len(alphabet))

	h := fnv.New32a()
	_, _ = h.Write([]byte(domain))
	n := h.Sum32()

	var buf [6]byte
	for i := len(buf) - 1; i >= 0; i-- {
		buf[i] = alphabet[n%base]
		n /= base
	}

	return string(buf[:]) + "." + s.fakeRootDomain
}
