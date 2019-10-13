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
	var chans []chan interface{}
	for i := 0; i < len(jobs); i++ {
		chans = append(chans, make(chan interface{}))
	}
	wg := &sync.WaitGroup{}
	for i, f := range jobs {
		wg.Add(1)
		go func(i int, f job) {
			defer wg.Done()
			if i == 0 {
				f(chans[i], chans[i])
			} else {
				f(chans[i-1], chans[i])
			}
			close(chans[i])
		}(i, f)
	}
	wg.Wait()
}

func SingleHash(in, out chan interface{}) {
	fmt.Println("SingleHash starts")
	var strs []string
	for v := range in {
		strs = append(strs, strconv.Itoa(v.(int)))
	}
	var chans []chan string
	for i := 0; i < len(strs)*3; i++ {
		chans = append(chans, make(chan string))
	}
	k := 0
	for i, s := range strs {
		go func(i int, s string, res chan string) {
			defer func() {
				fmt.Println("SingleHash th done!", i, "-", "1")
			}()
			fmt.Println("SingleHash th ", i, "-", 1)
			res <- DataSignerCrc32(s)
		}(i, s, chans[k])
		go func(i int, s string, out chan string) {
			defer func() {
				fmt.Println("SingleHash th done!", i, "-", "2")
			}()
			fmt.Println("SingleHash th ", i, "-", 2)
			out <- DataSignerMd5(s)
		}(i, s, chans[k+1])
		go func(i int, res chan string, in chan string) {
			defer func() {
				fmt.Println("SingleHash th done!", i, "-", "3")
			}()
			fmt.Println("SingleHash th ", i, "-", 3)
			s := <-in
			res <- DataSignerCrc32(s)
		}(i, chans[k+2], chans[k+1])
		k += 3
	}
	for i := 0; i < len(chans); i += 3 {
		res := <-chans[i] + "~" + <-chans[i+2]
		out <- res
	}
}

func MultiHash(in, out chan interface{}) {
	fmt.Println("MultiHash")
	var strs []string
	for v := range in {
		strs = append(strs, v.(string))
	}
	var chans []chan string
	for i := 0; i < len(strs)*6; i++ {
		chans = append(chans, make(chan string))
	}
	k := 0
	for i, s := range strs {
		for th := 0; th <= 5; th++ {
			go func(i int, s string, res chan string, th int) {
				defer func() {
					fmt.Println("MultiHash th done!", i, "-", th)
				}()
				fmt.Println("MultiHash th ", i, "-", th)
				res <- DataSignerCrc32(strconv.Itoa(th) + s)
			}(i, s, chans[k], th)
			k++
		}
	}
	res := ""
	for i := 0; i < len(chans); i += 6 {
		res = ""
		for j := i; j <= i+5; j++ {
			res += <-chans[j]
		}
		out <- res
	}
}

func CombineResults(in, out chan interface{}) {
	fmt.Println("CombineResults")
	var res []string
	for s := range in {
		res = append(res, s.(string))
	}
	sort.Strings(res)
	out <- strings.Join(res, "_")
}
