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
