package main

import (
	"fmt"
)

func main() {
	in := make(chan int, 2)

	go func(out chan<- int) {
		for i := 0; i <= 10; i++ {
			fmt.Println("before", i)
			out <- i
			fmt.Println("after", i)
		}
		close(out)
		fmt.Println("generator finish")
	}(in)

	for i := range in {
		fmt.Println("\tget", i, len(in))
	}

	return

	for {
		i, isClosed := <-in
		if !isClosed {
			fmt.Println("\tchan closed", i)
			break
		}
		fmt.Println("\tget", i)
	}

	// fmt.Scanln()
}
