package main

import (
	"fmt"
	//"os"
	"log"
	//log "github.com/sirupsen/logrus"

	"context"
	"strings"

	//"flag"

	//"github.com/y0ssar1an/q"

	"github.com/bobbae/q"
	"github.com/bobbae/q/pkg1"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type Config struct {
	username   string
	password   string
	configFile string
}

var qC Config

func main() {
	q.O = "stderr"
	q.P = ".*"

	//flag.StringVar(&q.O, "qo", "q", "q log output destination")
	//flag.StringVar(&q.P, "qp", "", "q package/function regexp pattern")
	//flag.Parse()
	//qqtest()

	ctx := context.WithValue(context.Background(), "config", &qC)

	var rootCmd = &cobra.Command{
		Use:  "qqtest",
		Long: `qq demo program`,
		Run: func(cmd *cobra.Command, args []string) {
			qqtest(ctx)
		},
	}

	rootCmd.PersistentFlags().StringVarP(&q.O, "qoutput", "O", q.O, "q output destination")
	rootCmd.PersistentFlags().StringVarP(&q.P, "qpattern", "P", q.P, "q regexp pattern")
	rootCmd.PersistentFlags().StringVarP(&qC.username, "username", "u", "", "username")
	rootCmd.PersistentFlags().StringVarP(&qC.password, "password", "p", "", "password")
	rootCmd.PersistentFlags().StringVarP(&qC.configFile, "configFile", "c", "", "configuration file")

	viper.BindPFlag("username", rootCmd.PersistentFlags().Lookup("username"))
	viper.BindPFlag("password", rootCmd.PersistentFlags().Lookup("password"))

	var cmdDemo = &cobra.Command{
		Use:   "demo [message string]",
		Short: "Demo Post Msg",
		Run: func(cmd *cobra.Command, args []string) {
			demo1(strings.Join(args, " "))
		},
	}

	var cmdTest = &cobra.Command{
		Use:   "test <topic> <msg>",
		Short: "Test Post Msg",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 2 {
				log.Fatalf("Error: missing topic or message")
			}
			test1(args[0], args[1])
		},
	}

	var cmdEcho = &cobra.Command{
		Use:   "echo [string to echo]",
		Short: "Echo anything to the screen",
		Long: `echo is for echoing anything back.
Echo works a lot like print, except it has a child command.`,
		Args: cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Echo: " + strings.Join(args, " "))
		},
	}

	var echoTimes int

	var cmdTimes = &cobra.Command{
		Use:   "times [string to echo]",
		Short: "Echo anything to the screen more times",
		Long: `echo things multiple times back to the user by providing
a count and a string.`,
		Args: cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			for i := 0; i < echoTimes; i++ {
				fmt.Println("Echo: " + strings.Join(args, " "))
			}
		},
	}

	cmdTimes.Flags().IntVarP(&echoTimes, "times", "t", 1, "times to echo the input")



	rootCmd.AddCommand(cmdDemo, cmdTest)
	cmdDemo.AddCommand(cmdEcho)
	cmdEcho.AddCommand(cmdTimes)


	rootCmd.Execute()
}

func qqtest(ctx context.Context) {
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

	q.P = "" // off all q.Q()

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

/*
func getHome() string {
	homeDir := os.Getenv("HOME")
	if homeDir == "" {
		homeDir = os.Getenv("USERPROFILE")
	}
	return homeDir
}
*/

func test1(topic, message string) {
	q.Q(topic, message)
}

func demo1(message string) {
	q.Q(message, qC.username, qC.password)
	q.Q(qC.configFile)

	if qC.configFile != "" {
		viper.SetConfigFile(qC.configFile)
	}

	/*
		viper.SetConfigName("config")
		viper.SetConfigType("json")

		 homeDir := getHome()
		viper.AddConfigPath(homeDir + "/.exctl")
		viper.AddConfigPath(".")
	*/

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}

	q.Q(viper.GetString("username"), viper.GetString("password"), viper.GetString("tag1.tag2.tag3"))
	log.Printf("username: %s, password: %s", qC.username, qC.password)
	q.Q(viper.GetString("a.b"))
}
