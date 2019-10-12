package main

// сюда писать код
import (
	"sort"
	"strconv"
	"strings"
	"sync"
)

//type job func(in, out chan interface{})

func ExecutePipeline(jobs ...job) {
	var chans []chan interface{}
	for i := 0; i < len(jobs); i++ {
		chans = append(chans, make(chan interface{}))
	}

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		jobs[0](chans[0], chans[0])
	}()
	go func() {
		defer wg.Done()
		jobs[1](chans[0], chans[1])
	}()
	wg.Wait()
}

func SingleHash(in, out chan interface{}) {
	x := <-in
	s := strconv.Itoa(x.(int))
	out <- DataSignerCrc32(s) + "~" + DataSignerMd5(s)
}
func MultiHash(in, out chan interface{}) {
	defer func() {
		out <- "finish"
		close(in)
	}()
	for th := 0; th < 5; th++ {
		s := <-in
		out <- DataSignerCrc32(strconv.Itoa(th) + s.(string))
	}
}

func CombineResults(in, out chan interface{}) {
	var res []string
	for {
		s := <-in
		res = append(res, s.(string))
		if s == "finish" {
			close(in)
			break
		}
	}
	sort.Strings(res)
	out <- strings.Join(res, "_")
}

/*
func main() {

	var ok = true
	var recieved uint32
	freeFlowJobs := []job{
		job(func(in, out chan interface{}) {
			fmt.Println("job before write")
			out <- 1
			time.Sleep(10 * time.Millisecond)
			currRecieved := atomic.LoadUint32(&recieved)
			// в чем тут суть
			// если вы накапливаете значения, то пока вся функция не отрабоатет - дальше они не пойдут
			// тут я проверяю, что счетчик увеличился в следующей функции
			// это значит что туда дошло значение прежде чем текущая функция отработала
			if currRecieved == 0 {
				ok = false
			}
			fmt.Println("job after write")
		}),
		job(func(in, out chan interface{}) {
			fmt.Println("job before read")
			for _ = range in {
				atomic.AddUint32(&recieved, 1)
			}
			fmt.Println("job after read")
		}),
	}
	ExecutePipeline(freeFlowJobs...)
	fmt.Println(recieved)
	if !ok || recieved == 0 {
		fmt.Println("no value free flow - dont collect them")
	}
}
*/
