package newton

import (
	"os"
	"strings"
	"unicode"
)

type EnvironmentSource struct {
	prefix    string
	separator string
}

func FromEnvironment() *EnvironmentSource {
	return FromEnvironmentCustomSeparator("", "|")
}
func FromEnvironmentWithPrefix(prefix string) *EnvironmentSource {
	return FromEnvironmentCustomSeparator(prefix, "|")
}
func FromEnvironmentCustomSeparator(prefix, separator string) *EnvironmentSource {
	return &EnvironmentSource{prefix: prefix, separator: separator}
}

func (this *EnvironmentSource) Name() string {
	return "environment"
}

func (this *EnvironmentSource) Strings(key string) ([]string, error) {
	key = this.prefix + sanitizeKey(key)

	if value := os.Getenv(key); len(value) > 0 {
		return strings.Split(value, this.separator), nil
	}

	if value := os.Getenv(strings.ToUpper(key)); len(value) > 0 {
		return strings.Split(value, this.separator), nil
	}

	if value := os.Getenv(strings.ToLower(key)); len(value) > 0 {
		return strings.Split(value, this.separator), nil
	}

	return nil, KeyNotFoundError
}
func sanitizeKey(key string) string {
	sanitized := ""

	for _, character := range key {
		if unicode.IsDigit(character) {
			sanitized += string(character)
		} else if unicode.IsLetter(character) {
			sanitized += string(character)
		} else {
			sanitized += "_"
		}
	}

	return sanitized
}
