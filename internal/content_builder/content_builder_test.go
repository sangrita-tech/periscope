package contentbuilder_test

import (
	"testing"

	contentbuilder "github.com/sangrita-tech/periscope/internal/content_builder"
	"github.com/stretchr/testify/assert"
)

func Test_AddBlock_WhitespaceContent_ReturnsEmptyResult(t *testing.T) {
	b := contentbuilder.New()

	b.AddBlock("a.txt", "   \n\t  ")
	b.AddBlock("b.txt", "")

	assert.Equal(t, "", b.Result())
}

func Test_AddBlock_NoTrailingNewline_ReturnsNormalizedBlock(t *testing.T) {
	b := contentbuilder.New()

	b.AddBlock("one", "hello")

	assert.Equal(t, "one\n\nhello\n", b.Result())
}

func Test_AddBlock_WithTrailingNewline_ReturnsSameBlock(t *testing.T) {
	b := contentbuilder.New()

	b.AddBlock("one", "hello\n")

	assert.Equal(t, "one\n\nhello\n", b.Result())
}

func Test_AddBlock_MultipleBlocks_ReturnsSeparatedBlocks(t *testing.T) {
	b := contentbuilder.New()

	b.AddBlock("one", "hello")
	b.AddBlock("two", "world\n")

	want :=
		"one\n\nhello\n" +
			"\n" +
			"two\n\nworld\n"

	assert.Equal(t, want, b.Result())
}
