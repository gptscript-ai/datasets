package dataset

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type Dataset interface {
	GetID() string
	Type() string
	Length() int
	Nth(i int) (string, error)
	Range(i, j int) ([]string, error)
}

// ArrayDataset represents an array of generic JSON data.
type ArrayDataset struct {
	ID   string
	Data []any
}

func (d *ArrayDataset) GetID() string {
	return d.ID
}

func (d *ArrayDataset) Type() string {
	return "array"
}

func (d *ArrayDataset) Length() int {
	return len(d.Data)
}

func (d *ArrayDataset) Nth(i int) (string, error) {
	if i < 0 || i >= len(d.Data) {
		return "", fmt.Errorf("index %d out of bounds for dataset %s", i, d.ID)
	}

	datum, err := json.Marshal(d.Data[i])
	if err != nil {
		return "", fmt.Errorf("error marshalling data at index %d in dataset %s: %v", i, d.ID, err)
	}

	return string(datum), nil
}

func (d *ArrayDataset) Range(i, j int) ([]string, error) {
	if i > j {
		return nil, fmt.Errorf("invalid range %d - %d for dataset %s", i, j, d.ID)
	}

	if i < 0 || j >= len(d.Data) {
		return nil, fmt.Errorf("range %d - %d out of bounds for dataset %s", i, j, d.ID)
	}

	var data []string
	for k := i; k <= j; k++ {
		datum, err := d.Nth(k)
		if err != nil {
			return nil, err
		}
		data = append(data, datum)
	}

	return data, nil
}

// FileDataset represents a single file in the workspace.
// This dataset supports three different iteration strategies:
// - LineMethod: each line in the file is a separate piece of data
// - SplitMethod: the file is split by a delimiter, specified in a metadata file
// - WholeMethod: the entire file is a single piece of data
type FileDataset struct {
	Method       IterationMethod
	ID, Splitter string
	Contents     []byte
}

func (d *FileDataset) GetID() string {
	return d.ID
}

func (d *FileDataset) Type() string {
	return "file"
}

func (d *FileDataset) Length() int {
	fileStr := string(d.Contents)
	switch d.Method {
	case LineMethod:
		return len(strings.Split(fileStr, "\n"))
	case SplitMethod:
		return len(strings.Split(fileStr, d.Splitter))
	case WholeMethod:
		return 1
	}
	return 0
}

func (d *FileDataset) Nth(i int) (string, error) {
	fileStr := string(d.Contents)
	switch d.Method {
	case LineMethod:
		lines := strings.Split(fileStr, "\n")
		if i < 0 || i >= len(lines) {
			return "", fmt.Errorf("index %d out of bounds for dataset %s", i, d.ID)
		}
		return lines[i], nil
	case SplitMethod:
		parts := strings.Split(fileStr, d.Splitter)
		if i < 0 || i >= len(parts) {
			return "", fmt.Errorf("index %d out of bounds for dataset %s", i, d.ID)
		}
		return parts[i], nil
	case WholeMethod:
		if i > 0 {
			return "", fmt.Errorf("index %d out of bounds for dataset %s", i, d.ID)
		}
		return fileStr, nil
	}
	return "", fmt.Errorf("unknown iteration strategy %s for dataset %s", d.Method, d.ID)
}

func (d *FileDataset) Range(i, j int) ([]string, error) {
	if i > j {
		return nil, fmt.Errorf("invalid range %d - %d for dataset %s", i, j, d.ID)
	}

	fileStr := string(d.Contents)
	switch d.Method {
	case LineMethod:
		lines := strings.Split(fileStr, "\n")
		if i < 0 || j >= len(lines) {
			return nil, fmt.Errorf("range %d - %d out of bounds for dataset %s", i, j, d.ID)
		}
		return lines[i : j+1], nil
	case SplitMethod:
		parts := strings.Split(fileStr, d.Splitter)
		if i < 0 || j >= len(parts) {
			return nil, fmt.Errorf("range %d - %d out of bounds for dataset %s", i, j, d.ID)
		}
		return parts[i : j+1], nil
	case WholeMethod:
		if i > 0 || j > 1 {
			return nil, fmt.Errorf("range %d - %d out of bounds for dataset %s", i, j, d.ID)
		}
		return []string{fileStr}, nil
	}
	return nil, fmt.Errorf("unknown iteration strategy %s for dataset %s", d.Method, d.ID)
}

// FolderDataset represents a folder in the workspace, where each file is a single piece of data.
type FolderDataset struct {
	ID    string
	Files []string
}

func (d *FolderDataset) GetID() string {
	return d.ID
}

func (d *FolderDataset) Type() string {
	return "folder"
}

func (d *FolderDataset) Length() int {
	return len(d.Files)
}

func (d *FolderDataset) Nth(i int) (string, error) {
	data, _, err := d.nthWithCurrentSize(i, 0)
	return data, err
}

func (d *FolderDataset) nthWithCurrentSize(i int, currentSize int64) (string, int64, error) {
	if i < 0 || i >= len(d.Files) {
		return "", 0, fmt.Errorf("index %d out of bounds for dataset %s", i, d.ID)
	}

	fileName := d.Files[i]
	fileStat, err := os.Stat(fileName)
	if err != nil {
		return "", 0, fmt.Errorf("error getting info for file %s: %v", fileName, err)
	}

	if fileStat.Size()+currentSize > 100*1024*1024 { // 100 MiB
		return "", 0, fmt.Errorf("dataset %s is too large to read (combined file size must be under 100 MiB)", d.ID)
	}
	currentSize += fileStat.Size()

	contents, err := os.ReadFile(fileName)
	if err != nil {
		return "", 0, fmt.Errorf("error reading file %s: %v", fileName, err)
	}

	return string(contents), currentSize, nil
}

func (d *FolderDataset) Range(i, j int) ([]string, error) {
	if i > j {
		return nil, fmt.Errorf("invalid range %d - %d for dataset %s", i, j, d.ID)
	}

	if i < 0 || j >= len(d.Files) {
		return nil, fmt.Errorf("range %d - %d out of bounds for dataset %s", i, j, d.ID)
	}

	var (
		data     []string
		contents string
		size     int64
		err      error
	)
	for k := i; k <= j; k++ {
		contents, size, err = d.nthWithCurrentSize(k, size)
		if err != nil {
			return nil, err
		}
		data = append(data, contents)
	}

	return data, nil
}
