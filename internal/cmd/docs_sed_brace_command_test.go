package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestResolveHeading_AllValues(t *testing.T) {
	assert.Equal(t, "TITLE", resolveHeading("t"))
	assert.Equal(t, "SUBTITLE", resolveHeading("s"))
	assert.Equal(t, "HEADING_1", resolveHeading("1"))
	assert.Equal(t, "HEADING_6", resolveHeading("6"))
	assert.Equal(t, "NORMAL_TEXT", resolveHeading("0"))
	assert.Equal(t, "HEADING_3", resolveHeading("3"))
	assert.Equal(t, "CUSTOM", resolveHeading("CUSTOM"))
}

func TestBuildBraceTextStyleRequests_ImplicitReset(t *testing.T) {
	expression := &braceExpr{Bold: boolPtr(true), Indent: indentNotSet}
	requests := buildBraceTextStyleRequests(expression, 1, 10)
	require.Len(t, requests, 2)
	assert.Contains(t, requests[0].UpdateTextStyle.Fields, "bold")
	assert.Contains(t, requests[0].UpdateTextStyle.Fields, "baselineOffset")
	assert.True(t, requests[1].UpdateTextStyle.TextStyle.Bold)

	additive := &braceExpr{Bold: boolPtr(true), NoReset: true, Indent: indentNotSet}
	additiveRequests := buildBraceTextStyleRequests(additive, 1, 10)
	require.Len(t, additiveRequests, 1)
	assert.True(t, additiveRequests[0].UpdateTextStyle.TextStyle.Bold)
}

func TestClassifyMatch_BraceImage(t *testing.T) {
	expression := sedExpr{
		brace: &braceExpr{
			ImgRef: "https://example.com/img.png",
			Width:  100,
			Height: 50,
			Indent: indentNotSet,
		},
	}
	match := classifyMatch(10, "hello", []int{0, 5}, "hello", "world", expression)
	assert.NotNil(t, match.image)
	assert.Equal(t, "https://example.com/img.png", match.image.URL)
	assert.Equal(t, 100, match.image.Width)
	assert.Equal(t, 50, match.image.Height)
	assert.Equal(t, int64(10), match.start)
	assert.Equal(t, int64(15), match.end)
}

func TestClassifyMatch_PlainText(t *testing.T) {
	match := classifyMatch(0, "foo", []int{0, 3}, "foo", "bar", sedExpr{})
	assert.Nil(t, match.image)
	assert.Equal(t, "bar", match.newText)
}
