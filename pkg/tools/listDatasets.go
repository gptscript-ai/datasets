package tools

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/gptscript-ai/datasets/pkg/dataset"
)

type listDatasetsRequest struct {
	WorkspaceID string `json:"workspaceID"`
}

func ListDatasets(w http.ResponseWriter, r *http.Request) {
	var req listDatasetsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if req.WorkspaceID == "" {
		http.Error(w, "workspaceID is required", http.StatusBadRequest)
		return
	}

	m, err := dataset.NewManager(req.WorkspaceID)
	if err != nil {
		fmt.Printf("failed to create dataset manager: %v\n", err)
		os.Exit(1)
	}

	datasets, err := m.ListDatasets(r.Context())
	if err != nil {
		fmt.Printf("failed to list datasets: %v\n", err)
		os.Exit(1)
	}

	if err := json.NewEncoder(w).Encode(datasets); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
