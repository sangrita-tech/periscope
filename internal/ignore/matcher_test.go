package ignore

import "testing"

func TestMatcherMatchesLiteralDirectory(t *testing.T) {
	matcher, err := NewMatcher([]string{"vendor"})
	if err != nil {
		t.Fatal(err)
	}

	if !matcher.Match("vendor", true) {
		t.Fatal("vendor directory should be ignored")
	}

	if !matcher.Match("internal/vendor", true) {
		t.Fatal("nested vendor directory should be ignored")
	}

	if matcher.Match("vendor.go", false) {
		t.Fatal("vendor.go should not be ignored by literal vendor rule")
	}
}

func TestMatcherMatchesGlobs(t *testing.T) {
	matcher, err := NewMatcher([]string{"*.yaml", "*_test.go", "**/*.yml"})
	if err != nil {
		t.Fatal(err)
	}

	cases := []string{
		"config.yaml",
		"configs/config.yaml",
		"internal/foo_test.go",
		"configs/local.yml",
	}

	for _, tc := range cases {
		if !matcher.Match(tc, false) {
			t.Fatalf("%s should be ignored", tc)
		}
	}
}

func TestMatcherMatchesGrepStyleAlternation(t *testing.T) {
	matcher, err := NewMatcher([]string{`some\|other`})
	if err != nil {
		t.Fatal(err)
	}

	if !matcher.Match("pkg/some/file.go", false) {
		t.Fatal("regexp alternative some should be ignored")
	}

	if !matcher.Match("pkg/other/file.go", false) {
		t.Fatal("regexp alternative other should be ignored")
	}

	if matcher.Match("pkg/third/file.go", false) {
		t.Fatal("third should not be ignored")
	}
}
