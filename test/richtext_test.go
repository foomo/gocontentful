package test

import (
	"github.com/foomo/gocontentful/test/testapi"
	"github.com/stretchr/testify/require"
	"testing"
)

var targetRichText = &testapi.RichTextNode{
	NodeType: "document",
	Content: []interface{}{testapi.RichTextNode{
		NodeType: "heading-1",
		Content: []interface{}{
			testapi.RichTextNodeTextNode{
				NodeType: "text",
				Value:    "A sample page",
				Marks:    []testapi.RichTextMark{},
			},
		},
	},
		testapi.RichTextNode{
			NodeType: "paragraph",
			Content: []interface{}{
				testapi.RichTextNodeTextNode{
					NodeType: "text",
					Value:    "The paragraph ",
					Marks:    []testapi.RichTextMark{},
				},
				testapi.RichTextNodeTextNode{
					NodeType: "text",
					Value:    "with bold stuff",
					Marks: []testapi.RichTextMark{
						{Type: "bold"},
					},
				},
			},
		},
	},
}

func TestHTMLToRichText(t *testing.T) {
	html := `
<html>
	<body>
		<h1>A sample page</h1>
		<p>The paragraph <b>with bold stuff</b></p>
	</body>
</html>
`
	richText := testapi.HtmlToRichText(html)
	require.NotEmpty(t, richText)
	require.EqualValues(t, targetRichText, richText)
}

func TestRichTextToHTML(t *testing.T) {
	html, err := testapi.RichTextToHtml(targetRichText, nil, nil, nil, nil, testapi.SpaceLocaleGerman)
	require.NoError(t, err)
	want := "<h1>A sample page</h1><p>The paragraph <b>with bold stuff</b></p>"
	require.Equal(t, want, html)

}
