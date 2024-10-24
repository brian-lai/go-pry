package main

import "github.com/brian-lai/go-pry/pry"

func a(b int) {
	c := 5
	pry.Pry()
	_ = c
}

func main() {
	a(5)
}
