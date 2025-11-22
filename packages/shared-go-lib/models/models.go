package models

// TemplateFile 表示一个完整的模板文件
type TemplateFile struct {
	// 模板元信息
	Name                string `yaml:"name" json:"name"`
	Description         string `yaml:"description" json:"description"`
	Version             string `yaml:"version" json:"version"`
	Author              string `yaml:"author" json:"author"`
	CliqTemplateVersion string `yaml:"cliq_template_version" json:"cliq_template_version"`

	// 命令列表
	Cmds []Command `yaml:"cmds" json:"cmds"`
}

// Command 表示一个命令模板
type Command struct {
	ID          string               `yaml:"id" json:"id"` // 添加 ID 字段
	Name        string               `yaml:"name" json:"name"`
	Description string               `yaml:"description" json:"description"`
	Command     string               `yaml:"command" json:"command"`
	Variables   []VariableDefinition `yaml:"variables" json:"variables"` // Changed from map to array
}

// VariableDefinition 表示命令中的一个变量定义（扁平化结构）
type VariableDefinition struct {
	Name        string                 `yaml:"name" json:"name"`
	Type        string                 `yaml:"type" json:"type"`
	ArgName     string                 `yaml:"arg_name,omitempty" json:"arg_name,omitempty"`
	Label       string                 `yaml:"label" json:"label"`
	Description string                 `yaml:"description" json:"description"`
	Required    bool                   `yaml:"required" json:"required"`
	Options     map[string]interface{} `yaml:"options,omitempty" json:"options,omitempty"`
}

// 变量类型常量
const (
	VarTypeText       = "string"
	VarTypeFileInput  = "file_input"
	VarTypeFileOutput = "file_output"
	VarTypeBoolean    = "boolean"
	VarTypeNumber     = "number"
	VarTypeSelect     = "select"
)
