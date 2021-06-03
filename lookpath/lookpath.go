package lookpath

type LoopPath struct {
	Path    string `yaml:"path,omitempty" json:"path,omitempty"`
	Command string `yaml:"command,omitempty" json:"command,omitempty"`
	Install string `yaml:"install,omitempty" json:"install,omitempty"`
}
