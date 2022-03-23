package main

import (
	"flag"
	"fmt"
	"os"

	//	"os/signal"
	"sync"
	//	"syscall"

	"github.com/CosmicBDry/prometheus/configAgent/options"
	"github.com/CosmicBDry/prometheus/configAgent/task"
)

func main() {

	var (
		h, help                              bool
		RemoteServer, server, promtool, conf string
		Wg                                   sync.WaitGroup
	)

	flag.BoolVar(&help, "help", false, "help manual")
	flag.BoolVar(&h, "h", false, "help manual")
	flag.StringVar(&RemoteServer, "RemoteServer", "http://localhost:9090", "Get RemoteServer Config Data")
	flag.StringVar(&conf, "conf", "/opt/prometheus/prometheus.yml", "Config of Local Promtheus-Server ")
	flag.StringVar(&promtool, "promtool", "/opt/prometheus/promtool", "The path of promtool")

	flag.Usage = func() {
		fmt.Println("Usage: ./configAgent [-RemoteServer http://localhost:9090] [-conf /opt/prometheus/prometheus.yml]" +
			" " + "[-promtool /opt/prometheus/promtool]")
		flag.PrintDefaults()
	}
	flag.Parse()

	if h || help {
		flag.Usage()
		os.Exit(0)
	}
	//server = RemoteServer
	server = "http://192.168.1.13:8001"
	options := options.NewOption(server, promtool, conf)
	Register := task.NewRegisterTask(options)
	Heartbeat := task.NewHeartBeatTask(options)
	ConfigTask := task.NewConfigTask(options)

	//通过Wg优雅的终止工作例程-------------------------------------------------------------->
	Wg.Add(3) //添加三个工作例程，即将被监听的工作例程
	go func() {
		Register.Run()
		Wg.Done()
	}()
	go func() {
		Heartbeat.Run()
		Wg.Done()
	}()
	go func() {
		ConfigTask.Run()
		Wg.Done()
	}()

	Wg.Wait()

	//也可以通过signal.Notify监听工作例程的系统进程终止信号，来防止工作例程未执行-------------------------------------------------------------->
	//Interrupt := make(chan os.Signal, 1)
	//signal.Notify(Interrupt, syscall.SIGINT, syscall.SIGTERM) //syscall.SIGINT类似于ctrl+c(-2)，syscall.SIGTERM(-15)为终止默认的程序
	//<-Interrupt

	//总结：Waitgroup和signal.Notify选择其中一个即可-------------------------------------------------------------->
}
