package workerq

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

const (
	URLTemplate = "https://age-of-empires-2-api.herokuapp.com/api/v1/civilization/%d"
)

type Job struct {
	ID int
}

type Worker struct {
	ID       int
	JobQueue chan Job
	Quit     chan int
	Wg       *sync.WaitGroup
}

func NewWorker(ID int, jobQueue chan Job, quit chan int, wg *sync.WaitGroup) *Worker {
	return &Worker{
		ID:       ID,
		JobQueue: jobQueue,
		Quit:     quit,
		Wg:       wg,
	}
}

func (w *Worker) Start() {
	c := &http.Client{Timeout: time.Millisecond * 15000}
	go func() {
		for {
			select {
			case job := <-w.JobQueue:

				res, _ := w.callAPI(job.ID, c)
				fmt.Printf("worker %d pick job %d: %s\n", w.ID, job.ID, res)
				w.Wg.Done()

			case <-w.Quit:
				break
			}
		}
	}()
}

func (w *Worker) callAPI(id int, c *http.Client) (string, error) {
	ur := fmt.Sprintf(URLTemplate, id)
	req, err := http.NewRequest(http.MethodGet, ur, nil)
	if err != nil {
		return "", err
	}
	res, err := c.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	output, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	return string(output), nil
}
