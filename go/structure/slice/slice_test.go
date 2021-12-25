package slice

import (
	"fmt"
	"sync"
	"testing"
	"unsafe"
)

// 引用类型

// go test slice_test.go -v
/**
* 切片初始化
 */
func TestSliceInit(t *testing.T) {
	// 初始化方式1：直接声明
	var slice1 []int
	t.Log(len(slice1), cap(slice1)) // 0, 0
	slice1 = append(slice1, 1)

	t.Log(len(slice1), cap(slice1), unsafe.Sizeof(slice1)) // 1, 1, 24

	// 初始化方式2：使用字面量
	slice2 := []int{1, 2, 3, 4}
	t.Log(len(slice2), cap(slice2), unsafe.Sizeof(slice2)) // 4, 4, 24

	// 初始化方式3：使用make创建slice
	slice3 := make([]int, 3, 5)            // make([]T, len, cap) cap不传则和len一样
	t.Log(len(slice3), cap(slice3))        // 3, 5
	t.Log(slice3[0], slice3[1], slice3[2]) // 0, 0, 0
	// t.Log(slice3[3], slice3[4]) // panic: runtime error: index out of range [3] with length 3
	slice3 = append(slice3, 1)
	t.Log(len(slice3), cap(slice3), unsafe.Sizeof(slice3)) // 4, 5, 24

	// 初始化方式4: 从切片或数组“截取”
	arr := [100]int{}
	for i := range arr {
		arr[i] = i
	}
	slcie4 := arr[1:3]
	// 如果我们只用到一个slice的一小部分，那么底层的整个数组也将继续保存在内存当中。当这个底层数组很大，或者这样的场景很多时，可能会造成内存急剧增加，造成崩溃。 所以在这样的场景下，我们可以将需要的分片复制到一个新的slice中去，减少内存的占用。例如一个很大的切片data里，我们需要的数据是data[m:n]，那么我们创建一个新的slice变量r，将数据复制到r中返回
	slice5 := make([]int, len(slcie4))
	copy(slice5, slcie4)
	t.Log(len(slcie4), cap(slcie4), unsafe.Sizeof(slcie4)) // 2，99，24
	t.Log(len(slice5), cap(slice5), unsafe.Sizeof(slice5)) // 2，2，24
}

/**
* 浅拷贝
 */
func TestSliceShadowCopy(t *testing.T) {
	slice1 := []int{1, 2, 3, 4, 5}
	slice2 := slice1     // 浅拷贝（注意 := 对于引用类型是浅拷贝，对于值类型是深拷贝）
	t.Logf("%p", slice1) // 0xc00001c120
	t.Logf("%p", slice2) // 0xc00001c120
	// 同时改变两个数组，这时就是浅拷贝，未扩容时，修改 slice1 的元素之后，slice2 的元素也会跟着修改
	slice1[0] = 10
	t.Log(slice1, len(slice1), cap(slice1)) // [10 2 3 4 5] 5 5
	t.Log(slice2, len(slice2), cap(slice2)) // [10 2 3 4 5] 5 5
	// 注意下：扩容后，slice1和slice2不再指向同一个数组，修改 slice1 的元素之后，slice2 的元素不会被修改了
	slice1 = append(slice1, 5, 6, 7, 8)
	slice1[0] = 11                          // 这里可以发现，slice1[0] 被修改为了 11, slice1[0] 还是10
	t.Log(slice1, len(slice1), cap(slice1)) // [11 2 3 4 5 5 6 7 8] 9 10
	t.Log(slice2, len(slice2), cap(slice2)) // [10 2 3 4 5] 5 5
}

/**
* 深拷贝
 */
func TestSliceDeepCopy(t *testing.T) {
	slice1 := []int{1, 2, 3, 4, 5}
	slice2 := make([]int, 5, 5)
	copy(slice2, slice1)                    // 深拷贝
	t.Log(slice1, len(slice1), cap(slice1)) // [1 2 3 4 5] 5 5
	t.Log(slice2, len(slice2), cap(slice2)) // [1 2 3 4 5] 5 5
	slice1[1] = 100                         //只会影响slice1
	t.Log(slice1, len(slice1), cap(slice1)) // [1 100 3 4 5] 5 5
	t.Log(slice2, len(slice2), cap(slice2)) // [1 2 3 4 5] 5 5
}

/**
* 切片截取
 */
func TestSliceSubstr(t *testing.T) {
	slice1 := []int{1, 2, 3, 4, 5}
	slice2 := slice1[:]
	// 截取 slice[left:right:max]
	// left：省略默认0；right：省略默认len(slice1)；max: 省略默认len(slice1)
	// len = right-left+1
	// cap = max-left
	t.Log(slice2, len(slice2), cap(slice2)) // 1 2 3 4 5] 5 5
	slice3 := slice1[1:]
	t.Log(slice3, len(slice3), cap(slice3)) // [2 3 4 5] 4 4
	slice4 := slice1[:2]
	t.Log(slice4, len(slice4), cap(slice4)) // [1 2] 2 5
	slice5 := slice1[1:2]
	t.Log(slice5, len(slice5), cap(slice5)) // [2] 1 4
	slice6 := slice1[:2:5]
	t.Log(slice6, len(slice6), cap(slice6)) // [1 2] 2 5
	slice7 := slice1[1:2:2]
	t.Log(slice7, len(slice7), cap(slice7)) // [2] 1 1
}

func TestSliceEmptyOrNil(t *testing.T) {
	var slice1 []int            // slice1 is nil slice
	slice2 := make([]int, 0)    // slcie2 is empty slice
	var slice3 = make([]int, 2) // slice3 is zero slice
	if slice1 == nil {
		t.Log("slice1 is nil.") // 会输出这行
	}
	if slice2 == nil {
		t.Log("slice2 is nil.") // 不会输出这行
	}
	t.Log(slice3) // [0 0]
}

/**
* 切片遍历
 */
func TestSliceTravel(t *testing.T) {
	slice1 := []int{1, 2, 3, 4}
	for i := 0; i < len(slice1); i++ {
		t.Log(slice1[i])
	}
	for idx, e := range slice1 {
		t.Log(idx, e)
	}
	for _, e := range slice1 {
		t.Log(e)
	}
}

/**
* 切片增长
 */
func TestSliceGrowing(t *testing.T) {
	slice1 := []int{}
	for i := 0; i < 10; i++ {
		slice1 = append(slice1, i)
		t.Log(len(slice1), cap(slice1))
	}
	// 1 1
	// 2 2
	// 3 4
	// 4 4
	// 5 8
	// 6 8
	// 7 8
	// 8 8
	// 9 16
	// 10 16
}

func TestSliceDelete(t *testing.T) {
	slice1 := []int{1, 2, 3, 4, 5}
	var x int
	x, slice1 = slice1[len(slice1)-1], slice1[:len(slice1)-1]
	t.Log(x, slice1, len(slice1), cap(slice1)) // 5 [1 2 3 4] 4 5

	slice1 = append(slice1[:2], slice1[3:]...) // 删除第2个元素
	t.Log(slice1, len(slice1), cap(slice1))    // [1 2 4] 3 5
}

func TestSliceReverse(t *testing.T) {
	a := []int{1, 2, 3, 4, 5}
	for left, right := 0, len(a)-1; left < right; left, right = left+1, right-1 {
		a[left], a[right] = a[right], a[left]
	}
	t.Log(a, len(a), cap(a)) // [5 4 3 2 1] 5 5
}

/**
* 切片共享存储空间
* 多个切片可能共享同一个底层数组，这种情况下，对其中一个切片或者底层数组的更改，会影响到其他切片
 */
func TestSliceShareMemory(t *testing.T) {
	slice1 := []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12"}
	Q2 := slice1[3:6]
	t.Log(Q2, len(Q2), cap(Q2)) // [4 5 6] 3 9
	Q3 := slice1[5:8]
	t.Log(Q3, len(Q3), cap(Q3)) // [6 7 8] 3 7
	Q3[0] = "Unkown"
	t.Log(Q2, Q3) // [4 5 Unkown] [Unkown 7 8]

	a := []int{1, 2, 3, 4, 5}
	shadow := a[1:3]
	t.Log(shadow, a)             // [2 3] [1 2 3 4 5]
	shadow = append(shadow, 100) // 会修改指向数组的所有切片
	t.Log(shadow, a)             // [2 3 100] [1 2 3 100 5]
}

/**
* 切片非并发安全
* 多次执行，每次得到的结果都不一样
* 可以考虑使用 channel 本身的特性 (阻塞) 来实现安全的并发读写
 */
func TestSliceConcurrencySafe(t *testing.T) {
	a := make([]int, 0)
	var wg sync.WaitGroup
	for i := 0; i < 10000; i++ {
		wg.Add(1)
		go func(i int) {
			a = append(a, i)
			wg.Done()
		}(i)
	}
	wg.Wait()
	t.Log(len(a)) // not equal 10000
}

func TestSliceConcurrencySafeByChanel(t *testing.T) {
	buffer := make(chan int)
	a := make([]int, 0)
	// 消费者
	go func() {
		for v := range buffer {
			a = append(a, v)
		}
	}()
	// 生产者
	var wg sync.WaitGroup
	for i := 0; i < 10000; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			buffer <- i
		}(i)
	}
	wg.Wait()
	t.Log(len(a)) // 10000
}

func TestSliceConcurrencySafeByMutex(t *testing.T) {
	var lock sync.Mutex //互斥锁
	a := make([]int, 0)
	var wg sync.WaitGroup
	for i := 0; i < 10000; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			lock.Lock()
			defer lock.Unlock()
			a = append(a, i)
		}(i)
	}
	wg.Wait()
	t.Log(len(a)) // equal 10000
}

func sliceModify(s []int) {
	s[0] = 100
}

func sliceAppend(s []int) []int {
	// 这里 s 虽然改变了，但并不会影响外层函数的 s
	s = append(s, 100)
	// 内层slice的len,cap，但不影响外层slice的len,cap
	// 如果未触发扩容，内存slice指向的数组会被改变，会影响外层slice
	// 如果触发扩容，内存slice指向的数组不会被改变，不影响外层slice
	fmt.Println(s, len(s), cap(s))
	return s
}

func sliceAppendPtr(s *[]int) {
	// 会改变外层 s 本身
	*s = append(*s, 100)
	return
}

// 注意：Go语言中所有的传参都是值传递（传值），都是一个副本，一个拷贝。
// 拷贝的内容是非引用类型（int、string、struct等这些），在函数中就无法修改原内容数据；
// 拷贝的内容是引用类型（interface、指针、map、slice、chan等这些），这样就可以修改原内容数据。
func TestSliceFn(t *testing.T) {
	// 参数为引用类型slice：外层slice的len/cap不会改变，指向的底层数组会改变
	s := []int{1, 1, 1}
	newS := sliceAppend(s)            // 函数内发生了扩容
	t.Log(s, len(s), cap(s))          // [1 1 1] 3 3
	t.Log(newS, len(newS), cap(newS)) // [1 1 1 100] 4 6

	s2 := make([]int, 0, 5)
	newS = sliceAppend(s2)                       // 函数内未发生扩容
	t.Log(s2, s2[0:5], len(s2), cap(s2))         // [] [100 0 0 0 0] 0 5
	t.Log(newS, newS[0:5], len(newS), cap(newS)) // [100] [100 0 0 0 0] 1 5

	// 参数为引用类型slice的指针：外层slice的len/cap会改变，指向的底层数组会改变
	sliceAppendPtr(&s)
	t.Log(s, len(s), cap(s)) // [1 1 1 100] 4 6
	sliceModify(s)
	t.Log(s, len(s), cap(s)) // [100 1 1 100] 4 6
}

func BenchmarkSliceConcurrencySafeByMutex(b *testing.B) {
	var lock sync.Mutex //互斥锁
	a := make([]int, 0)
	var wg sync.WaitGroup
	for i := 0; i < b.N; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			lock.Lock()
			defer lock.Unlock()
			a = append(a, i)
		}(i)
	}
	wg.Wait()
}

//  go test -bench='Benchmark$' -cpu=2,4 -benchtime=5s .
func BenchmarkSliceConcurrencySafeByChanel(b *testing.B) {
	buffer := make(chan int)
	a := make([]int, 0)
	go func() {
		for v := range buffer {
			a = append(a, v)
		}
	}()
	for i := 0; i < b.N; i++ {
		buffer <- i
	}
}
