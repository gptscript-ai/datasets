package dataset

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDatasets(t *testing.T) {
	ctx := context.Background()
	m, err := NewManager()
	require.NoError(t, err)

	workspaceID, err := m.gptscriptClient.CreateWorkspace(ctx, "directory")
	require.NoError(t, err)
	t.Logf("workspace ID: %s", workspaceID)

	m.workspaceID = workspaceID

	t.Cleanup(func() {
		_ = m.gptscriptClient.DeleteWorkspace(ctx)
	})

	dataset, err := m.NewDataset(ctx, "test dataset", "our lovely test dataset")
	require.NoError(t, err)
	require.Equal(t, "test dataset", dataset.GetName())
	require.Equal(t, "our lovely test dataset", dataset.GetDescription())
	require.Equal(t, 0, dataset.GetLength())

	// Let's add a couple elements.
	_, err = dataset.AddElement(ctx, "file1", "The first file", []byte("This is dataset file 1.\n"))
	require.NoError(t, err)
	require.Equal(t, 1, dataset.GetLength())

	_, err = dataset.AddElement(ctx, "file2", "The second file", []byte("This is dataset file 2.\n"))
	require.NoError(t, err)
	require.Equal(t, 2, dataset.GetLength())

	// Now test for file name collision. "file!" will take file_. "file@" will try file_, and then ultimately take file__1.
	// All we need to test for is that the behavior still works well, as this is an implementation detail that doesn't
	// concern the user.
	_, err = dataset.AddElement(ctx, "file!", "The third file", []byte("This is dataset file 3.\n"))
	require.NoError(t, err)
	require.Equal(t, 3, dataset.GetLength())

	_, err = dataset.AddElement(ctx, "file@", "The fourth file", []byte("This is dataset file 4.\n"))
	require.NoError(t, err)
	require.Equal(t, 4, dataset.GetLength())

	// Let's read it back.
	dataset, err = m.GetDataset(ctx, dataset.GetID())
	require.NoError(t, err)
	require.Equal(t, "test dataset", dataset.GetName())
	require.Equal(t, "our lovely test dataset", dataset.GetDescription())

	metas := dataset.ListElements()
	require.Len(t, metas, 4)

	oneBytes, _, err := dataset.GetElement(ctx, "file1")
	require.NoError(t, err)
	require.Equal(t, "This is dataset file 1.\n", string(oneBytes))

	twoBytes, _, err := dataset.GetElement(ctx, "file2")
	require.NoError(t, err)
	require.Equal(t, "This is dataset file 2.\n", string(twoBytes))

	threeBytes, _, err := dataset.GetElement(ctx, "file!")
	require.NoError(t, err)
	require.Equal(t, "This is dataset file 3.\n", string(threeBytes))

	fourBytes, _, err := dataset.GetElement(ctx, "file@")
	require.NoError(t, err)
	require.Equal(t, "This is dataset file 4.\n", string(fourBytes))

	datasets, err := m.ListDatasets(ctx)
	require.NoError(t, err)
	require.Len(t, datasets, 1)
}
