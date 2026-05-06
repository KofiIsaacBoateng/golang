package greet

import "strings"

// export functions begin with a capital letter. Go automatically understands that.
func Hello(name string) string {
	cleaned := normaliseName(name);

	return "Hello, " + cleaned + ".";
}

// local function begins with a small letter.
func normaliseName(name string) string {
	trimmed := strings.TrimSpace(name);

	if trimmed == "" {
		return "GUEST";
	}

	return strings.ToUpper(trimmed);
}