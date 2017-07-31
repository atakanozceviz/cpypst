package model

import "sync"

type Tmp struct {
	sync.Mutex
	tmp string
}

func (t *Tmp) Write(s string) {
	t.Lock()
	t.tmp = s
	t.Unlock()
}
func (t *Tmp) Read() string {
	t.Lock()
	s := t.tmp
	t.Unlock()
	return s
}
