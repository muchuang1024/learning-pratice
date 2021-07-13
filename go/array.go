package main

import "fmt"

func main()  {
	res := make([]int, 0)
	result := [][]int{}
	for i := 1; i < 15; i++ {
		res = append(res, i)
		result = append(result, res)
	}

	fmt.Println("res", res)
	fmt.Println("result", result)
}