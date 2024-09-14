package main

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

// Task definition
type Task interface {
	Process(i int)
}

// Email task definition
type EmailTask struct {
	Email       string
	Subject     string
	MessageBody string
}

// Way to process the Email task
func (t *EmailTask) Process(i int) {
	fmt.Printf("Worker %d is processing the email %s\n", i, t.Email)
	// Simulate a time consuming process
	time.Sleep(2 * time.Second)
}

// Image processing task
type ImageProcessingTask struct {
	ImageUrl string
}

// Way to process the Image
func (t *ImageProcessingTask) Process(i int) {
	fmt.Printf("Worker %d is processing the image %s\n", i, t.ImageUrl)
	// Simulate a time consuming process
	time.Sleep(5 * time.Second)
}

// Worker pool definition
type WorkerPool struct {
	Tasks       []Task
	concurrency int
	tasksChan   chan Task
	wg          sync.WaitGroup
}

// Functions to execute the worker pool

func (wp *WorkerPool) worker(i int) {
	for task := range wp.tasksChan {
		task.Process(i)
		wp.wg.Done()
	}
	fmt.Printf("Worker %d has finished\n", i)
	fmt.Println(strings.Repeat("-", 10))
}

func (wp *WorkerPool) Run() {
	// Initialize the tasks channel
	wp.tasksChan = make(chan Task, len(wp.Tasks))

	// Start workers
	for i := 0; i < wp.concurrency; i++ {
		go wp.worker(i+1)
	}

	// Send tasks to the tasks channel
	wp.wg.Add(len(wp.Tasks))
	for _, task := range wp.Tasks {
		wp.tasksChan <- task
	}
	close(wp.tasksChan)

	// Wait for all tasks to finish
	wp.wg.Wait()
}
