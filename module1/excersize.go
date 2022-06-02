package module1

import (
	"fmt"
	"time"
)

func ItemOne() []string {
	a := []string{"I", "am", "stupid", "and", "weak"}
	for i, v := range a {
		if v == "stupid" {
			a[i] = "smart"
		}
		if v == "weak" {
			a[i] = "strong"
		}
	}
	return a
}

func ItemTwo() {
	queue := make(chan int, 10)
	defer close(queue)
	// 生产者协程
	go func() {
		for {
			select {
			case <-time.After(1 * time.Second):
				queue <- time.Now().Nanosecond()
			}
		}
	}()
	// 消费者协程
	go func() {
		for {
			select {
			case <-time.After(1 * time.Second):
				fmt.Println(<-queue)
			}
		}
	}()
	time.Sleep(1 * time.Hour)
}
