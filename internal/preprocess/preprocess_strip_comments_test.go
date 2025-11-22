package preprocess_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/sangrita-tech/periscope/internal/preprocess"
)

func Test_PreprocessStripComments_RemovesPureCommentLines(t *testing.T) {
	c := preprocess.New().AddStripComments()

	in := "// c1\nx // keep\n# c2\n-- c3\n; c4\ny\n"
	out, res, err := c.Process("file", in)

	assert.NoError(t, err)
	assert.Equal(t, "x // keep\ny\n", out)
	assert.Empty(t, res)
}
