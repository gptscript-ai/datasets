package tools

import (
	"fmt"
	"os"

	"github.com/gptscript-ai/datasets/pkg/dataset"
)

func GetAllElements(workspace, datasetID string) {
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
	for _, e := range elements {
		eBytes, _, err := d.GetElement(e.Name)
		if err != nil {
			fmt.Printf("failed to get element: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("%s: %s\n", e.Name, string(eBytes))
	}
}
