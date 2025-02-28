package worker

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Result struct {
	Line       string
	LineNumber int
	Path       string
}

type Results struct {
	Inner []Result
}

func NewResult(line string, lineNUmber int, path string) Result {
	return Result{Line: line, LineNumber: lineNUmber, Path: path}
}

func FindInFile(path string, find string) *Results {
	file, err := os.Open(path)
	if err != nil {
		fmt.Println("error", err)
		return nil
	}
	results := &Results{make([]Result, 0)}
	scanner := bufio.NewScanner(file)
	lineNun := 1
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), find) {
			r := NewResult(scanner.Text(), lineNun, path)
			results.Inner = append(results.Inner, r)
		}
		lineNun += 1
	}
	if len(results.Inner) == 0 {
		return nil
	} else {
		return results
	}
}
