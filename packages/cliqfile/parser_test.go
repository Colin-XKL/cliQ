package cliqfile

import (
	"reflect"
	"testing"
)

func TestParse(t *testing.T) {
	validYAML := `
name: My Template
description: A description
version: 1.0
author: Me
cliq_template_version: 1.0
cmds:
  - name: cmd1
    description: desc1
    command: echo {{var1}}
    variables:
      - name: var1
        type: string
        label: Var 1
`
	template, err := Parse([]byte(validYAML))
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	expected := &TemplateFile{
		Name:                "My Template",
		Description:         "A description",
		Version:             "1.0",
		Author:              "Me",
		CliqTemplateVersion: "1.0",
		Cmds: []Command{
			{
				Name:        "cmd1",
				Description: "desc1",
				Command:     "echo {{var1}}",
				Variables: []VariableDefinition{
					{
						Name:  "var1",
						Type:  "string",
						Label: "Var 1",
					},
				},
			},
		},
	}

	if !reflect.DeepEqual(template, expected) {
		t.Errorf("Parsed template does not match expected.\nGot: %+v\nExpected: %+v", template, expected)
	}
}

func TestParseInvalidYAML(t *testing.T) {
	invalidYAML := `
name: My Template
  description: Indentation error
`
	_, err := Parse([]byte(invalidYAML))
	if err == nil {
		t.Error("Expected error for invalid YAML, got nil")
	}
}
