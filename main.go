package main

import (
	"fmt"
	"os"

	"github.com/gptscript-ai/datasets/pkg/tools"
)

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
	default:
		fmt.Printf("unknown command: %s\n", os.Args[1])
		os.Exit(1)
	}
}
