package concurrent

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
)

// JobPool defines the operations to interact with a job pool containing multiple workers working together concurrently.
type JobPool interface {
	// AddWorker adds a worker to process the items in the job pool using f.
	AddWorker(ctx context.Context, wg *sync.WaitGroup, db *sql.DB, f func(string, *sql.DB) error)
	// Start starts the workers in the job pool.
	Start(ctx context.Context)
	// Enqueue adds the given input to the job pool to be processed by its workers.
	Enqueue(input string)
}

type jobPool struct {
	inputChan  chan string
	workerChan chan string
}

// AddWorker adds a worker to process the items in job pool using f.
func (jp jobPool) AddWorker(ctx context.Context, wg *sync.WaitGroup, db *sql.DB, f func(string, *sql.DB) error) {
	go func() {
		defer wg.Done()
		for {
			select {
			case url := <-jp.workerChan:
				if err := f(url, db); err != nil {
					fmt.Printf("Error while processing input %s : %s", url, err)
				}
			case <-ctx.Done():
				return
			}
		}
	}()
}

// Start starts the workers in the job pool.
func (jp jobPool) Start(ctx context.Context) {
	go func() {
		for {
			select {
			case url := <-jp.inputChan:
				if ctx.Err() != nil {
					close(jp.workerChan)
					return
				}
				jp.workerChan <- url
			case <-ctx.Done():
				close(jp.workerChan)
				return
			}
		}
	}()
}

// Enqueue adds the given input to the job pool to be processed by its workers.
func (jp jobPool) Enqueue(input string) {
	jp.inputChan <- input
}

// NewJobPool creates and returns a JobPool.
func NewJobPool(numWorkers int) JobPool {
	return &jobPool{
		inputChan:  make(chan string),
		workerChan: make(chan string, numWorkers),
	}
}
