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
        if e.Field == "command" {
            foundUndefined = true
        }
        if e.Field == "variables" {
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
