package tools

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/gptscript-ai/datasets/pkg/dataset"
	"github.com/gptscript-ai/datasets/pkg/util"
)

type outputFilter struct {
	Output string `json:"output,omitempty"`
}

var idRegex = regexp.MustCompile(`gds://[a-z0-9]{5}`)

func findDatasetIds(content string) []string {
	return idRegex.FindAllString(content, -1)
}

func OutputFilter(w http.ResponseWriter, r *http.Request) {
	var req outputFilter
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err := w.Write([]byte(req.Output))
	if err != nil {
		return
	}

	datasetIDs := findDatasetIds(req.Output)
	if len(datasetIDs) == 0 {
		return
	}

	workspaceID, err := util.GetWorkspaceID(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	m, err := dataset.NewManager(workspaceID)
	if err != nil {
		fmt.Printf("failed to create dataset manager: %v\n", err)
		os.Exit(1)
	}

	budget := 30_000
outerFor:
	for _, id := range datasetIDs {
		d, err := m.GetDataset(r.Context(), id)
		if err != nil {
			if strings.Contains(err.Error(), "not found") {
				http.Error(w, "dataset not found", http.StatusNotFound)
				return
			}
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		toWrite := d
		toWrite.Elements = make(map[string]dataset.Element)

		for _, element := range d.Elements {
			budget -= len(element.Contents)
			budget -= len(element.BinaryContents)
			if budget < 0 {
				toWrite.Truncated = true
				if err := json.NewEncoder(w).Encode(toWrite); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				data := fmt.Sprintf("\nDataset truncated, %d of %d items not returned\n", len(d.Elements)-len(toWrite.Elements), len(d.Elements))
				_, _ = w.Write([]byte(data))
				continue outerFor
			}
		}

		if err := json.NewEncoder(w).Encode(toWrite); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	return
}
