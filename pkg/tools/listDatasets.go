package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/gptscript-ai/datasets/pkg/dataset"
)

func ListDatasets() {
	m, err := dataset.NewManager()
	if err != nil {
		fmt.Printf("failed to create dataset manager: %v\n", err)
		os.Exit(1)
	}

	datasets, err := m.ListDatasets(context.Background())
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
