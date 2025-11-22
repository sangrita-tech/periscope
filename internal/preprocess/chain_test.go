package preprocess_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/sangrita-tech/periscope/internal/preprocess"
)

func Test_Chain_Process_NoSteps_ReturnsOriginalContentAndEmptyResults(t *testing.T) {
	c := preprocess.New()

	out, res, err := c.Process("x", "hello")

	assert.NoError(t, err)
	assert.Equal(t, "hello", out)
	assert.Empty(t, res)
}

func Test_Chain_Add_NilStep_DoesNothing(t *testing.T) {
	c := preprocess.New()

	c.Add(nil)

	out, res, err := c.Process("x", "hello")

	assert.NoError(t, err)
	assert.Equal(t, "hello", out)
	assert.Empty(t, res)
}

func Test_Chain_MultipleSteps_ExecutedInOrder(t *testing.T) {
	c := preprocess.New().
		AddMaskURL().
		AddStripComments().
		AddCollapseEmptyLines()

	in := "http://google.com\n\n// comment\n\nhi\n\n"
	out, res, err := c.Process("path", in)

	assert.NoError(t, err)
	assert.NotContains(t, out, "google.com")
	assert.NotContains(t, out, "// comment")
	assert.Contains(t, out, ".example.com")
	assert.Contains(t, out, "hi\n")
	assert.Contains(t, res, "mask_url")
}
