package yaml

import (
    "encoding/base64"
    "regexp"
    "strings"

    "gopkg.in/yaml.v3"

    "repo/shared-go-lib/models"
)

var thinkTagRegex = regexp.MustCompile(`(?s){{think}}.*?{{/think}}`)

func StripThinkTags(s string) string {
	return thinkTagRegex.ReplaceAllString(s, "")
}

var fencedBlockRegex = regexp.MustCompile("(?s)^\\s*" + "```" + "(?:yaml|yml)?\\s*\\n(.*?)\\n\\s*" + "```" + "\\s*$")
func StripFences(s string) string {
	out := StripThinkTags(s)
	out = strings.TrimSpace(out)

	// Extract content only if the entire string is a single fenced code block.
	if m := fencedBlockRegex.FindStringSubmatch(out); m != nil && len(m) > 1 {
		return strings.TrimSpace(m[1])
	}
	return out
}

func UnmarshalTemplate(s string) (*models.TemplateFile, error) {
    var t models.TemplateFile
    if err := yaml.Unmarshal([]byte(s), &t); err != nil {
        return nil, err
    }
    return &t, nil
}

func MarshalTemplate(t *models.TemplateFile) (string, error) {
    b, err := yaml.Marshal(t)
    if err != nil {
        return "", err
    }
    return string(b), nil
}

func Base64Encode(s string) string {
	return base64.StdEncoding.EncodeToString([]byte(s))
}
