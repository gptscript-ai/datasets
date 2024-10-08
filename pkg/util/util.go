package util

import (
	"os"
	"path/filepath"
)

// ToFileName converts a name to be alphanumeric plus underscores.
func ToFileName(name string) string {
	for i := 0; i < len(name); i++ {
		if (name[i] < 'a' || name[i] > 'z') && (name[i] < 'A' || name[i] > 'Z') && (name[i] < '0' || name[i] > '9') {
			name = name[:i] + "_" + name[i+1:]
		}
	}
	return name
}

func EnsureUniqueFilename(base, name string) string {
	for {
		if _, err := os.Stat(filepath.Join(base, name)); err == nil {
			name += "_"
		} else if !os.IsNotExist(err) {
			return ""
		} else {
			break
		}
	}
	return name
}
