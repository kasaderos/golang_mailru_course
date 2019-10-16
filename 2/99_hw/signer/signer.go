package main

// сюда писать код
import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync"
)

func ExecutePipeline(jobs ...job) {
	wg := &sync.WaitGroup{}
	var pred chan interface{}
	for i, f := range jobs {
		wg.Add(1)
		out := make(chan interface{})
		in := pred
		pred = out
		go func(i int, f job) {
			defer wg.Done()
			f(in, out)
			close(out)
		}(i, f)
	}
	wg.Wait()
}

func SingleHash(in, out chan interface{}) {
	const goroutinesnum = 7 // можно получить путем подсчета из канала
	out1 := make(chan string)
	out2 := make(chan string)
	mu := &sync.Mutex{}
	for v := range in {
		s := strconv.Itoa(v.(int))
		bw := make(chan string)
		go func(s string) {
			out1 <- DataSignerCrc32(s)
		}(s)
		go func(s string) {
			mu.Lock()
			bw <- DataSignerMd5(s)
			mu.Unlock()
		}(s)
		go func() {
			out2 <- DataSignerCrc32(<-bw)
		}()
	}
	for i := 0; i < goroutinesnum; i++ {
		out <- <-out1 + "~" + <-out2
	}
}

func MultiHash(in, out chan interface{}) {
	out2 := make(chan string)
	workerInput := make(chan int)
	for v := range in {
		mu := &sync.Mutex{}
		res := ""
		for i := 0; i <= 5; i++ {
			go func(s string, th chan int) {
				v := <-th
				mu.Lock()
				res += DataSignerCrc32(string(v) + s)
				if v == 5 {
					out2 <- res
				}
				mu.Unlock()
			}(v.(string), workerInput)
		}
	}
	for i := 0; i < 7; i++ {
		for k := 0; k <= 5; k++ {
			fmt.Println(k)
			workerInput <- k
		}
	}
	for i := 0; i < 7; i++ {
		out <- <-out2
	}
}

func CombineResults(in, out chan interface{}) {
	var res []string
	for s := range in {
		res = append(res, s.(string))
	}
	sort.Strings(res)
	out <- strings.Join(res, "_")
}
