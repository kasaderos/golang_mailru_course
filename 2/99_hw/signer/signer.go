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
	var prev chan interface{}
	for i, f := range jobs {
		wg.Add(1)
		out := make(chan interface{})
		in := prev
		prev = out
		go func(i int, f job, in, out chan interface{}) {
			defer wg.Done()
			f(in, out)
			close(out)
		}(i, f, in, out)
	}
	wg.Wait()
}

func SingleHash(in, out chan interface{}) {
	mu := &sync.Mutex{}
<<<<<<< HEAD
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
=======
	wg := &sync.WaitGroup{}
	for v := range in {
		betweenCh := make(chan string)
		value := strconv.Itoa(v.(int))
		go func(s string, out chan string) {
			out <- DataSignerCrc32(s)
		}(value, betweenCh)
		wg.Add(1)
		go func(s string, out chan interface{}, previos <-chan string, mu *sync.Mutex) {
			defer wg.Done()
			mu.Lock()
			md5 := DataSignerMd5(s)
			mu.Unlock()
			result := DataSignerCrc32(md5)
			out <- <-previos + "~" + result
		}(value, out, betweenCh, mu)
>>>>>>> home3
	}
	wg.Wait()
}

func MultiHash(in, out chan interface{}) {
	wg := &sync.WaitGroup{}
	for v := range in {
		var temp chan string
		for th := 0; th <= 5; th++ {
			outTh := make(chan string)
			inTh := temp
			temp = outTh
			wg.Add(1)
			go func(s string, th int, inTh <-chan string, outTh chan string, out chan interface{}) {
				defer wg.Done()
				result := DataSignerCrc32(strconv.Itoa(th) + s)
				if inTh == nil {
					outTh <- result
				} else if th < 5 {
					outTh <- <-inTh + result
				} else if th == 5 {
					out <- <-inTh + result
				}
			}(v.(string), th, inTh, outTh, out)
		}
	}
	wg.Wait()
}

func CombineResults(in, out chan interface{}) {
	var result []string
	for s := range in {
		result = append(result, s.(string))
	}
	sort.Strings(result)
	out <- strings.Join(result, "_")
}
