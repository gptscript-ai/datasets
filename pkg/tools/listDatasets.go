package tools

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/gptscript-ai/datasets/pkg/dataset"
)

func ListDatasets(workspaceDir string) {
	m, err := dataset.NewManager(workspaceDir)
	if err != nil {
		fmt.Printf("failed to create dataset manager: %v\n", err)
		os.Exit(1)
	}

	datasets, err := m.ListDatasets()
	if err != nil {
		fmt.Printf("failed to list datasets: %v\n", err)
		os.Exit(1)
	}

	datasetsJSON, err := json.Marshal(datasets)
	if err != nil {
		fmt.Printf("failed to marshal datasets: %v\n", err)
		os.Exit(1)
	}

	fmt.Print(string(datasetsJSON))
}
