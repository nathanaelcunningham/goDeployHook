package main

import (
	"log"
	"os/exec"
)

type Job struct {
	Repository string
	Script     string
}

func (j Job) Process() {
	log.Println(j.Repository)
	out, err := exec.Command(j.Script).CombinedOutput()
	if err != nil {
		log.Println(err)
	}

	if len(out) > 0 {
		log.Println(string(out))
	}
}

type Queue struct {
	Jobs chan Job
}

func (q *Queue) Push(job Job) {
	q.Jobs <- job
}

func (q *Queue) Start() {
	go q.loop()
}

func (q *Queue) loop() {
	for {
		job := <-q.Jobs
		job.Process()
	}
}
