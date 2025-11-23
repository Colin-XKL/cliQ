package cliqfile

import (
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

// GenerateFromCommand creates a TemplateFile from a command string.
func GenerateFromCommand(commandStr string) (*TemplateFile, error) {
	if commandStr == "" {
		return nil, fmt.Errorf("command string cannot be empty")
	}

	// Extract variables from command string
	variables := extractVariables(commandStr) // using the function from validator.go

	// Generate template
	templateFile := &TemplateFile{
		Name:                "Generated Template",
		Description:         "Automatically generated template from command",
		Version:             "1.0",
		Author:              "cliQ",
		CliqTemplateVersion: "1.0",
		Cmds: []Command{
			{
				ID:          "generated_cmd_1",
				Name:        "Generated Command",
				Description: "Automatically generated command",
				Command:     commandStr,
				Variables:   []VariableDefinition{},
			},
		},
	}

	// Create variable definitions
	for _, varName := range variables {
		varType := determineVariableType(varName)

		varDef := VariableDefinition{
			Name:        varName,
			Type:        varType,
			Label:       getLabelFromVariableName(varName),
			Description: fmt.Sprintf("The %s parameter", varName),
			Required:    true,
		}

		switch varType {
		case VarTypeFileInput, VarTypeFileOutput:
			varDef.Options = map[string]interface{}{
				"file_types": []string{".*"},
			}
		case VarTypeNumber:
			varDef.Options = map[string]interface{}{
				"default": 1,
				"min":     0,
				"max":     100,
			}
		case VarTypeBoolean:
			varDef.Options = map[string]interface{}{
				"default": false,
			}
		}

		templateFile.Cmds[0].Variables = append(templateFile.Cmds[0].Variables, varDef)
	}

	return templateFile, nil
}

// GenerateYAML serializes the TemplateFile to YAML string.
func GenerateYAML(template *TemplateFile) (string, error) {
	if template == nil {
		return "", fmt.Errorf("template cannot be nil")
	}

	data, err := yaml.Marshal(template)
	if err != nil {
		return "", fmt.Errorf("failed to marshal template: %w", err)
	}

	yamlStr := string(data)
	if yamlStr == "" || yamlStr == "{}\n" {
		return "", fmt.Errorf("generated YAML is empty or invalid")
	}

	return yamlStr, nil
}

func determineVariableType(varName string) string {
	if strings.HasSuffix(varName, "_file") || strings.HasSuffix(varName, "_path") ||
		strings.Contains(varName, "file") || strings.Contains(varName, "path") {
		if strings.Contains(varName, "input") || strings.Contains(varName, "src") {
			return VarTypeFileInput
		} else if strings.Contains(varName, "output") || strings.Contains(varName, "dest") {
			return VarTypeFileOutput
		} else {
			return VarTypeFileInput
		}
	} else if strings.Contains(varName, "number") || strings.Contains(varName, "size") ||
		strings.Contains(varName, "width") || strings.Contains(varName, "height") {
		return VarTypeNumber
	} else if strings.Contains(varName, "enable") || strings.Contains(varName, "use") ||
		strings.Contains(varName, "flag") || strings.Contains(varName, "show") {
		return VarTypeBoolean
	} else {
		return VarTypeText
	}
}

func getLabelFromVariableName(varName string) string {
	label := strings.ReplaceAll(varName, "_", " ")
	label = strings.ReplaceAll(label, "-", " ")
	if len(label) > 0 {
		label = strings.ToUpper(string(label[0])) + label[1:]
	}
	return label
}
