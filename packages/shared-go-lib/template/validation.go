package template

import (
    "fmt"
    "strings"

    "repo/shared-go-lib/models"
)

var allowedTypes = map[string]struct{}{
    "string": {},
    "file_input": {},
    "file_output": {},
    "number": {},
    "boolean": {},
    "select": {},
}

func ValidateTemplate(t *models.TemplateFile) error {
    if t.Name == "" || t.Description == "" || t.Version == "" || t.Author == "" || t.CliqTemplateVersion == "" {
        return fmt.Errorf("missing required metadata fields")
    }
    if len(t.Cmds) == 0 {
        return fmt.Errorf("cmds must contain at least one command")
    }
    for _, c := range t.Cmds {
        if c.Name == "" || c.Description == "" || c.Command == "" {
            return fmt.Errorf("command '%s' missing required fields", c.Name)
        }
        if len(c.Variables) == 0 {
            return fmt.Errorf("command '%s' must define variables", c.Name)
        }
        names := map[string]struct{}{}
        for _, v := range c.Variables {
            if v.Name == "" || v.Type == "" || v.Label == "" {
                return fmt.Errorf("variable missing required fields in command '%s'", c.Name)
            }
            if _, ok := allowedTypes[v.Type]; !ok {
                return fmt.Errorf("variable '%s' has unsupported type '%s'", v.Name, v.Type)
            }
            if _, dup := names[v.Name]; dup {
                return fmt.Errorf("duplicate variable name '%s'", v.Name)
            }
            names[v.Name] = struct{}{}
        }
        // placeholder consistency: each var should appear in command string
        for name := range names {
            ph := "{{" + name + "}}"
            if !strings.Contains(c.Command, ph) {
                return fmt.Errorf("variable '%s' not referenced in command", name)
            }
        }
    }
    return nil
}