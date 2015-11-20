package main

/*
Run as a tcp port proxy if there are multi datacentors in your
production, Receive the traffic and redirect to tranfer.sh server.

cz-20151119
*/

import (
	"flag"
	"github.com/VividCortex/godaemon"
	"log"
	"log/syslog"
	"os"
	"os/signal"
	"syscall"
	"time"
)

//ignore signal
func waitSignal() {
	var sigChan = make(chan os.Signal, 1)
	signal.Notify(sigChan)
	for sig := range sigChan {
		if sig == syscall.SIGINT || sig == syscall.SIGTERM {
			log.Printf("terminated by signal %v\n", sig)
		} else {
			log.Printf("received signal: %v, ignore\n", sig)
		}
	}
}

//net timeout
const timeout = time.Second * 2

func main() {
	// options
	var bind, backend, logTo string
	var buffer uint
	var daemon bool
	flag.StringVar(&bind, "bind", ":8002", "locate ip and port")
	flag.StringVar(&backend, "backend", "127.0.0.1:8003", "backend server ip and port")
	flag.StringVar(&logTo, "logTo", "stdout", "stdout or syslog")
	flag.UintVar(&buffer, "buffer", 4096, "buffer size")
	flag.BoolVar(&daemon, "daemon", false, "run as daemon process")
	flag.Parse()

	log.SetOutput(os.Stdout)
	if logTo == "syslog" {
		w, err := syslog.New(syslog.LOG_INFO, "portproxy")
		if err != nil {
			log.Fatal(err)
		}
		log.SetOutput(w)
	}

	if daemon == true {
		godaemon.MakeDaemon(&godaemon.DaemonAttr{})
	}

	p := New(bind, backend, uint32(buffer))
	log.Println("portproxy started.")
	go p.Start()
	waitSignal()
}
