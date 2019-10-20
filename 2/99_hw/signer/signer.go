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
