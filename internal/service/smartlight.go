package service

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	netrpc "net/rpc"
	"os"
	"strings"

	"github.com/moson-mo/smartlight/internal/rpc"
)

func Run() {
	// setup and parse arguments
	d := flag.Uint("d", 10, "duration after disabling the keyboard backlight in seconds")
	loff := flag.Int("loff", 0, "brightness level when switching backlight off (default 0)")
	lon := flag.Int("lon", 1, "brightness level when switching backlight on")
	v := flag.Bool("v", true, "be verbose")
	flag.Parse()

	// create new service
	s, err := new(*d, *loff, *lon, *v)
	if err != nil {
		fmt.Println(err)
		return
	}

	// setup & start rpc server
	h := rpc.Server{
		Start: s.Start,
		Stop:  s.Stop,
		Quit: func() {
			s.Stop()
			os.Exit(0)
		},
		IsRunning: &s.IsRunning,
	}
	startRPCServer(&h)

	// start service
	go s.Start()

	// control service
	reader := bufio.NewReader(os.Stdin)
	for {
		input, _ := reader.ReadString('\n')
		input = strings.Replace(input, "\n", "", -1)
		if input == "stop" {
			s.Stop()
		}
		if input == "start" {
			go s.Start()
		}
		if input == "quit" {
			break
		}
	}
}

func startRPCServer(srv *rpc.Server) error {
	sock, err := net.Listen("tcp", ":31987")
	if err != nil {
		return err
	}

	netrpc.RegisterName("smartlight", srv)
	go netrpc.Accept(sock)
	return nil
}
