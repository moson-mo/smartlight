package service

import (
	"fmt"
	"sync"
	"time"

	"github.com/godbus/dbus"
	"github.com/moson-mo/smartlight/internal/helper"
	hook "github.com/robotn/gohook"
)

type service struct {
	bus    *dbus.Conn
	obj    dbus.BusObject
	timer  *time.Timer
	done   chan bool
	events chan hook.Event
	mut    *sync.Mutex
	active bool

	Duration  time.Duration
	OffLevel  int
	OnLevel   int
	Verbose   bool
	IsRunning bool
}

func new(duration uint, offLevel, onLevel int, verbose bool) (*service, error) {
	d := time.Duration(duration) * time.Second
	fmt.Println(d)
	s := service{
		Duration:  d,
		OffLevel:  offLevel,
		OnLevel:   onLevel,
		Verbose:   verbose,
		IsRunning: false,

		done:   make(chan bool),
		timer:  time.NewTimer(d),
		mut:    &sync.Mutex{},
		active: false,
	}

	var err error
	s.bus, err = dbus.SystemBus()
	if err != nil {
		return nil, err
	}
	s.obj = s.bus.Object("org.freedesktop.UPower", "/org/freedesktop/UPower/KbdBacklight")

	return &s, nil
}

func (s *service) Start() {
	if s.IsRunning {
		helper.PrintErrorString("already running!")
		return
	}
	s.IsRunning = true
	s.setBacklightValue(0) // dirty hack to fix problems with some models ?!
	s.waitForEvents()
}

func (s *service) Stop() {
	s.done <- true
	s.done <- true
	s.IsRunning = false
	s.mut.Lock()
	s.active = false
	s.mut.Unlock()
}

func (s *service) waitForEvents() {
	s.events = hook.Start()
	defer func() {
		s.timer.Stop()
		hook.End()
		fmt.Println("stopped")
	}()
	wg := sync.WaitGroup{}
	wg.Add(1)
	defer wg.Wait()
	go func() {
		defer wg.Done()
		for {
			select {
			case <-s.done:
				return
			case <-s.timer.C:
				s.setBacklightOff()
			}
		}
	}()

	for {
		select {
		case <-s.done:
			return
		case <-s.events:
			s.setBacklightOn()
		}
	}
}

func (s *service) setBacklightOff() {
	if s.Verbose {
		fmt.Println("Switching backlight off")
	}

	err := s.setBacklightValue(s.OffLevel)
	if err != nil {
		helper.PrintError(err)
		return
	}

	s.mut.Lock()
	s.active = false
	s.mut.Unlock()
}

func (s *service) setBacklightOn() {
	s.timer.Reset(s.Duration)
	if !s.active {
		if s.Verbose {
			fmt.Println("Switching backlight on")
		}

		err := s.setBacklightValue(s.OnLevel)
		if err != nil {
			helper.PrintError(err)
			return
		}

		s.mut.Lock()
		s.active = true
		s.mut.Unlock()
	}
}

func (s *service) setBacklightValue(level int) error {
	call := s.obj.Call("org.freedesktop.UPower.KbdBacklight.SetBrightness", 0, level)
	if call.Err != nil {
		return call.Err
	}
	return nil
}
