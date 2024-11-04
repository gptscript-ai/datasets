package tools

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gptscript-ai/datasets/pkg/dataset"
)

type addElementsRequest struct {
	WorkspaceID string            `json:"workspaceID"`
	DatasetID   string            `json:"datasetID"`
	Elements    []dataset.Element `json:"elements"`
}

func AddElements(w http.ResponseWriter, r *http.Request) {
	var req addElementsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if req.WorkspaceID == "" {
		http.Error(w, "workspaceID is required", http.StatusBadRequest)
	} else if len(req.Elements) == 0 {
		http.Error(w, "elements is required", http.StatusBadRequest)
		return
	}

	m, err := dataset.NewManager(req.WorkspaceID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var d dataset.Dataset
	if req.DatasetID == "" {
		d, err = m.NewDataset(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		d, err = m.GetDataset(r.Context(), req.DatasetID)
		if err != nil {
			if strings.Contains(err.Error(), "not found") {
				http.Error(w, "dataset not found", http.StatusNotFound)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	for _, element := range req.Elements {
		if err := d.AddElement(element); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	if err := d.Save(r.Context()); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if _, err = w.Write([]byte(d.ID)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
