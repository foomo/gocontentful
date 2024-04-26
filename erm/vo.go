package erm

// Locale VO
type Locale struct {
	Name         string `json:"name,omitempty"`
	Code         string `json:"code,omitempty"`
	FallbackCode string `json:"fallbackCode,omitempty"`
	Default      bool   `json:"default,omitempty"`
}

// ContentTypeSysAttributes VO
type ContentTypeSysAttributes struct {
	ID       string `json:"id,omitempty"`
	Type     string `json:"type,omitempty"`
	LinkType string `json:"linkType,omitempty"`
}

// ContentTypeSys VO
type ContentTypeSys struct {
	Sys ContentTypeSysAttributes `json:"sys,omitempty"`
}

// ContentfulSys VO
type ContentfulSys struct {
	ID          string         `json:"id,omitempty"`
	Type        string         `json:"type,omitempty"`
	LinkType    string         `json:"linkType,omitempty"`
	ContentType ContentTypeSys `json:"contentType,omitempty"`
	CreatedAt   string         `json:"createdAt,omitempty"`
	UpdatedAt   string         `json:"updatedAt,omitempty"`
	Revision    float64        `json:"revision,omitempty"`
	Version     float64        `json:"version,omitempty"`
}

// ContentTypeFieldItemsValidation VO
type ContentTypeFieldItemsValidation struct {
	LinkContentType []string `json:"linkContentType,omitempty"`
}

// ContentTypeFieldItems VO
type ContentTypeFieldItems struct {
	Type        string                            `json:"type,omitempty"`
	Validations []ContentTypeFieldItemsValidation `json:"validations,omitempty"`
	LinkType    string                            `json:"linkType,omitempty"`
}

// ContentTypeField VO
type ContentTypeField struct {
	ID              string                 `json:"id,omitempty"`
	Name            string                 `json:"name,omitempty"`
	Type            string                 `json:"type,omitempty"`
	Items           *ContentTypeFieldItems `json:"items,omitempty"`
	LinkType        string                 `json:"linkType,omitempty"`
	Omitted         bool                   `json:"omitted,omitempty"`
	ReferencedTypes []string               `json:"referencedTypes,omitempty"`
}

// ContentType VO
type ContentType struct {
	Sys         ContentfulSys      `json:"sys,omitempty"`
	Name        string             `json:"name,omitempty"`
	Description string             `json:"description,omitempty"`
	Fields      []ContentTypeField `json:"fields,omitempty"`
}

type ExportFile struct {
	ContentTypes []ContentType `json:"contentTypes"`
	Locales      []Locale      `json:"locales"`
}
