package main

import (
	"fmt"
	"log"
)

type Worker struct {
	ID         int
	WorkerPool chan chan Job
	JobChannel chan Job
	quit       chan bool
}

func NewWorker(workerPool chan chan Job, id int) Worker {
	return Worker{
		ID:         id,
		WorkerPool: workerPool,
		JobChannel: make(chan Job),
		quit:       make(chan bool),
	}
}

func (w Worker) Start() {
	fmt.Printf("starting worker %d \n", w.ID)
	go func() {
		for {
			// register the current worker into the worker queue.
			w.WorkerPool <- w.JobChannel

			select {
			case job := <-w.JobChannel:
				fmt.Printf("worker %d <- request\n", w.ID)

				if err := job.Request.predict(); err != nil {
					log.Fatal(err.Error())
				}

			case <-w.quit:
				fmt.Printf("worker %d <- quit\n", w.ID)
				return
			}
		}
	}()
}

func (w Worker) Stop() {
	go func() {
		w.quit <- true
	}()
}
