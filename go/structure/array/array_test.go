package array

import (
	"testing"
	"time"
)

// 值类型

// go test array_test.go -v
/**
* 数组初始化
 */
func TestArrayInit(t *testing.T) {
	t.Log(time.Now().UnixNano()/1e6, time.Now().Unix())
	var arr1 [3]int
	arr2 := [4]int{1, 2, 3, 4}
	arr3 := [...]int{1, 2, 3}
	t.Log(arr1, arr2, arr3)
}

/**
* 数组遍历
 */
func TestArrayTravel(t *testing.T) {
	arr1 := [...]int{1, 2, 3, 4}
	for i := 0; i < len(arr1); i++ {
		t.Log(arr1[i])
	}
	for idx, e := range arr1 {
		t.Log(idx, e)
	}
	for _, e := range arr1 {
		t.Log(e)
	}
}

/**
* 数组截取 arr[开始索引（包含）， 结束索引（不包含）]
* 不支持 arr[-1]
 */
func TestArraySection(t *testing.T) {
	arr1 := [...]int{1, 2, 3, 4, 5}
	t.Log(arr1[1:2])
	t.Log(arr1[1:len(arr1)])
	t.Log(arr1[:])
}
