package dataset

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"unicode"

	"github.com/gptscript-ai/go-gptscript"
)

type ElementMeta struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type Element struct {
	ElementMeta `json:",inline"`
	File        string `json:"file"`
	Index       int    `json:"index"`
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

func (d *Dataset) GetName() string {
	return d.Name
}

func (d *Dataset) GetDescription() string {
	return d.Description
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

func (d *Dataset) GetElement(ctx context.Context, name string) ([]byte, Element, error) {
	e, exists := d.Elements[name]
	if !exists {
		return nil, Element{}, fmt.Errorf("element %s not found", name)
	}

	contents, err := d.m.gptscriptClient.ReadFileInWorkspace(ctx, e.File, gptscript.ReadFileInWorkspaceOptions{
		WorkspaceID: d.m.workspaceID,
	})
	if err != nil {
		return nil, Element{}, fmt.Errorf("failed to read element %s: %w", name, err)
	}

	return contents, e, nil
}

func (d *Dataset) AddElement(ctx context.Context, name, description string, contents []byte) (Element, error) {
	if _, exists := d.Elements[name]; exists {
		return Element{}, fmt.Errorf("element %s already exists", name)
	}

	fileName, err := d.m.EnsureUniqueElementFilename(ctx, d.ID, toFileName(name))
	if err != nil {
		return Element{}, fmt.Errorf("failed to generate unique file name: %w", err)
	}

	if err := d.m.gptscriptClient.WriteFileInWorkspace(ctx, fileName, contents, gptscript.WriteFileInWorkspaceOptions{
		WorkspaceID: d.m.workspaceID,
	}); err != nil {
		return Element{}, fmt.Errorf("failed to write element %s: %w", name, err)
	}

	e := Element{
		ElementMeta: ElementMeta{
			Name:        name,
			Description: description,
		},
		Index: len(d.Elements),
		File:  fileName,
	}

	d.Elements[name] = e
	return e, d.save(ctx)
}

func (d *Dataset) save(ctx context.Context) error {
	datasetJSON, err := json.Marshal(d)
	if err != nil {
		return fmt.Errorf("failed to marshal dataset: %w", err)
	}

	if err := d.m.gptscriptClient.WriteFileInWorkspace(ctx, datasetMetaFolder+"/"+d.ID, datasetJSON, gptscript.WriteFileInWorkspaceOptions{
		WorkspaceID: d.m.workspaceID,
	}); err != nil {
		return fmt.Errorf("failed to write dataset file: %w", err)
	}
	return nil
}

// toFileName converts a name to be alphanumeric plus underscores.
func toFileName(name string) string {
	return strings.Map(func(c rune) rune {
		if !unicode.IsLetter(c) && !unicode.IsDigit(c) {
			return '_'
		}
		return c
	}, name)
}
