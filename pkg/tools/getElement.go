package tools

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gptscript-ai/datasets/pkg/dataset"
	"github.com/gptscript-ai/datasets/pkg/util"
)

type getElementRequest struct {
	DatasetID string `json:"datasetID"`
	Name      string `json:"name"`
}

func GetElement(w http.ResponseWriter, r *http.Request) {
	var req getElementRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if req.DatasetID == "" {
		http.Error(w, "datasetID is required", http.StatusBadRequest)
		return
	} else if req.Name == "" {
		http.Error(w, "name is required", http.StatusBadRequest)
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
		http.Error(w, fmt.Sprintf("failed to get dataset: %v\n", err), http.StatusInternalServerError)
		return
	}

	element, err := d.GetElement(req.Name)
	if err != nil {
		http.Error(w, "element not found", http.StatusNotFound)
		return
	}

	// Remove the index from the element before returning it to the user.
	eNoIndex := dataset.ElementNoIndex{
		ElementMeta:    element.ElementMeta,
		Contents:       element.Contents,
		BinaryContents: element.BinaryContents,
	}

	if err := json.NewEncoder(w).Encode(eNoIndex); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
