package matcher_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/sangrita-tech/periscope/internal/matcher"
)

func Test_New_TrimsAndDropsEmptyMasks(t *testing.T) {
	m := matcher.New([]string{"  a* ", "   ", "", "\n\tb"})

	assert.True(t, m.Match("a123"))
	assert.True(t, m.Match("b"))
	assert.False(t, m.Match("c"))
}

func Test_Match_NoMasks_ReturnsTrueForAnyValue(t *testing.T) {
	m := matcher.New(nil)

	assert.True(t, m.Match(""))
	assert.True(t, m.Match("anything"))
	assert.True(t, m.Match("foo/bar/baz"))
}

func Test_Match_Wildcards_WorkAtEndsAndMiddle(t *testing.T) {
	m := matcher.New([]string{
		"foo*",
		"*bar",
		"a*b*c",
		"one*two",
	})

	assert.True(t, m.Match("foo"))
	assert.True(t, m.Match("foobar"))
	assert.True(t, m.Match("xxbar"))
	assert.True(t, m.Match("a___b___c"))
	assert.True(t, m.Match("oneZZtwo"))

	assert.False(t, m.Match("fo"))
	assert.False(t, m.Match("ba"))
	assert.False(t, m.Match("a_b_c_d"))
	assert.False(t, m.Match("oneZZtw"))
}

func Test_Match_LiteralRegexChars_AreQuoted(t *testing.T) {
	m := matcher.New([]string{"a.+b", "file(1).txt"})

	assert.True(t, m.Match("a.+b"))
	assert.True(t, m.Match("file(1).txt"))

	assert.False(t, m.Match("axxxb"))
	assert.False(t, m.Match("file1.txt"))
}

func Test_Match_MultipleMasks_OrSemantics(t *testing.T) {
	m := matcher.New([]string{"*.go", "*.md"})

	assert.True(t, m.Match("main.go"))
	assert.True(t, m.Match("README.md"))
	assert.False(t, m.Match("image.png"))
}
