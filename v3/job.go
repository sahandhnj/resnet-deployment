package main

type Job struct {
	Request Request
}

var JobQueue chan Job
