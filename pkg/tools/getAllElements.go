package tools

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/gptscript-ai/datasets/pkg/dataset"
)

func GetAllElementsLLM(datasetID string) {
	elems := getAllElements(datasetID)

	var elemStrings []string
	for _, e := range elems {
		rawContents, err := base64.StdEncoding.DecodeString(e.Contents)
		if err != nil {
			rawContents = []byte(e.Contents)
		}
		elemStrings = append(elemStrings, fmt.Sprintf(`{"name": %q, "description": %q, "contents": %q}`, e.Name, e.Description, string(rawContents)))
	}

	fmt.Printf("[%s]", strings.Join(elemStrings, ","))
}

func GetAllElementsSDK(datasetID string) {
	elems := getAllElements(datasetID)
	elemsJSON, err := json.Marshal(elems)
	if err != nil {
		fmt.Printf("failed to marshal elements: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(string(elemsJSON))
}

func getAllElements(datasetID string) []elem {
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

	return elems
}
