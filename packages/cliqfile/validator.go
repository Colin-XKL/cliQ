package cliqfile

import (
	"fmt"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

type ValidationError struct {
	Line    int    `json:"line"`
	Column  int    `json:"column"`
	Field   string `json:"field"`
	Message string `json:"message"`
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("Line %d: %s - %s", e.Line, e.Field, e.Message)
}

// Validate parses and validates the cliqfile YAML content.
func Validate(data []byte) ([]ValidationError, error) {
	var node yaml.Node
	if err := yaml.Unmarshal(data, &node); err != nil {
		return nil, fmt.Errorf("YAML syntax error: %w", err)
	}

	if node.Kind == 0 && len(node.Content) == 0 {
		return []ValidationError{{Message: "Empty document"}}, nil
	}

	if node.Kind == yaml.DocumentNode {
		if len(node.Content) == 0 {
			return []ValidationError{{Message: "Empty document"}}, nil
		}
		return validateRoot(node.Content[0]), nil
	}
	return validateRoot(&node), nil
}

func validateRoot(node *yaml.Node) []ValidationError {
	var errors []ValidationError
	if node.Kind != yaml.MappingNode {
		errors = append(errors, ValidationError{Line: node.Line, Message: "Root must be a mapping (object)"})
		return errors
	}

	requiredFields := map[string]bool{
		"name":                  false,
		"description":           false,
		"version":               false,
		"author":                false,
		"cliq_template_version": false,
		"cmds":                  false,
	}

	// Iterate over keys
	for i := 0; i < len(node.Content); i += 2 {
		if i+1 >= len(node.Content) {
			break
		}
		keyNode := node.Content[i]
		valNode := node.Content[i+1]
		key := keyNode.Value

		if _, ok := requiredFields[key]; ok {
			requiredFields[key] = true
			if valNode.Value == "" && valNode.Kind == yaml.ScalarNode {
				errors = append(errors, ValidationError{Line: keyNode.Line, Field: key, Message: "Field cannot be empty"})
			}
		}

		if key == "cmds" {
			cmdErrors := validateCmds(valNode)
			errors = append(errors, cmdErrors...)
		}
	}

	for field, found := range requiredFields {
		if !found {
			errors = append(errors, ValidationError{Line: node.Line, Field: field, Message: "Missing required field"})
		}
	}

	return errors
}

func validateCmds(node *yaml.Node) []ValidationError {
	var errors []ValidationError
	if node.Kind != yaml.SequenceNode {
		errors = append(errors, ValidationError{Line: node.Line, Field: "cmds", Message: "cmds must be a list"})
		return errors
	}

	if len(node.Content) == 0 {
		errors = append(errors, ValidationError{Line: node.Line, Field: "cmds", Message: "cmds list cannot be empty"})
	}

	for _, cmdNode := range node.Content {
		errors = append(errors, validateCommand(cmdNode)...)
	}
	return errors
}

func validateCommand(node *yaml.Node) []ValidationError {
	var errors []ValidationError
	if node.Kind != yaml.MappingNode {
		errors = append(errors, ValidationError{Line: node.Line, Message: "Command must be an object"})
		return errors
	}

	requiredFields := map[string]bool{
		"name":        false,
		"description": false,
		"command":     false,
	}

	var commandStr string
	var varsNode *yaml.Node

	for i := 0; i < len(node.Content); i += 2 {
		if i+1 >= len(node.Content) {
			break
		}
		keyNode := node.Content[i]
		valNode := node.Content[i+1]
		key := keyNode.Value

		if _, ok := requiredFields[key]; ok {
			requiredFields[key] = true
			if valNode.Value == "" && valNode.Kind == yaml.ScalarNode {
				errors = append(errors, ValidationError{Line: keyNode.Line, Field: key, Message: "Field cannot be empty"})
			}
		}

		if key == "command" {
			commandStr = valNode.Value
		}
		if key == "variables" {
			varsNode = valNode
		}
	}

	for field, found := range requiredFields {
		if !found {
			errors = append(errors, ValidationError{Line: node.Line, Field: field, Message: "Missing required field in command"})
		}
	}

	definedVars := make(map[string]bool)
	if varsNode != nil {
		varErrors, dVars := validateVariables(varsNode)
		errors = append(errors, varErrors...)
		definedVars = dVars
	}

	if commandStr != "" {
		// Check usage
		usedVars := extractVariables(commandStr)
		for _, v := range usedVars {
			if !definedVars[v] {
				errors = append(errors, ValidationError{Line: node.Line, Field: "command", Message: fmt.Sprintf("Variable '%s' used in command string but not defined in variables", v)})
			}
		}
		// Check unused vars
		for v := range definedVars {
			found := false
			for _, uv := range usedVars {
				if uv == v {
					found = true
					break
				}
			}
			if !found {
				// Find line for this variable? It's hard to point to the variable definition here without storing more info.
				// We'll point to the command block start for now.
				errors = append(errors, ValidationError{Line: node.Line, Field: "variables", Message: fmt.Sprintf("Variable '%s' defined but not used in command string", v)})
			}
		}
	}

	return errors
}

func validateVariables(node *yaml.Node) ([]ValidationError, map[string]bool) {
	var errors []ValidationError
	definedVars := make(map[string]bool)

	if node.Kind != yaml.SequenceNode {
		errors = append(errors, ValidationError{Line: node.Line, Field: "variables", Message: "variables must be a list"})
		return errors, definedVars
	}

	for _, varNode := range node.Content {
		vErrors, name := validateVariable(varNode)
		errors = append(errors, vErrors...)
		if name != "" {
			if definedVars[name] {
				errors = append(errors, ValidationError{Line: varNode.Line, Field: "name", Message: fmt.Sprintf("Duplicate variable name '%s'", name)})
			}
			definedVars[name] = true
		}
	}
	return errors, definedVars
}

func validateVariable(node *yaml.Node) ([]ValidationError, string) {
	var errors []ValidationError
	name := ""

	if node.Kind != yaml.MappingNode {
		errors = append(errors, ValidationError{Line: node.Line, Message: "Variable must be an object"})
		return errors, name
	}

	requiredFields := map[string]bool{
		"name":  false,
		"type":  false,
		"label": false,
	}

	var typeVal string

	for i := 0; i < len(node.Content); i += 2 {
		if i+1 >= len(node.Content) {
			break
		}
		keyNode := node.Content[i]
		valNode := node.Content[i+1]
		key := keyNode.Value

		if _, ok := requiredFields[key]; ok {
			requiredFields[key] = true
			if valNode.Value == "" && valNode.Kind == yaml.ScalarNode {
				errors = append(errors, ValidationError{Line: keyNode.Line, Field: key, Message: "Field cannot be empty"})
			}
		}

		if key == "name" {
			name = valNode.Value
		}
		if key == "type" {
			typeVal = valNode.Value
		}
	}

	for field, found := range requiredFields {
		if !found {
			errors = append(errors, ValidationError{Line: node.Line, Field: field, Message: "Missing required field in variable"})
		}
	}

	if typeVal != "" {
		allowed := map[string]bool{
			VarTypeText: true, VarTypeFileInput: true, VarTypeFileOutput: true,
			VarTypeNumber: true, VarTypeBoolean: true, VarTypeSelect: true,
		}
		if !allowed[typeVal] {
			errors = append(errors, ValidationError{Line: node.Line, Field: "type", Message: fmt.Sprintf("Invalid type '%s'", typeVal)})
		}
	}

	return errors, name
}

func extractVariables(commandStr string) []string {
	// Basic regex for {{variable_name}}
	re := regexp.MustCompile(`\{\{\s*([a-zA-Z0-9_-]+)\s*\}\}`)
	matches := re.FindAllStringSubmatch(commandStr, -1)
	var vars []string
	seen := make(map[string]bool)
	for _, match := range matches {
		if len(match) > 1 {
			v := strings.TrimSpace(match[1])
			if !seen[v] {
				vars = append(vars, v)
				seen[v] = true
			}
		}
	}
	return vars
}
