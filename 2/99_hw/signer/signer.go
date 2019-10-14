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
	ch1 := make(chan interface{}, len(strs)*2)
	ch2 := make(chan interface{}, len(strs))
	ch3 := make(chan interface{}, len(strs))
	//wg := &sync.WaitGroup{}
	for i, s := range strs {
		go func(i int, s string, res chan string, k int) {
			defer func() {
				//wg.Done()
				fmt.Println("SingleHash th done!", i, "-", "1")
			}()
			<-ch2
			fmt.Println("SingleHash th ", i, "-", 1)
			res <- DataSignerCrc32(s)
		}(i, s, chans[k], k)

		//wg.Add(1)
		go func(i int, s string, out chan string, k int) {
			defer func() {
				//wg.Done()
				fmt.Println("SingleHash th done!", i, "-", "2")
			}()
			<-ch1
			fmt.Println("SingleHash th ", i, "-", 2)
			out <- DataSignerMd5(s)
		}(i, s, chans[k+1], k+1)
		go func(i int, res chan string, in chan string, k int) {
			defer func() {
				//wg.Done()
				fmt.Println("SingleHash th done!", i, "-", "3")
			}()
			<-ch3
			s := <-in
			fmt.Println("SingleHash th ", i, "-", 3)
			res <- DataSignerCrc32(s)
		}(i, chans[k+2], chans[k+1], k+2)
		k += 3
	}
	//close(ch)

	for i := 0; i < len(strs); i++ {
		go func() {
			ch1 <- 1
		}()
	}
	for i := 0; i < len(strs); i++ {
		ch3 <- 1
		ch2 <- 1
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
	ch := make(chan interface{}, len(strs)*6)
	for i, s := range strs {
		for th := 0; th <= 5; th++ {
			go func(i int, s string, res chan string, th int, k int) {
				defer func() {
					fmt.Println("MultiHash th done!", i, "-", th)
				}()
				<-ch
				fmt.Println("MultiHash th ", i, "-", th)
				res <- DataSignerCrc32(strconv.Itoa(th) + s)
			}(i, s, chans[k], th, k)
			k++
		}
	}
	res := ""
	for j := 0; j < len(strs)*6; j++ {
		ch <- 1
	}
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
