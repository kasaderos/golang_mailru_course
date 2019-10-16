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

	for v := range in {
		var temp chan string
		out1 := make(chan string)
		for th := 0; th <= 5; th++ {
			next := make(chan string)
			inn := temp
			temp = next
			go func(s string, th int) {
				if inn != nil {
					<-inn
				}
				out1 <- DataSignerCrc32(strconv.Itoa(th) + s)
				fmt.Println(s, th)
				if th < 5 {
					fmt.Println(s, th, "2")
					next <- "go"
				} else {
					res := ""
					for j := 0; j <= 5; j++ {
						res += <-out1
					}
					out2 <- res
					fmt.Println("#")
				}
			}(v.(string), th)
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
