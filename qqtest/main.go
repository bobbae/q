package main

import (
	"flag"

	"github.com/bobbae/q"
)

func main() {
	flag.Parse()
	q.Q("level", q.Level)
	one := 1
	two := 2
	three := 3
	q.Q(one, two, three)
	maintest1()
}

func maintest1() {
	maintestvar := "aaa"
	q.Q(maintestvar)
}
