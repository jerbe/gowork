package go_pool_test

import (
	"fmt"
	"sort"
	"sync/atomic"
	"testing"
	"time"
	"unsafe"
)

type C struct {
	q chan bool
}

func (c *C) start() {
	for {
		select {
		case <-c.q:
			fmt.Println("quit")
			return
		default:
			//fmt.Println("default")
		}
	}
}

func (c *C) Start() {
	go func() {
		c.start()
	}()
}

func (c *C) Stop() {
	fmt.Println("c.Stop()")
	c.q <- true
}

func Test_ChanC(t *testing.T) {
	CC := C{
		q: make(chan bool),
	}

	CC.Start()

	time.Sleep(time.Second)

	CC.Stop()

	time.Sleep(time.Hour)

}

func Test_Chan(t *testing.T) {
	var ci = make(chan int32, 100)
	//var cci = make(chan chan int)
	//go func() {
	//	for{
	//		cci <- ci
	//	}
	//}()

	var x = int32(0)
	go func() {
		for {
			i := <-ci
			fmt.Println(i)
		}
	}()

	go func() {
		for i := 0; i < 1000; i++ {
			go func() {
				atomic.AddInt32(&x, 1)
			}()
		}
	}()

	go func() {
		for i := int32(0); i < 2000; i++ {
			go func() {
				atomic.AddInt32(&x, 1)
			}()
		}
	}()

	go func() {
		for i := int32(0); i < 1000; i++ {
			go func() {
				atomic.AddInt32(&x, 1)
			}()
		}
	}()

	go func() {
		for i := 0; i < 4000; i++ {
			y := atomic.LoadInt32(&x)
			fmt.Println(y)
		}
	}()

	time.Sleep(time.Second * 20)
	fmt.Println("done", x)

	time.Sleep(time.Hour)

	//<-ci

}

type user struct {
	name string

	age int

	i int32
}

func Test_Unsafe(t *testing.T) {
	var u user

	uName := (*string)(unsafe.Pointer(&u))
	*uName = "name"

	uAge := (*int)(unsafe.Pointer(uintptr(unsafe.Pointer(&u)) + unsafe.Offsetof((&u).age)))
	*uAge = 10000

	i := (*int32)(unsafe.Pointer(uintptr(unsafe.Pointer(&u)) + unsafe.Offsetof((&u).i)))
	*i = 99

	fmt.Println(u)

}

func Test_Sort(t *testing.T) {
	keys := []string{
		"3-4",
		"1",
		"3-1",
		"2-2",
		"4-2",
		"10-1-2",
		"6-2",
		"3-3",
		"3-5",
		"4-1",
		"10-1-2",
	}

	sort.Strings(keys)

	fmt.Println(keys)

}
