package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/gptscript-ai/datasets/pkg/dataset"
)

func GetAllElements(datasetID string) {
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

	elements := d.ListElements()
	var elems []elem
	for _, e := range elements {
		eBytes, _, err := d.GetElement(context.Background(), e.Name)
		if err != nil {
			fmt.Printf("failed to get element: %v\n", err)
			os.Exit(1)
		}

		elems = append(elems, elem{
			Contents:    string(eBytes),
			Name:        e.Name,
			Description: e.Description,
		})
	}

	elemsJSON, err := json.Marshal(elems)
	if err != nil {
		fmt.Printf("failed to marshal elements: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(string(elemsJSON))
}
