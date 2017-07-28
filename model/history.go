package model

import (
	"sync"
)

type HistItem struct {
	Ip      string
	Content string
}

type History struct {
	sync.Mutex
	History []HistItem
}

func (h *History) Add(item HistItem) {
	h.Lock()
	h.History = append(h.History, item)
	h.Unlock()
}

func (h *History) Remove(i int) {
	h.Lock()
	h.History = append(h.History[:i], h.History[i+1:]...)
	h.Unlock()
}
