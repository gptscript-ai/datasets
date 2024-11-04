package tools

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gptscript-ai/datasets/pkg/dataset"
)

type listElementsRequest struct {
	WorkspaceID string `json:"workspaceID"`
	DatasetID   string `json:"datasetID"`
}

func ListElements(w http.ResponseWriter, r *http.Request) {
	var req listElementsRequest
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
		fmt.Printf("failed to create dataset manager: %v\n", err)
		os.Exit(1)
	}

	d, err := m.GetDataset(r.Context(), req.DatasetID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			http.Error(w, "dataset not found", http.StatusNotFound)
			return
		}
		fmt.Printf("failed to get dataset: %v\n", err)
		os.Exit(1)
	}

	if err := json.NewEncoder(w).Encode(d.ListElements()); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
