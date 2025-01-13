package main

import (
	"fmt"
	"runtime"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	fmt.Println("Hello, Gophers ")

	d := 1
	fmt.Println(&d)

	mi := make(map[string]int)
	mi["tomek"] = 15
	mi["tomek1"] = 25
	mi["tomek2"] = 35
	mi["tomek3"] = 145

	delete(mi, "tomek2")
	for v := range mi {
		fmt.Println(v)
	}
	sem := make(chan int)
	wg.Add(1)
	go func() {
		defer wg.Done()
		sem <- 12
		sem <- 38

		fmt.Println("works")

	}()

	fmt.Println(<-sem)
	fmt.Println(<-sem)

	fmt.Println(runtime.NumCPU())

	f := incrementor()
	fmt.Println(f())
	wg.Wait()
}

func incrementor() func() int {
	x := 0
	return func() int {
		x++
		return x
	}
}
