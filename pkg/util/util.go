package util

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"unicode"
)

// ToFileName converts a name to be alphanumeric plus underscores.
func ToFileName(name string) string {
	return strings.Map(func(c rune) rune {
		if !unicode.IsLetter(c) && !unicode.IsDigit(c) {
			return '_'
		}
		return c
	}, name)
}

func EnsureUniqueFilename(base, name string) (string, error) {
	var counter int
	uniqueName := name
	for {
		if _, err := os.Stat(filepath.Join(base, uniqueName)); err == nil {
			counter++
			uniqueName = fmt.Sprintf("%s_%d", name, counter)
		} else if !os.IsNotExist(err) {
			return "", err
		} else {
			return uniqueName, nil
		}
	}
}
