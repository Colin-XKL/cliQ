package cliqfile

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

// Parse parses the YAML content into a TemplateFile struct.
// It does NOT validate the content. Use Validate() for that.
func Parse(data []byte) (*TemplateFile, error) {
	var template TemplateFile
	if err := yaml.Unmarshal(data, &template); err != nil {
		return nil, fmt.Errorf("failed to parse YAML: %w", err)
	}
	return &template, nil
}
