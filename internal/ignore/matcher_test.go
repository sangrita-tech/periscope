package ignore

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func mustMatcher(t *testing.T, patterns []string) *Matcher {
	t.Helper()

	matcher, err := NewMatcher(patterns)
	assert.NoError(t, err)

	if err != nil {
		t.FailNow()
	}

	return matcher
}

func Test_EmptyPatterns_ReturnsFalse(t *testing.T) {
	matcher := mustMatcher(t, nil)

	matched := matcher.Match("file.go")

	assert.False(t, matched)
}

func Test_EmptyTarget_ReturnsFalse(t *testing.T) {
	matcher := mustMatcher(t, []string{"file.go"})

	matched := matcher.Match("")

	assert.False(t, matched)
}

func Test_CurrentDirectoryTarget_ReturnsFalse(t *testing.T) {
	matcher := mustMatcher(t, []string{"."})

	matched := matcher.Match(".")

	assert.False(t, matched)
}

func Test_PlainFileName_MatchesRootFile(t *testing.T) {
	matcher := mustMatcher(t, []string{"debug.log"})

	matched := matcher.Match("debug.log")

	assert.True(t, matched)
}

func Test_PlainFileName_MatchesNestedFile(t *testing.T) {
	matcher := mustMatcher(t, []string{"debug.log"})

	matched := matcher.Match("tmp/debug.log")

	assert.True(t, matched)
}

func Test_PlainDirectoryName_MatchesNestedFile(t *testing.T) {
	matcher := mustMatcher(t, []string{"vendor"})

	matched := matcher.Match("internal/vendor/module.go")

	assert.True(t, matched)
}

func Test_PathPattern_MatchesExactPathPrefix(t *testing.T) {
	matcher := mustMatcher(t, []string{"build/output"})

	matched := matcher.Match("build/output/app.bin")

	assert.True(t, matched)
}

func Test_PathPattern_DoesNotMatchSamePathUnderOtherDirectory(t *testing.T) {
	matcher := mustMatcher(t, []string{"build/output"})

	matched := matcher.Match("tmp/build/output/app.bin")

	assert.False(t, matched)
}

func Test_StarGlob_MatchesFileName(t *testing.T) {
	matcher := mustMatcher(t, []string{"*.log"})

	matched := matcher.Match("app.log")

	assert.True(t, matched)
}

func Test_StarGlob_DoesNotMatchPathSeparator(t *testing.T) {
	matcher := mustMatcher(t, []string{"logs/*.log"})

	matched := matcher.Match("logs/archive/app.log")

	assert.False(t, matched)
}

func Test_DoubleStarWithSlash_MatchesFileAtAnyDepth(t *testing.T) {
	matcher := mustMatcher(t, []string{"logs/**/*.log"})

	matched := matcher.Match("logs/2026/06/app.log")

	assert.True(t, matched)
}

func Test_DoubleStarWithoutSlash_MatchesAnyRemainingPath(t *testing.T) {
	matcher := mustMatcher(t, []string{"tmp/**"})

	matched := matcher.Match("tmp/a/b/c/cache.bin")

	assert.True(t, matched)
}

func Test_QuestionMark_MatchesSingleCharacter(t *testing.T) {
	matcher := mustMatcher(t, []string{"file-?.go"})

	matched := matcher.Match("file-a.go")

	assert.True(t, matched)
}

func Test_QuestionMark_DoesNotMatchMultipleCharacters(t *testing.T) {
	matcher := mustMatcher(t, []string{"file-?.go"})

	matched := matcher.Match("file-ab.go")

	assert.False(t, matched)
}

func Test_CharacterClass_MatchesAllowedCharacter(t *testing.T) {
	matcher := mustMatcher(t, []string{"file-[ab].go"})

	matched := matcher.Match("file-a.go")

	assert.True(t, matched)
}

func Test_UnclosedCharacterClass_IsTreatedAsPlainText(t *testing.T) {
	matcher := mustMatcher(t, []string{"file-[ab.go"})

	matched := matcher.Match("file-[ab.go")

	assert.True(t, matched)
}

func Test_InvalidGlob_ReturnsError(t *testing.T) {
	matcher, err := NewMatcher([]string{"file-[z-a].go"})

	assert.Error(t, err)
	assert.Nil(t, matcher)
}

func Test_MultiplePatterns_MatchesAnySuitablePattern(t *testing.T) {
	matcher := mustMatcher(t, []string{"*.tmp", "cache", "dist/**"})

	matched := matcher.Match("pkg/cache/data.bin")

	assert.True(t, matched)
}

func Test_CleanedTargetPath_MatchesNormalizedPath(t *testing.T) {
	matcher := mustMatcher(t, []string{"logs/app.log"})

	matched := matcher.Match("logs/./archive/../app.log")

	assert.True(t, matched)
}

func Test_NonMatchingTarget_ReturnsFalse(t *testing.T) {
	matcher := mustMatcher(t, []string{"*.log"})

	matched := matcher.Match("app.txt")

	assert.False(t, matched)
}
