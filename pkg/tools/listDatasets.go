package tools

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/gptscript-ai/datasets/pkg/dataset"
	"github.com/gptscript-ai/datasets/pkg/util"
)

func ListDatasets(w http.ResponseWriter, r *http.Request) {
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

	datasets, err := m.ListDatasets(r.Context())
	if err != nil {
		fmt.Printf("failed to list datasets: %v\n", err)
		os.Exit(1)
	}

	if err := json.NewEncoder(w).Encode(datasets); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
