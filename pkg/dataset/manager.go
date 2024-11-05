package dataset

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/gptscript-ai/go-gptscript"
)

const (
	datasetFolder = "datasets"
)

type Manager struct {
	gptscriptClient *gptscript.GPTScript
	workspaceID     string
}

func NewManager(workspaceID string) (Manager, error) {
	g, err := gptscript.NewGPTScript()
	if err != nil {
		return Manager{}, fmt.Errorf("failed to create GPTScript: %w", err)
	}

	return Manager{gptscriptClient: g, workspaceID: workspaceID}, nil
}

func (m *Manager) ListDatasets(ctx context.Context) ([]DatasetMeta, error) {
	files, err := m.gptscriptClient.ListFilesInWorkspace(ctx, gptscript.ListFilesInWorkspaceOptions{
		Prefix:      datasetFolder,
		WorkspaceID: m.workspaceID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list dataset files: %w", err)
	}

	var datasets []DatasetMeta
	for _, file := range files {
		contents, err := m.gptscriptClient.ReadFileInWorkspace(ctx, file, gptscript.ReadFileInWorkspaceOptions{
			WorkspaceID: m.workspaceID,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to read dataset file %s: %w", file, err)
		}

		var d Dataset
		if err = json.Unmarshal(contents, &d); err != nil {
			return nil, fmt.Errorf("failed to read dataset file %s: %w", file, err)
		}

		datasets = append(datasets, d.DatasetMeta)
	}

	return datasets, nil
}

func (m *Manager) NewDataset(ctx context.Context, name, description string) (Dataset, error) {
	randBytes := make([]byte, 3)
	if _, err := rand.Read(randBytes); err != nil {
		return Dataset{}, fmt.Errorf("failed to generate random bytes: %w", err)
	}

	id := fmt.Sprintf("gds://%x", randBytes)[:11]
	d := Dataset{
		DatasetMeta: DatasetMeta{
			ID:          id,
			Name:        name,
			Description: description,
		},
		Elements: make(map[string]Element),
	}

	// Now convert to JSON and save it to the workspace
	datasetJSON, err := json.Marshal(d)
	if err != nil {
		return Dataset{}, fmt.Errorf("failed to marshal dataset: %w", err)
	}

	if err := m.gptscriptClient.WriteFileInWorkspace(ctx, datasetFolder+"/"+idToFileName(id), datasetJSON, gptscript.WriteFileInWorkspaceOptions{
		WorkspaceID: m.workspaceID,
	}); err != nil {
		return Dataset{}, fmt.Errorf("failed to write dataset file: %w", err)
	}

	d.m = m
	return d, nil
}

func (m *Manager) GetDataset(ctx context.Context, id string) (Dataset, error) {
	fileName := idToFileName(id)
	data, err := m.gptscriptClient.ReadFileInWorkspace(ctx, datasetFolder+"/"+fileName, gptscript.ReadFileInWorkspaceOptions{
		WorkspaceID: m.workspaceID,
	})
	if err != nil {
		if !isNotFoundInWorkspaceError(err) {
			return Dataset{}, fmt.Errorf("dataset %s not found", id)
		}
		return Dataset{}, fmt.Errorf("failed to read dataset file: %w", err)
	}

	var d Dataset
	if err = json.Unmarshal(data, &d); err != nil {
		return Dataset{}, fmt.Errorf("failed to unmarshal dataset file %s: %w", datasetFolder+"/"+fileName, err)
	}

	d.m = m
	return d, nil
}

func idToFileName(id string) string {
	return id[6:] + ".gds"
}

func isNotFoundInWorkspaceError(err error) bool {
	var notFoundErr *gptscript.NotFoundInWorkspaceError
	return errors.As(err, &notFoundErr)
}
