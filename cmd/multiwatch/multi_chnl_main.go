package main

import (
	"fmt"
	"log"
	"time"

	"github.com/amukherj/gokraft/sync/multiwatch"
)

type Entry struct {
	id   int32
	data string
}

func main() {
	mw := multiwatch.NewMultiWatcher()
	mw.Watch()

	ch1 := make(chan interface{})
	mw.EnqueueWatcher(ch1)

	ch2 := make(chan interface{})
	mw.EnqueueWatcher(ch2)

	ch3 := make(chan interface{})
	mw.EnqueueWatcher(ch3)

	ch4 := make(chan interface{})
	mw.EnqueueWatcher(ch4)

	ch5 := make(chan interface{})
	mw.EnqueueWatcher(ch5)

	// Concurrent updates on input channels
	go func() {
		time.Sleep(10 * time.Millisecond)
		ch1 <- Entry{
			id:   1,
			data: "Hello",
		}

		time.Sleep(1 * time.Second)
		ch1 <- Entry{
			id:   1,
			data: "Hello Hello",
		}

		time.Sleep(1 * time.Second)
		ch1 <- Entry{
			id:   1,
			data: "Hello Hello Hello",
		}
		close(ch1)
		log.Println("Goroutine1 exiting")
	}()

	go func() {
		time.Sleep(10 * time.Millisecond)
		ch2 <- Entry{
			id:   2,
			data: "Hi",
		}

		time.Sleep(1 * time.Second)
		ch2 <- Entry{
			id:   2,
			data: "Hi Hi",
		}

		time.Sleep(1 * time.Second)
		ch2 <- Entry{
			id:   2,
			data: "Hi Hi Hi",
		}
		close(ch2)
		log.Println("Goroutine2 exiting")
	}()

	go func() {
		time.Sleep(10 * time.Millisecond)
		ch3 <- Entry{
			id:   3,
			data: "Hola",
		}

		time.Sleep(1 * time.Second)
		ch3 <- Entry{
			id:   3,
			data: "Hola Hola",
		}

		time.Sleep(1 * time.Second)
		ch3 <- Entry{
			id:   3,
			data: "Hola Hola Hola",
		}
		close(ch3)
		log.Println("Goroutine3 exiting")
	}()

	go func() {
		time.Sleep(10 * time.Millisecond)
		ch4 <- Entry{
			id:   4,
			data: "Servus",
		}

		time.Sleep(1 * time.Second)
		ch4 <- Entry{
			id:   4,
			data: "Servus Servus",
		}

		time.Sleep(1 * time.Second)
		ch4 <- Entry{
			id:   4,
			data: "Servus Servus Servus",
		}
		close(ch4)
		log.Println("Goroutine4 exiting")
	}()

	go func() {
		time.Sleep(10 * time.Millisecond)
		ch5 <- Entry{
			id:   5,
			data: "Bonjour",
		}

		time.Sleep(1 * time.Second)
		ch5 <- Entry{
			id:   5,
			data: "Bonjour Bonjour",
		}

		time.Sleep(1 * time.Second)
		ch5 <- Entry{
			id:   5,
			data: "Bonjour Bonjour Bonjour",
		}
		close(ch5)
		log.Println("Goroutine5 exiting")
	}()

	// Read from the multiwatcher
	done := false
	for !done {
		select {
		case data, ok := <-mw.Channel():
			if ok {
				fmt.Printf("Got: %v\n", data)
			} else {
				mw.Close()
			}
		case <-time.After(10 * time.Second):
			done = true
			mw.Stop()
		}
	}
	time.Sleep(10 * time.Second)
	mw.Stop()
}
