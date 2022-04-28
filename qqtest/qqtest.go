package main

import (
	//"fmt"
	//"os"

	//"flag"

	//"github.com/y0ssar1an/q"
	"github.com/bobbae/q"
	"github.com/bobbae/q/pkg1"

	"github.com/spf13/cobra"
)

func main() {
	//flag.StringVar(&q.O, "qo", "q", "q log output destination")
	//flag.StringVar(&q.P, "qp", "", "q package/function regexp pattern")
	//flag.Parse()
	//qqtest()

	var rootCmd = &cobra.Command{
		Use:  "qqtest",
		Long: `qq demo program`,
		Run: func(cmd *cobra.Command, args []string) {
			qqtest()
		},
	}

	rootCmd.PersistentFlags().StringVarP(&q.O, "qoutput", "O", q.O, "q output destination")
	rootCmd.PersistentFlags().StringVarP(&q.P, "qpattern", "P", q.P, "q regexp pattern")

	rootCmd.Execute() 

}

func qqtest() {
	q.Q(q.P, q.O)
	//q.P = "*"
	//q.O = "stdout"

	q.P = ".*"
	q.O = "stderr"
	q.Q("hello")
	one := 1
	two := 2
	three := 3
	q.Q(one, two, three)

	maintest1()

	q.P = ""

	maintest2()
	q.P = "pkg1.*"
	maintest3()
	pkg1.Pkg1_func1()
	q.P = "test.*"
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
