package config

type Config struct {
	SpaceID           string   `yaml:"spaceId,omitempty"`
	Environment       string   `yaml:"environment,omitempty"`
	ExportFile        string   `yaml:"exportFile,omitempty"`
	ContentTypes      []string `yaml:"contentTypes,omitempty"`
	PathTargetPackage string   `yaml:"pathTargetPackage,omitempty"`
	RequireVersion    string   `yaml:"requireVersion,omitempty"`
}
