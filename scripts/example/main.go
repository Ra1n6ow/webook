package main

import (
	"fmt"
	"time"
)

func main() {
	now1 := time.Now()
	now2 := now1.UnixMilli()
	now3 := time.UnixMilli(now2)

	fmt.Println(now1)
	fmt.Println(now2)
	fmt.Println(now3)
}
