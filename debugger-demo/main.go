package main

import (
	"fmt"
	"time"
)

func add(a ,b int) int
func test(a int) int {
	d :=3
	a = d+3
	return a
}
func main()  {
		a ,b:=1,3
		for i:=0;i<5;i++ {
			a++
		}
		test(a)
		add(a,b)
		go func() {
			a +=1
			b -=1
			time.Sleep(10e9)
			b +=1
		}()
		fmt.Println(a,b)
		time.Sleep(100e9)
}