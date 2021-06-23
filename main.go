package main

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/Jeffail/tunny"
)

func main() {
	now := time.Now()

	q := make(chan event)
	f, err := os.Create("text.txt")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	go pushToQueue(q, 200000)
	processQueue(q, f)

	fmt.Println(time.Since(now).String())
}

type event struct {
	index int
}

func pushToQueue(q chan<- event, max int) {
	for i := 0; i < max; i++ {
		q <- event{
			index: i,
		}
	}

	close(q)
	fmt.Println("pushed")
}

func processQueue(qe <-chan event, f *os.File) {
	wg := &sync.WaitGroup{}
	pool := tunny.New(100, createWorker(wg, f))
	defer pool.Close()

	for e := range qe {
		wg.Add(1)
		go pool.Process(e)
	}

	wg.Wait()
	println("processed all event")
}

func createWorker(wg *sync.WaitGroup, f *os.File) func() tunny.Worker {
	return func() tunny.Worker {
		return worker{wg: wg, f: f}
	}
}

type worker struct {
	wg *sync.WaitGroup
	f  *os.File
}

func (w worker) Process(payload interface{}) interface{} {
	defer w.wg.Done()
	var err error
	e := payload.(event)

	fmt.Println(e)
	_, err = w.f.WriteString(fmt.Sprintln(e))

	<-time.After(10 * time.Millisecond)

	return err
}

func (w worker) BlockUntilReady() {}
func (w worker) Interrupt()       {}
func (w worker) Terminate()       {}
