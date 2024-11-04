package dataset

import (
	"context"
	"testing"

	"github.com/gptscript-ai/go-gptscript"
	"github.com/stretchr/testify/require"
)

func TestDatasets(t *testing.T) {
	ctx := context.Background()

	g, err := gptscript.NewGPTScript()
	require.NoError(t, err)

	workspaceID, err := g.CreateWorkspace(ctx, "directory")
	require.NoError(t, err)
	t.Logf("workspace ID: %s", workspaceID)

	defer g.DeleteWorkspace(ctx, gptscript.DeleteWorkspaceOptions{
		WorkspaceID: workspaceID,
	})

	m, err := NewManager(workspaceID)
	require.NoError(t, err)

	m.workspaceID = workspaceID

	t.Cleanup(func() {
		_ = m.gptscriptClient.DeleteWorkspace(ctx)
	})

	dataset, err := m.NewDataset(ctx)
	require.NoError(t, err)
	require.Equal(t, 0, dataset.GetLength())

	// Let's add a couple elements.
	err = dataset.AddElement(Element{
		ElementMeta: ElementMeta{
			Name:        "file1",
			Description: "The first file",
		},
		Contents: "This is dataset file 1",
	})
	require.NoError(t, err)
	require.Equal(t, 1, dataset.GetLength())

	err = dataset.AddElement(Element{
		ElementMeta: ElementMeta{
			Name:        "file2",
			Description: "The second file",
		},
		Contents: "This is dataset file 2",
	})
	require.NoError(t, err)
	require.Equal(t, 2, dataset.GetLength())

	err = dataset.AddElement(Element{
		ElementMeta: ElementMeta{
			Name:        "binary file",
			Description: "has binary contents",
		},
		BinaryContents: []byte("binary contents"),
	})

	require.NoError(t, dataset.Save(ctx))

	// Let's read it back.
	dataset, err = m.GetDataset(ctx, dataset.GetID())
	require.NoError(t, err)

	metas := dataset.ListElements()
	require.Len(t, metas, 3)

	oneElement, err := dataset.GetElement("file1")
	require.NoError(t, err)
	require.Equal(t, "This is dataset file 1", oneElement.Contents)
	require.Equal(t, 0, oneElement.Index)

	twoElement, err := dataset.GetElement("file2")
	require.NoError(t, err)
	require.Equal(t, "This is dataset file 2", twoElement.Contents)
	require.Equal(t, 1, twoElement.Index)

	binaryElement, err := dataset.GetElement("binary file")
	require.NoError(t, err)
	require.Equal(t, []byte("binary contents"), binaryElement.BinaryContents)
	require.Equal(t, 2, binaryElement.Index)

	// Test to make sure the order was maintained
	elementMetas := dataset.ListElements()
	require.Equal(t, "file1", elementMetas[0].Name)
	require.Equal(t, "file2", elementMetas[1].Name)
	require.Equal(t, "binary file", elementMetas[2].Name)

	datasets, err := m.ListDatasets(ctx)
	require.NoError(t, err)
	require.Len(t, datasets, 1)
}
