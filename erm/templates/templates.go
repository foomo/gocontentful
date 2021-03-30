package templates

import _ "embed"

//go:embed "contentful_vo.gotmpl"
var TemplateVo []byte

//go:embed "contentful_vo_base.gotmpl"
var TemplateVoBase []byte

//go:embed "contentful_vo_lib.gotmpl"
var TemplateVoLib []byte

//go:embed "contentful_vo_lib_contenttype.gotmpl"
var TemplateVoLibContentType []byte
