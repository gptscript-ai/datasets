package tools

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/gptscript-ai/datasets/pkg/dataset"
)

func CreateDataset(workspace, name, description string) {
	m, err := dataset.NewManager(workspace)
	if err != nil {
		fmt.Printf("failed to create dataset manager: %v\n", err)
		os.Exit(1)
	}

	d, err := m.NewDataset(name, description)
	if err != nil {
		fmt.Printf("failed to create dataset: %v\n", err)
		os.Exit(1)
	}

	d.BaseDir = "" // This field is an implementation detail that we don't need to care about

	datasetJSON, err := json.Marshal(d)
	if err != nil {
		fmt.Printf("failed to marshal dataset: %v\n", err)
		os.Exit(1)
	}

	fmt.Print(string(datasetJSON))
}
