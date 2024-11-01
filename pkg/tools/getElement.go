package tools

import (
	"context"
	"encoding/base64"
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

func GetElementLLM(datasetID, elementName string) {
	element := getElement(datasetID, elementName)
	// Attempt to base64 decode the contents.
	rawContents, err := base64.StdEncoding.DecodeString(element.Contents)
	if err != nil {
		// If it's not base64, just use the contents.
		rawContents = []byte(element.Contents)
	}

	fmt.Printf(`{"name": %q, "description": %q, "contents": %q}`, element.Name, element.Description, string(rawContents))
}

func GetElementSDK(datasetID, elementName string) {
	element := getElement(datasetID, elementName)
	elementJSON, err := json.Marshal(element)
	if err != nil {
		fmt.Printf("failed to marshal element: %v\n", err)
		os.Exit(1)
	}

	fmt.Print(string(elementJSON))
}

func getElement(datasetID, elementName string) elem {
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

	return elem{
		Contents:    string(elementContents),
		Name:        e.Name,
		Description: e.Description,
	}
}
