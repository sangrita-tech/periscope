package preprocess_test

import (
	"testing"

	"github.com/sangrita-tech/periscope/internal/preprocess"
	"github.com/stretchr/testify/assert"
)

func Test_Helpers_AddStripComments_WiresStepIntoChain(t *testing.T) {
	c := preprocess.New().AddStripComments()

	out, _, err := c.Process("x", "// a\nb\n")
	assert.NoError(t, err)
	assert.Equal(t, "b\n", out)
}

func Test_Helpers_AddMaskURL_WiresStepIntoChain(t *testing.T) {
	c := preprocess.New().AddMaskURL()

	out, res, err := c.Process("x", "http://google.com\n")
	assert.NoError(t, err)
	assert.NotContains(t, out, "google.com")
	assert.Contains(t, res, "mask_url")
}

func Test_Helpers_AddCollapseEmptyLines_WiresStepIntoChain(t *testing.T) {
	c := preprocess.New().AddCollapseEmptyLines()

	out, _, err := c.Process("x", "a\n\n\nb\n")
	assert.NoError(t, err)
	assert.Equal(t, "a\n\nb\n", out)
}
