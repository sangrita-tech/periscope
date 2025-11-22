package matcher_test

import (
	"testing"

	"github.com/sangrita-tech/periscope/internal/matcher"
	"github.com/stretchr/testify/assert"
)

func Test_Match_NoMasks_ReturnsTrue(t *testing.T) {
	m := matcher.New(nil)

	assert.True(t, m.Match("anything"))
}

func Test_Match_EmptyAndSpacedMasks_AreIgnored_ReturnsTrueForNoEffectiveMasks(t *testing.T) {
	m := matcher.New([]string{"   ", "", "\n\t"})

	assert.True(t, m.Match("anything"))
}

func Test_Match_GlobAgainstWholeValue_ReturnsTrue(t *testing.T) {
	m := matcher.New([]string{"dir/*.go"})

	assert.True(t, m.Match("dir/main.go"))
	assert.False(t, m.Match("dir/main.txt"))
}

func Test_Match_ContainsFallback_ReturnsTrue(t *testing.T) {
	m := matcher.New([]string{"vendor/"})

	assert.True(t, m.Match("src/vendor/github.com/x/y.go"))
	assert.False(t, m.Match("src/internal/y.go"))
}

func Test_Match_NormalizesSlashesInMasksAndValues_ReturnsTrue(t *testing.T) {
	m := matcher.New([]string{`dir\*.go`})

	assert.True(t, m.Match(`dir\main.go`))
	assert.True(t, m.Match("dir/main.go"))
}

func Test_Match_RandomString_ReturnsTrue(t *testing.T) {
	m := matcher.New([]string{"h*world"})

	assert.True(t, m.Match("hello world"))
}
