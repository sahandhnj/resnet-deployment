package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	tf "github.com/tensorflow/tensorflow/tensorflow/go"
)

type Status string

const (
	Queued     = "Queued"
	Processing = "Processing"
	Done       = "Done"
	Failed     = "Failed"
)

type Request struct {
	ID        string
	ImageFile multipart.File
	ImageName string
	w         *http.ResponseWriter
	Status    Status
	StartedAt time.Time
}

type Resp struct {
	RequestedAt    time.Time   `json:"requestedAt"`
	ProcessingTime int64       `json:"processingTime"`
	Status         Status      `json:"status"`
	Prediction     *Prediction `json:"prediction"`
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

type Prediction struct {
	Name       string  `json:"name"`
	Probabiliy float32 `json:"probability"`
}

func RandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func (p *Request) predict() error {
	startProcessing := time.Now()

	fmt.Println("Handling request " + time.Now().String())
	modeldir := "meta/model"

	defer p.ImageFile.Close()

	inputName := "meta/input/" + p.ImageName + "-input-" + p.ID + ".jpg"

	// Load a frozen graph to use for queries
	modelpath := filepath.Join(modeldir, "resnet50.pb")
	model, err := ioutil.ReadFile(modelpath)
	if err != nil {
		return err
	}

	inputFile, err := os.Create(inputName)

	if err != nil {
		return err
	}
	defer inputFile.Close()

	_, err = io.Copy(inputFile, p.ImageFile)
	if err != nil {
		return err
	}

	// DecodeJpeg uses a scalar String-valued tensor as input.
	tensor, err := makeTensorFromImage(inputName)
	if err != nil {
		return err
	}
	// Construct an in-memory graph from the serialized form.
	graph := tf.NewGraph()
	if err := graph.Import(model, ""); err != nil {
		return err
	}

	// Create a session for inference over graph.
	session, err := tf.NewSession(graph, nil)
	if err != nil {
		return err
	}
	defer session.Close()

	// for _, obj := range graph.Operations() {
	// 	fmt.Printf("%s : %s -> %d\n", obj.Name(), obj.Type(), obj.NumOutputs())

	// }
	// Get all the input and output operations
	inputop := graph.Operation("input_1")
	// Output ops
	o1 := graph.Operation("fc1000/Softmax")
	// fmt.Println("HIII")
	// fmt.Printf("%s\n", o1.Name())
	// Execute COCO Graph

	outputs, err := session.Run(
		map[tf.Output]*tf.Tensor{
			inputop.Output(0): tensor,
		},
		[]tf.Output{
			o1.Output(0),
		},
		nil)
	if err != nil {
		return err
	}

	probabilities := outputs[0].Value().([][]float32)[0]
	pred, err := normalizeResult(probabilities, "meta/labels.json")
	if err != nil {
		return err
	}

	done = done + 1
	fmt.Printf("Finished predicting #%d with Id:%s\n", done, p.ID)

	if !QueuedResult {
		processingTime := time.Since(startProcessing).Nanoseconds()
		reqservice.Add(startProcessing, processingTime)
		resp := &Resp{
			RequestedAt:    p.StartedAt,
			ProcessingTime: processingTime,
			Status:         Done,
			Prediction:     pred,
		}

		json, err := json.Marshal(resp)
		if err != nil {
			return err
		}

		l := *p.w
		l.Header().Set("Content-Type", "application/json")
		l.Write(json)
	}

	return nil
}

func normalizeResult(probabilities []float32, labelsFile string) (*Prediction, error) {
	bestIdx := 0
	for i, p := range probabilities {
		if p > probabilities[bestIdx] {
			bestIdx = i
		}
	}

	file, err := os.Open(labelsFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	var labels []string
	for scanner.Scan() {
		labels = append(labels, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	name := labels[bestIdx]
	name = strings.Replace(name, "\"", "", -1)
	name = strings.TrimSpace(name)
	fmt.Printf("Most likely to be a %s (%2.0f)\n", name, probabilities[bestIdx]*100.0)

	return &Prediction{
		Name:       name,
		Probabiliy: probabilities[bestIdx] * 100.0,
	}, nil
}
