package main

import (
	"fmt"
	"os"

	"github.com/gptscript-ai/datasets/pkg/tools"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("missing argument")
		os.Exit(1)
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
		tools.ListElements(workspace, os.Getenv("DATASET_ID"))
	case "getElement":
		tools.GetElement(workspace, os.Getenv("DATASET_ID"), os.Getenv("ELEMENT"))
	case "createDataset":
		tools.CreateDataset(workspace, os.Getenv("DATASET_NAME"), os.Getenv("DATASET_DESCRIPTION"))
	case "addElement":
		tools.AddElement(workspace, os.Getenv("DATASET_ID"), os.Getenv("ELEMENT_NAME"), os.Getenv("ELEMENT_DESCRIPTION"), os.Getenv("ELEMENT_TYPE"), []byte(os.Getenv("ELEMENT_CONTENT")))
	}
}
