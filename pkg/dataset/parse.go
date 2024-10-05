package dataset

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/tidwall/gjson"
)

type metadata struct {
	GptscriptMetadata bool   `json:"gptscriptMetadata"`
	File              string `json:"file"`
	Method            string `json:"method,omitempty"`
	Splitter          string `json:"splitter,omitempty"`
}

func ParseDataset(id, workspace string) (Dataset, error) {
	idFile := workspace + string(os.PathSeparator) + id

	idInfo, err := os.Stat(idFile)
	if err != nil {
		return nil, fmt.Errorf("error getting info for dataset %s: %v", id, err)
	}

	if idInfo.IsDir() {
		return parseDir(id, workspace)
	} else if idInfo.Size() > 100*1024*1024 { // 100 MiB
		return nil, fmt.Errorf("dataset %s is too large (over 100 MiB)", id)
	}

	idContents, err := os.ReadFile(idFile)
	if err != nil {
		return nil, fmt.Errorf("error reading data from file %s: %v", id, err)
	}

	if strings.HasSuffix(id, ".json") {
		// This is either a metadata file, or a data array. Check for the array first.
		res := gjson.Get(string(idContents), "data")
		if res.Exists() && res.IsArray() {
			return parseArray(id, []byte(res.Raw))
		}

		// Now check to see whether it's a metadata file to point to a file dataset.
		var meta metadata
		if err := json.Unmarshal(idContents, &meta); err == nil && meta.GptscriptMetadata {
			return parseMeta(id, workspace, meta)
		}
	}

	// This is just a generic file. Treat it as a file dataset with default parameters.
	return &FileDataset{
		Method:   LineMethod,
		ID:       id,
		Contents: normalizeLineEndings(idContents),
	}, nil
}

func parseArray(id string, data []byte) (Dataset, error) {
	var dataArray []any
	if err := json.Unmarshal(data, &dataArray); err != nil {
		return nil, fmt.Errorf("error unmarshalling data for dataset %s: %v", id, err)
	}

	return &ArrayDataset{
		ID:   id,
		Data: dataArray,
	}, nil
}

func parseMeta(id, workspace string, meta metadata) (Dataset, error) {
	// Set the default method to "line" if it's not set.
	if meta.Method == "" {
		meta.Method = string(LineMethod)
	}

	if !strings.HasPrefix(meta.File, workspace) {
		meta.File = workspace + string(os.PathSeparator) + meta.File
	}

	fileInfo, err := os.Stat(meta.File)
	if err != nil {
		return nil, fmt.Errorf("error getting info for dataset %s: %v", id, err)
	}

	if fileInfo.IsDir() {
		return nil, fmt.Errorf("file dataset %s points to a directory", id)
	}

	contents, err := os.ReadFile(meta.File)
	if err != nil {
		return nil, fmt.Errorf("error reading file dataset %s: %v", id, err)
	}

	return &FileDataset{
		Method:   IterationMethod(meta.Method),
		ID:       id,
		Splitter: meta.Splitter,
		Contents: normalizeLineEndings(contents),
	}, nil
}

func parseDir(id, workspace string) (Dataset, error) {
	files, err := recursiveGetFilenames(id, workspace)
	if err != nil {
		return nil, err
	}

	return &FolderDataset{
		ID:    id,
		Files: files,
	}, nil
}

func recursiveGetFilenames(id, workspace string) ([]string, error) {
	dirContents, err := os.ReadDir(workspace + string(os.PathSeparator) + id)
	if err != nil {
		return nil, fmt.Errorf("error reading directory contents for dataset %s: %v", id, err)
	}

	var filenames []string
	for _, entry := range dirContents {
		if entry.IsDir() {
			subFiles, err := recursiveGetFilenames(id+string(os.PathSeparator)+entry.Name(), workspace)
			if err != nil {
				return nil, err
			}
			filenames = append(filenames, subFiles...)
		} else {
			filenames = append(filenames, workspace+string(os.PathSeparator)+id+string(os.PathSeparator)+entry.Name())
		}
	}

	return filenames, nil
}

func normalizeLineEndings(contents []byte) []byte {
	return bytes.ReplaceAll(contents, []byte("\r\n"), []byte("\n"))
}
