package main

import (
	"encoding/binary"
	"fmt"
)

func main() {

	key := [10]byte{}

	// 2个字节 （2 * 8 = 16）
	bytes := make([]byte, 2)
	binary.LittleEndian.PutUint16(bytes, uint16(8))
	for i := 0; i < 2; i++ {
		key[i] = bytes[i]
	}

	fmt.Printf("bytes: %b\n", bytes)

	// 8个字节 （8 * 8 = 64）
	bytes = make([]byte, 8)
	for _, v := range []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10} {
		bytes[v/8] = bytes[v/8] | (1 << uint(v%8))
	}
	fmt.Printf("bytes: %b\n", bytes)

	// 合并
	for i := 0; i < len(bytes); i++ {
		key[i+2] = bytes[i]
	}

	fmt.Printf("key: %b\n", key)

	var str1 string = "111"
	var str2 string = "5e91c4a425e1f75e608b5730"
	fmt.Printf("str1: %b\n", []byte(str1))
	fmt.Printf("str2: %b\n", []byte(str2))
	bytes = make([]byte, 0, 48)
	//bytes = append(bytes, []byte(str1)...)
	//bytes = append(bytes, []byte(str2)...)
	//bytes := Key{}
	count := 0
	for i, b := range []byte(str1) {
		fmt.Println(11, i, b)
		bytes = append(bytes, b)
		count += 1
	}
	fmt.Println(111, count)
	for _, b := range []byte(str2) {
		bytes = append(bytes, b)
	}

	fmt.Printf("bytes: %b, %+v\n", bytes, len(bytes))
}

