package rpc

import (
	"errors"
	"time"
)

// Response is the message returned from service
type Response struct {
	Message string
}

// Server is the RPC server handler (which functions are available)
type Server struct {
	IsRunning *bool
	Stop      func()
	Start     func()
	Quit      func()
}

// Execute is the method that can be called from the rpc clients
func (s *Server) Execute(req string, res *Response) error {
	if req == "disable" {
		if *s.IsRunning {
			s.Stop()
			res.Message = "disabled"
			return nil
		}
		res.Message = "service is already disabled"
		return nil
	}
	if req == "enable" {
		if !*s.IsRunning {
			go s.Start()
			res.Message = "enabled"
			return nil
		}
		res.Message = "service already enabled/running"
		return nil
	}
	if req == "status" {
		if *s.IsRunning {
			res.Message = "enabled"
			return nil
		}
		res.Message = "disabled"
		return nil
	}
	if req == "quit" {
		go func() {
			time.Sleep(100 * time.Millisecond)
			s.Quit()
		}()
		res.Message = "quit"
		return nil
	}
	res.Message = "invalid"
	return errors.New("invalid command")
}
