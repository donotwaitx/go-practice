// 12 - Goroutines & Channels
// Concurrency là siêu năng lực của Go.
package main

import (
	"fmt"
	"sync"
	"time"
)

func worker(id int, wg *sync.WaitGroup) {
	defer wg.Done() // báo cho WaitGroup khi xong
	fmt.Printf("worker %d bắt đầu\n", id)
	time.Sleep(100 * time.Millisecond)
	fmt.Printf("worker %d xong\n", id)
}

func producer(ch chan<- int) {
	for i := 1; i <= 5; i++ {
		ch <- i * 10 // gửi vào channel
	}
	close(ch) // đóng channel khi xong
}

func main() {
	// 1) Goroutine cơ bản với WaitGroup
	var wg sync.WaitGroup
	for i := 1; i <= 3; i++ {
		wg.Add(1)
		go worker(i, &wg)
	}
	wg.Wait()
	fmt.Println("--- tất cả workers xong ---")

	// 2) Channel — kênh giao tiếp giữa goroutines
	ch := make(chan int)
	go producer(ch)
	for v := range ch { // tự dừng khi channel close
		fmt.Println("nhận:", v)
	}

	// 3) Buffered channel
	bch := make(chan string, 2)
	bch <- "a"
	bch <- "b"
	close(bch)
	for v := range bch {
		fmt.Println("buf:", v)
	}

	// 4) select — chờ nhiều channels
	c1, c2 := make(chan string, 1), make(chan string, 1)
	go func() { time.Sleep(50 * time.Millisecond); c1 <- "from c1" }()
	go func() { time.Sleep(30 * time.Millisecond); c2 <- "from c2" }()

	for i := 0; i < 2; i++ {
		select {
		case m := <-c1:
			fmt.Println(m)
		case m := <-c2:
			fmt.Println(m)
		case <-time.After(200 * time.Millisecond):
			fmt.Println("timeout")
		}
	}
}
