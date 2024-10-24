package main

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/gptscript-ai/datasets/pkg/dataset"
	"github.com/gptscript-ai/datasets/pkg/tools"
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
		tools.GetElement(os.Getenv("DATASETID"), os.Getenv("ELEMENT"))
	case "createDataset":
		tools.CreateDataset(os.Getenv("DATASETNAME"), os.Getenv("DATASETDESCRIPTION"))
	case "addElement":
		tools.AddElement(os.Getenv("DATASETID"), os.Getenv("ELEMENTNAME"), os.Getenv("ELEMENTDESCRIPTION"), []byte(os.Getenv("ELEMENTCONTENT")))
	case "addElements":
		elements, err := handleGzip(os.Getenv("ELEMENTS"))
		if err != nil {
			fmt.Printf("failed to decompress elements: %v\n", err)
			os.Exit(1)
		}

		var elementInputs []elementInput
		if err := json.Unmarshal([]byte(elements), &elementInputs); err != nil {
			fmt.Printf("failed to unmarshal elements: %v\n", err)
			os.Exit(1)
		}

		addElements(os.Getenv("DATASETID"), elementInputs)
	case "getAllElements":
		tools.GetAllElements(os.Getenv("DATASETID"))
	default:
		fmt.Printf("unknown command: %s\n", os.Args[1])
		os.Exit(1)
	}
}

func handleGzip(elements string) (string, error) {
	var gz struct {
		Content string `json:"_gz"`
	}
	if err := json.Unmarshal([]byte(elements), &gz); err != nil {
		// If it didn't unmarshal, then it's not gzipped
		return elements, nil
	}

	decoded, err := base64.StdEncoding.DecodeString(gz.Content)
	if err != nil {
		return "", err
	}

	reader, err := gzip.NewReader(bytes.NewReader(decoded))
	if err != nil {
		return "", err
	}

	result, err := io.ReadAll(reader)
	if err != nil {
		return "", err
	}

	return string(result), nil
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
