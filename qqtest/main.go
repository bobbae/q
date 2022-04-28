package main

import (
	//"flag"

	//"github.com/y0ssar1an/q"
	"github.com/bobbae/q"
	"github.com/bobbae/q/pkg1"
)

func main() {
	//flag.Parse()
	//q.Level = "all"
	//q.Output = ""
	q.Level = "all"
	q.Output = "stderr"
	q.Q("hello")
	//q.Q("level", q.Level)
	one := 1
	two := 2
	three := 3
	//q.Level = ""
	q.Q(one, two, three)
	maintest1()
	maintest2()
	q.Level = "main.main"
	q.Q("reset level to all")
	maintest1()
	maintest2()
	pkg1.pkg1_func1()
}

func maintest1() {
	maintestvar := "aaa"
	q.Q("test1", maintestvar)
}

func maintest2() {
	q.Q("test2")
}
