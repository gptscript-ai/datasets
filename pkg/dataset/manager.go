package dataset

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/gptscript-ai/datasets/pkg/util"
)

const (
	datasets  = "datasets"
	extension = ".dataset.json"
)

type Manager struct {
	datasetDir string
}

func NewManager(workspaceDir string) (Manager, error) {
	datasetDir := filepath.Join(workspaceDir, datasets)
	if _, err := os.Stat(datasetDir); os.IsNotExist(err) {
		if err := os.Mkdir(datasetDir, 0755); err != nil {
			return Manager{}, fmt.Errorf("failed to create dataset directory: %w", err)
		}
	}

	return Manager{datasetDir: datasetDir}, nil
}

func (m *Manager) ListDatasets() ([]DatasetMeta, error) {
	files, err := filepath.Glob(filepath.Join(m.datasetDir, "*"+extension))
	if err != nil {
		return nil, fmt.Errorf("failed to list dataset files: %w", err)
	}

	var datasets []DatasetMeta
	for _, file := range files {
		var d Dataset
		data, err := os.ReadFile(file)
		if err != nil {
			return nil, fmt.Errorf("failed to read dataset file %s: %w", file, err)
		}

		if err = json.Unmarshal(data, &d); err != nil {
			return nil, fmt.Errorf("failed to unmarshal dataset file %s: %w", file, err)
		}

		datasets = append(datasets, d.DatasetMeta)
	}

	return datasets, nil
}

func (m *Manager) NewDataset(name, description string) (Dataset, error) {
	randBytes := make([]byte, 16)
	if _, err := rand.Read(randBytes); err != nil {
		return Dataset{}, fmt.Errorf("failed to generate random bytes: %w", err)
	}

	id := fmt.Sprintf("%x", randBytes)
	dirName := util.EnsureUniqueFilename(m.datasetDir, util.ToFileName(name))
	baseDir := filepath.Join(m.datasetDir, dirName)
	if err := os.Mkdir(baseDir, 0755); err != nil {
		return Dataset{}, fmt.Errorf("failed to create dataset directory: %w", err)
	}

	d := Dataset{
		DatasetMeta: DatasetMeta{
			ID:          id,
			Name:        name,
			Description: description,
		},
		BaseDir:  baseDir,
		Elements: make(map[string]Element),
	}

	// Now convert to JSON and save it to the workspace
	datasetJSON, err := json.Marshal(d)
	if err != nil {
		return Dataset{}, fmt.Errorf("failed to marshal dataset: %w", err)
	}

	if err := os.WriteFile(baseDir+extension, datasetJSON, 0644); err != nil {
		return Dataset{}, fmt.Errorf("failed to write dataset file: %w", err)
	}

	return d, nil
}

func (m *Manager) GetDataset(id string) (Dataset, error) {
	files, err := filepath.Glob(filepath.Join(m.datasetDir, "*"+extension))
	if err != nil {
		return Dataset{}, fmt.Errorf("failed to list dataset files: %w", err)
	}

	for _, file := range files {
		var d Dataset
		data, err := os.ReadFile(file)
		if err != nil {
			return Dataset{}, fmt.Errorf("failed to read dataset file %s: %w", file, err)
		}

		if err = json.Unmarshal(data, &d); err != nil {
			return Dataset{}, fmt.Errorf("failed to unmarshal dataset file %s: %w", file, err)
		}

		if d.GetID() == id {
			return d, nil
		}
	}

	return Dataset{}, fmt.Errorf("dataset with ID %s not found", id)
}
