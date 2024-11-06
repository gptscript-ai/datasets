package util

import (
	"fmt"
	"net/http"
	"strings"
)

func GetWorkspaceID(r *http.Request) (string, error) {
	for _, kv := range r.Header.Values("X-GPTScript-Env") {
		if value, ok := strings.CutPrefix(kv, "GPTSCRIPT_WORKSPACE_ID="); ok {
			return value, nil
		}
	}

	return "", fmt.Errorf("GPTSCRIPT_WORKSPACE_ID not found in environment header")
}
