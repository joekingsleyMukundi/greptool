package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/alexflint/go-arg"
	"github.com/joekingsleyMukundi/greptool/worker"
	"github.com/joekingsleyMukundi/greptool/worklist"
)

func discoverDirs(wl *worklist.Worklist, path string) {
	entries, err := os.ReadDir(path)
	if err != nil {
		fmt.Println("Reader error:", err)
	}
	for _, entry := range entries {
		if entry.IsDir() {
			nextPath := filepath.Join(path, entry.Name())
			discoverDirs(wl, nextPath)
		} else {
			wl.Add(worklist.NewJob(filepath.Join(path, entry.Name())))
		}
	}
}

type Args struct {
	SearchTerm string `arg:"positional,required"`
	SearchDir  string `arg:"positional"`
}

func main() {
	var args Args
	arg.MustParse(&args)
	var workersWg sync.WaitGroup
	wl := worklist.New(100)
	results := make(chan worker.Result, 100)
	numWorkers := 10
	workersWg.Add(1)
	go func() {
		defer workersWg.Done()
		discoverDirs(&wl, args.SearchDir)
		wl.Finalize(numWorkers)
	}()
	for i := 0; i < numWorkers; i++ {
		workersWg.Add(1)
		go func() {
			defer workersWg.Done()
			for {
				workEntry := wl.Next()
				if workEntry.Path != "" {
					workerResult := worker.FindInFile(workEntry.Path, args.SearchTerm)
					if workerResult != nil {
						for _, r := range workerResult.Inner {
							results <- r
						}
					}
				} else {
					return
				}
			}
		}()
	}
	blockWorkersWg := make(chan struct{})
	go func() {
		workersWg.Wait()
		close(blockWorkersWg)
	}()
	var dispayWg sync.WaitGroup
	dispayWg.Add(1)
	go func() {
		for {
			select {
			case r := <-results:
				fmt.Printf("%v[%v]:%v\n", r.Line, r.LineNumber, r.Path)
			case <-blockWorkersWg:
				if len(results) == 0 {
					dispayWg.Done()
					return
				}
			}
		}
	}()
	dispayWg.Wait()
}
