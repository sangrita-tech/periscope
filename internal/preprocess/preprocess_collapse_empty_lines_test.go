package preprocess_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/sangrita-tech/periscope/internal/preprocess"
)

func Test_PreprocessCollapseEmptyLines_CollapsesSequences(t *testing.T) {
	c := preprocess.New().AddCollapseEmptyLines()

	in := "a\n\n\n \n\t\nb\n\nc\n"
	out, res, err := c.Process("x", in)

	assert.NoError(t, err)
	assert.Equal(t, "a\n\nb\n\nc\n", out)
	assert.Empty(t, res)
}
