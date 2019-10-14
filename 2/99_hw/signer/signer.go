package main

// сюда писать код
import (
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
	mu := &sync.Mutex{}
	const goroutinesNum = 7
	workerInput1 := make(chan string)
	workerInput2 := make(chan string)
	ch1 := make(chan string)
	ch3 := make(chan string)
	for i := 0; i < goroutinesNum; i++ {
		ch2 := make(chan string)
		go func(s chan string, res chan string) {
			res <- DataSignerCrc32(<-s)
		}(workerInput1, ch1)
		go func(s chan string, out chan string, mu *sync.Mutex) {
			mu.Lock()
			out <- DataSignerMd5(<-s)
			mu.Unlock()
		}(workerInput2, ch2, mu)
		go func(res, in chan string) {
			res <- DataSignerCrc32(<-in)
		}(ch3, ch2)
	}
	for v := range in {
		s := strconv.Itoa(v.(int))
		workerInput1 <- s
		workerInput2 <- s
	}
	for i := 0; i < goroutinesNum; i++ {
		out <- <-ch1 + "~" + <-ch3
	}
}

func MultiHash(in, out chan interface{}) {
	var chans []chan string
	k := 0
	for v := range in {
		for th := 0; th <= 5; th++ {
			chans = append(chans, make(chan string))
			go func(s string, res chan string, th int) {
				res <- DataSignerCrc32(strconv.Itoa(th) + s)
			}(v.(string), chans[k], th)
			k++
		}
	}
	for i := 0; i < len(chans); i += 6 {
		res := ""
		for j := i; j <= i+5; j++ {
			res += <-chans[j]
		}
		out <- res
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
