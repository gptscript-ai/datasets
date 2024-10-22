package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/gptscript-ai/datasets/pkg/dataset"
)

func AddElement(datasetID, name, description string, content []byte) {
	m, err := dataset.NewManager()
	if err != nil {
		fmt.Printf("failed to create dataset manager: %v\n", err)
		os.Exit(1)
	}

	d, err := m.GetDataset(context.Background(), datasetID)
	if err != nil {
		fmt.Printf("failed to get dataset: %v\n", err)
		os.Exit(1)
	}

	e, err := d.AddElement(context.Background(), name, description, content)
	if err != nil {
		fmt.Printf("failed to create element: %v\n", err)
		os.Exit(1)
	}

	elementJSON, err := json.Marshal(e.ElementMeta)
	if err != nil {
		fmt.Printf("failed to marshal element: %v\n", err)
		os.Exit(1)
	}

	fmt.Print(string(elementJSON))
}
