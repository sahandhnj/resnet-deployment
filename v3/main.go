package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/sahandhnj/ml-deployment-benchmarks/v3/db"
	"github.com/sahandhnj/ml-deployment-benchmarks/v3/service"
)

const (
	MaxWorker    = 0
	MaxQueue     = 0
	Address      = ":3002"
	QueuedResult = false
)

var (
	works = 0
	done  = 0
)

var dbhandler *db.DBStore
var reqservice *service.ReqService

func main() {
	http.HandleFunc("/resnet/v1/predict", requestHandler)
	http.HandleFunc("/stat", reqDataHandler)

	fmt.Printf("size of queue %d\n", MaxQueue)
	dbhandler, err := db.NewDBStore()
	reqservice = service.NewReqService(dbhandler)

	if err != nil {
		log.Fatal(err)
	}

	JobQueue = make(chan Job, MaxQueue)

	dispatcher := NewDispatcher(MaxWorker)
	dispatcher.Run()

	fmt.Println("Listening on " + Address)
	err = http.ListenAndServe(Address, nil)
	if err != nil {
		log.Fatal(err)
	}
}
