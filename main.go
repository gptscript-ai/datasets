package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/gptscript-ai/datasets/pkg/dataset"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("missing argument")
		os.Exit(1)
	}

	workspace := os.Getenv("GPTSCRIPT_WORKSPACE_DIR")
	if workspace == "" {
		fmt.Println("missing GPTSCRIPT_WORKSPACE_DIR")
		os.Exit(1)
	}

	arg := os.Args[1]

	var (
		result string
		err    error
	)
	switch arg {
	case "info":
		result, err = info(os.Getenv("ID"), workspace)
	case "load_one":
		result, err = loadOne(os.Getenv("ID"), os.Getenv("INDEX"), workspace)
	case "load_range":
		result, err = loadRange(os.Getenv("ID"), os.Getenv("START"), os.Getenv("END"), workspace)
	case "load_all":
		result, err = loadAll(os.Getenv("ID"), workspace)
	}

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println(result)
}

func info(id, workspace string) (string, error) {
	set, err := dataset.ParseDataset(id, workspace)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("Dataset ID: %s, length: %d", set.GetID(), set.Length()), nil
}

func loadOne(id, index, workspace string) (string, error) {
	set, err := dataset.ParseDataset(id, workspace)
	if err != nil {
		return "", err
	}

	indexInt, err := strconv.Atoi(index)
	if err != nil {
		return "", fmt.Errorf("invalid index: %v", err)
	}

	data, err := set.Nth(indexInt)
	if err != nil {
		return "", err
	}

	return data, nil
}

func loadRange(id, start, end, workspace string) (string, error) {
	set, err := dataset.ParseDataset(id, workspace)
	if err != nil {
		return "", err
	}

	startInt, err := strconv.Atoi(start)
	if err != nil {
		return "", fmt.Errorf("invalid start: %v", err)
	}
	endInt, err := strconv.Atoi(end)
	if err != nil {
		return "", fmt.Errorf("invalid end: %v", err)
	}

	data, err := set.Range(startInt, endInt)
	if err != nil {
		return "", err
	}

	return strings.Join(data, "\n\n"), nil
}

func loadAll(id, workspace string) (string, error) {
	set, err := dataset.ParseDataset(id, workspace)
	if err != nil {
		return "", err
	}

	data, err := set.Range(0, set.Length()-1)
	if err != nil {
		return "", err
	}

	return strings.Join(data, "\n\n"), nil
}
