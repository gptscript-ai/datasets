package dataset

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"

	"github.com/gptscript-ai/go-gptscript"
)

type ElementMeta struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type Element struct {
	ElementMeta    `json:",inline"`
	Index          int    `json:"index"`
	Contents       string `json:"contents"`
	BinaryContents []byte `json:"binaryContents"`
}

// ElementNoIndex is used for returning data to the user, since the user does not care about the index.
type ElementNoIndex struct {
	ElementMeta    `json:",inline"`
	Contents       string `json:"contents"`
	BinaryContents []byte `json:"binaryContents"`
}

type DatasetMeta struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type Dataset struct {
	m           *Manager
	DatasetMeta `json:",inline"`
	Elements    map[string]Element `json:"elements"`
}

func (d *Dataset) GetID() string {
	return d.ID
}

func (d *Dataset) GetLength() int {
	return len(d.Elements)
}

func (d *Dataset) ListElements() []ElementMeta {
	var elements []Element
	for _, element := range d.Elements {
		elements = append(elements, element)
	}
	sort.Slice(elements, func(i, j int) bool {
		return elements[i].Index < elements[j].Index
	})

	var elementMetas []ElementMeta
	for _, element := range elements {
		elementMetas = append(elementMetas, element.ElementMeta)
	}
	return elementMetas
}

func (d *Dataset) GetAllElements() []ElementNoIndex {
	var elements []Element
	for _, element := range d.Elements {
		elements = append(elements, element)
	}
	sort.Slice(elements, func(i, j int) bool {
		return elements[i].Index < elements[j].Index
	})

	var noIndex []ElementNoIndex
	for _, element := range elements {
		noIndex = append(noIndex, ElementNoIndex{
			ElementMeta:    element.ElementMeta,
			Contents:       element.Contents,
			BinaryContents: element.BinaryContents,
		})
	}

	return noIndex
}

func (d *Dataset) GetElement(name string) (Element, error) {
	e, exists := d.Elements[name]
	if !exists {
		return Element{}, fmt.Errorf("element %s not found", name)
	}

	return e, nil
}

func (d *Dataset) AddElement(e Element) error {
	if _, exists := d.Elements[e.Name]; exists {
		return fmt.Errorf("element %s already exists", e.Name)
	}

	e.Index = len(d.Elements)
	d.Elements[e.Name] = e
	return nil
}

func (d *Dataset) Save(ctx context.Context) error {
	datasetJSON, err := json.Marshal(d)
	if err != nil {
		return fmt.Errorf("failed to marshal dataset: %w", err)
	}

	if err := d.m.gptscriptClient.WriteFileInWorkspace(ctx, datasetFolder+"/"+idToFileName(d.ID), datasetJSON, gptscript.WriteFileInWorkspaceOptions{
		WorkspaceID: d.m.workspaceID,
	}); err != nil {
		return fmt.Errorf("failed to write dataset file: %w", err)
	}
	return nil
}
