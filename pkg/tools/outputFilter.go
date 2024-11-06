package tools

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"sort"
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

		output := struct {
			ID          string            `json:"id,omitempty"`
			Name        string            `json:"name,omitempty"`
			Description string            `json:"description,omitempty"`
			Items       []dataset.Element `json:"items,omitempty"`
			Length      int               `json:"length,omitempty"`
			Truncated   bool              `json:"truncated,omitempty"`
		}{
			ID:          d.ID,
			Name:        d.Name,
			Description: d.Description,
			Length:      len(d.Elements),
		}

		var elementList []dataset.Element
		for _, element := range d.Elements {
			elementList = append(elementList, element)
		}
		sort.Slice(elementList, func(i, j int) bool {
			return elementList[i].Index < elementList[j].Index
		})

		for _, element := range elementList {
			budget -= len(element.Contents)
			budget -= len(element.BinaryContents)
			if budget < 0 {
				output.Truncated = true
				if err := json.NewEncoder(w).Encode(output); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				data := fmt.Sprintf("\nDataset %s truncated, %d of %d items not returned\n", id, len(elementList)-len(output.Items), len(elementList))
				_, _ = w.Write([]byte(data))
				continue outerFor
			}
			output.Items = append(output.Items, element)
		}

		if err := json.NewEncoder(w).Encode(output); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	return
}
