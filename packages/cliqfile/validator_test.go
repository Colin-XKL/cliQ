package cliqfile

import (
	"testing"
)

func TestValidate(t *testing.T) {
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
	errors, err := Validate([]byte(validYAML))
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if len(errors) > 0 {
		t.Errorf("Expected no errors, got %v", errors)
	}

	invalidYAML := `
name:
description:
version:
author:
cliq_template_version:
cmds: []
`
	errors, err = Validate([]byte(invalidYAML))
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if len(errors) == 0 {
		t.Error("Expected errors, got none")
	}
}

func TestValidateRoot(t *testing.T) {
	tests := []struct {
		name          string
		yaml          string
		expectedError string
	}{
		{
			name: "Empty document",
			yaml: "",
			expectedError: "Empty document",
		},
		{
			name: "Root not a mapping",
			yaml: "- item1",
			expectedError: "Root must be a mapping (object)",
		},
		{
			name: "Missing required fields",
			yaml: `
cmds:
  - name: cmd1
    description: desc1
    command: echo hello
`,
			expectedError: "Missing required field", // Will match multiple errors
		},
		{
			name: "Empty required field",
			yaml: `
name: ""
description: A description
version: 1.0
author: Me
cliq_template_version: 1.0
cmds:
  - name: cmd1
    description: desc1
    command: echo hello
`,
			expectedError: "Field cannot be empty",
		},
		{
			name: "Cmds not a list",
			yaml: `
name: My Template
description: A description
version: 1.0
author: Me
cliq_template_version: 1.0
cmds: "not a list"
`,
			expectedError: "cmds must be a list",
		},
		{
			name: "Cmds list empty",
			yaml: `
name: My Template
description: A description
version: 1.0
author: Me
cliq_template_version: 1.0
cmds: []
`,
			expectedError: "cmds list cannot be empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errors, err := Validate([]byte(tt.yaml))
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
			found := false
			for _, e := range errors {
				if e.Message == tt.expectedError {
					found = true
					break
				}
				// For "Missing required field", it might be multiple, just check if any matches
				if tt.expectedError == "Missing required field" && e.Message == "Missing required field" {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("Expected error containing '%s', got %v", tt.expectedError, errors)
			}
		})
	}
}

func TestValidateCommand(t *testing.T) {
	tests := []struct {
		name          string
		yaml          string
		expectedError string
	}{
		{
			name: "Command not an object",
			yaml: `
name: My Template
description: A description
version: 1.0
author: Me
cliq_template_version: 1.0
cmds:
  - "not an object"
`,
			expectedError: "Command must be an object",
		},
		{
			name: "Missing command fields",
			yaml: `
name: My Template
description: A description
version: 1.0
author: Me
cliq_template_version: 1.0
cmds:
  - name: cmd1
`,
			expectedError: "Missing required field in command",
		},
		{
			name: "Empty command field",
			yaml: `
name: My Template
description: A description
version: 1.0
author: Me
cliq_template_version: 1.0
cmds:
  - name: cmd1
    description: desc1
    command: ""
`,
			expectedError: "Field cannot be empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errors, err := Validate([]byte(tt.yaml))
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
			found := false
			for _, e := range errors {
				if e.Message == tt.expectedError {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("Expected error containing '%s', got %v", tt.expectedError, errors)
			}
		})
	}
}

func TestValidateVariables(t *testing.T) {
	tests := []struct {
		name          string
		yaml          string
		expectedError string
	}{
		{
			name: "Variables not a list",
			yaml: `
name: My Template
description: A description
version: 1.0
author: Me
cliq_template_version: 1.0
cmds:
  - name: cmd1
    description: desc1
    command: echo hello
    variables: "not a list"
`,
			expectedError: "variables must be a list",
		},
		{
			name: "Variable not an object",
			yaml: `
name: My Template
description: A description
version: 1.0
author: Me
cliq_template_version: 1.0
cmds:
  - name: cmd1
    description: desc1
    command: echo hello
    variables:
      - "not an object"
`,
			expectedError: "Variable must be an object",
		},
		{
			name: "Missing variable fields",
			yaml: `
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
`,
			expectedError: "Missing required field in variable",
		},
		{
			name: "Invalid variable type",
			yaml: `
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
        type: invalid_type
        label: Var 1
`,
			expectedError: "Invalid type 'invalid_type'",
		},
		{
			name: "Duplicate variable name",
			yaml: `
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
      - name: var1
        type: number
        label: Var 1 Duplicate
`,
			expectedError: "Duplicate variable name 'var1'",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errors, err := Validate([]byte(tt.yaml))
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
			found := false
			for _, e := range errors {
				if e.Message == tt.expectedError {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("Expected error containing '%s', got %v", tt.expectedError, errors)
			}
		})
	}
}

func TestValidateVariableUsage(t *testing.T) {
    yamlStr := `
name: Test
description: Test
version: 1.0
author: Test
cliq_template_version: 1.0
cmds:
  - name: test
    description: test
    command: echo {{unused_var}}
    variables:
      - name: defined_var
        type: string
        label: Defined
`
    errors, err := Validate([]byte(yamlStr))
    if err != nil {
        t.Fatalf("Unexpected error: %v", err)
    }
    foundUnused := false
    foundUndefined := false

    for _, e := range errors {
        if e.Field == "command" && e.Message == "Variable 'unused_var' used in command string but not defined in variables" {
            foundUndefined = true
        }
        if e.Field == "variables" && e.Message == "Variable 'defined_var' defined but not used in command string" {
            foundUnused = true
        }
    }

    if !foundUndefined {
        t.Error("Expected error about undefined variable usage")
    }
    if !foundUnused {
        t.Error("Expected error about unused variable definition")
    }
}

func TestValidateNoVariables(t *testing.T) {
	validYAML := `
name: My Template
description: A description
version: 1.0
author: Me
cliq_template_version: 1.0
cmds:
  - name: cmd1
    description: desc1
    command: echo hello
`
	errors, err := Validate([]byte(validYAML))
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if len(errors) > 0 {
		t.Errorf("Expected no errors, got %v", errors)
	}

    // Also test with variables: []
    validYAML2 := `
name: My Template
description: A description
version: 1.0
author: Me
cliq_template_version: 1.0
cmds:
  - name: cmd1
    description: desc1
    command: echo hello
    variables: []
`
	errors, err = Validate([]byte(validYAML2))
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if len(errors) > 0 {
		t.Errorf("Expected no errors, got %v", errors)
	}
}
