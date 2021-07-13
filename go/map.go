package main

import "fmt"

func main() {
	m := map[int]int{
		1: 1,
		2: 2,
		3: 3,
		5: 5,
	}
	for k, v := range m {
		fmt.Println(k, v)
	}

	var next = make([]float64, 5)
	next[1] = 1
	next[4] = 3

	fmt.Println(next)

}
