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
	datasetFolder     = "datasets"
	datasetMetaFolder = "datasets/meta"
)

type Manager struct {
	gptscriptClient *gptscript.GPTScript
	workspaceID     string
}

func NewManager() (Manager, error) {
	g, err := gptscript.NewGPTScript()
	if err != nil {
		return Manager{}, fmt.Errorf("failed to create GPTScript: %w", err)
	}

	return Manager{gptscriptClient: g}, nil
}

func (m *Manager) ListDatasets(ctx context.Context) ([]DatasetMeta, error) {
	files, err := m.gptscriptClient.ListFilesInWorkspace(ctx, gptscript.ListFilesInWorkspaceOptions{
		Prefix:      datasetMetaFolder,
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
	randBytes := make([]byte, 16)
	if _, err := rand.Read(randBytes); err != nil {
		return Dataset{}, fmt.Errorf("failed to generate random bytes: %w", err)
	}

	id := fmt.Sprintf("%x", randBytes)
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

	if err := m.gptscriptClient.WriteFileInWorkspace(ctx, datasetMetaFolder+"/"+id, datasetJSON, gptscript.WriteFileInWorkspaceOptions{
		WorkspaceID: m.workspaceID,
	}); err != nil {
		return Dataset{}, fmt.Errorf("failed to write dataset file: %w", err)
	}

	d.m = m
	return d, nil
}

func (m *Manager) GetDataset(ctx context.Context, id string) (Dataset, error) {
	data, err := m.gptscriptClient.ReadFileInWorkspace(ctx, datasetMetaFolder+"/"+id, gptscript.ReadFileInWorkspaceOptions{
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
		return Dataset{}, fmt.Errorf("failed to unmarshal dataset file %s: %w", datasetMetaFolder+"/"+id, err)
	}

	d.m = m
	return d, nil
}

func (m *Manager) EnsureUniqueElementFilename(ctx context.Context, datasetID, name string) (string, error) {
	var counter int
	uniqueName := name
	for {
		if _, err := m.gptscriptClient.ReadFileInWorkspace(ctx, datasetFolder+"/"+datasetID+"/"+uniqueName, gptscript.ReadFileInWorkspaceOptions{
			WorkspaceID: m.workspaceID,
		}); err == nil {
			counter++
			uniqueName = fmt.Sprintf("%s_%d", name, counter)
		} else if !isNotFoundInWorkspaceError(err) {
			return "", fmt.Errorf("failed to check if file exists: %w", err)
		} else {
			return datasetFolder + "/" + datasetID + "/" + uniqueName, nil
		}
	}
}

func isNotFoundInWorkspaceError(err error) bool {
	var notFoundErr *gptscript.NotFoundInWorkspaceError
	return errors.As(err, &notFoundErr)
}
