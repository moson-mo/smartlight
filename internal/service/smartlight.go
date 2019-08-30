package service

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	netrpc "net/rpc"
	"os"
	"strings"

	"github.com/moson-mo/smartlight/internal/helper"

	"github.com/moson-mo/smartlight/internal/rpc"
)

// Run loads configuration (from file or arguments), start's rpc server and service
func Run() {
	// load config (if available)
	c, err := loadConfig()
	if err != nil {
		fmt.Println("Problem loading config file: " + err.Error())
		fmt.Println("Trying to save default config file...")
		nc := config{
			Duration: 15,
			OffLevel: 0,
			OnLevel:  1,
			Comments: "Duration = timeout in seconds after switching to 'off'; OffLevel = brightness level when switched off; OnLevel = your guessed it!?!",
		}
		c = &nc
		err = createConfig(nc)
		if err != nil {
			helper.PrintError(err)
		}
	}

	// setup and parse arguments
	d := flag.Uint("d", 15, "duration after disabling the keyboard backlight in seconds")
	loff := flag.Int("loff", 0, "brightness level when switching backlight off (default 0)")
	lon := flag.Int("lon", 1, "brightness level when switching backlight on")
	flag.Parse()

	// check if cmdline args have been given and override config if so
	if isFlagPassed("d") {
		c.Duration = *d
	}
	if isFlagPassed("loff") {
		c.OffLevel = *loff
	}
	if isFlagPassed("lon") {
		c.OnLevel = *lon
	}

	// create new service
	s, err := new(*c)
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
	err = startRPCServer(&h)
	if err != nil {
		helper.PrintError(err)
		return
	}

	// start service
	go s.Start()

	// control service (actually this was meant for debugging purposes, but ok, let's keep it :) )
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

// start the rpc server which is accessed by the tray and cli applications
func startRPCServer(srv *rpc.Server) error {
	sock, err := net.Listen("tcp", ":31987")
	if err != nil {
		return err
	}

	netrpc.RegisterName("smartlight", srv)
	go netrpc.Accept(sock)
	return nil
}

// helper function which is checking if a flag has been entered as cmdline args. Thanks stackoverflow :)
func isFlagPassed(name string) bool {
	found := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == name {
			found = true
		}
	})
	return found
}
