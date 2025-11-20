package transformers

import (
	"hash/fnv"
	"strings"

	"github.com/sangrita-tech/periscope/internal/transformer"
	"mvdan.cc/xurls/v2"
)

type maskURLTransformer struct {
	fakeRootDomain string
	mapping        map[string]string
	counter        uint32
}

func MaskURL() transformer.Transformer {
	return &maskURLTransformer{
		fakeRootDomain: "example.com",
		mapping:        make(map[string]string, 32),
		counter:        0,
	}
}

func (t *maskURLTransformer) Transform(path, content string) (string, transformer.Result, error) {
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

		fakeDomain, ok := t.mapping[domain]
		if !ok {
			t.mapping[domain] = t.generateFakeDomain(domain)
		}

		return u[:authorityStart] + strings.Replace(u[authorityStart:], domain, fakeDomain, 1)
	})

	res := transformer.Result{
		Name:    "mask_url",
		Mapping: t.mapping,
	}

	return out, res, nil
}

func (t *maskURLTransformer) generateFakeDomain(domain string) string {
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

	return string(buf[:]) + "." + t.fakeRootDomain
}
