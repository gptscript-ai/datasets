package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/gptscript-ai/datasets/pkg/dataset"
	"github.com/gptscript-ai/datasets/pkg/tools"
)

type elementInput struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Content     string `json:"content"`
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println(`usage: gptscript-go-tool <command>
subcommands: listDatasets, listElements, getElement, createDataset, addElement
env vars: GPTSCRIPT_WORKSPACE_DIR`)
	}

	workspace := os.Getenv("GPTSCRIPT_WORKSPACE_DIR")
	if workspace == "" {
		fmt.Println("missing GPTSCRIPT_WORKSPACE_DIR")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "listDatasets":
		tools.ListDatasets(workspace)
	case "listElements":
		tools.ListElements(workspace, os.Getenv("DATASETID"))
	case "getElement":
		tools.GetElement(workspace, os.Getenv("DATASETID"), os.Getenv("ELEMENT"))
	case "createDataset":
		tools.CreateDataset(workspace, os.Getenv("DATASETNAME"), os.Getenv("DATASETDESCRIPTION"))
	case "addElement":
		tools.AddElement(workspace, os.Getenv("DATASETID"), os.Getenv("ELEMENTNAME"), os.Getenv("ELEMENTDESCRIPTION"), []byte(os.Getenv("ELEMENTCONTENT")))
	case "addElements":
		var elements []elementInput
		if err := json.Unmarshal([]byte(os.Getenv("ELEMENTS")), &elements); err != nil {
			fmt.Printf("failed to unmarshal elements: %v\n", err)
			os.Exit(1)
		}
		addElements(workspace, os.Getenv("DATASETID"), elements)
	case "getAllElements":
		tools.GetAllElements(workspace, os.Getenv("DATASETID"))
	default:
		fmt.Printf("unknown command: %s\n", os.Args[1])
		os.Exit(1)
	}
}

func addElements(workspace, datasetID string, elements []elementInput) {
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

	for _, e := range elements {
		content := []byte(e.Content)
		_, err := d.AddElement(e.Name, e.Description, content)
		if err != nil {
			fmt.Printf("failed to create element: %v\n", err)
			os.Exit(1)
		}
	}

	fmt.Println("elements added successfully")
}
