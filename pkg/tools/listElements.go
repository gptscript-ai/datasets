package tools

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/gptscript-ai/datasets/pkg/dataset"
)

func ListElements(workspace, datasetID string) {
	m, err := dataset.NewManager(workspace)
	if err != nil {
		fmt.Printf("failed to create dataset manager: %v\n", err)
		os.Exit(1)
	}

	d, err := m.GetDataset(datasetID)
	if err != nil {
		fmt.Printf("failed to get dataset: %v\n", err)
		os.Exit(1)
	}

	elements := d.ListElements()
	elementsJSON, err := json.Marshal(elements)
	if err != nil {
		fmt.Printf("failed to marshal elements: %v\n", err)
		os.Exit(1)
	}

	fmt.Print(string(elementsJSON))
}
