package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"
)

func студент(ctx context.Context, студентNum int, out chan<- int) {
	waitTime := time.Duration(rand.Intn(100)+10) * time.Millisecond
	fmt.Println(студентNum, "sleep", waitTime)
	select {
	// case <-end:
	// 	fmt.Println("студент", студентNum, "finished by ctx")
	// 	return
	case <-ctx.Done():
		fmt.Println("студент", студентNum, "finished by ctx")
		return
	case <-time.After(waitTime):
		fmt.Println("студент", студентNum, "придумал вопрос")
		out <- студентNum
	}
}

func main() {
	ctx, finish := context.WithCancel(context.Background())
	result := make(chan int, 1)
	// end := make(chan int, 10)

	for i := 0; i <= 10; i++ {
		go студент(ctx, i, result)
	}

	foundBy := <-result
	fmt.Println("запрос задан", foundBy)
	finish()

	time.Sleep(time.Second)
}
