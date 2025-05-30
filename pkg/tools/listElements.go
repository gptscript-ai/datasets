package tools

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gptscript-ai/datasets/pkg/dataset"
	"github.com/gptscript-ai/datasets/pkg/util"
)

type listElementsRequest struct {
	DatasetID string `json:"datasetID"`
}

func ListElements(w http.ResponseWriter, r *http.Request) {
	var req listElementsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if req.DatasetID == "" {
		http.Error(w, "datasetID is required", http.StatusBadRequest)
		return
	}

	workspaceID, err := util.GetWorkspaceID(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	m, err := dataset.NewManager(workspaceID)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to create dataset manager: %v\n", err), http.StatusInternalServerError)
		return
	}

	d, err := m.GetDataset(r.Context(), req.DatasetID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			http.Error(w, "dataset not found", http.StatusNotFound)
			return
		}
		http.Error(w, fmt.Sprintf("failed to get dataset: %v\n", err), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(d.ListElements()); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
