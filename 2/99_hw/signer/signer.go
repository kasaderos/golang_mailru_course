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
	wg := &sync.WaitGroup{}
	for v := range in {
		s := strconv.Itoa(v.(int))
		out1 := make(chan string)
		go func(s string, out1 chan string) {
			out1 <- DataSignerCrc32(s)
		}(s, out1)
		wg.Add(1)
		go func(s string, out chan interface{}, mu *sync.Mutex) {
			defer wg.Done()
			mu.Lock()
			v := DataSignerMd5(s)
			mu.Unlock()
			out <- <-out1 + "~" + DataSignerCrc32(v)
		}(s, out, mu)
	}
	wg.Wait()
}

func MultiHash(in, out chan interface{}) {
	wg := &sync.WaitGroup{}
	for v := range in {
		var temp chan string
		for i := 0; i <= 5; i++ {
			to := make(chan string)
			from := temp
			temp = to
			wg.Add(1)
			go func(s string, i int, from, to chan string, out chan interface{}) {
				defer wg.Done()
				v := DataSignerCrc32(strconv.Itoa(i) + s)
				if from == nil {
					to <- v
				} else if i < 5 {
					to <- <-from + v
				} else if i == 5 {
					out <- <-from + v
				}
			}(v.(string), i, from, to, out)
		}
	}
	wg.Wait()
}

func CombineResults(in, out chan interface{}) {
	var res []string
	for s := range in {
		res = append(res, s.(string))
	}
	sort.Strings(res)
	out <- strings.Join(res, "_")
}
