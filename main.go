package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/gptscript-ai/datasets/pkg/dataset"
	"github.com/gptscript-ai/datasets/pkg/tools"
	"github.com/gptscript-ai/go-gptscript"
)

type elementInput struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Contents    string `json:"contents"`
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println(`usage: gptscript-go-tool <command>
subcommands: listDatasets, listElements, getElement, createDataset, addElement
env vars: GPTSCRIPT_WORKSPACE_DIR`)
	}

	switch os.Args[1] {
	case "listDatasets":
		tools.ListDatasets()
	case "listElements":
		tools.ListElements(os.Getenv("DATASETID"))
	case "getElement":
		tools.GetElementLLM(os.Getenv("DATASETID"), os.Getenv("ELEMENT"))
	case "getAllElements":
		tools.GetAllElementsLLM(os.Getenv("DATASETID"))
	case "getElementSDK":
		tools.GetElementSDK(os.Getenv("DATASETID"), os.Getenv("ELEMENT"))
	case "getAllElementsSDK":
		tools.GetAllElementsSDK(os.Getenv("DATASETID"))
	case "createDataset":
		tools.CreateDataset(os.Getenv("DATASETNAME"), os.Getenv("DATASETDESCRIPTION"))
	case "addElement":
		tools.AddElement(os.Getenv("DATASETID"), os.Getenv("ELEMENTNAME"), os.Getenv("ELEMENTDESCRIPTION"), []byte(os.Getenv("ELEMENTCONTENT")))
	case "addElements":
		var elementInputs []elementInput
		if err := json.Unmarshal([]byte(gptscript.GetEnv("ELEMENTS", "")), &elementInputs); err != nil {
			fmt.Printf("failed to unmarshal elements: %v\n", err)
			os.Exit(1)
		}

		addElements(os.Getenv("DATASETID"), elementInputs)
	default:
		fmt.Printf("unknown command: %s\n", os.Args[1])
		os.Exit(1)
	}
}

func addElements(datasetID string, elements []elementInput) {
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

	for _, e := range elements {
		content := []byte(e.Contents)
		_, err := d.AddElement(context.Background(), e.Name, e.Description, content)
		if err != nil {
			fmt.Printf("failed to create element: %v\n", err)
			os.Exit(1)
		}
	}

	fmt.Println("elements added successfully")
}
