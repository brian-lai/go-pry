package main

import (
	"fmt"

	"github.com/brian-lai/go-pry/pry"
)

func main() {
	testMap := map[string]int{
		"duck": 1,
		"blue": 2,
		"5":    0xDEAD,
	}
	for k, v := range testMap {
		fmt.Println(k, v)
	}
	pry.Pry()
}
