package replacement

import (
	"testing"

	"github.com/sangrita-tech/periscope/internal/model"
	"github.com/stretchr/testify/assert"
)

func Test_NewReplacer_ReturnsReplacer(t *testing.T) {
	replacer := NewReplacer()

	assert.NotNil(t, replacer)
}

func Test_Apply_WithEmptyRules_ReturnsOriginalText(t *testing.T) {
	replacer := NewReplacer()
	text := "hello world"

	result := replacer.Apply(text, nil)

	assert.Equal(t, text, result)
}

func Test_Apply_WithSingleRule_ReplacesText(t *testing.T) {
	replacer := NewReplacer()
	rules := []model.Replacement{
		{
			Pattern: "world",
			Value:   "gopher",
		},
	}

	result := replacer.Apply("hello world", rules)

	assert.Equal(t, "hello gopher", result)
}

func Test_Apply_WithRepeatedPattern_ReplacesAllOccurrences(t *testing.T) {
	replacer := NewReplacer()
	rules := []model.Replacement{
		{
			Pattern: "go",
			Value:   "rust",
		},
	}

	result := replacer.Apply("go is go and go", rules)

	assert.Equal(t, "rust is rust and rust", result)
}

func Test_Apply_WithMissingPattern_ReturnsOriginalText(t *testing.T) {
	replacer := NewReplacer()
	text := "hello world"
	rules := []model.Replacement{
		{
			Pattern: "python",
			Value:   "go",
		},
	}

	result := replacer.Apply(text, rules)

	assert.Equal(t, text, result)
}

func Test_Apply_WithEmptyPattern_IgnoresRule(t *testing.T) {
	replacer := NewReplacer()
	text := "hello world"
	rules := []model.Replacement{
		{
			Pattern: "",
			Value:   "broken",
		},
	}

	result := replacer.Apply(text, rules)

	assert.Equal(t, text, result)
}

func Test_Apply_WithMultipleRules_AppliesRulesInOrder(t *testing.T) {
	replacer := NewReplacer()
	rules := []model.Replacement{
		{
			Pattern: "hello",
			Value:   "hi",
		},
		{
			Pattern: "hi",
			Value:   "bye",
		},
	}

	result := replacer.Apply("hello world", rules)

	assert.Equal(t, "bye world", result)
}

func Test_Apply_WithEmptyValue_RemovesPattern(t *testing.T) {
	replacer := NewReplacer()
	rules := []model.Replacement{
		{
			Pattern: "bad ",
			Value:   "",
		},
	}

	result := replacer.Apply("bad hello bad world", rules)

	assert.Equal(t, "hello world", result)
}

func Test_Apply_WithRuleAfterEmptyPattern_AppliesValidRule(t *testing.T) {
	replacer := NewReplacer()
	rules := []model.Replacement{
		{
			Pattern: "",
			Value:   "broken",
		},
		{
			Pattern: "world",
			Value:   "gopher",
		},
	}

	result := replacer.Apply("hello world", rules)

	assert.Equal(t, "hello gopher", result)
}
