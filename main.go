package main

import (
	"fmt"
	"os"
	"sync"
	"time"
)

func main() {
	now := time.Now()

	q := make(chan event, 100)
	f, err := os.Create("text.txt")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	go pushToQueue(q, 200_000)
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
	subQ := make(chan event, 10)
	go func() {
		for e := range qe {
			wg.Add(1)
			subQ <- e
		}
		close(subQ)
		println("end sub queue")
	}()

	for e := range subQ {
		go func(ee event) {
			fmt.Println(ee)
			_, err := f.WriteString(fmt.Sprintln(ee))
			if err != nil {
				fmt.Println(err, "write error")
			}
			<-time.After(10 * time.Second)
			wg.Done()
		}(e)
	}

	wg.Wait()
	println("processed all event")
}
