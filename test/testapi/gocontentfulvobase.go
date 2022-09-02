// Code generated by https://github.com/foomo/gocontentful v1.0.12-beta - DO NOT EDIT.
package testapi

import "github.com/foomo/contentful"

type ContentTypeSysAttributes struct {
	ID       string `json:"id,omitempty"`
	Type     string `json:"type,omitempty"`
	LinkType string `json:"linkType,omitempty"`
}

type ContentTypeSys struct {
	Sys ContentTypeSysAttributes `json:"sys,omitempty"`
}

type ContentfulSys struct {
	ID               string         `json:"id,omitempty"`
	Type             string         `json:"type,omitempty"`
	LinkType         string         `json:"linkType,omitempty"`
	ContentType      ContentTypeSys `json:"contentType,omitempty"`
	Environment      ContentTypeSys `json:"environment,omitempty"`
	Space            ContentTypeSys `json:"space,omitempty"`
	CreatedAt        string         `json:"createdAt,omitempty"`
	UpdatedAt        string         `json:"updatedAt,omitempty"`
	Revision         float64        `json:"revision,omitempty"`
	Version          float64        `json:"version,omitempty"`
	PublishedCounter float64        `json:"publishedCounter,omitempty"`
	PublishedVersion float64        `json:"publishedVersion,omitempty"`
}

type ContentfulReferencedEntry struct {
	Entry *contentful.Entry
	Col   *contentful.Collection
	LogFn func(
		contentType string,
		entryID string,
		method string,
		err error,
	)
}

type ContentTypeFieldLocation struct {
	Lat float64 `json:"lat,omitempty"`
	Lon float64 `json:"lon,omitempty"`
}

type RichTextNode struct {
	NodeType string        `json:"nodeType"`
	Content  []interface{} `json:"content"`
	Data     RichTextData  `json:"data"`
}

type RichTextNodeTextNode struct {
	NodeType string         `json:"nodeType"`
	Data     RichTextData   `json:"data"`
	Value    string         `json:"value"`
	Marks    []RichTextMark `json:"marks"`
}

type RichTextData struct {
	URI    string          `json:"uri,omitempty"`
	Target *ContentTypeSys `json:"target,omitempty"`
}

type RichTextGenericNode struct {
	NodeType string                 `json:"nodeType"`
	Content  []*RichTextGenericNode `json:"content"`
	Data     map[string]interface{} `json:"data"`
	Value    string                 `json:"value"`
	Marks    []RichTextMark         `json:"marks"`
}

type richTextHtmlTag struct {
	attrs      map[string]string
	name       string
	customHTML string
}

type richTextHtmlTags []richTextHtmlTag

type RichTextMark struct {
	Type string `json:"type,omitempty"`
}

const (
	FieldTypeLink      string = "Link"
	FieldLinkTypeEntry string = "Entry"
	FieldLinkTypeAsset string = "Asset"

	HtmlHeading1        string = "h1"
	HtmlHeading2        string = "h2"
	HtmlHeading3        string = "h3"
	HtmlHeading4        string = "h4"
	HtmlHeading5        string = "h5"
	HtmlHeading6        string = "h6"
	HtmlParagraph       string = "p"
	HtmlItalic          string = "i"
	HtmlEm              string = "em"
	HtmlBold            string = "b"
	HtmlStrong          string = "strong"
	HtmlUnderline       string = "u"
	HtmlAnchor          string = "a"
	HtmlImage           string = "img"
	HtmlBlockquote      string = "blockquote"
	HtmlCode            string = "code"
	HtmlUnorderedList   string = "ul"
	HtmlOrderedList     string = "ol"
	HtmlListItem        string = "li"
	HtmlBreak           string = "br"
	HtmlHorizontalRule  string = "hr"
	HtmlAttributeHref   string = "href"
	HtmlTable           string = "table"
	HtmlTableRow        string = "tr"
	HtmlTableHeaderCell string = "th"
	HtmlTableCell       string = "td"

	RichTextNodeDocument        string = "document"
	RichTextNodeParagraph       string = "paragraph"
	RichTextNodeHeading1        string = "heading-1"
	RichTextNodeHeading2        string = "heading-2"
	RichTextNodeHeading3        string = "heading-3"
	RichTextNodeHeading4        string = "heading-4"
	RichTextNodeHeading5        string = "heading-5"
	RichTextNodeHeading6        string = "heading-6"
	RichTextNodeHyperlink       string = "hyperlink"
	RichTextNodeEntryHyperlink  string = "entry-hyperlink"
	RichTextNodeAssetHyperlink  string = "asset-hyperlink"
	RichTextNodeEmbeddedAsset   string = "embedded-asset-block"
	RichTextNodeEmbeddedEntry   string = "embedded-entry-block"
	RichTextNodeText            string = "text"
	RichTextNodeUnorderedList   string = "unordered-list"
	RichTextNodeOrderedList     string = "ordered-list"
	RichTextNodeListItem        string = "list-item"
	RichTextNodeBlockquote      string = "blockquote"
	RichTextNodeHR              string = "hr"
	RichTextNodeTable           string = "table"
	RichTextNodeTableRow        string = "table-row"
	RichTextNodeTableHeaderCell string = "table-header-cell"
	RichTextNodeTableCell       string = "table-cell"

	RichTextMarkBold      string = "bold"
	RichTextMarkItalic    string = "italic"
	RichTextMarkUnderline string = "underline"
	RichTextMarkCode      string = "code"
)

const (
	StatusDraft     = "draft"
	StatusChanged   = "changed"
	StatusPublished = "published"
)
