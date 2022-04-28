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
	q.Q(one, two, three)

	maintest1()

	q.Level = ""

	maintest2()
	q.Level = "pkg1"
	maintest3()
	pkg1.Pkg1_func1()
	q.Level = "test"
	maintest4()
}

func maintest1() {
	maintestvar := "aaa"
	q.Q("test1", maintestvar)
}

func maintest2() {
	q.Q("test2")
}

func maintest3() {
	q.Q("test3")
}

func maintest4() {
	q.Q("test4")
}
