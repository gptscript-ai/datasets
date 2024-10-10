package dataset

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/gptscript-ai/datasets/pkg/util"
)

type ElementMeta struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type Element struct {
	ElementMeta `json:",inline"`
	File        string `json:"file"`
}

type DatasetMeta struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type Dataset struct {
	DatasetMeta `json:",inline"`
	BaseDir     string             `json:"baseDir,omitempty"`
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
	var elements []ElementMeta
	for _, element := range d.Elements {
		elements = append(elements, element.ElementMeta)
	}
	return elements
}

func (d *Dataset) GetElement(name string) ([]byte, Element, error) {
	e, exists := d.Elements[name]
	if !exists {
		return nil, Element{}, fmt.Errorf("element %s not found", name)
	}

	contents, err := os.ReadFile(d.BaseDir + string(os.PathSeparator) + e.File)
	if err != nil {
		return nil, Element{}, fmt.Errorf("failed to read element %s: %w", name, err)
	}

	return contents, e, nil
}

func (d *Dataset) AddElement(name, description string, contents []byte) (Element, error) {
	if _, exists := d.Elements[name]; exists {
		return Element{}, fmt.Errorf("element %s already exists", name)
	}

	fileName, err := util.EnsureUniqueFilename(d.BaseDir, util.ToFileName(name))
	if err != nil {
		return Element{}, fmt.Errorf("failed to generate unique file name: %w", err)
	}

	loc := filepath.Join(d.BaseDir, fileName)
	if err := os.WriteFile(loc, contents, 0644); err != nil {
		return Element{}, fmt.Errorf("failed to write element %s: %w", name, err)
	}

	e := Element{
		ElementMeta: ElementMeta{
			Name:        name,
			Description: description,
		},
		File: fileName,
	}

	d.Elements[name] = e
	return e, d.save()
}

func (d *Dataset) save() error {
	datasetJSON, err := json.Marshal(d)
	if err != nil {
		return fmt.Errorf("failed to marshal dataset: %w", err)
	}

	if err := os.WriteFile(d.BaseDir+extension, datasetJSON, 0644); err != nil {
		return fmt.Errorf("failed to write dataset file: %w", err)
	}

	return nil
}
