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

	strs := make([]string, 10, 100)
	var i int
	mu := &sync.RWMutex{}
	//mu2 := &sync.Mutex{}
	for v := range in {
		s := strconv.Itoa(v.(int))
		ch := make(chan string)
		go func(s string, i int) {
			defer mu.Unlock()
			mu.Lock()
			strs[i] = DataSignerCrc32(s)
			i++
		}(s, i)
		i++
		go func(s string, out chan string) {
			defer mu.Unlock()
			mu.Lock()
			out <- DataSignerMd5(s)
		}(s, ch)
		go func(in chan string, i int) {
			defer mu.Unlock()
			mu.Lock()
			strs[i] = DataSignerCrc32(<-in)
		}(ch, i)
		i++
	}
	mu.RLock()
	for _, res := range strs {
		out <- res
	}
	mu.RUnlock()

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
