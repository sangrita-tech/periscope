package preprocess_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/sangrita-tech/periscope/internal/preprocess"
)

func Test_PreprocessMaskURL_ReplacesDomainAndReturnsMapping(t *testing.T) {
	c := preprocess.New().AddMaskURL()

	in := "http://google.com/a\nhttps://google.com/b\n"
	out, res, err := c.Process("x", in)

	assert.NoError(t, err)

	assert.NotContains(t, out, "google.com")
	assert.Contains(t, out, ".example.com")

	assert.Contains(t, res, "mask_url")
	assert.Equal(t, 1, len(res["mask_url"]))
	assert.Contains(t, res["mask_url"], "google.com")
	assert.Contains(t, res["mask_url"]["google.com"], ".")
}
