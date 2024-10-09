package dataset

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

const testWorkspace = "testworkspace"

func TestDatasetsRead(t *testing.T) {
	wd, err := os.Getwd()
	require.NoError(t, err)

	workspaceDir := filepath.Join(wd, testWorkspace)
	m, err := NewManager(workspaceDir)
	require.NoError(t, err)

	datasetMetas, err := m.ListDatasets()
	require.NoError(t, err)
	require.Len(t, datasetMetas, 2)

	datasetOne, err := m.GetDataset("one")
	require.NoError(t, err)
	require.Equal(t, "one", datasetOne.GetName())
	require.Equal(t, "The first test dataset", datasetOne.GetDescription())
	require.Equal(t, 2, datasetOne.GetLength())

	oneMetas := datasetOne.ListElements()
	require.Len(t, oneMetas, 2)

	oneOneBytes, _, err := datasetOne.GetElement("file1")
	require.NoError(t, err)
	require.Equal(t, "This is dataset 1, file 1.\n", string(oneOneBytes))

	oneTwoBytes, _, err := datasetOne.GetElement("file2")
	require.NoError(t, err)
	require.Equal(t, "This is dataset 1, file 2.\n", string(oneTwoBytes))

	datasetTwo, err := m.GetDataset("two")
	require.NoError(t, err)
	require.Equal(t, "two", datasetTwo.GetName())
	require.Equal(t, "The second test dataset", datasetTwo.GetDescription())
	require.Equal(t, 2, datasetTwo.GetLength())

	twoMetas := datasetTwo.ListElements()
	require.Len(t, twoMetas, 2)

	twoOneBytes, _, err := datasetTwo.GetElement("file1")
	require.NoError(t, err)
	require.Equal(t, "This is dataset 2, file 1.\n", string(twoOneBytes))

	twoTwoBytes, _, err := datasetTwo.GetElement("file2")
	require.NoError(t, err)
	require.Equal(t, "This is dataset 2, file 2.\n", string(twoTwoBytes))
}

func TestDatasetWrite(t *testing.T) {
	wd, err := os.Getwd()
	require.NoError(t, err)

	workspaceDir := filepath.Join(wd, testWorkspace)
	m, err := NewManager(workspaceDir)
	require.NoError(t, err)

	t.Cleanup(func() {
		threeFiles, _ := filepath.Glob(filepath.Join(workspaceDir, "datasets", "three", "*"))

		for _, file := range threeFiles {
			_ = os.Remove(file)
		}

		_ = os.Remove(filepath.Join(workspaceDir, "datasets", "three"))
		_ = os.Remove(filepath.Join(workspaceDir, "datasets", "three.dataset.json"))
	})

	datasetThree, err := m.NewDataset("three", "The third test dataset")
	require.NoError(t, err)
	require.Equal(t, "three", datasetThree.GetName())
	require.Equal(t, "The third test dataset", datasetThree.GetDescription())
	require.Equal(t, 0, datasetThree.GetLength())

	// Let's add a couple elements.
	_, err = datasetThree.AddElement("file1", "The first file", []byte("This is dataset 3, file 1.\n"))
	require.NoError(t, err)
	require.Equal(t, 1, datasetThree.GetLength())

	_, err = datasetThree.AddElement("file2", "The second file", []byte("This is dataset 3, file 2.\n"))
	require.NoError(t, err)
	require.Equal(t, 2, datasetThree.GetLength())

	// Let's read it back.
	datasetThree, err = m.GetDataset(datasetThree.GetID())
	require.NoError(t, err)
	require.Equal(t, "three", datasetThree.GetName())
	require.Equal(t, "The third test dataset", datasetThree.GetDescription())

	threeMetas := datasetThree.ListElements()
	require.Len(t, threeMetas, 2)

	threeOneBytes, _, err := datasetThree.GetElement("file1")
	require.NoError(t, err)
	require.Equal(t, "This is dataset 3, file 1.\n", string(threeOneBytes))

	threeTwoBytes, _, err := datasetThree.GetElement("file2")
	require.NoError(t, err)
	require.Equal(t, "This is dataset 3, file 2.\n", string(threeTwoBytes))
}
