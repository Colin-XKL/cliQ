package template

type Template struct {
    Name               string  `yaml:"name"`
    Description        string  `yaml:"description"`
    Version            string  `yaml:"version"`
    Author             string  `yaml:"author"`
    CliqTemplateVersion string `yaml:"cliq_template_version"`
    Cmds               []Cmd   `yaml:"cmds"`
}

type Cmd struct {
    Name        string     `yaml:"name"`
    Description string     `yaml:"description"`
    Command     string     `yaml:"command"`
    Variables   []Variable `yaml:"variables"`
}

type Variable struct {
    Name        string                 `yaml:"name"`
    Type        string                 `yaml:"type"`
    Label       string                 `yaml:"label"`
    Description string                 `yaml:"description"`
    Required    bool                   `yaml:"required"`
    Options     map[string]interface{} `yaml:"options,omitempty"`
}