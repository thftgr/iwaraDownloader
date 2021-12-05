package pool

func worker(jobs <-chan func() interface{}, results chan<- interface{}) {
	for j := range jobs {
		results <- j()
	}
}

func StartPool(fs Jobs, poolSize int) (v []interface{}) {
	jobLength := len(fs)
	jobs := make(chan func() interface{}, jobLength)
	res := make(chan interface{}, jobLength)

	for w := 0; w < poolSize; w++ {
		go worker(jobs, res)
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

//func main() {
//	var o []func() interface{}
//	for i := 0; i < 10; i++ {
//		ii := i
//		o = append(o, func() interface{} {
//			time.Sleep(time.Second)
//			fmt.Println(ii * 2)
//			return ii * 2
//		})
//	}
//	fmt.Println(StartPool(o, 10))
//
//}
