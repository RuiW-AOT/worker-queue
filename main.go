package main

import (
	"fmt"
	"sync"

	"github.com/RuiW-AOT/worker-queue/workerq"
)

func main() {
	ids := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	wg := &sync.WaitGroup{}

	dispatcher := workerq.NewDispatcher(3, wg)
	dispatcher.Start()

	for _, id := range ids {
		job := workerq.Job{ID: id}
		wg.Add(1)
		dispatcher.Submit(job)
	}
	wg.Wait()
	fmt.Println("finished")
}
