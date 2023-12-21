package main

import (
	"fmt"
	"net/textproto"
	"strconv"
	"strings"
)

type Signal int

const (
	Reload = iota
	Shutdown
	Dump
	Debug
	Halt
	NewCircuit
	ClearCircuit
	Heartbeat
	Dormant
	Active
)

func (s Signal) String() string {
	return [...]string{
		"RELOAD",
		"SHUTDOWN",
		"DUMP",
		"DEBUG",
		"HALT",
		"NEWNYM",
		"CLEARDNSCACHE",
		"HEARTBEAT",
		"DORMANT",
		"ACTIVE",
	}[s]
}

type Controller struct {
	*textproto.Conn
}

func NewController(addr string) (*Controller, error) {
	conn, err := textproto.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}
	return &Controller{conn}, nil
}

func (c *Controller) Signal(signal Signal) error {
	_, _, err := c.makeRequest("SIGNAL" + signal.String())
	if err != nil {
		return err
	}
	return nil
}

func (c *Controller) GetAddress() (string, error) {
	return c.getInfo("address")
}

func (c *Controller) GetBytesRead() (int, error) {
	return c.getInfoInt("traffic/read")
}

func (c *Controller) GetBytesWritten() (int, error) {
	return c.getInfoInt("traffic/written")
}

func (c *Controller) GetVersion() (string, error) {
	return c.getInfo("version")
}

// AuthenticateNone authenticate to controller without password or cookie.
func (c *Controller) AuthenticateNone() error {
	_, _, err := c.makeRequest("AUTHENTICATE")
	if err != nil {
		return err
	}
	return nil
}

func (c *Controller) makeRequest(request string) (int, string, error) {
	id, err := c.Cmd(request)
	if err != nil {
		return 0, "", err
	}
	c.StartResponse(id)
	defer c.EndResponse(id)
	return c.ReadResponse(250)
}

func (c *Controller) getInfo(key string) (string, error) {
	_, msg, err := c.makeRequest("GETINFO " + key)
	if err != nil {
		return "", err
	}
	lines := strings.Split(msg, "\n")
	for _, line := range lines {
		parts := strings.SplitN(line, "=", 2)
		if parts[0] == key {
			return parts[1], nil
		}
	}
	return "", fmt.Errorf(key + " not found")
}

func (c *Controller) getInfoInt(key string) (int, error) {
	s, err := c.getInfo(key)
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(s)
}
