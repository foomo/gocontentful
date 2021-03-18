package templates

import _ "embed"

//go:embed "contentful_vo.tmpl"
var TemplateVo []byte

//go:embed "contentful_vo_base.tmpl"
var TemplateVoBase []byte

//go:embed "contentful_vo_lib.tmpl"
var TemplateVoLib []byte

//go:embed "contentful_vo_lib_contenttype.tmpl"
var TemplateVoLibContentType []byte
