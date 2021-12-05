package pool

import (
	"github.com/pterm/pterm"
)

func worker(number int, jobs <-chan func() interface{}, results chan<- interface{}) {
	pterm.Info.Println("goroutine worker ", number, "started")
	for j := range jobs {
		results <- j()
	}
	pterm.Info.Println("goroutine worker ", number, "ended")
}

func StartPool(fs Jobs, poolSize int) (v []interface{}) {
	jobLength := len(fs)
	jobs := make(chan func() interface{}, jobLength)
	res := make(chan interface{}, jobLength)

	for w := 0; w < poolSize; w++ {
		go worker(w, jobs, res)
	}

	for j := 0; j < jobLength; j++ {
		jobs <- fs[j]
	}
	close(jobs) // jobs 에 더이상 쌓인게 없으면 worker 가 종료된다.

	for a := 0; a < jobLength; a++ {
		v = append(v, <-res)
	}
	return v

}

type Jobs []func() interface{}
