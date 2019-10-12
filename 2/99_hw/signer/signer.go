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
			if i < len(jobs) {
				close(chans[i])
			}
		}(i, f)
	}
	wg.Wait()
}

func SingleHash(in, out chan interface{}) {
	fmt.Println("starting SingleHash")
	var strs []string
	for v := range in {
		x := v
		strs = append(strs, strconv.Itoa(x.(int)))
	}
	for _, s := range strs {
		res := DataSignerCrc32(s) + "~" + DataSignerCrc32(DataSignerMd5(s))
		fmt.Println("SingleHash read from (in) chan", res)
		out <- res
	}
	fmt.Println("ending SingleHash")
}

func MultiHash(in, out chan interface{}) {
	fmt.Println("starting MultiHash")
	for s := range in {
		res := ""
		for th := 0; th <= 5; th++ {
			res += DataSignerCrc32(strconv.Itoa(th) + s.(string))
			fmt.Println("MultiHash read from (in) chan", res)
		}
		out <- res
	}
	fmt.Println("ending MultiHash")
}

func CombineResults(in, out chan interface{}) {
	fmt.Println("starting Combine")
	var res []string
	for s := range in {
		res = append(res, s.(string))
		fmt.Println("Combine read from (in) chan", res)
	}
	sort.Strings(res)
	fmt.Println("HASH : ", strings.Join(res, "_"))
	out <- strings.Join(res, "_")
	fmt.Println("ending Combine")
}
