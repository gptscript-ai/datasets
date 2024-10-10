package util

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestToFileName(t *testing.T) {
	tests := []struct {
		name   string
		output string
	}{
		{
			name:   "test",
			output: "test",
		},
		{
			name:   "test@",
			output: "test_",
		},
		{
			name:   "test@123",
			output: "test_123",
		},
		{
			name:   "TEST_TEST-TEST",
			output: "TEST_TEST_TEST",
		},
		{
			name:   "123",
			output: "123",
		},
		{
			name:   "_",
			output: "_",
		},
		{
			name:   "!!!!!!!!!!!!!!!!!!!",
			output: "___________________",
		},
	}

	for _, test := range tests {
		if ToFileName(test.name) != test.output {
			t.Errorf("ToFileName(%s) != %s", test.name, test.output)
		}
	}
}

func TestEnsureUniqueFilename(t *testing.T) {
	wd, err := os.Getwd()
	require.NoError(t, err)

	base, err := os.MkdirTemp(wd, "test")
	require.NoError(t, err)

	t.Cleanup(func() {
		_ = os.RemoveAll(base)
	})

	name, err := EnsureUniqueFilename(base, "test")
	require.NoError(t, err)
	require.Equal(t, "test", name)
	_ = os.WriteFile(filepath.Join(base, "test"), []byte{}, 0644)

	name, err = EnsureUniqueFilename(base, "test")
	require.NoError(t, err)
	require.Equal(t, "test_1", name)
	_ = os.WriteFile(filepath.Join(base, "test_1"), []byte{}, 0644)

	name, err = EnsureUniqueFilename(base, "test")
	require.NoError(t, err)
	require.Equal(t, "test_2", name)
	_ = os.WriteFile(filepath.Join(base, "test_2"), []byte{}, 0644)

	name, err = EnsureUniqueFilename(base, "test")
	require.NoError(t, err)
	require.Equal(t, "test_3", name)
}
