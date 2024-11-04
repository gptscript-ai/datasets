package tools

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gptscript-ai/datasets/pkg/dataset"
)

type getAllElementsRequest struct {
	WorkspaceID string `json:"workspaceID"`
	DatasetID   string `json:"datasetID"`
}

func GetAllElements(w http.ResponseWriter, r *http.Request) {
	var req getAllElementsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if req.WorkspaceID == "" {
		http.Error(w, "workspaceID is required", http.StatusBadRequest)
		return
	} else if req.DatasetID == "" {
		http.Error(w, "datasetID is required", http.StatusBadRequest)
		return
	}

	m, err := dataset.NewManager(req.WorkspaceID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	d, err := m.GetDataset(r.Context(), req.DatasetID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			http.Error(w, "dataset not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	elements := d.GetAllElements()

	if err := json.NewEncoder(w).Encode(elements); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
