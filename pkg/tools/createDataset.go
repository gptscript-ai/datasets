package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/gptscript-ai/datasets/pkg/dataset"
)

func CreateDataset(name, description string) {
	m, err := dataset.NewManager()
	if err != nil {
		fmt.Printf("failed to create dataset manager: %v\n", err)
		os.Exit(1)
	}

	d, err := m.NewDataset(context.Background(), name, description)
	if err != nil {
		fmt.Printf("failed to create dataset: %v\n", err)
		os.Exit(1)
	}

	datasetJSON, err := json.Marshal(d)
	if err != nil {
		fmt.Printf("failed to marshal dataset: %v\n", err)
		os.Exit(1)
	}

	fmt.Print(string(datasetJSON))
}
