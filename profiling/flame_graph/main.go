package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"time"
)

func fn1() {
	b := fn2()
	fmt.Println(fn3(b))
}

func fn2() (b []byte) {
	b, _ = json.Marshal(map[string]int{
		"a":  22,
		"bb": 333,
	})
	return
}

func fn3(b []byte) int {
	var m map[string]int
	json.Unmarshal(b, &m)

	if len(m) > 0 {
		return m["a"]
	}
	return 0
}

func main() {
	go func() {
		for {
			fn1()
			time.Sleep(1e8)
		}
	}()

	panic(http.ListenAndServe(":8080", nil))
}
