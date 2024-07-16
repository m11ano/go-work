package main

import (
	"fmt"
	"sort"
	"strings"
	"sync"
)

func SingleHash(in, out chan interface{}) {
	wg := &sync.WaitGroup{}
	md5Mu := &sync.Mutex{}
	for val := range in {
		strVal := fmt.Sprintf("%v", val)
		wg.Add(1)
		go func(strVal string) {
			defer wg.Done()

			firstTaskCh := make(chan string, 1)
			defer close(firstTaskCh)
			secondTaskCh := make(chan string, 1)
			defer close(secondTaskCh)

			go func(result chan string) {
				result <- DataSignerCrc32(strVal)
			}(firstTaskCh)

			go func(result chan string) {
				md5Mu.Lock()
				md5Result := DataSignerMd5(strVal)
				md5Mu.Unlock()
				result <- DataSignerCrc32(md5Result)
			}(secondTaskCh)

			result := <-firstTaskCh + "~" + <-secondTaskCh

			out <- result

		}(strVal)
	}
	wg.Wait()
}

const thCount = 6

func MultiHash(in, out chan interface{}) {
	wg := &sync.WaitGroup{}
	for val := range in {
		strVal := fmt.Sprintf("%v", val)
		wg.Add(1)
		go func(strVal string) {
			defer wg.Done()

			wgTh := &sync.WaitGroup{}
			result := [thCount]string{}
			for i := 0; i < thCount; i++ {
				wgTh.Add(1)
				go func(i int) {
					defer wgTh.Done()
					result[i] = DataSignerCrc32(fmt.Sprintf("%v%v", i, strVal))
				}(i)
			}
			wgTh.Wait()
			out <- strings.Join(result[:], "")
		}(strVal)
	}
	wg.Wait()
}

func CombineResults(in, out chan interface{}) {
	data := make([]string, 0)
	for val := range in {
		data = append(data, fmt.Sprintf("%v", val))
	}
	sort.Strings(data)
	out <- strings.Join(data, "_")
}

func ExecutePipeline(jobs ...job) {
	wg := &sync.WaitGroup{}
	in := make(chan interface{})

	for i, jobFunc := range jobs {
		wg.Add(1)
		out := make(chan interface{})
		go func(job job, in, out chan interface{}, i int) {
			defer wg.Done()
			defer close(out)
			job(in, out)
		}(jobFunc, in, out, i)
		in = out
	}
	wg.Wait()
}
