package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/gptscript-ai/datasets/pkg/dataset"
)

type elem struct {
	Contents    string `json:"contents,omitempty"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

func GetElement(datasetID, elementName string) {
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

	elementContents, e, err := d.GetElement(context.Background(), elementName)
	if err != nil {
		fmt.Printf("failed to get element: %v\n", err)
		os.Exit(1)
	}

	element := elem{
		Contents:    string(elementContents),
		Name:        e.Name,
		Description: e.Description,
	}

	elementJSON, err := json.Marshal(element)
	if err != nil {
		fmt.Printf("failed to marshal element: %v\n", err)
		os.Exit(1)
	}

	fmt.Print(string(elementJSON))
}
