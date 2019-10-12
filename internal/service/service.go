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

// create new service instance
func new(c config) (*service, error) {
	d := time.Duration(c.Duration) * time.Second
	s := service{
		Duration:  d,
		OffLevel:  c.OffLevel,
		OnLevel:   c.OnLevel,
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

// Start is starting the service which checks for keyboard / mouse events and sets the backlight accordingly.
func (s *service) Start() {
	s.mut.Lock()
	if s.IsRunning {
		helper.PrintErrorString("already running!")
		return
	}
	s.IsRunning = true
	s.mut.Unlock()
	s.setBacklightValue(0) // dirty hack to fix problems with some models ?!
	s.waitForEvents()
}

// Stop stops the service
func (s *service) Stop() {
	s.done <- true
	s.done <- true
	s.IsRunning = false
	s.mut.Lock()
	s.active = false
	s.mut.Unlock()
}

// Event loop that checks for keyboard / mouse input and switches backlight on/off accordingly
// the core of this whole thing :)
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

// Method that is switching off the backlight (according to OffLevel value)
func (s *service) setBacklightOff() {
	fmt.Println("Switching backlight off")

	err := s.setBacklightValue(s.OffLevel)
	if err != nil {
		helper.PrintError(err)
		return
	}

	s.mut.Lock()
	s.active = false
	s.mut.Unlock()
}

// Method that is switching on the backlight (according to OnLevel value)
func (s *service) setBacklightOn() {
	s.timer.Reset(s.Duration)
	if !s.active {
		fmt.Println("Switching backlight on")

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

// Methods that set's the backlight value via dbus
func (s *service) setBacklightValue(level int) error {
	call := s.obj.Call("org.freedesktop.UPower.KbdBacklight.SetBrightness", 0, level)
	if call.Err != nil {
		return call.Err
	}
	return nil
}
