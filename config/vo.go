package config

type Config struct {
	SpaceID           string   `json:"spaceId,omitempty" yaml:"spaceId,omitempty"`
	Environment       string   `json:"environment,omitempty" yaml:"environment,omitempty"`
	ExportFile        string   `json:"exportFile,omitempty" yaml:"exportFile,omitempty"`
	ContentTypes      []string `json:"contentTypes,omitempty" yaml:"contentTypes,omitempty"`
	PathTargetPackage string   `json:"pathTargetPackage,omitempty" yaml:"pathTargetPackage,omitempty"`
	RequireVersion    string   `json:"requireVersion,omitempty" yaml:"requireVersion,omitempty"`
}
