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
		go func(s string, bw chan string) {
			mu.Lock()
			bw <- DataSignerMd5(s)
			mu.Unlock()
		}(s, bw)
		go func(bw chan string) {
			out2 <- DataSignerCrc32(<-bw)
		}(bw)
	}
	for i := 0; i < goroutinesnum; i++ {
		out <- <-out1 + "~" + <-out2
	}
}

func MultiHash(in, out chan interface{}) {
	out2 := make(chan string)
	for v := range in {
		var temp chan string
		for i := 0; i <= 5; i++ {
			outt := make(chan string)
			inn := temp
			temp = outt
			go func(s string, i int, in, out chan string) {
				v := DataSignerCrc32(strconv.Itoa(i) + s)
				if in == nil {
					out <- v
				} else if i < 5 {
					out <- <-in + v
				} else if i == 5 {
					out2 <- <-in + v
				}
			}(v.(string), i, inn, outt)
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
