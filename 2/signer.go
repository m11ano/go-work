package main

import "C"
import (
	"fmt"
	"sync"
)

func SingleHash(in, out chan interface{}) {
	for val := range in {
		strVal := fmt.Sprintf("%v", val)
		fmt.Println(DataSignerCrc32(strVal))
		out <- strVal
	}
}

func MultiHash(in, out chan interface{}) {
	for val := range in {
		out <- val
	}
}

func CombineResults(in, out chan interface{}) {
	for _ = range in {
		//
	}
	out <- "TEST"
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
