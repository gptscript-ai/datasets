package tools

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gptscript-ai/datasets/pkg/dataset"
	"github.com/gptscript-ai/datasets/pkg/util"
)

type addElementsRequest struct {
	DatasetID   string            `json:"datasetID"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Elements    []dataset.Element `json:"elements"`
}

func AddElements(w http.ResponseWriter, r *http.Request) {
	var req addElementsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if len(req.Elements) == 0 {
		http.Error(w, "elements is required", http.StatusBadRequest)
		return
	}

	workspaceID, err := util.GetWorkspaceID(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	m, err := dataset.NewManager(workspaceID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var d dataset.Dataset
	if req.DatasetID == "" {
		d, err = m.NewDataset(r.Context(), req.Name, req.Description)
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
