package main

import (
	"fmt"
	"time"
)

func main() {
	now := time.Now()

	q := make(chan event)

	go pushToQueue(q, 1000)

	processQueue(q)
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
}

func processQueue(event <-chan event) {
	for e := range event {
		processEvent(e)
	}
}

func processEvent(e event) {
	fmt.Println(e)
	<-time.After(10 * time.Millisecond)
}
