package main

import (
	"time"
)

func main() {
	ch := make(chan int, 10)

	go func() {
		var i = 1
		for {
			i++
			ch <- i
		}
	}()

	for {
		select {
		case x := <-ch:
			println(x)
		case <-time.After(10 * time.Second):
			println(time.Now().Unix())
		}
	}
}
