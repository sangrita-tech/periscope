package preprocess_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/sangrita-tech/periscope/internal/preprocess"
)

func Test_StepFunc_Apply_CallsWrappedFunction(t *testing.T) {
	step := preprocess.NewStep(
		"demo",
		func(path, content string) (string, map[string]string, error) {
			return content + "_x", map[string]string{"k": "v"}, nil
		},
	)

	out, res, err := step.Apply("p", "hi")

	assert.NoError(t, err)
	assert.Equal(t, "hi_x", out)
	assert.Equal(t, map[string]string{"k": "v"}, res)
	assert.Equal(t, "demo", step.Name())
}
