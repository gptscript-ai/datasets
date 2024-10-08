package tools

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"

	"github.com/gptscript-ai/datasets/pkg/dataset"
)

type elem struct {
	Contents       string `json:"contents,omitempty"`
	ContentsBase64 string `json:"contents_base64,omitempty"`
	Type           string `json:"type"`
	Name           string `json:"name"`
	Description    string `json:"description,omitempty"`
}

func GetElement(workspace, datasetID, elementName string) {
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

	elementContents, e, err := d.GetElement(elementName)
	if err != nil {
		fmt.Printf("failed to get element: %v\n", err)
		os.Exit(1)
	}

	element := elem{
		Type:        string(e.Type),
		Name:        e.Name,
		Description: e.Description,
	}

	if e.Type == dataset.DataTypeBytes {
		element.ContentsBase64 = base64.StdEncoding.EncodeToString(elementContents)
	} else {
		element.Contents = string(elementContents)
	}

	elementJSON, err := json.Marshal(element)
	if err != nil {
		fmt.Printf("failed to marshal element: %v\n", err)
		os.Exit(1)
	}

	print(string(elementJSON))
}
