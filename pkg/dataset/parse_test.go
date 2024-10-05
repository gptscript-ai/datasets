package dataset

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

const testWorkspace = "testworkspace"

func TestParseArrayDataset(t *testing.T) {
	wd, err := os.Getwd()
	require.NoError(t, err)
	dataset, err := ParseDataset("array.json", wd+string(os.PathSeparator)+testWorkspace)
	require.NoError(t, err)

	require.Equal(t, "array.json", dataset.GetID())
	require.Equal(t, "array", dataset.Type())
	require.Equal(t, 6, dataset.Length())
	indexZero, err := dataset.Nth(0)
	require.NoError(t, err)
	require.Equal(t, "\"one\"", indexZero)
	oneToFour, err := dataset.Range(1, 5)
	require.NoError(t, err)
	require.Equal(t, []string{`"two"`, `["three"]`, `{"four":true}`, `["five",5]`, `6`}, oneToFour)
}

func TestParseFileDataset(t *testing.T) {
	wd, err := os.Getwd()
	require.NoError(t, err)
	dataset, err := ParseDataset("file.txt", wd+string(os.PathSeparator)+testWorkspace)
	require.NoError(t, err)

	require.Equal(t, "file.txt", dataset.GetID())
	require.Equal(t, "file", dataset.Type())
	require.Equal(t, 5, dataset.Length())
	indexZero, err := dataset.Nth(0)
	require.NoError(t, err)
	require.Equal(t, "This is the first line.", indexZero)
	oneToFour, err := dataset.Range(1, 4)
	require.NoError(t, err)
	require.Equal(t,
		[]string{"This is the second line.",
			"This is the third line.",
			"This is the fourth line.",
			"This is the fifth line."},
		oneToFour)
}

func TestParseFileSplitterDataset(t *testing.T) {
	wd, err := os.Getwd()
	require.NoError(t, err)
	dataset, err := ParseDataset("file_meta.json", wd+string(os.PathSeparator)+testWorkspace)
	require.NoError(t, err)

	require.Equal(t, "file_meta.json", dataset.GetID())
	require.Equal(t, "file", dataset.Type())
	require.Equal(t, 4, dataset.Length())
	indexZero, err := dataset.Nth(0)
	require.NoError(t, err)
	require.Equal(t, "This is the first datum.", indexZero)
	oneToThree, err := dataset.Range(1, 3)
	require.NoError(t, err)
	require.Equal(t,
		[]string{"This is the second datum.",
			"This is the third datum.",
			"This is the fourth datum."},
		oneToThree)
}

func TestParseFolderDataset(t *testing.T) {
	wd, err := os.Getwd()
	require.NoError(t, err)
	dataset, err := ParseDataset("dataset_dir", wd+string(os.PathSeparator)+testWorkspace)
	require.NoError(t, err)

	require.Equal(t, "dataset_dir", dataset.GetID())
	require.Equal(t, "folder", dataset.Type())
	require.Equal(t, 2, dataset.Length())
	indexZero, err := dataset.Nth(0)
	require.NoError(t, err)
	require.Equal(t, "This is file 1, line 1.\nThis is file 1, line 2.", indexZero)
	zeroAndOne, err := dataset.Range(0, 1)
	require.NoError(t, err)
	require.Equal(t,
		[]string{"This is file 1, line 1.\nThis is file 1, line 2.",
			"This is file 2, line 1.\nThis is file 2, line 2."},
		zeroAndOne)
}
