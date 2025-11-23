package cliqfile

import (
	"encoding/base64"
	"regexp"
	"strings"
)

var thinkTagRegex = regexp.MustCompile(`(?s){{think}}.*?{{/think}}`)

func StripThinkTags(s string) string {
	return thinkTagRegex.ReplaceAllString(s, "")
}

var fencedBlockRegex = regexp.MustCompile("(?s)^\\s*```(?:yaml|yml)?\\s*\\n(.*?)\\n\\s*```\\s*$")

func StripFences(s string) string {
	out := StripThinkTags(s)
	out = strings.TrimSpace(out)

	// Extract content only if the entire string is a single fenced code block.
	if m := fencedBlockRegex.FindStringSubmatch(out); m != nil && len(m) > 1 {
		return strings.TrimSpace(m[1])
	}
	return out
}

func Base64Encode(s string) string {
	return base64.StdEncoding.EncodeToString([]byte(s))
}
