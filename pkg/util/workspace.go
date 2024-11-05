package util

import (
	"fmt"
	"net/http"
	"strings"
)

func GetWorkspaceID(r *http.Request) (string, error) {
	for _, pair := range strings.Split(r.Header.Get("X-GPTScript-Env"), ",") {
		key, value, _ := strings.Cut(pair, "=")
		if key == "GPTSCRIPT_WORKSPACE_ID" {
			return value, nil
		}
	}

	return "", fmt.Errorf("GPTSCRIPT_WORKSPACE_ID not found in environment header")
}
