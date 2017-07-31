package model

import (
	"sync"
)

type Connection struct {
	Ip     string
	Name   string
	Active bool
}

type Connections struct {
	sync.Mutex
	Connections map[string]*Connection
}

func (c *Connections) Add(con Connection) {
	if c.Connections == nil {
		c.Connections = make(map[string]*Connection)
	}
	c.Lock()
	c.Connections[con.Ip] = &con
	c.Unlock()
}

func (c *Connections) Remove(con Connection) {
	c.Lock()
	delete(c.Connections, con.Ip)
	c.Unlock()
}

func (c *Connections) Disable(con Connection) {
	c.Lock()
	c.Connections[con.Ip].Active = false
	c.Unlock()
}

func (c *Connections) Enable(con Connection) {
	c.Lock()
	c.Connections[con.Ip].Active = true
	c.Unlock()
}
